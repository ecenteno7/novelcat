package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"booksapp/internal/models"

	fp "github.com/amonsat/fullname_parser"
	"github.com/joho/godotenv"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	pbModels "github.com/pocketbase/pocketbase/models"
)

// Fetches a list of book data from NYT Books API
func FetchBookList(fetchDate string, listName string) ([]models.Book, error) {
	// 2008-07-01 - fetch date example
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	nytBaseUrl := "https://api.nytimes.com/svc/books/v3"
	nytApiKey := os.Getenv("NYT_API_KEY")
	fmt.Printf("Date: %s\n", fetchDate)
	requestURL := fmt.Sprintf(
		"%s/lists/%s/%s.json?api-key=%s",
		nytBaseUrl,
		fetchDate,
		listName,
		nytApiKey,
	)

	fmt.Printf("Request URL: %s\n", requestURL)

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

	return response.Results.Books, err
}

// getFirstAndLastName accepts a full name string
// and returns the first name and last name
// as the first and second return values
func getFirstAndLastName(fullname string) (string, string, string) {
	// splitting the fullname single string by whitespace
	// names := strings.Split(fullname, " ")
	parsedFullname := fp.ParseFullname(fullname)
	first := parsedFullname.First
	middle := parsedFullname.Middle
	last := parsedFullname.Last
	return last, first, middle
}

// Print information about each book
func printBookResultsInfo(books []models.Book) {
	for _, book := range books {
		fmt.Println("Title:", book.Title)
		fmt.Println("Description:", book.Description)
		firstName, lastName, _ := getFirstAndLastName(book.Author)
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
		authorRecordId := []string{}
		if strings.Contains(book.Author, "and") {
			authors := strings.Split(book.Author, "and")
			for _, author := range authors {
				id, err := InsertAuthor(db, author)
				if err != nil {
					return err
				}
				authorRecordId = append(authorRecordId, id)
			}
		} else {
			id, err := InsertAuthor(db, book.Author)
			if err != nil {
				return err
			}
			authorRecordId = append(authorRecordId, id)
		}

		bookRecord, _ := db.Dao().FindFirstRecordByFilter(
			"books", "title = {:title} && author_id = {:author_id}",
			dbx.Params{"title": book.Title, "author_id": authorRecordId[0]},
		)

		if bookRecord != nil {
			log.Printf("Book already exists in db: %s", book.Title)
			return nil
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

	firstName, lastName, middleName := getFirstAndLastName(authorName)

	author, err := db.Dao().FindFirstRecordByFilter(
		"authors", "last_name = {:lastName} && firstName = {:firstName}",
		dbx.Params{"lastName": lastName, "firstName": firstName},
	)

	if author != nil {
		return author.Id, nil
	}

	record := pbModels.NewRecord(collection)

	record.Set("last_name", lastName)
	record.Set("first_name", firstName)
	record.Set("middle_name", middleName)

	if err := db.Dao().SaveRecord(record); err != nil {
		return "", err
	}

	return record.Id, nil
}

func fetchBooksBetweenDates(db *pocketbase.PocketBase, start time.Time, end time.Time, listName string) error {

	for date := start; date.Before(end); date = date.AddDate(0, 1, 0) {
		fetchDate := date.Format("2006-01-02")
		books, err := FetchBookList(fetchDate, listName)
		if err != nil {
			log.Printf("Error fetching book data for date %s: %v", fetchDate, err)
			continue
		}

		// Print information about stored books
		printBookResultsInfo(books)

		// Insert fetched book data into the database
		if err := InsertBooks(db, books); err != nil {
			log.Printf("Error inserting books into the database for date %s: %v", fetchDate, err)
			continue
		}
		time.Sleep(12 * time.Second)
	}

	return nil
}

func PopulateBooksInDB(db *pocketbase.PocketBase) error {
	// Fetch book data from NYT API for each month in between dates
	start := time.Date(2008, time.July, 1, 0, 0, 0, 0, time.UTC)
	today := time.Now()
	end := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, time.UTC)
	listName := "hardcover-fiction"

	err := fetchBooksBetweenDates(db, start, end, listName)

	if err != nil {
		log.Printf("Error fetching %s books between dates %v %v", listName, start, end)
	}

	return nil
}

func AddBookToUser(db *pocketbase.PocketBase, userId string, bookId string, bookshelf string) error {
	log.Printf("Called addbooktouser %s %s", userId, bookId)

	collection, err := db.Dao().FindCollectionByNameOrId("book_user")
	if err != nil {
		return err
	}
	record := pbModels.NewRecord(collection)

	record.Set("user_id", userId)
	record.Set("book_id", bookId)
	record.Set("bookshelf", bookshelf)

	if err := db.Dao().SaveRecord(record); err != nil {
		return err
	}

	return nil
}

func CreateUser(db *pocketbase.PocketBase, username string, name string, email string, password string) error {

	collection, err := db.Dao().FindCollectionByNameOrId("users")
	if err != nil {
		return err
	}
	record := pbModels.NewRecord(collection)

	record.Set("username", username)
	record.Set("name", name)
	record.Set("email", email)
	record.SetPassword(password)

	if err := db.Dao().SaveRecord(record); err != nil {
		return err
	}

	return nil
}
