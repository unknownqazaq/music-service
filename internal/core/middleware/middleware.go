package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"music-service/internal/core/logger"
	"music-service/internal/core/response"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Claims struct {
	UserID           int64  `json:"user_id"`
	Email            string `json:"email"`
	Role             string `json:"role"`
	SubscriptionType string `json:"subscription_type"`
	jwt.RegisteredClaims
}

type contextKey string

const (
	ClaimsKey       contextKey = "user_claims"
	requestIDHeader            = "X-Request-ID"
)

func CORS() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func RequestID() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)
			if requestID == "" {
				requestID = uuid.NewString()
			}
			r.Header.Set(requestIDHeader, requestID)
			w.Header().Set(requestIDHeader, requestID)
			next.ServeHTTP(w, r)
		})
	}
}

func Logger(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)
			l := log.With(
				zap.String("request_id", requestID),
				zap.String("url", r.URL.String()),
			)
			ctx := logger.ToContext(r.Context(), l)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Trace() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.FromContext(r.Context())
			rw := response.NewResponseWriter(w)

			before := time.Now()
			log.Debug(">>> incoming HTTP request",
				zap.String("http_method", r.Method),
				zap.Time("time", before.UTC()),
			)

			next.ServeHTTP(rw, r)

			log.Debug("<<< done HTTP request",
				zap.Int("status_code", rw.GetStatusCode()),
				zap.Duration("latency", time.Since(before)),
			)
		})
	}
}

func Recoverer() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.FromContext(r.Context())
			responseHandler := response.NewHTTPResponseHandler(log, w)

			defer func() {
				if p := recover(); p != nil {
					responseHandler.PanicResponse(p, "during handle HTTP request got unexpected panic")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func AuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.FromContext(r.Context())
			responseHandler := response.NewHTTPResponseHandler(log, w)

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				responseHandler.JSONResponse(http.StatusUnauthorized, response.ErrorResponse{Error: "unauthorized"})
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				responseHandler.JSONResponse(http.StatusUnauthorized, response.ErrorResponse{Error: "invalid auth token format"})
				return
			}

			tokenStr := parts[1]
			claims := &Claims{}

			token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				responseHandler.JSONResponse(http.StatusUnauthorized, response.ErrorResponse{Error: "invalid or expired token"})
				return
			}

			ctx := context.WithValue(r.Context(), ClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRole(allowedRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.FromContext(r.Context())
			responseHandler := response.NewHTTPResponseHandler(log, w)

			claims, ok := r.Context().Value(ClaimsKey).(*Claims)
			if !ok || claims.Role != allowedRole {
				responseHandler.JSONResponse(http.StatusForbidden, response.ErrorResponse{Error: "forbidden"})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func UserClaimsFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(ClaimsKey).(*Claims)
	return claims, ok
}
