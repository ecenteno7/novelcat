package main

import (
	"booksapp/routes"
	"booksapp/services"
	"log"
	"os"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	// Create a new Pocketbase application instance
	app := pocketbase.New()

	// Add middleware to serve static files from the provided public directory
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// Serve webpage
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("../../client/build"), false))

		// Register book routes
		routes.RegisterBookRoutes(e, app)

		services.PopulateBooksInDB(app)
		return nil
	})

	// Start the Pocketbase application
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
