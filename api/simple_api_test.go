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

func TestCreateReturns(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/simple", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, Create(c)) {
		assert.NotEmpty(t, rec.Body.String())
	}
}

func TestCreateValidRequestReturns201(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	testBody := map[string]interface{}{
		"name": "john",
		"number": 1234,
	} 
	testBodyBytes, _ := json.Marshal(testBody)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/simple", bytes.NewBuffer(testBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, Create(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

func TestCreateSavesToDB(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	testBody := map[string]interface{}{
		"name": "john",
		"number": 1234,
	} 
	testBodyBytes, _ := json.Marshal(testBody)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/simple", bytes.NewBuffer(testBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	assert.NoError(t, Create(c))
	
	var simple models.Simple
	db.DbManager().First(&simple)
	
	assert.Equal(t, testBody["name"], simple.Name)
	assert.Equal(t, testBody["number"], simple.Number)
}

func TestCreateInvalidRequestReturns400(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	testBody := map[string]interface{}{
		"name": "john",
	} 
	testBodyBytes, _ := json.Marshal(testBody)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/simple", bytes.NewBuffer(testBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, Create(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}	
}

func TestCreateInvalidRequestNotSavedToDB(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	testBody := map[string]interface{}{
		"name": "john",
	} 
	testBodyBytes, _ := json.Marshal(testBody)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/simple", bytes.NewBuffer(testBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	assert.NoError(t, Create(c))
	
	var simple models.Simple
	result := db.DbManager().First(&simple)
	assert.Error(t, result.Error)
}

func TestGetAllReturns200(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/simple", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	assert.NoError(t, GetAll(c))
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetAllReturnsEmptyList(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/simple", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	assert.NoError(t, GetAll(c))

	stringResult := rec.Body.String()

	var simples []models.Simple
	assert.NoError(t, json.Unmarshal([]byte(stringResult), &simples))

	assert.Len(t, simples, 0)
}

func TestGetAllSingleEntry(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	testBody := map[string]interface{}{
		"name": "john",
		"number": 1234,
	} 
	testBodyBytes, _ := json.Marshal(testBody)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/simple", bytes.NewBuffer(testBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	assert.NoError(t, Create(c))

	req = httptest.NewRequest(http.MethodGet, "/api/simple", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	assert.NoError(t, GetAll(c))

	stringResult := rec.Body.String()

	var simples []models.Simple
	assert.NoError(t, json.Unmarshal([]byte(stringResult), &simples))

	assert.Len(t, simples, 1)
	assert.Equal(t, testBody["name"], simples[0].Name)
	assert.Equal(t, testBody["number"], simples[0].Number)
}

func TestGetAllThreeEntries(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	type testBody map[string]interface{}
	testBodySlice := []testBody{}

	tB1 := testBody{
		"name": "john",
		"number": 1234,
	}
	tB2 := testBody{
		"name": "jane",
		"number": 1,
	}
	tB3 := testBody{
		"name": "SAM_01234",
		"number": -123,
	}

	testBodySlice = append(testBodySlice, tB1)
	testBodySlice = append(testBodySlice, tB2)
	testBodySlice = append(testBodySlice, tB3)

	e := echo.New()
	for _, testBody := range testBodySlice {
		testBodyBytes, _ := json.Marshal(testBody)

		req := httptest.NewRequest(http.MethodPost, "/api/simple", bytes.NewBuffer(testBodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		assert.NoError(t, Create(c))
	}

	req := httptest.NewRequest(http.MethodGet, "/api/simple", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	assert.NoError(t, GetAll(c))

	stringResult := rec.Body.String()

	var simples []models.Simple
	assert.NoError(t, json.Unmarshal([]byte(stringResult), &simples))

	assert.Len(t, simples, 3)
	
	for i, testBody := range testBodySlice {
		assert.Equal(t, testBody["name"], simples[i].Name)
		assert.Equal(t, testBody["number"], simples[i].Number)
	}
}

func TestGetReturns200(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	testBody := map[string]interface{}{
		"name": "john",
		"number": 1234,
	} 
	testBodyBytes, _ := json.Marshal(testBody)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/simple", bytes.NewBuffer(testBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	assert.NoError(t, Create(c))

	req = httptest.NewRequest(http.MethodGet, "/api/simple/1", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	
	c.SetParamNames("id")
	c.SetParamValues("1")

	assert.NoError(t, Get(c))
	assert.Equal(t, http.StatusOK, rec.Code)	
}

func TestInvalidGetReturns404(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	testBody := map[string]interface{}{
		"name": "john",
		"number": 1234,
	} 
	testBodyBytes, _ := json.Marshal(testBody)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/simple", bytes.NewBuffer(testBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	assert.NoError(t, Create(c))

	req = httptest.NewRequest(http.MethodGet, "/api/simple/1", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	
	c.SetParamNames("id")
	c.SetParamValues("2")

	assert.NoError(t, Get(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)	
}

func TestGetByIDReturnsEntry(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	testBody := map[string]interface{}{
		"name": "john",
		"number": 1234,
	} 
	testBodyBytes, _ := json.Marshal(testBody)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/simple", bytes.NewBuffer(testBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	assert.NoError(t, Create(c))

	req = httptest.NewRequest(http.MethodGet, "/api/simple/1", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	
	c.SetParamNames("id")
	c.SetParamValues("1")

	assert.NoError(t, Get(c))

	stringResult := rec.Body.String()

	var simple models.Simple
	assert.NoError(t, json.Unmarshal([]byte(stringResult), &simple))

	assert.Equal(t, testBody["name"], simple.Name)
	assert.Equal(t, testBody["number"], simple.Number)
}
