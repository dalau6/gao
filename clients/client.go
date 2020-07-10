package clients

import (
	client "clients/gen/client"
	"context"
	"log"

	jwt "github.com/dgrijalva/jwt-go"
	"goa.design/goa/v3/security"
)

var (
	// Key is the key used in JWT authentication
	Key = []byte("secret")
)

// client service example implementation.
// The example methods log the requests and return zero values.
type clientsrvc struct {
	logger *log.Logger
}

// NewClient returns the client service implementation.
func NewClient(logger *log.Logger) client.Service {
	return &clientsrvc{logger}
}

// JWTAuth implements the authorization logic for service "client" for the
// "jwt" security scheme.
func (s *clientsrvc) JWTAuth(ctx context.Context,
	token string, scheme *security.JWTScheme) (context.Context,
	error) {

	claims := make(jwt.MapClaims)

	// authorize request
	// 1. parse JWT token, token key is hardcoded to "secret" in this example
	_, err := jwt.ParseWithClaims(token,
		claims, func(_ *jwt.Token) (interface{},
			error) {
			return Key, nil
		})
	if err != nil {
		s.logger.Print("Unable to obtain claim from token, it's invalid")
		return ctx, client.Unauthorized("invalid token")
	}

	s.logger.Print("claims retrieved, validating against scope")
	s.logger.Print(claims)

	// 2. validate provided "scopes" claim
	if claims["scopes"] == nil {
		s.logger.Print("Unable to get scope since the scope is empty")
		return ctx, client.InvalidScopes("invalid scopes in token")
	}
	scopes, ok := claims["scopes"].([]interface{})
	if !ok {
		s.logger.Print("An error occurred when retrieving the scopes")
		s.logger.Print(ok)
		return ctx, client.InvalidScopes("invalid scopes in token")
	}
	scopesInToken := make([]string, len(scopes))
	for _, scp := range scopes {
		scopesInToken = append(scopesInToken, scp.(string))
	}
	if err := scheme.Validate(scopesInToken); err != nil {
		s.logger.Print("Unable to parse token, check error below")
		return ctx, client.InvalidScopes(err.Error())
	}
	return ctx, nil

}

// Add implements add.
func (s *clientsrvc) Add(ctx context.Context,
	p *client.AddPayload) (err error) {
	s.logger.Print("client.add started")
	newClient := client.ClientManagement{
		ClientID:      p.ClientID,
		ClientName:    p.ClientName,
		ContactName:   p.ContactName,
		ContactEmail:  p.ContactEmail,
		ContactMobile: p.ContactMobile,
	}
	err = CreateClient(&newClient)
	if err != nil {
		s.logger.Print("An error occurred...")
		s.logger.Print(err)
		return
	}
	s.logger.Print("client.add completed")
	return
}

// Get implements get.
func (s *clientsrvc) Get(ctx context.Context,
	p *client.GetPayload) (res *client.ClientManagement,
	err error) {
	s.logger.Print("client.get started")
	result, err := GetClient(p.ClientID)
	if err != nil {
		s.logger.Print("An error occurred...")
		s.logger.Print(err)
		return
	}
	s.logger.Print("client.get completed")
	return &result, err
}

// Show implements show.
func (s *clientsrvc) Show(ctx context.Context,
	p *client.ShowPayload) (res client.ClientManagementCollection,
	err error) {
	s.logger.Print("client.show started")
	res, err = ListClients()
	if err != nil {
		s.logger.Print("An error occurred...")
		s.logger.Print(err)
		return
	}
	s.logger.Print("client.show completed")
	return
}
