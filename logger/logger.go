package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"superindo/v1/model"
	"time"
)

var dedicated *logrus.Logger

type CustomStringFormatter struct{}

func (f *CustomStringFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(entry.Message), nil
}

func Init_Logger(customlog *logrus.Logger) *logrus.Logger {
	var err error
	var file *os.File
	dedicated = customlog
	var logFileName = get_log_file_name()
	dedicated.SetLevel(logrus.InfoLevel)
	dedicated.SetFormatter(&CustomStringFormatter{})
	if file, err = os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err != nil {
		fmt.Println("Failed to open log file:", err)
		return dedicated
	}
	//defer file.Close()

	app_level := os.Getenv("APP_LEVEL")
	switch app_level {
	case "production":
		// hanya menulis log di file
		dedicated.SetOutput(file)
	default:
		// menulis log di console dan file
		dedicated.SetOutput(io.MultiWriter(file, os.Stdout))
	}

	return dedicated
}

func Logger_middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			request_data := get_log_request(c)
			request_log := Create_api_loger_string(request_data)
			dedicated.Info(request_log)
			return next(c)
		}
	}
}

func get_log_request(c echo.Context) model.LogModel {
	var result model.LogModel
	request_id := c.Response().Header().Get(echo.HeaderXRequestID)
	uri := c.Request().RequestURI
	bytesIn := c.Request().Header.Get(echo.HeaderContentLength)
	timezone, _ := time.Now().Local().Zone()
	method := c.Request().Method
	request_headers := make(map[string]string)
	for key, values := range c.Request().Header {
		request_headers[key] = strings.Join(values, ", ")
	}

	request_header_data := c.Request().Header
	request_header_json, _ := json.Marshal(request_header_data)

	contentType := c.Request().Header.Get("Content-Type")

	var request_body_string string
	if strings.Contains(contentType, echo.MIMEApplicationJSON) {
		request_body_string = get_request_body_string(c)
	}
	if strings.Contains(contentType, echo.MIMEMultipartForm) {
		request_body_string = get_request_body_form_data(c)
	}

	var request_detail_model model.RequestDetail
	request_detail_model.Request_date = time.Now().Format("2006-01-02 15:04:05.999 -07:00")
	request_detail_model.Request_id = request_id
	request_detail_model.Method = method
	request_detail_model.Url = uri
	request_detail_model.Ip = c.RealIP()
	request_detail_model.Api_type = "Request"
	request_detail_model.Level_name = "Debug"
	request_detail_model.Channel = "development"
	request_detail_model.Timezone = timezone
	request_detail_model.Bytes_in = bytesIn

	request_detail_json, _ := json.Marshal(request_detail_model)

	result.Mode = "Request"
	result.Detail = string(request_detail_json)
	result.Header = string(request_header_json)
	result.Body = request_body_string
	return result
}

func Get_log_response(c echo.Context, resBody []byte) model.LogModel {
	start := time.Now()
	bytesOut := fmt.Sprint(c.Response().Size)
	latency := time.Since(start)
	latencyHuman := latency.String()
	status := c.Response().Status
	request_id := c.Response().Header().Get(echo.HeaderXRequestID)
	uri := c.Request().RequestURI
	bytesIn := c.Request().Header.Get(echo.HeaderContentLength)
	timezone, _ := time.Now().Local().Zone()
	method := c.Request().Method

	var response_detail_data model.ResponseDetail
	response_detail_data.Request_date = time.Now().Format("2006-01-02 15:04:05.999 -07:00")
	response_detail_data.Request_id = request_id
	response_detail_data.Method = method
	response_detail_data.Url = uri
	response_detail_data.Ip = c.RealIP()
	response_detail_data.Api_type = "Response"
	response_detail_data.Level_name = "Debug"
	response_detail_data.Channel = "development"
	response_detail_data.Timezone = timezone
	response_detail_data.Bytes_in = bytesIn
	response_detail_data.Bytes_out = bytesOut
	//response_detail_data.Error = err.Error()
	response_detail_data.Latency = latency
	response_detail_data.Latency_human = latencyHuman
	response_detail_data.Status = status

	response_detail_json, _ := json.Marshal(response_detail_data)

	response_header_data := c.Request().Header
	response_header_json, _ := json.Marshal(response_header_data)

	compactResBody := new(bytes.Buffer)
	if err := json.Compact(compactResBody, resBody); err != nil {
		// Handle error jika diperlukan
		fmt.Println("Error compacting response body:", err)
	}

	var result model.LogModel
	result.Mode = "Response"
	result.Detail = string(response_detail_json)
	result.Header = string(response_header_json)
	result.Body = compactResBody.String()
	return result
}

