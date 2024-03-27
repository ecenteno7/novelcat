package main

import (
	"fmt"
	"log"
	"populate-db/models"
	"populate-db/services"
)

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

	data, err := services.FetchBookData()
	if err != nil {
		log.Fatal(err)
	}

	books := data.Results.Books

	printBookResultsInfo(books)
}
