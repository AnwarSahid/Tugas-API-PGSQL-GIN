package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "pgsql_go"
)

type Book struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
}

var (
	db  *sql.DB
	err error
)
var router = gin.Default()

func main() {
	psqlInfo := fmt.Sprintf("host=  %s port=%d user=%s password=%s dbname= %s sslmode=disable", host, port, user, password, dbname)
	db, err = sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("successfuly connect database")

	getAllBook()
	getBookById()
	createBook()
	updateBook()
	deleteBook()

	router.Run(":4000")
}

func getAllBook() {
	router.GET("/buku", func(c *gin.Context) {
		var results = []Book{}
		rows, err := db.Query(`select * from books`)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		for rows.Next() {
			var book = Book{}

			err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.Description)
			if err != nil {
				panic(err)
			}
			results = append(results, book)
		}
		c.JSON(http.StatusOK, gin.H{
			"book": results,
		})
	})
}

func getBookById() {
	router.GET("/buku/:bookID", func(c *gin.Context) {
		bookID := c.Param("bookID")
		// var book = []Book{}
		stmt, err := db.Prepare("SELECT id, title, author, description FROM books WHERE id=$1")
		if err != nil {
			panic(err)
		}
		var book Book
		err = stmt.QueryRow(bookID).Scan(&book.ID, &book.Title, &book.Author, &book.Description)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("Book not found")
			} else {
				panic(err)
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"id":          book.ID,
			"title":       book.Title,
			"author":      book.Author,
			"description": book.Description,
		})
	})
}

func createBook() {
	router.POST("/buku/create", func(c *gin.Context) {
		var book = Book{}
		var newBook Book

		if err := c.ShouldBindJSON((&newBook)); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		query := `insert into books(title,author, description) values($1, $2, $3) Returning *`
		err = db.QueryRow(query, newBook.Title, newBook.Author, newBook.Description).Scan(&book.ID, &book.Title, &book.Author, &book.Description)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Data Created Successfuly",
			"data":    newBook,
		})
	})

}

func updateBook() {
	router.PUT("/buku/update/:ID", func(c *gin.Context) {
		bookID := c.Param("ID")
		var updatedBook Book

		if err := c.ShouldBindJSON(&updatedBook); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		query := `update books set title = $2 ,author =$3, description=$4 where id= $1`
		res, err := db.Exec(query, bookID, updatedBook.Title, updatedBook.Author, updatedBook.Description)
		if err != nil {
			panic(err)
		}
		count, err := res.RowsAffected()
		c.JSON(http.StatusOK, gin.H{
			"message":     "Data updated Successfuly",
			"banyak data": count,
		})
	})
}

func deleteBook() {
	router.DELETE("buku/delete/:ID", func(c *gin.Context) {
		bookID := c.Param("ID")
		query := `DELETE from books WHERE id = $1`
		_, err := db.Exec(query, bookID)
		if err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Data Deleted succesfully",
		})
	})
}
