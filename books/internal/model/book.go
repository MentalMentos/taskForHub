package model

import "github.com/kamva/mgm/v3"

// Book - структура книги
type Book struct {
	mgm.DefaultModel `bson:",inline"` // Встроенные ID и таймстампы
	Author           string           `json:"author" bson:"author"`
	Title            string           `json:"title" bson:"title"`
}
