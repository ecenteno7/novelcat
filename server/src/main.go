package main

import (
	"log"
	// "fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"

    "github.com/pocketbase/dbx"

)

func main() {
	// Create a new Pocketbase application instance
	app := pocketbase.New()

	// Add middleware to serve static files from the provided public directory
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("../../client/build"), false))

		return nil
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/getBooksByUserId/:userId", func(c echo.Context) error {
			userId := c.PathParam("userId")
			log.Printf("Called get books by user id: %v", userId)

			// retrieve multiple "articles" collection records by a custom dbx expression(s)
			records, err := app.Dao().FindRecordsByExpr("book_user",
				dbx.NewExp("user_id = {:userId}", dbx.Params{"userId": userId}),
			)

			if err != nil {
				return c.JSON(500, map[string]string{"message": "ERROR"})
			} else {

				for _, element := range records {
					log.Println(element.Get("book_id"))
					if errs := app.Dao().ExpandRecord(element, []string{"book_id", "bookshelf_type_id"}, nil); len(errs) > 0 {
						log.Printf("failed to expand: %v", errs)
					}
					log.Println(element.ExpandedOne("book_id").Get("title"))
				}
				return c.JSON(http.StatusOK, map[string]interface{}{"message": "Hello " + userId, "books": records })
			}

			return nil
		})

		return nil
	})

	// Start the Pocketbase application
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}

}

