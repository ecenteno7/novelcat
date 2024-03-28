package routes

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

// RegisterBookRoutes registers book-related routes
func RegisterBookRoutes(e *core.ServeEvent, app *pocketbase.PocketBase) error {
	e.Router.GET("/getBooksByUserId/:userId", func(c echo.Context) error {
		userId := c.PathParam("userId")
		log.Printf("Called get books by user id: %v", userId)

		// Retrieve multiple "articles" collection records by a custom dbx expression(s)
		records, err := app.Dao().FindRecordsByExpr("book_user",
			dbx.NewExp("user_id = {:userId}", dbx.Params{"userId": userId}),
		)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "ERROR"})
		}

		// Process retrieved records
		for _, element := range records {
			log.Println(element.Get("book_id"))
			if errs := app.Dao().ExpandRecord(element, []string{"book_id", "bookshelf_type_id"}, nil); len(errs) > 0 {
				log.Printf("failed to expand: %v", errs)
			}
			log.Println(element.ExpandedOne("book_id").Get("title"))
		}

		return c.JSON(http.StatusOK, map[string]interface{}{"message": "Hello " + userId, "books": records})
	})

	return nil
}
