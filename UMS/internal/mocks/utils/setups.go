package mocks

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
)

// SetUpGinTest creates a gin test context with a JSON request.
func SetUpGinTest(method, url string, body string, authWorks bool) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest(method, url, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	c.Request = req

	if authWorks {
		c.Set(config.UserIDMetadataKey, "test-user-id")
	}

	return c, w
}

// SetUpGinTestMultipart creates a gin test context with a multipart file upload request.
func SetUpGinTestMultipart(fileContent, fileName string, authWorks bool) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, _ := writer.CreateFormFile("file", fileName)
	io.WriteString(part, fileContent)
	writer.Close()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodPost, "/api/files/upload", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	c.Request = req
	if authWorks {
		c.Set(config.UserIDMetadataKey, "test-user-id")
	}
	return c, w
}
