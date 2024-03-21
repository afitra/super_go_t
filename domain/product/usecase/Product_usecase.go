package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"strconv"
	"superindo/v1/connection"
	"superindo/v1/domain/product"
	"superindo/v1/helper"
	"superindo/v1/model"
	"superindo/v1/response"
	"time"
)

type Product_usecase struct {
	source             string
	product_repository product.Repository
	redis_cache        connection.Redis_cache
	redis_client       *redis.Client
}

func NewProduct_usecase(product_repository product.Repository, redis connection.Redis_cache, redis_con *redis.Client) product.Usecase {
	return &Product_usecase{"Product_usecase", product_repository, redis, redis_con}
}

func (IN *Product_usecase) Use_post_product_register(c echo.Context, request model.Req_register_product) (interface{}, error) {
	var resp interface{}
	var err error
	var query_processing model.Query_proccessing

	if query_processing.Tx, err = IN.product_repository.Rep_begin_transaction(); err != nil {
		resp := response.Set_error_response(err, response.Code_error_general, false, response.Message_error_general)
		return resp, err
	}

	payload := set_payload_product_insert(request)

	if query_processing.Tx, err = IN.product_repository.Rep_insert_product(query_processing.Tx, payload); err != nil {
		resp := response.Set_error_response(err, response.Code_error_general, false, response.Message_error_general)
		return resp, err
	}

	if err = IN.product_repository.Rep_commit_transaction(query_processing.Tx); err != nil {
		resp := response.Set_error_response(err, response.Code_error_general, false, response.Message_error_general)
		return resp, err
	}

	resp = response.Set_success_response(response.Code_success_general, true, response.Message_success_data_added, nil)
	return resp, err
}

func set_payload_product_insert(request model.Req_register_product) model.Product {
	var payload model.Product
	payload.Name = request.Name
	payload.Product_code = helper.Generate_code("PR", 8)
	payload.Product_type = request.Product_type
	payload.Description = request.Description
	payload.Price = request.Price
	return payload
}

func (IN *Product_usecase) Use_get_product_list(c echo.Context) (interface{}, error) {
	var resp interface{}
	var err error
	var data []model.Product
	var offset int
	var limit int

	if offset, err = strconv.Atoi(c.QueryParam("offset")); err != nil {
		resp := response.Set_error_response(err, response.Code_error_general, false, response.Message_error_parameter)
		return resp, err
	}

	if limit, err = strconv.Atoi(c.QueryParam("limit")); err != nil {
		resp := response.Set_error_response(err, response.Code_error_general, false, response.Message_error_parameter)
		return resp, err
	}

	if data, err = IN.product_repository.Rep_get_product_list(offset, limit); err != nil {
		resp := response.Set_error_response(err, response.Code_error_general, false, response.Message_error_general)
		return resp, err
	}

	resp = response.Set_success_response(response.Code_success_general, true, response.Message_success_general, data)
	return resp, err
}

func (IN *Product_usecase) Use_get_product_search(c echo.Context) (interface{}, error) {
	var resp interface{}
	var err error
	var data model.Product

	if data, err = IN.product_repository.Rep_get_product_search_name_or_product_code(c.QueryParam("data")); err != nil {
		resp := response.Set_error_response(err, response.Code_error_general, false, response.Message_success_data_not_found)
		return resp, err
	}

	resp = response.Set_success_response(response.Code_success_general, true, response.Message_success_general, data)
	return resp, err
}

