package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

var db *sqlx.DB

func main() {
	var err error
	dsn := "user=jamal password=1234 dbname=db_employee sslmode=disable" // Замените на свои данные
	db, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}

	e := echo.New()

	e.GET("/users", getUsers)
	e.POST("/users", createUser)
	e.GET("/users/:id", getUser)
	e.PUT("/users/:id", updateUser)
	e.DELETE("/users/:id", deleteUser)

	e.Logger.Fatal(e.Start(":8080"))
}

func getUsers(c echo.Context) error {
	users := []User{}

	err := db.Select(&users, "SELECT * FROM users")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, users)
}

func createUser(c echo.Context) error {
	user := User{}

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	query := "INSERT INTO users (name, phone, address) VALUES ($1, $2, $3) RETURNING id"

	err := db.QueryRow(query, user.Name, user.Phone, user.Address).Scan(&user.ID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, user)
}

func getUser(c echo.Context) error {
	id := c.Param("id")

	user := User{}

	query := "SELECT * FROM users WHERE id = $1"

	err := db.Get(&user, query, id)

	if err != nil {
		// Возвращаем 404 если пользователь не найден
		return c.JSON(http.StatusNotFound, fmt.Sprintf("User with ID %s not found", id))
	}

	return c.JSON(http.StatusOK, user)
}

func updateUser(c echo.Context) error {
	id := c.Param("id")

	user := User{}

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	query := "UPDATE users SET name = $1, phone = $2, address = $3 WHERE id = $4"

	_, err := db.Exec(query, user.Name, user.Phone, user.Address, id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

func deleteUser(c echo.Context) error {
	id := c.Param("id")

	query := "DELETE FROM users WHERE id = $1"

	_, err := db.Exec(query, id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
