package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	TMSpb "github.com/rakshithrajs/cloud/UMS/gen/TMS/v1"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	grpcClient "github.com/rakshithrajs/cloud/UMS/internal/grpcClient"
	handlerErrors "github.com/rakshithrajs/cloud/UMS/internal/handlers/errors"
	"github.com/rakshithrajs/cloud/UMS/internal/mocks"
	mockUtils "github.com/rakshithrajs/cloud/UMS/internal/mocks/utils"
)

func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name                string
		authorizationHeader string
		validateErr         mocks.GrpcOperationError
		claims              *TMSpb.TokenClaims
		expectedStatusCode  int
		expectedError       any
		expectedUserID      string
	}{
		{
			name:               "authorization fails because of missing authorization header",
			expectedStatusCode: http.StatusUnauthorized,
			expectedError:      handlerErrors.ErrMissingAuthHeader.Error(),
		},
		{
			name:                "authorization fails because token is unauthenticated",
			authorizationHeader: "Bearer invalid-token",
			validateErr:         mocks.GrpcOpInvalidToken,
			expectedStatusCode:  http.StatusUnauthorized,
			expectedError:       handlerErrors.ErrUnauthorized.Error(),
		},
		{
			name:                "authorization fails because of internal error from TMS",
			authorizationHeader: "Bearer valid-token",
			validateErr:         mocks.GrpcOpInternalError,
			expectedStatusCode:  http.StatusInternalServerError,
			expectedError:       handlerErrors.ErrSomethingWentWrong.Error(),
		},
		{
			name:                "authorization succeeds with valid token",
			authorizationHeader: "Bearer valid-token",
			claims: &TMSpb.TokenClaims{
				UserID: "test-user-id",
			},
			expectedStatusCode: http.StatusOK,
			expectedUserID:     "test-user-id",
		},
	}

	gin.SetMode(gin.TestMode)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokensClient := &mocks.MockTokensClient{
				Claims:      tt.claims,
				ValidateErr: tt.validateErr,
			}
			tmsClient := grpcClient.NewTMSClient(tokensClient)

			router := gin.New()
			router.Use(AuthMiddleware(tmsClient))
			router.GET("/test", func(c *gin.Context) {
				userID, _ := c.Get(config.UserIDMetadataKey)
				c.JSON(http.StatusOK, gin.H{"userID": userID})
			})

			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			if tt.authorizationHeader != config.NullString {
				req.Header.Set("Authorization", tt.authorizationHeader)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatusCode {
				t.Errorf("expected status code %d, got %d", tt.expectedStatusCode, w.Code)
			}

			if tt.expectedError != nil {
				mockUtils.CheckError(t, w, tt.expectedError)
			}

			if tt.expectedUserID != config.NullString {
				mockUtils.CheckData(t, w, map[string]any{"userID": tt.expectedUserID})
			}
		})
	}
}