func (IN *Product_usecase) Use_get_product_filter(c echo.Context) (interface{}, error) {
	var resp interface{}
	var err error
	var data []model.Product
	var offset int
	var limit int
	ctx := context.Background()
	var cache_result interface{}
	var products []model.Product

	if offset, err = strconv.Atoi(c.QueryParam("offset")); err != nil {
		resp := response.Set_error_response(err, response.Code_error_general, false, response.Message_error_parameter)
		return resp, err
	}

	if limit, err = strconv.Atoi(c.QueryParam("limit")); err != nil {
		resp := response.Set_error_response(err, response.Code_error_general, false, response.Message_error_parameter)
		return resp, err
	}

	cache_key := fmt.Sprintf("%s_%d_%d", c.Param("product_type"), offset, limit)

	if cache_result, err = IN.redis_cache.Redis_get_key(ctx, cache_key); err != nil {
		resp := response.Set_error_response(err, response.Code_error_general, false, response.Message_error_general)
		return resp, err
	}

	if cache_result != nil {
		products, err = mapping_array_product_from_redis(cache_result)
	}

	if err != nil {
		resp := response.Set_error_response(err, response.Code_error_general, false, response.Message_error_general)
		return resp, err
	}

	if len(products) > 0 {
		resp = response.Set_success_response(response.Code_success_general, true, response.Message_success_general, products)
		return resp, err
	}

	if data, err = IN.product_repository.Rep_get_product_filter_by_product_code(c.Param("product_type"), offset, limit); err != nil {
		resp := response.Set_error_response(err, response.Code_error_general, false, response.Message_success_data_not_found)
		return resp, err
	}

	// Save data to cache
	if err = IN.set_data_array_product_to_redis(cache_key, data); err != nil {
		resp := response.Set_error_response(err, response.Code_error_general, false, response.Message_success_data_not_found)
		return resp, err
	}

	resp = response.Set_success_response(response.Code_success_general, true, response.Message_success_general, data)
	return resp, err
}

func (IN *Product_usecase) set_data_array_product_to_redis(cache_key string, data []model.Product) error {
	var err error
	var data_byte []byte
	ctx := context.Background()
	if data_byte, err = json.Marshal(data); err != nil {
		return err
	}

	var payload model.Chache_model
	payload.Key = cache_key
	payload.Data = data_byte
	payload.Expired = 30 * time.Minute

	if err = IN.redis_cache.Redis_set_key(ctx, payload); err != nil {
		return err
	}

	return nil
}

func mapping_array_product_from_redis(cache_result interface{}) ([]model.Product, error) {
	var err error

	cachedData, ok := cache_result.([]interface{})
	if !ok {
		err = errors.New(response.Message_error_something_wrong)
		return nil, err
	}

	var products []model.Product
	for _, item := range cachedData {
		productMap, ok := item.(map[string]interface{})
		if !ok {
			err = errors.New(response.Message_error_something_wrong)
			return nil, err
		}

		priceFloat, ok := productMap["price"].(float64)
		if !ok {
			err = errors.New(response.Message_error_something_wrong)
			return nil, err
		}
		price := int(priceFloat)

		product := model.Product{

			Name:          productMap["name"].(string),
			Product_code:  productMap["product_code"].(string),
			Product_type:  productMap["product_type"].(string),
			Description:   productMap["description"].(string),
			Price:         price,
			Register_date: productMap["register_date"].(string),
		}

		// Append product to products slice
		products = append(products, product)
	}
	return products, err
}

func (IN *Product_usecase) Use_get_product_sort(c echo.Context) (interface{}, error) {
	var resp interface{}
	var err error
	var data []model.Product
	var offset int
	var limit int
	var key = c.Param("key")
	var validate_param bool

	switch key {
	case "name":
		validate_param = true
		break
	case "price":
		validate_param = true
		break
	case "register_date":
		validate_param = true
		break
	default:
		validate_param = false
		break
	}

	if !validate_param {
		err = errors.New("Unkown data sort")
		resp := response.Set_error_response(err, response.Code_error_general, false, response.Message_error_parameter)
		return resp, err
	}

	if offset, err = strconv.Atoi(c.QueryParam("offset")); err != nil {
		resp := response.Set_error_response(err, response.Code_error_general, false, response.Message_error_parameter)
		return resp, err
	}

	if limit, err = strconv.Atoi(c.QueryParam("limit")); err != nil {
		resp := response.Set_error_response(err, response.Code_error_general, false, response.Message_error_parameter)
		return resp, err
	}
	if data, err = IN.product_repository.Rep_get_product_sort_by_key(c.Param("key"), offset, limit); err != nil {
		resp := response.Set_error_response(err, response.Code_error_general, false, response.Message_success_data_not_found)
		return resp, err
	}

	resp = response.Set_success_response(response.Code_success_general, true, response.Message_success_general, data)
	return resp, err
}
