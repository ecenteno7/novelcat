package main

import (
	"booksapp/internal/routes"
	"booksapp/internal/services"
	_ "embed"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	// Create a new Pocketbase application instance
	app := pocketbase.New()

	// os.Setenv()

	// Add middleware to serve static files from the provided public directory
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// Serve webpage
		hostname, err := os.Hostname()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		webPath := ""
		if strings.Contains(hostname, ".local") {
			webPath = "../client/build"
		} else {
			webPath = "/pb/web"
		}
		fmt.Printf("Webpath: %s", webPath)
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS(webPath), false))

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
