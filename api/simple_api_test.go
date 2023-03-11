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
