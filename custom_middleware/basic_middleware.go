package custom_middleware

import (
	"github.com/go-playground/validator/v10"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"os"
	"superindo/v1/logger"
	"superindo/v1/model"
)

type customMiddleware struct {
	e *echo.Echo
}

type customValidator struct {
	Validator *validator.Validate
}

func (cv *customValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

var (
	echGroup  model.EchoGroup
	customlog *logrus.Logger
)

func InitMiddleware(ech *echo.Echo, echoGroup model.EchoGroup, cusLog *logrus.Logger) {
	cm := &customMiddleware{ech}
	echGroup = echoGroup

	customlog = cusLog
	cm.cors()
	ech.Use(middleware.RequestID())
	ech.Use(logger.Logger_middleware())
	ech.Use(middleware.BodyDumpWithConfig(cm.customDumpBody()))
	ech.Use(middleware.Recover())
	//cm.basicAuth()
	//cm.jwtAuth()
	cm.customValidation()
}

func (cm customMiddleware) cors() {
	cm.e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"Access-Control-Allow-Origin"},
		//AllowMethods: []string{"*"},
		AllowMethods: []string{"GET", "POST"},
	}))
}

func (cm customMiddleware) customDumpBody() middleware.BodyDumpConfig {
	return middleware.BodyDumpConfig{
		Handler: func(c echo.Context, reqBody, resBody []byte) {
			response_data := logger.Get_log_response(c, resBody)
			response_log := logger.Create_api_loger_string(response_data)
			customlog.Info(response_log)
			// Simpan response body dalam variabel di context
			c.Set("responseBody", resBody)
		},
	}
}

func (cm customMiddleware) basicAuth() {
	echGroup.Private.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == os.Getenv(`BASIC_AUTH_USERNAME`) && password == os.Getenv(`BASIC_AUTH_PASSWORD`) {
			return true, nil
		}
		return false, nil
	}))
}

func (cm customMiddleware) jwtAuth() {
	echGroup.API.Use(echojwt.WithConfig(echojwt.Config{
		SigningMethod: "HS512",
		SigningKey:    []byte(os.Getenv(`JWT_SECRET`)),
	}))
}

func (cm *customMiddleware) customValidation() {
	validate := validator.New()
	customValidator := &customValidator{Validator: validate}
	cm.e.Validator = customValidator
}
