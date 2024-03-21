package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	_redis_connection "superindo/v1/connection"
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
	custom_log            *logrus.Logger
	echo_group            model.EchoGroup
	redis_client          *redis.Client
)

func init() {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	time.Local = loc
	ech = echo.New()
	ech.Debug = true
	loadEnv()

	custom_log = logrus.New()
	logger.Init_Logger(custom_log)

	getDBConn()
	init_redis()

}

func main() {

	echo_group = model.EchoGroup{
		API:     ech.Group("/api"),
		Private: ech.Group("/private"),
		Public:  ech.Group("/public"),
	}

	custom_middleware.InitMiddleware(ech, echo_group, custom_log)
	setDependencyInjection()

	ech.GET("/ping", ping)
	var err error
	if err = ech.Start(":" + os.Getenv(`APP_PORT`)); err != nil {
		logger.MakeError(nil, err).Debug(err)
	}
}

func setDependencyInjection() {

	redis_con := _redis_connection.NewRedis_Connection(redis_client)

	// Repository

	employee_repository := _product_repository.NewProduct_repository(basic_sqlx_connection)

	// Usecase
	employee_usecase := _product_usecase.NewProduct_usecase(employee_repository, redis_con, redis_client)

	// Handler
	_product_delievery.NewProduct_delievery(echo_group, employee_usecase)

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

func init_redis() {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}

	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	redis_client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	ctx := context.Background()

	_, err := redis_client.Ping(ctx).Result()
	if err != nil {
		// Handle error
		panic(err)
	}

}
