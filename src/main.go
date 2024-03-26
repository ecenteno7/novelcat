package main

import (
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
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./pb_public"), false))

		// // Create a new book collection if it doesn't exist
		// if err := createBookCollectionIfNotExist(app); err != nil {
		//     log.Fatal(err)
		// }

		// // Create a sample book
		// if err := createSampleBook(app); err != nil {
		//     log.Fatal(err)
		// }

		// Log a message indicating that the server has started
		log.Println("Server started and sample book created.")
		return nil
	})

	// Start the Pocketbase application
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}

}

// createBookCollectionIfNotExist creates a new collection named "books" if it doesn't already exist
// func createBookCollectionIfNotExist(app *pocketbase.PocketBase) error {
//     // Check if the "books" collection exists
//     exists, err := app.Dao().HasCollection("books")
//     if err != nil {
//         return err
//     }
//     // If not, create it
//     if !exists {
//         if err := app.Dao().CreateCollection("books"); err != nil {
//             return err
//         }
//     }
//     return nil
// }

// createSampleBook creates a sample book in the "books" collection
// func createSampleBook(app *pocketbase.PocketBase) error {
//     // Create a new book instance
//     book := &models.Book{
//         Title:       "Sample Book",
//         AuthorID:    1, // Update with appropriate author ID
//         Description: "This is a sample book created upon serving Pocketbase.",
//         // Add any additional fields as needed
//     }

//     // Save the book to the "books" collection
//     if err := app.Dao().Collection("books").Save(book); err != nil {
//         return err
//     }

//     return nil
// }
