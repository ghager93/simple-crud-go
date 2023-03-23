package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"simple-crud/go/db"
	"simple-crud/go/models"
)

func setupTest(t *testing.T) func() {
	db.Init("test.db")

	return func() {
		db.DropTables()
	}
}

func createPayload(name string, number int) map[string]interface{} {
	return map[string]interface{}{
		"name": name,
		"number": number,
	}
}

func postSimple(payload map[string]interface{}) (*httptest.ResponseRecorder, error) {
	payloadBytes, _ := json.Marshal(payload)	

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/simple", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := Create(c)

	return rec, err
}

func getAllSimple() (*httptest.ResponseRecorder, error) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/simple", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := GetAll(c)

	return rec, err
}

func getSimple(id string) (*httptest.ResponseRecorder, error) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/simple/" + id, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(id)

	err := Get(c)

	return rec, err	
}

func deleteSimple(id string) (*httptest.ResponseRecorder, error) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/simple/" + id, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(id)

	err := Delete(c)

	return rec, err
}

func TestCreateReturns(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	payload := createPayload("", 0)
	rec, err := postSimple(payload)

	assert.NoError(t, err)
	assert.NotEmpty(t, rec.Body.String())
}

func TestCreateValidRequestReturns201(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	payload := createPayload("john", 1234)
	rec, err := postSimple(payload)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestCreateSavesToDB(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	payload := createPayload("john", 1234)
	_, err := postSimple(payload)
	assert.NoError(t, err)
	
	var simple models.Simple
	db.DbManager().First(&simple)
	
	assert.Equal(t, payload["name"], simple.Name)
	assert.Equal(t, payload["number"], simple.Number)
}

func TestCreateInvalidRequestReturns400(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	payload := map[string]interface{}{
		"name": "john",
	} 
	
	rec, err := postSimple(payload)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateInvalidRequestNotSavedToDB(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	payload := map[string]interface{}{
		"name": "john",
	} 
	
	_, err := postSimple(payload)
	assert.NoError(t, err)
	
	var simple models.Simple
	result := db.DbManager().First(&simple)
	assert.Error(t, result.Error)
}

func TestGetAllReturns200(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	rec, err := getAllSimple()

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetAllReturnsEmptyList(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	rec, err := getAllSimple()

	assert.NoError(t, err)

	stringResult := rec.Body.String()

	var simples []models.Simple
	assert.NoError(t, json.Unmarshal([]byte(stringResult), &simples))

	assert.Len(t, simples, 0)
}

func TestGetAllSingleEntry(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	payload := createPayload("john", 1234)

	_, err := postSimple(payload)

	assert.NoError(t, err)

	rec, err := getAllSimple()

	assert.NoError(t, err)

	stringResult := rec.Body.String()

	var simples []models.Simple
	assert.NoError(t, json.Unmarshal([]byte(stringResult), &simples))

	assert.Len(t, simples, 1)
	assert.Equal(t, payload["name"], simples[0].Name)
	assert.Equal(t, payload["number"], simples[0].Number)
}

func TestGetAllThreeEntries(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	type testBody map[string]interface{}
	payloadSlice := []testBody{}

	payload1 := createPayload("john", 1234)
	payload2 := createPayload("jane", 1)
	payload3 := createPayload("SAM_01234", -123)

	payloadSlice = append(payloadSlice, payload1)
	payloadSlice = append(payloadSlice, payload2)
	payloadSlice = append(payloadSlice, payload3)

	for _, payload := range payloadSlice {
		_, err := postSimple(payload)
		assert.NoError(t, err)
	}

	rec, err := getAllSimple()

	assert.NoError(t, err)

	stringResult := rec.Body.String()

	var simples []models.Simple
	assert.NoError(t, json.Unmarshal([]byte(stringResult), &simples))

	assert.Len(t, simples, 3)
	
	for i, payload := range payloadSlice {
		assert.Equal(t, payload["name"], simples[i].Name)
		assert.Equal(t, payload["number"], simples[i].Number)
	}
}

func TestGetReturns200(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	payload := createPayload("john", 1234)
	_, err := postSimple(payload)

	assert.NoError(t, err)

	rec, err := getSimple("1")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)	
}

func TestInvalidGetReturns404(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	payload := createPayload("john", 1234)
	_, err := postSimple(payload)

	assert.NoError(t, err)

	rec, err := getSimple("2")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)	
}

func TestGetByIDReturnsEntry(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	payload := createPayload("john", 1234)
	_, err := postSimple(payload)

	assert.NoError(t, err)

	rec, err := getSimple("1")

	assert.NoError(t, err)

	stringResult := rec.Body.String()

	var simple models.Simple
	assert.NoError(t, json.Unmarshal([]byte(stringResult), &simple))

	assert.Equal(t, payload["name"], simple.Name)
	assert.Equal(t, payload["number"], simple.Number)
}

func TestDeleteReturns200(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	payload := createPayload("john", 1234)
	_, err := postSimple(payload)

	assert.NoError(t, err)

	rec, err := deleteSimple("1")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestInvalidDeleteReturns404(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	rec, err := deleteSimple("1")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDeleteRemovesEntry(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	type testBody map[string]interface{}
	payloadSlice := []testBody{}

	payload1 := createPayload("john", 1234)
	payload2 := createPayload("jane", 1)
	payload3 := createPayload("SAM_01234", -123)

	payloadSlice = append(payloadSlice, payload1)
	payloadSlice = append(payloadSlice, payload2)
	payloadSlice = append(payloadSlice, payload3)

	for _, payload := range payloadSlice {
		_, err := postSimple(payload)
		assert.NoError(t, err)
	}

	_, err := deleteSimple("2")

	assert.NoError(t, err)

	rec, err := getSimple("1")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	rec, err = getSimple("2")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	rec, err = getSimple("3")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}


