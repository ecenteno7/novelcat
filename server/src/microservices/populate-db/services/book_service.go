package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"populate-db/models"
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
