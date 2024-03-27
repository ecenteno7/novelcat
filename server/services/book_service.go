package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"booksapp/models"

	"github.com/joho/godotenv"
	"github.com/pocketbase/pocketbase"
	pbModels "github.com/pocketbase/pocketbase/models"
)

// Fetches a list of book data from NYT Books API
func FetchBookData() (models.APIResponse, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	nytBaseUrl := "https://api.nytimes.com/svc/books/v3"
	nytApiKey := os.Getenv("NYT_API_KEY")

	requestURL := fmt.Sprintf(
		"%s/lists/2008-07-01/hardcover-fiction.json?api-key=%s",
		nytBaseUrl,
		nytApiKey,
	)

	res, err := http.Get(requestURL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close() // Close the response body after processing

	// Parse JSON response
	var response models.APIResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	return response, err
}

// Print information about each book
func printBookResultsInfo(books []models.Book) {
	for _, book := range books {
		fmt.Println("Title:", book.Title)
		fmt.Println("Description:", book.Description)
	}
}

// InsertBooks inserts the given books into the database
func InsertBooks(db *pocketbase.PocketBase, books []models.Book) error {
	// Iterate over books and save each one

	collection, err := db.Dao().FindCollectionByNameOrId("books")
	if err != nil {
		return err
	}

	for _, book := range books {
		record := pbModels.NewRecord(collection)

		record.Set("title", book.Title)
		record.Set("description", book.Description)

		if err := db.Dao().SaveRecord(record); err != nil {
			return err
		}
	}
	return nil
}

func PopulateBooksInDB(db *pocketbase.PocketBase) {
	// Fetch book data from NYT API
	log.Printf("Starting populate books in db")
	data, err := FetchBookData()
	if err != nil {
		log.Fatal(err)
	}

	// Extract books from API response
	books := data.Results.Books

	// Print information about stored books
	printBookResultsInfo(books)

	// Insert fetched book data into the database
	if err := InsertBooks(db, books); err != nil {
		log.Fatal(err)
	}
}
