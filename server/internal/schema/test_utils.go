package schema

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func bindMockRequest[schema any](t *testing.T, payload map[string]any) (schema, error) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	b, err := json.Marshal(payload)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("POST", "/company", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	var s schema
	err = c.ShouldBindJSON(&s)
	return s, err
}
