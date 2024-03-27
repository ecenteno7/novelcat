package models

// Book defined by NYT API response
type Book struct {
	Rank            int       `json:"rank"`
	RankLastWeek    int       `json:"rank_last_week"`
	WeeksOnList     int       `json:"weeks_on_list"`
	Asterisk        int       `json:"asterisk"`
	Dagger          int       `json:"dagger"`
	PrimaryISBN10   string    `json:"primary_isbn10"`
	PrimaryISBN13   string    `json:"primary_isbn13"`
	Publisher       string    `json:"publisher"`
	Description     string    `json:"description"`
	Price           string    `json:"price"`
	Title           string    `json:"title"`
	Author          string    `json:"author"`
	Contributor     string    `json:"contributor"`
	ContributorNote string    `json:"contributor_note"`
	BookImage       string    `json:"book_image"`
	BuyLinks        []BuyLink `json:"buy_links"`
}

// BuyLink defined by NYT API response - gives a link to purchase book.
type BuyLink struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
