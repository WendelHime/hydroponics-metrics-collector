package middlewares

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/WendelHime/hydroponics-metrics-collector/internal/shared/errors"
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// EnsureValidToken is a middleware that will check the validity of our JWT.
func EnsureValidToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		issuerURL, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/")
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse the issuer url")
		}

		r.Header.Set("Authorization", r.Header.Get("X-Endpoint-API-UserInfo"))

		provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

		jwtValidator, err := validator.New(
			provider.KeyFunc,
			validator.RS256,
			issuerURL.String(),
			[]string{os.Getenv("AUTH0_AUDIENCE")},
			validator.WithCustomClaims(
				func() validator.CustomClaims {
					return &CustomClaims{}
				},
			),
			validator.WithAllowedClockSkew(time.Minute),
		)
		if err != nil {
			log.Error().Err(err).Msg("Failed to set up the jwt validator")
		}

		errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
			log.Warn().Err(err).Msg("Encountered error while validating JWT")
			errors.RenderErr(w, r, errors.UnauthorizedErr)
		}

		middleware := jwtmiddleware.New(
			jwtValidator.ValidateToken,
			jwtmiddleware.WithErrorHandler(errorHandler),
		)

		middleware.CheckJWT(next).ServeHTTP(w, r)
	})
}

// CustomClaims contains custom data we want from the token.
type CustomClaims struct {
	Issuer string `json:"iss"`
	Scope  string `json:"scope"`
}

// Validate does nothing for this example, but we need
// it to satisfy validator.CustomClaims interface.
func (c CustomClaims) Validate(ctx context.Context) error {
	return nil
}

// HasScope checks whether our claims have a specific scope.
func (c CustomClaims) HasScope(expectedScope string) bool {
	result := strings.Split(c.Scope, " ")
	for i := range result {
		if result[i] == expectedScope {
			return true
		}
	}

	return false
}

func HasScope(scopes string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
			claims := token.CustomClaims.(*CustomClaims)
			if !claims.HasScope(scopes) {
				errors.RenderErr(w, r, errors.ForbiddenErr)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func UserMatches(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")

		token := r.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
		claims := token.CustomClaims.(*CustomClaims)
		if claims.Issuer != userID {
			errors.RenderErr(w, r, errors.ForbiddenErr)
			return
		}
		next.ServeHTTP(w, r)
	})
}
