package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"simple-crud/go/db"
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

	db.Init("test.db")

	if assert.NoError(t, Create(c)) {
		assert.NotEmpty(t, rec.Body.String())
	}
}