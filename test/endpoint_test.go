package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"superindo/v1/model"
	"testing"
	"time"
)

var (
	baseUrl      = "http://localhost:3000"
	product_type = []string{"sayuran", "protein", "buah", "snack"}
	old_product  model.Product
)

func generate_name() string {
	var userName string
	err := faker.FakeData(&userName)
	if err != nil {
		panic(err)
	}
	return userName
}

func random_product_type() string {
	var num int
	rand.Seed(time.Now().UnixNano())
	num = rand.Intn(3)
	return product_type[num]
}

func makeString(length int) string {
	return strings.Repeat("a", length)
}

func randomBetween100kAnd1M() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(900001) + 100000
}

func createRequest(method string, url string, requestBody map[string]interface{}) (newMap map[string]interface{}, err error) {

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, baseUrl+url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response, err
}

func Test_success_register(t *testing.T) {
	name := generate_name()
	product_type := random_product_type()
	description := makeString(100)
	price := randomBetween100kAnd1M()

	old_product.Name = name
	old_product.Product_type = product_type
	old_product.Description = description
	old_product.Price = price

	requestBody := map[string]interface{}{
		"name":         name,
		"product_type": product_type,
		"description":  description,
		"price":        price,
	}

	response, err := createRequest(http.MethodPost, "/api/product/register", requestBody)
	if err != nil {
		t.Errorf("Terjadi kesalahan,  %s", err.Error())
	}

	assert.NotEmpty(t, response["code"])
	assert.NotEmpty(t, response["status"])
	assert.NotEmpty(t, response["title"])
	assert.NotEmpty(t, response["message"])
	assert.Nil(t, response["data"], "Response data harus bernilai nil")

	assert.Equal(t, response["code"], "00")
	assert.Equal(t, response["status"], true)
	assert.Equal(t, response["title"], "Success")
	assert.Equal(t, response["message"], "Add Data in Process")
	assert.Equal(t, response["data"], nil)

}

func Test_error_register_duplicate(t *testing.T) {

	requestBody := map[string]interface{}{
		"name":         old_product.Name,
		"product_type": old_product.Product_type,
		"description":  old_product.Description,
		"price":        old_product.Price,
	}

	response, err := createRequest(http.MethodPost, "/api/product/register", requestBody)
	if err != nil {
		t.Errorf("Terjadi kesalahan,  %s", err.Error())
	}

	assert.NotEmpty(t, response["code"])
	assert.NotEmpty(t, response["title"])
	assert.NotEmpty(t, response["message"])

	assert.Equal(t, response["code"], "500")
	assert.IsType(t, true, response["status"], "Response status harus bertipe bool")
	assert.False(t, response["status"].(bool), "Response status harus bernilai false")
	assert.Equal(t, response["title"], "Error")
	assert.Equal(t, response["message"], "Internal Server Error")
}

func Test_error_register_unknown_product_type(t *testing.T) {
	name := generate_name()
	description := makeString(100)
	price := randomBetween100kAnd1M()

	requestBody := map[string]interface{}{
		"name":         name,
		"product_type": "some tipe",
		"description":  description,
		"price":        price,
	}

	response, err := createRequest(http.MethodPost, "/api/product/register", requestBody)
	if err != nil {
		t.Errorf("Terjadi kesalahan,  %s", err.Error())
	}

	assert.NotEmpty(t, response["code"])
	assert.NotEmpty(t, response["title"])
	assert.NotEmpty(t, response["message"])

	assert.Equal(t, response["code"], "500")
	assert.IsType(t, true, response["status"], "Response status harus bertipe bool")
	assert.False(t, response["status"].(bool), "Response status harus bernilai false")
	assert.Equal(t, response["title"], "Error")
	assert.Equal(t, response["message"], "Internal Server Error")
}

func TestSearchProduct(t *testing.T) {
	requestBody := map[string]interface{}{}
	code_product := "PR-4TTXTJZ7"
	url := fmt.Sprintf("/api/product/search?data=%s", code_product)

	response, err := createRequest(http.MethodGet, url, requestBody)
	if err != nil {
		t.Fatalf("Gagal membuat permintaan HTTP: %v", err)
	}

	assert.Equal(t, "00", response["code"], "Kode harus 00")
	assert.True(t, response["status"].(bool), "Status harus true")
	assert.Equal(t, "Success", response["title"], "Title harus 'Success'")
	assert.Equal(t, "Data Already Sent", response["message"], "Message harus 'Data Already Sent'")

	data := response["data"].(map[string]interface{})

	assert.Equal(t, "produk10", data["name"], "Nama produk harus 'produk10'")
	assert.Equal(t, "PR-4TTXTJZ7", data["product_code"], "Kode produk harus 'PR-4TTXTJZ7'")
	assert.Equal(t, "buah", data["product_type"], "Tipe produk harus 'buah'")
	assert.Equal(t, "description satu", data["description"], "Deskripsi produk harus 'description satu'")
	assert.Equal(t, 100000.0, data["price"], "Harga produk harus 100000")
	assert.Equal(t, "2024-03-20", data["register_date"], "Tanggal pendaftaran harus '2024-03-20'")
}

func TestSearchProductNotFound(t *testing.T) {
	requestBody := map[string]interface{}{}

	code_product := "PR-4ABCD1234Z7WW"
	url := fmt.Sprintf("/api/product/search?data=%s", code_product)

	response, err := createRequest(http.MethodGet, url, requestBody)
	if err != nil {
		t.Fatalf("Gagal membuat permintaan HTTP: %v", err)
	}

	assert.Equal(t, "500", response["code"], "Kode harus 500")
	assert.False(t, response["status"].(bool), "Status harus false")
	assert.Equal(t, "Error", response["title"], "Title harus 'Error'")
	assert.Equal(t, "Data not found for the specified entity.", response["message"], "Message harus 'Data not found for the specified entity.'")
}
