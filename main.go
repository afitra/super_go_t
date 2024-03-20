package main

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"superindo/v1/custom_middleware"
	_product_delievery "superindo/v1/domain/product/delievery"
	_product_repository "superindo/v1/domain/product/repository"
	_product_usecase "superindo/v1/domain/product/usecase"
	"superindo/v1/logger"
	"superindo/v1/model"
	"time"
)

var (
	ech                   *echo.Echo
	basic_sql_connection  *sql.DB
	basic_sqlx_connection *sqlx.DB
	customlog             *logrus.Logger
	echoGroup             model.EchoGroup
)

func init() {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	time.Local = loc
	ech = echo.New()
	ech.Debug = true
	loadEnv()

	customlog = logrus.New()
	logger.Init_Logger(customlog)

	getDBConn()

}

func main() {

	echoGroup = model.EchoGroup{
		API:     ech.Group("/api"),
		Private: ech.Group("/private"),
		Public:  ech.Group("/public"),
	}

	custom_middleware.InitMiddleware(ech, echoGroup, customlog)
	setDependencyInjection()

	ech.GET("/ping", ping)
	var err error
	if err = ech.Start(":" + os.Getenv(`APP_PORT`)); err != nil {
		logger.MakeError(nil, err).Debug(err)
	}
}

func setDependencyInjection() {
	// Repository

	employee_repository := _product_repository.NewProduct_repository(basic_sqlx_connection)

	// Usecase
	employee_usecase := _product_usecase.NewProduct_usecase(employee_repository)

	// Handler
	_product_delievery.NewProduct_delievery(echoGroup, employee_usecase)

}

func ping(echTx echo.Context) error {
	response := map[string]interface{}{
		"status": "online",
		"month":  "Server Actived!!",
	}

	return echTx.JSON(http.StatusOK, response)
}

func loadEnv() {
	var err error
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		return
	}

	if err = godotenv.Load(); err != nil {
		logger.MakeError(nil, err).Fatal(err)
	}
}

func getDBConn() {
	var err error
	var sqlxConn *sqlx.DB
	var sqlConn *sql.DB

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")

	// SQLX Connection
	sqlxConnection := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)
	if sqlxConn, err = sqlx.Connect("postgres", sqlxConnection); err != nil {
		logger.MakeError(nil, err).Fatal(err)
		os.Exit(1)
	}

	// SQL Connection
	sqlConnection := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)
	if sqlConn, err = sql.Open("postgres", sqlConnection); err != nil {
		logger.MakeError(nil, err).Fatal(err)
		os.Exit(1)
	}

	// Ping SQLX
	if err = sqlxConn.Ping(); err != nil {
		logger.MakeError(nil, err).Fatal(err)
		os.Exit(1)
	}

	// Ping SQL
	if err = sqlConn.Ping(); err != nil {
		logger.MakeError(nil, err).Fatal(err)
		os.Exit(1)
	}

	// Set global variables
	basic_sql_connection = sqlConn
	basic_sqlx_connection = sqlxConn

	defer basic_sql_connection.Close()
}
