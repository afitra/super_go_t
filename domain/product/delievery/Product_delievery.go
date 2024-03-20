package delievery

import (
	"fmt"
	"superindo/v1/domain/product"
	"superindo/v1/model"
)

type Product_delievery struct {
	Product_usecase product.Usecase
}

func NewProduct_delievery(echoGroup model.EchoGroup, emus product.Usecase) {
	handler := &Product_delievery{
		Product_usecase: emus,
	}
	fmt.Println(handler)

	//echoGroup.API.GET("/employee/:kota", handler.del_get_employee_by_id)
	//echoGroup.Public.POST("/employee/register", handler.del_register_employee)
	//echoGroup.Public.POST("/employee/login", handler.del_login_employee)
}
