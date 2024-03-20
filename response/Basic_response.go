package response

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"superindo/v1/helper"
)

type Base_error_response struct {
	Code    string   `json:"code"`
	Status  bool     `json:"status"`
	Title   string   `json:"title"`
	Message string   `json:"message"`
	Detail  []string `json:"detail,omitempty"`
}
type Base_response struct {
	Code    string      `json:"code"`
	Status  bool        `json:"status"`
	Title   string      `json:"title"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// SetTitle title of Response errors
func Reverse_error_response(c echo.Context, err error) error {
	result := Base_error_response{
		Code:    Code_error_general,
		Status:  false,
		Title:   "Error",
		Message: Message_error_general,
	}

	env_string := os.Getenv("ENABLE_ERROR_DEBUG_MESSAGE")
	flag := helper.Convert_string_to_bool(env_string)
	if flag {
		result.Detail = helper.Split_string_to_array(err.Error())
	}
	return c.JSON(http.StatusOK, result)
}

func Set_error_response(err error, code string, status bool, message string) Base_error_response {
	result := Base_error_response{
		Code:    code,
		Status:  status,
		Title:   "Error",
		Message: message,
	}

	env_string := os.Getenv("ENABLE_ERROR_DEBUG_MESSAGE")
	flag := helper.Convert_string_to_bool(env_string)
	if flag {
		result.Detail = helper.Split_string_to_array(err.Error())
	}

	return result
}

func Set_success_response(code string, status bool, message string, data interface{}) Base_response {
	var result Base_response
	result.Code = code
	result.Status = status
	result.Title = "Success"
	result.Message = message
	result.Data = data
	if !status {
		result.Title = "Error"
	}

	return result
}
