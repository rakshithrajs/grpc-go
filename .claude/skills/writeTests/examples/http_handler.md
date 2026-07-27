# example only -  http handler 

```go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/todo/internal/handlers"
	users "github.com/rakshithrajs/todo/internal/handlers/users"
	"github.com/rakshithrajs/todo/internal/models"
	logger "github.com/rakshithrajs/todo/internal/utils/logger"
	validation "github.com/rakshithrajs/todo/internal/utils/validation"
)

const (
	funcNameCreateTodoHandler = "CreateTodoHandler"
)

func (h *TodoHandler) CreateTodoHandler(c *gin.Context) {
	userID, err := handlers.GetUser(c, funcNameCreateTodoHandler)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{logger.ErrorKey: err.Error()})
		return
	}

	var req models.CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{logger.ErrorKey: users.ErrInvalidJSON.Error()})
		return
	}

	if err := validation.Validate.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{logger.ErrorKey: validation.FieldErrors(err)})
		return
	}

	todo, err := h.service.CreateTodo(c.Request.Context(), req.Title, userID)
	if err != nil {
		if handlers.HandleDomainError(c, err) {
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{logger.ErrorKey: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"todo": todo})
}
```

```go
package handlers

import (
	"net/http"
	"testing"

	"github.com/rakshithrajs/todo/internal/handlers"
	users "github.com/rakshithrajs/todo/internal/handlers/users"
	"github.com/rakshithrajs/todo/internal/mocks"
	"github.com/rakshithrajs/todo/internal/models"
	"github.com/rakshithrajs/todo/internal/storage"
	validation "github.com/rakshithrajs/todo/internal/utils/validation"
)

func TestCreateTodoHandler(t *testing.T) {
	tests := []struct {
		name          string
		body          string
		auth          bool
		mockErr       mocks.DbOperationError
		expectedCode  int
		expectedError any
		expectedData  any
	}{
		{
			name:          "todo creation fails due to missing auth",
			body:          `{"title": "Test Todo"}`,
			auth:          false,
			expectedCode:  http.StatusUnauthorized,
			expectedError: handlers.ErrUnathorized.Error(),
		},
		{
			name:          "todo creation fails due to invalid json",
			body:          `{`,
			auth:          true,
			expectedCode:  http.StatusBadRequest,
			expectedError: users.ErrInvalidJSON.Error(),
		},
		{
			name:         "empty title",
			body:         `{"title": ""}`,
			auth:         true,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"title": validation.ErrTitleRequired.Error(),
			},
		},
		{
			name:         "title with only spaces",
			body:         `{"title": "   "}`,
			auth:         true,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"title": validation.ErrTitleRequired.Error(),
			},
		},
		{
			name:         "title with a invalid length",
			body:         `{"title": "a"}`,
			auth:         true,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"title": validation.ErrInvalidLength.Error(),
			},
		},
		{
			name:          "title already exists",
			body:          `{"title": "Test Todo"}`,
			auth:          true,
			mockErr:       mocks.DbOpDuplicateTitle,
			expectedCode:  http.StatusConflict,
			expectedError: storage.ErrTodoTitleAlreadyExists.Error(),
		},
		{
			name:          "internal server error",
			body:          `{"title": "Test Todo"}`,
			auth:          true,
			mockErr:       mocks.DbOpInternalError,
			expectedCode:  http.StatusInternalServerError,
			expectedError: storage.ErrFailedToCreateTodo.Error(),
		},
		{
			name:         "success",
			body:         `{"title": "Test Todo"}`,
			auth:         true,
			mockErr:      mocks.DbOpSuccess,
			expectedCode: http.StatusCreated,
			expectedData: map[string]models.Todo{
				"todo": {
					ID:     new(1),
					Title:  new("Test Todo"),
					UserID: new("test-user-id"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := mocks.SetUpGinTest(http.MethodPost, "/todos", tt.body, tt.auth)

			svc := &mocks.MockTodoService{MockErr: tt.mockErr}
			handler := NewTodoHandler(svc)

			handler.CreateTodoHandler(c)

			if w.Code != tt.expectedCode {
				t.Errorf("expected %d, got %d", tt.expectedCode, w.Code)
			}

			if tt.expectedCode == http.StatusCreated {
				mocks.CheckData(t, w, tt.expectedData)
				return
			}

			if tt.expectedError != nil {
				mocks.CheckError(t, w, tt.expectedError)
			}
		})
	}
}
```