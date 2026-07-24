package mocks

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error any `json:"error"`
}

func CheckError(t *testing.T, w *httptest.ResponseRecorder, expected any) {
	t.Helper()

	body := w.Body.Bytes()

	var resp ErrorResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	switch want := expected.(type) {
	case string:
		got, ok := resp.Error.(string)
		if !ok {
			t.Fatalf("expected single error %q, got keyed error %v", want, resp.Error)
		}
		if got != want {
			t.Errorf("expected %q, got %q", want, got)
		}
	case map[string]string:
		got, ok := resp.Error.(map[string]any)
		if !ok {
			t.Fatalf("expected keyed validation errors, got %v", resp.Error)
		}
		if len(got) != len(want) {
			t.Errorf("expected %d errors, got %d (%v)", len(want), len(got), got)
		}
		for key, wantErr := range want {
			gotVal, exists := got[key]
			if !exists {
				t.Errorf("expected error for field %q", key)
				continue
			}
			gotStr, ok := gotVal.(string)
			if !ok {
				t.Errorf("expected string error for field %q, got %v", key, gotVal)
				continue
			}
			if gotStr != wantErr {
				t.Errorf("expected field %q to have error %q, got %q", key, wantErr, gotStr)
			}
		}
	default:
		t.Fatalf("unsupported expected error type %T", expected)
	}
}

func SetUpGinTest(method, url string, body string, authWorks bool) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest(method, url, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	c.Request = req

	if authWorks {
		c.Set("userID", "test-user-id")
	}

	return c, w
}

func SetUpGinTestMultipart(fileContent string, authWorks bool) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, _ := writer.CreateFormFile("file", "test.txt")
	io.WriteString(part, fileContent)
	writer.Close()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodPost, "/api/files/upload", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	c.Request = req
	if authWorks {
		c.Set("userID", "test-user-id")
	}
	return c, w
}

func CheckData(t *testing.T, w *httptest.ResponseRecorder, expected any) {
	t.Helper()

	body := w.Body.Bytes()

	switch want := expected.(type) {
	case map[string]any:
		var got map[string]any
		if err := json.Unmarshal(body, &got); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if !reflect.DeepEqual(want, got) {
			t.Errorf("expected %v, got %v", want, got)
		}
	case map[string]string:
		var got map[string]string
		if err := json.Unmarshal(body, &got); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if !reflect.DeepEqual(want, got) {
			t.Errorf("expected %v, got %v", want, got)
		}
	default:
		t.Fatalf("unsupported expected data type %T", expected)
	}
}
