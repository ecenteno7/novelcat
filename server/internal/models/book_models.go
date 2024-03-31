// book.go
package models

import (
	"github.com/pocketbase/pocketbase/models"
)

// ensures that the Article struct satisfy the models.Model interface
var _ models.Model = (*Book)(nil)

type Book struct {
	models.BaseModel

	Title       string `db:"title" json:"title"`
	AuthorID    int    `db:"author_id" json:"authorId"`
	Description string `db:"description" json:"description"`
	Author      string `json:"author"`
	// Published   types.DateTime `db:"published" json:"published"`
}

type Author struct {
	LastName   string `db:"last_name"`
	FirstName  string `db:"first_name"`
	MiddleName string `db:"middle_name"`
}

func (m *Book) TableName() string {
	return "books" // the name of your collection
}
