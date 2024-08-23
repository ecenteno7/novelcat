package routes

import (
	"booksapp/internal/services"
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

// RegisterBookRoutes registers book-related API routes
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
			for _, record := range element.ExpandedAll("book_id") {
				log.Println(record.Get("title"))
			}
		}

		return c.JSON(http.StatusOK, map[string]interface{}{"message": "Returning records for user " + userId, "books": records})
	})

	e.Router.GET("/seedDb", func(c echo.Context) error {
		err := services.PopulateBooksInDB(app)
		if err == nil {
			return c.JSON(http.StatusOK, map[string]interface{}{"message": "DB seeded successfully"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Encountered error while seeding DB: " + err.Error()})
	})

	e.Router.POST("/addBookToUser", func(c echo.Context) error {
		data := struct {
			UserId    string `json:"userId"`
			BookId    string `json:"bookId"`
			Bookshelf string `json:"bookshelf"`
		}{}
		if err := c.Bind(&data); err != nil {
			return apis.NewBadRequestError("Failed to read request data", err)
		}

		userId := data.UserId
		bookId := data.BookId
		bookshelf := data.Bookshelf

		err := services.AddBookToUser(app, userId, bookId, bookshelf)
		if err == nil {
			return c.JSON(http.StatusOK, map[string]interface{}{"message": "Succesfully added book to user."})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Encountered error in add book to user: " + err.Error()})
	})

	e.Router.POST("/addUser", func(c echo.Context) error {
		data := struct {
			Username string `json:"username" db:"username"`
			Name     string `json:"name" db:"name"`
			Email    string `json:"email" db:"email"`
			Password string `json:"password" db:"password"`
		}{}
		if err := c.Bind(&data); err != nil {
			return apis.NewBadRequestError("Failed to read request data", err)
		}

		username := data.Username
		name := data.Name
		email := data.Email
		pwd := data.Password

		err := services.CreateUser(app, username, name, email, pwd)
		if err == nil {
			return c.JSON(http.StatusOK, map[string]interface{}{"message": "Successfully added user to db."})
		}

		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Encountered error in add user to db: " + err.Error()})

	})

	return nil
}
