package main

import (
	"database/sql"
	"fmt"
	"log"

	"populate-db/models"
	"populate-db/services"

	_ "github.com/mattn/go-sqlite3"
)

func handleDbConnection() {
	// Initialize PocketBase database connection
	db, err := sql.Open("sqlite3", "../pb_data/database.db")
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	defer db.Close()

	// // Perform CRUD operations
	// // Example: Insert a new book record into the database
	// book := models.Book{
	// 	Title:       "Sample Book",
	// 	Author:      "John Doe",
	// 	Description: "A sample book description",
	// 	// Add other fields as needed
	// }
	// if err := db.Create(&book).Error; err != nil {
	// 	log.Fatal("Error inserting book record:", err)
	// }

	// // Example: Query books from the database
	// var books []models.Book
	// if err := db.Find(&books).Error; err != nil {
	// 	log.Fatal("Error querying books:", err)
	// }

	// // Process retrieved books
	// for _, b := range books {
	// 	log.Println("Book Title:", b.Title)
	// 	// Print other book details as needed
	// }
}

// Print information about each book
func printBookResultsInfo(books []models.Book) {
	for _, book := range books {
		fmt.Println("Title:", book.Title)
		fmt.Println("Author:", book.Author)
		fmt.Println("Description:", book.Description)
		fmt.Println("Publisher:", book.Publisher)
		fmt.Println("Price:", book.Price)
		fmt.Println("Book Image URL:", book.BookImage)
		fmt.Println("Buy Links:")
		for _, link := range book.BuyLinks {
			fmt.Println("  -", link.Name+":", link.URL)
		}
		fmt.Println()
	}
}

func main() {

	handleDbConnection()

	data, err := services.FetchBookData()
	if err != nil {
		log.Fatal(err)
	}

	books := data.Results.Books

	printBookResultsInfo(books)
}
