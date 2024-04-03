package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"booksapp/internal/models"

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

// getFirstAndLastName accepts a full name string
// and returns the first name and last name
// as the first and second return values
func getFirstAndLastName(fullname string) (string, string) {
	// splitting the fullname single string by whitespace
	names := strings.Split(fullname, " ")

	return names[0], names[1]
}

// Print information about each book
func printBookResultsInfo(books []models.Book) {
	for _, book := range books {
		fmt.Println("Title:", book.Title)
		fmt.Println("Description:", book.Description)
		firstName, lastName := getFirstAndLastName(book.Author)
		fmt.Println("Author:", lastName, ", ", firstName)
	}
}

// InsertBooks inserts the given books into the database
func InsertBooks(db *pocketbase.PocketBase, books []models.Book) error {
	// Iterate over books and save each one

	bookCollection, err := db.Dao().FindCollectionByNameOrId("books")
	if err != nil {
		return err
	}

	for _, book := range books {

		authorRecordId, err := InsertAuthor(db, book.Author)
		if err != nil {
			return err
		}

		record := pbModels.NewRecord(bookCollection)

		record.Set("title", book.Title)
		record.Set("description", book.Description)
		record.Set("author_id", authorRecordId)

		if err := db.Dao().SaveRecord(record); err != nil {
			return err
		}
	}
	return nil
}

func InsertAuthor(db *pocketbase.PocketBase, authorName string) (string, error) {
	collection, err := db.Dao().FindCollectionByNameOrId("authors")
	if err != nil {
		return "", err
	}

	firstName, lastName := getFirstAndLastName(authorName)

	record := pbModels.NewRecord(collection)

	record.Set("last_name", lastName)
	record.Set("first_name", firstName)

	if err := db.Dao().SaveRecord(record); err != nil {
		return "", err
	}

	return record.Id, nil
}

func PopulateBooksInDB(db *pocketbase.PocketBase) error {
	// Fetch book data from NYT API
	log.Printf("Starting populate books in db")
	data, err := FetchBookData()
	if err != nil {
		db.Logger().Error("Error fetching book data from API")
		return err
	}

	// Extract books from API response
	books := data.Results.Books

	// Print information about stored books
	printBookResultsInfo(books)

	// Insert fetched book data into the database
	if err := InsertBooks(db, books); err != nil {
		db.Logger().Error("Error inserting books into the database")
		return err
	}
	return nil
}
