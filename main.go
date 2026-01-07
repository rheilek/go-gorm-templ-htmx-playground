package main

import (
	"fmt"
	"net/http"

	decl "github.com/rheilek/go-gorm-templ-htmx-playground/declarative"
)

type Book struct {
	ISBN   string `gorm:"column:ISBN;primaryKey" label:"ISBN"`
	Title  string `gorm:"column:Title" label:"Title"`
	Author string `gorm:"column:Author" label:"Author"`
}

func main() {

	db := decl.Setup(Book{})

	db.Create(&Book{ISBN: "978-3-16-148410-0", Title: "Go Programming", Author: "John Doe"})
	db.Create(&Book{ISBN: "978-1-23-456789-7", Title: "Advanced Go", Author: "Jane Smith"})
	db.Create(&Book{ISBN: "978-0-12-345678-9", Title: "Go Concurrency", Author: "Alice Johnson"})
	db.Create(&Book{ISBN: "978-9-87-654321-0", Title: "Go Web Development", Author: "Bob Brown"})

	fmt.Println("Listening on :8080")
	fmt.Println(http.ListenAndServe(":8080", nil))
}
