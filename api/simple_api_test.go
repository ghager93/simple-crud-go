package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateReturns(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/simple", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)



	if assert.NoError(t, Create(c)) {
		assert.NotEmpty(t, rec.Body.String())
	}
}