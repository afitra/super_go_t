package usecase

import (
	"errors"
	"github.com/labstack/echo/v4"
	"strconv"
	"superindo/v1/domain/product"
	"superindo/v1/helper"
	"superindo/v1/model"
	"superindo/v1/response"
)

type Product_usecase struct {
	source             string
	product_repository product.Repository
}

func NewProduct_usecase(product_repository product.Repository) product.Usecase {
	return &Product_usecase{"Product_usecase", product_repository}
}

func (IN *Product_usecase) Use_product_register(c echo.Context, request model.Req_register_product) (interface{}, error) {
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

	if offset, err = strconv.Atoi(c.QueryParam("offset")); err != nil {
		resp := response.Set_error_response(err, response.Code_error_general, false, response.Message_error_parameter)
		return resp, err
	}

	if limit, err = strconv.Atoi(c.QueryParam("limit")); err != nil {
		resp := response.Set_error_response(err, response.Code_error_general, false, response.Message_error_parameter)
		return resp, err
	}
	if data, err = IN.product_repository.Rep_get_product_filter_by_product_code(c.Param("product_type"), offset, limit); err != nil {
		resp := response.Set_error_response(err, response.Code_error_general, false, response.Message_success_data_not_found)
		return resp, err
	}

	resp = response.Set_success_response(response.Code_success_general, true, response.Message_success_general, data)
	return resp, err
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
