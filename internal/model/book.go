package model

type Book struct {
	ID     string `json:"id" gorm:"primaryKey"`
	Title  string `json:"title"`
	Author string `json:"author"`
}
