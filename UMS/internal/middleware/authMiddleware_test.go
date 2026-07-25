package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	handlerUtils "github.com/rakshithrajs/cloud/UMS/internal/handlers/utils"
	middlewareUtils "github.com/rakshithrajs/cloud/UMS/internal/middleware/utils"
	mockUtils "github.com/rakshithrajs/cloud/UMS/internal/mocks/utils"
)

func init() {
	config.SetConfig(&config.Config{JWTSecret: "test-secret"})
}

func GenerateJwt(iss, sub any, iat, exp int64, jwtSecret string, signingMethod jwt.SigningMethod) (string, error) {
	claims := jwt.MapClaims{
		"iss": iss,
		"sub": sub,
		"iat": iat,
		"exp": exp,
	}

	token := jwt.NewWithClaims(signingMethod, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name                string
		AuthorizationHeader string
		expectedStatusCode  int
		expectedError       any
	}{
		{
			name:                "authorization fails because of Missing Authorization Header",
			AuthorizationHeader: "",
			expectedStatusCode:  http.StatusUnauthorized,
			expectedError:       handlerUtils.ErrMissingAuthHeader.Error(),
		},
		{
			name:                "authorization fails because of Missing Bearer",
			AuthorizationHeader: "okenabcdefg",
			expectedStatusCode:  http.StatusUnauthorized,
			expectedError:       middlewareUtils.ErrMissingBearerToken.Error(),
		},
		{
			name:                "authorization fails because of Missing Token",
			AuthorizationHeader: "Bearer ",
			expectedStatusCode:  http.StatusUnauthorized,
			expectedError:       middlewareUtils.ErrMissingBearerToken.Error(),
		},
		{
			name: "authorization fails because of Invalid Signing Method",
			AuthorizationHeader: func() string {
				token, _ := GenerateJwt("cloud-app", "user123", time.Now().Unix(), time.Now().Add(1*time.Hour).Unix(), "test-secret", jwt.SigningMethodHS384)

				return "Bearer " + token
			}(),
			expectedStatusCode: http.StatusUnauthorized,
			expectedError:      middlewareUtils.ErrInvalidToken.Error(),
		},
		{
			name: "authorization fails because of Invalid Token",
			AuthorizationHeader: func() string {
				token, _ := GenerateJwt("cloud-app", "user123", time.Now().Unix(), time.Now().Add(1*time.Hour).Unix(), "wrong-secret", jwt.SigningMethodHS256)
				return "Bearer " + token
			}(),
			expectedStatusCode: http.StatusUnauthorized,
			expectedError:      middlewareUtils.ErrInvalidToken.Error(),
		},
		{
			name: "authorization fails because of Invalid Issuer",
			AuthorizationHeader: func() string {
				token, _ := GenerateJwt("invalid-issuer", "user123", time.Now().Unix(), time.Now().Add(1*time.Hour).Unix(), "test-secret", jwt.SigningMethodHS256)
				return "Bearer " + token
			}(),
			expectedStatusCode: http.StatusUnauthorized,
			expectedError:      middlewareUtils.ErrInvalidToken.Error(),
		},
		{
			name: "authorization fails because of Invalid Subject",
			AuthorizationHeader: func() string {
				token, _ := GenerateJwt("cloud-app", 12345, time.Now().Unix(), time.Now().Add(1*time.Hour).Unix(), "test-secret", jwt.SigningMethodHS256)
				return "Bearer " + token
			}(),
			expectedStatusCode: http.StatusUnauthorized,
			expectedError:      middlewareUtils.ErrInvalidToken.Error(),
		},
		{
			name: "authorization fails because of Zero Issue Time",
			AuthorizationHeader: func() string {
				token, _ := GenerateJwt("cloud-app", "user123", 0000000, time.Now().Add(1*time.Hour).Unix(), "test-secret", jwt.SigningMethodHS256)
				return "Bearer " + token
			}(),
			expectedStatusCode: http.StatusUnauthorized,
			expectedError:      middlewareUtils.ErrInvalidToken.Error(),
		},
		{
			name: "authorization fails because of Future Issue Time",
			AuthorizationHeader: func() string {
				futureTime := time.Now().Add(1 * time.Hour).Unix()
				token, _ := GenerateJwt("cloud-app", "user123", futureTime, futureTime+3600, "test-secret", jwt.SigningMethodHS256)
				return "Bearer " + token
			}(),
			expectedStatusCode: http.StatusUnauthorized,
			expectedError:      middlewareUtils.ErrInvalidToken.Error(),
		},
		{
			name: "authorization fails because of Expired Token",
			AuthorizationHeader: func() string {
				pastTime := time.Now().Add(-1 * time.Hour).Unix()
				token, _ := GenerateJwt("cloud-app", "user123", pastTime-3600, pastTime, "test-secret", jwt.SigningMethodHS256)
				return "Bearer " + token
			}(),
			expectedStatusCode: http.StatusUnauthorized,
			expectedError:      middlewareUtils.ErrTokenExpired.Error(),
		},
		{
			name: "authorization succeeds with Valid Token",
			AuthorizationHeader: func() string {
				currentTime := time.Now().Unix()
				token, _ := GenerateJwt("cloud-app", "user123", currentTime, currentTime+3600, "test-secret", jwt.SigningMethodHS256)
				return "Bearer " + token
			}(),
			expectedStatusCode: http.StatusOK,
		},
	}

	// arrange
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(AuthMiddleware())

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// act
			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			if tt.AuthorizationHeader != "" {
				req.Header.Set("Authorization", tt.AuthorizationHeader)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// assert
			if w.Code != tt.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatusCode, w.Code)
			}

			if tt.expectedError != nil {
				mockUtils.CheckError(t, w, tt.expectedError)
			}
		})
	}
}
