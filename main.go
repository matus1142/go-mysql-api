package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type Book struct {
	Id     int    `json: "id"`
	Name   string `json: "name"`
	Author string `json: author`
}

type BookInsert struct {
	Name   string `json: "name"`
	Author string `json: author`
}

func GetHelloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"error":        false,
		"message":      "Welcome to RESTful CRUD API with NodeJS, Express, MYSQL",
		"written_by":   "Matus",
		"published_on": "https://matus.dev",
	})
}

// rerieve all books
func GetBooks(c *gin.Context) {
	// username := c.Query("username")
	// password := c.Query("password")
	var empty_check int = 0
	db, err := sql.Open("mysql", "root:@tcp(localhost)/nodejs_api")

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}
	rows, err := db.Query("SELECT * FROM books")
	data := []Book{}
	for rows.Next() {
		var id int
		var name string
		var author string
		var created_at string
		var updated_at string
		err = rows.Scan(&id, &name, &author, &created_at, &updated_at)
		data = append(data, Book{Id: id, Name: name, Author: author})
		fmt.Println(data)
		checkErr(err)
		empty_check = 1
	}
	if empty_check == 0 {
		c.JSON(http.StatusOK, gin.H{
			"error":   true,
			"message": "Books table is empty",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"error":   false,
			"data":    data,
			"message": "Successfully retrieved all books",
		})
	}

	defer db.Close()
}

//add a new book
func PostBook(c *gin.Context) {
	var data BookInsert
	c.BindJSON(&data)
	db, err := sql.Open("mysql", "root:@tcp(localhost)/nodejs_api")

	if err != nil {
		panic(err.Error())
	}

	if data.Name == "" || data.Author == "" {
		c.JSON(http.StatusOK, gin.H{
			"error":   true,
			"message": "Pleaase provide book name and author",
		})
	} else {
		query := "INSERT INTO books (name, author) VALUES(?,?)"
		_, err := db.ExecContext(context.Background(), query, data.Name, data.Author)
		if err != nil {
			panic(err.Error())
		}
		c.JSON(http.StatusOK, gin.H{"data": data})
	}

	fmt.Println(data)

	defer db.Close()
}

// retrieve book by id
func GetBookById(c *gin.Context) {

	id := c.Param("id")
	// fmt.Println("String =", id)
	// fmt.Printf("Type =%T\n", id)
	if id != "" {
		var empty_check int = 0

		intId, _ := strconv.Atoi(id) //convert string to integer
		// fmt.Println("Int =", intId)
		// fmt.Printf("Type =%T\n", intId)
		db, err := sql.Open("mysql", "root:@tcp(localhost)/nodejs_api")

		// if there is an error opening the connection, handle it
		if err != nil {
			panic(err.Error())
		}

		rows, err := db.Query("SELECT * FROM books WHERE id = ?", intId)
		// fmt.Println(rows)
		data := []Book{}
		for rows.Next() {
			var id int
			var name string
			var author string
			var created_at string
			var updated_at string
			err = rows.Scan(&id, &name, &author, &created_at, &updated_at)
			data = append(data, Book{Id: id, Name: name, Author: author})
			fmt.Println(data)
			checkErr(err)
			empty_check = 1
		}
		if empty_check == 0 {
			c.JSON(http.StatusOK, gin.H{
				"error":   true,
				"message": "Book not found",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"error":   false,
				"data":    data,
				"message": "Successfully retrieved book data",
			})
		}

		defer db.Close()
	}
}

// updated book by id
func UpdateBook(c *gin.Context) {
	var updateData Book
	c.ShouldBindJSON(&updateData)
	id := updateData.Id
	name := updateData.Name
	author := updateData.Author
	fmt.Println("id:", updateData.Id, "name:", updateData.Name, "author:", updateData.Author)

	db, err := sql.Open("mysql", "root:@tcp(localhost)/nodejs_api")
	if err != nil {
		panic(err.Error())
	}
	if id != 0 {
		fmt.Println("OK")
		if name != "" && author != "" {
			_, err = db.Query("UPDATE books SET name = ?, author = ? WHERE id = ?", name, author, id)
		} else if name != "" {
			_, err = db.Query("UPDATE books SET name = ? WHERE id = ?", name, id)
		} else if author != "" {
			_, err = db.Query("UPDATE books SET author = ? WHERE id = ?", author, id)
		} else {
			c.JSON(http.StatusOK, gin.H{
				"error":   true,
				"message": "Pleaase provide book name or author",
			})
		}
		checkErr(err)
	}

	rows, err := db.Query("SELECT * FROM books WHERE id = ?", id)
	data := []Book{}
	for rows.Next() {
		var id int
		var name string
		var author string
		var created_at string
		var updated_at string
		err = rows.Scan(&id, &name, &author, &created_at, &updated_at)
		data = append(data, Book{Id: id, Name: name, Author: author})
		checkErr(err)
		c.JSON(http.StatusOK, gin.H{
			"error":   false,
			"data":    data,
			"message": "Successfully updated book data",
		})
	}

	defer db.Close()

}

//deleted book by id
func DeleteBook(c *gin.Context) {
	var id_exist_status = 0
	strid := c.Param("id")
	id, _ := strconv.Atoi(strid)

	db, err := sql.Open("mysql", "root:@tcp(localhost)/nodejs_api")
	if err != nil {
		panic(err.Error())
	}
	if strid != "" {
		_, err = db.Query("DELETE FROM books WHERE id = ?", id)
		checkErr(err)

		rows, err := db.Query("SELECT * FROM books WHERE id = ?", id)
		checkErr(err)
		for rows.Next() {
			id_exist_status = 1
		}
		if id_exist_status == 0 {
			c.JSON(http.StatusOK, gin.H{
				"error":   false,
				"message": "Successfully deleted book data",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"error":   true,
				"message": "Cannot deleted book data",
			})
		}

	}

	defer db.Close()
}

func main() {
	r := gin.Default()
	r.GET("/", GetHelloWorld)
	r.GET("/books", GetBooks)
	r.GET("/book/:id", GetBookById)
	r.POST("/book", PostBook)
	r.PUT("/book", UpdateBook)
	r.DELETE("/book/:id", DeleteBook)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