func Create_api_loger_string(logModel model.LogModel) string {
	var result string
	result += fmt.Sprintf("\n========= %s Information Start =========\n", logModel.Mode)
	result += fmt.Sprintf("\n%s Detail -=> %s\n", logModel.Mode, string(logModel.Detail))
	result += fmt.Sprintf("\n%s Header -=> %s\n", logModel.Mode, string(logModel.Header))
	result += fmt.Sprintf("\n%s Body -=> %s\n", logModel.Mode, string(logModel.Body))
	result += fmt.Sprintf("\n========= %s Information END =========\n", logModel.Mode)
	return result
}

func get_log_file_name() string {
	appName := os.Getenv("APP_NAME")
	currentTime := time.Now()
	logDir := "logs"
	logFileName := fmt.Sprintf("%s/%s_%s.log", logDir, appName, currentTime.Format("20060102"))
	_ = os.Mkdir(logDir, os.ModePerm) // Buat direktori logs jika belum ada
	return logFileName
}

func get_request_body_form_data(c echo.Context) string {
	form, err := c.MultipartForm()
	if err != nil {
		fmt.Println("Error parsing form data:", err)
		return ""
	}

	// Menampilkan informasi files
	files := make([]map[string]interface{}, 0)
	for _, headers := range form.File {
		for _, header := range headers {
			fileInfo := map[string]interface{}{
				"FileName":    header.Filename,
				"ContentType": header.Header.Get("Content-Type"),
				"Size":        header.Size,
			}
			files = append(files, fileInfo)
		}
	}

	// Menampilkan informasi form fields
	formFields := make(map[string]interface{})
	for key, values := range form.Value {
		if len(values) > 0 {
			switch key {
			case "password":
				formFields[key] = "****"
			default:
				formFields[key] = values[0]
			}
		}
	}

	// Menampilkan hasil akhir dalam format yang diinginkan
	result := map[string]interface{}{
		"Files":      files,
		"FormFields": formFields,
	}

	resultString, err := json.Marshal(result)
	if err != nil {
		fmt.Println("Error marshalling form data:", err)
		return ""
	}

	return string(resultString)
}

func get_request_body_string(c echo.Context) string {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		fmt.Println("Error reading request body:", err)
		return ""
	}

	var bodyJSON map[string]interface{}
	if err := json.Unmarshal(body, &bodyJSON); err != nil {
		fmt.Println("Error unmarshalling JSON body:", err)
		return ""
	}
	modify_sensitive_fields(bodyJSON)
	bodyString, err := json.Marshal(bodyJSON)
	if err != nil {
		fmt.Println("Error marshalling JSON body:", err)
		return ""
	}

	// Mengembalikan body ke body aslinya
	c.Request().Body = ioutil.NopCloser(strings.NewReader(string(body)))
	return string(bodyString)
}

func modify_sensitive_fields(bodyJSON map[string]interface{}) {
	for key, value := range bodyJSON {
		switch key {
		case "password":
			bodyJSON[key] = "****"
		default:
			bodyJSON[key] = value
		}

	}
}

func Custom_logger(source string, title string, message string) {
	var result string
	result += fmt.Sprintf("\n--- *** --- %s Logger Information Start --- *** ---\n", source)
	result += fmt.Sprintf("\n%s - %s -=> %s\n", source, title, message)
	result += fmt.Sprintf("\n--- *** --- %s Logger Information End --- *** ---\n", source)
	dedicated.Info(result)
}

func MakeError(c echo.Context, data error) *logrus.Entry {
	entry := logrus.NewEntry(logrus.StandardLogger())
	if c != nil && c.Response() != nil {
		var rl model.RequestLogger
		rl.RequestID = c.Response().Header().Get(echo.HeaderXRequestID)
		entry = dedicated.WithFields(logrus.Fields{"params": rl})
		return entry
	}

	entry = entry.WithError(data)
	//entry.Error(data.Error()) // Menambahkan pesan error ke log
	dedicated.Error(fmt.Sprintf("\n%s", data.Error()))

	return entry
}
