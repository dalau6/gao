package clients

import (
	signin "clients/gen/signin"
	"context"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"goa.design/goa/v3/security"
)

// signin service example implementation.
// The example methods log the requests and return zero values.
type signinsrvc struct {
	logger *log.Logger
}

// NewSignin returns the signin service implementation.
func NewSignin(logger *log.Logger) signin.Service {
	return &signinsrvc{logger}
}

// BasicAuth implements the authorization logic for service "signin" for the
// "basic" security scheme.
func (s *signinsrvc) BasicAuth(ctx context.Context,
	user, pass string, scheme *security.BasicScheme) (context.Context,
	error) {

	if user != "gopher" && pass != "academy" {
		return ctx, signin.
			Unauthorized("invalid username and password combination")
	}

	return ctx, nil
}

// Creates a valid JWT
func (s *signinsrvc) Authenticate(ctx context.Context,
	p *signin.AuthenticatePayload) (res *signin.Creds,
	err error) {

	// create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"nbf":    time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
		"iat":    time.Now().Unix(),
		"exp":    time.Now().Add(time.Duration(9) * time.Minute).Unix(),
		"scopes": []string{"api:read", "api:write"},
	})

	s.logger.Printf("user '%s' logged in", p.Username)

	// note that if "SignedString" returns an error then it is returned as
	// an internal error to the client
	t, err := token.SignedString(Key)
	if err != nil {
		return nil, err
	}

	res = &signin.Creds{
		JWT: t,
	}

	return
}
