package delievery

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"superindo/v1/domain/product"
	"superindo/v1/model"
	"superindo/v1/response"
)

type Product_delievery struct {
	Product_usecase product.Usecase
}

func NewProduct_delievery(echoGroup model.EchoGroup, emus product.Usecase) {
	handler := &Product_delievery{
		Product_usecase: emus,
	}

	echoGroup.API.POST("/product/register", handler.Del_post_product_register)
	echoGroup.API.GET("/product/list", handler.Del_get_product_list)
	echoGroup.API.GET("/product/search", handler.Del_get_product_search)
	echoGroup.API.GET("/product/filter-product_type/:product_type", handler.Del_get_product_filter)
	echoGroup.API.GET("/product/sort/:key", handler.Del_get_product_sort)

}

func (emus *Product_delievery) Del_post_product_register(c echo.Context) error {
	var resp interface{}
	var request model.Req_register_product

	if err := c.Bind(&request); err != nil {
		return response.Reverse_error_response(c, err)
	}

	if err := c.Validate(request); err != nil {
		return response.Reverse_error_response(c, err)
	}

	resp, _ = emus.Product_usecase.Use_product_register(c, request)
	return c.JSON(http.StatusOK, resp)
}

func (emus *Product_delievery) Del_get_product_list(c echo.Context) error {
	var resp interface{}
	resp, _ = emus.Product_usecase.Use_get_product_list(c)
	return c.JSON(http.StatusOK, resp)
}

func (emus *Product_delievery) Del_get_product_search(c echo.Context) error {
	var resp interface{}
	resp, _ = emus.Product_usecase.Use_get_product_search(c)
	return c.JSON(http.StatusOK, resp)
}

func (emus *Product_delievery) Del_get_product_filter(c echo.Context) error {
	var resp interface{}
	resp, _ = emus.Product_usecase.Use_get_product_filter(c)
	return c.JSON(http.StatusOK, resp)
}

func (emus *Product_delievery) Del_get_product_sort(c echo.Context) error {
	var resp interface{}
	resp, _ = emus.Product_usecase.Use_get_product_sort(c)
	return c.JSON(http.StatusOK, resp)
}
