// book.go
package models

import (
    "github.com/pocketbase/pocketbase/models"
    "github.com/pocketbase/pocketbase/tools/types"

)

// ensures that the Article struct satisfy the models.Model interface
var _ models.Model = (*Book)(nil)

type Book struct {
    models.BaseModel

    Title        string         `db:"title" json:"title"`
    Genre       string         `db:"genre" json:"genre"`
    AuthorID    int            `db:"author_id" json:"authorId"`
    Description string         `db:"description" json:"description"`
    Published   types.DateTime `db:"published" json:"published"`
}

func (m *Book) TableName() string {
    return "books" // the name of your collection
}
