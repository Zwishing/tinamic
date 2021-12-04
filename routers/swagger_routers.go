package routers

import (
	swagger "github.com/arsmn/fiber-swagger/v2"
	_ "github.com/arsmn/fiber-swagger/v2/example/docs"
	"github.com/gofiber/fiber/v2"
)

// SwaggerRoute func for describe group of API Docs routes.
func SwaggerRoute(a *fiber.App) {
	// Create routes group.
	route := a.Group("/swagger")

	// Routes for GET method:
	route.Get("*", swagger.Handler) // get one user by ID
	route.Get("/swagger/*", swagger.New(swagger.Config{ // custom
		URL: "http://example.com/doc.json",
		DeepLinking: false,
		// Expand ("list") or Collapse ("none") tag groups by default
		DocExpansion: "none",
		// Prefill OAuth ClientId on Authorize popup
		OAuth: &swagger.OAuthConfig{
			AppName:  "OAuth Provider",
			ClientId: "21bb4edc-05a7-4afc-86f1-2e151e4ba6e2",
		},
		// Ability to change OAuth2 redirect uri location
		OAuth2RedirectUrl: "http://localhost:8080/swagger/oauth2-redirect.html",
	}))
}
