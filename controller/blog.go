package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"hagavi-blog-go/models"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

const (
	blogPostTable = "blog_posting"
)

func CreateBlogPostHandler(c *fiber.Ctx, db *sql.DB) error {
	var newPost models.BlogPost
	err := c.BodyParser(&newPost)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}
	location := time.FixedZone("GMT+7", 7*60*60)
	createdTime := time.Now().In(location)
	sqlStatementCreate := fmt.Sprintf(`INSERT INTO %s (title, content, created_at, updated_at, author) VALUES ($1, $2, $3, $4, $5)`, blogPostTable)
	result, err := db.Exec(sqlStatementCreate, newPost.Title, newPost.Content, createdTime, createdTime, newPost.Author)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	if rowsAffected, err := result.RowsAffected(); err == nil && rowsAffected > 0 {
		newPostJSON, err := json.Marshal(newPost)
		if err == nil {
			return c.JSON(http.StatusOK, string(newPostJSON))
		}
	}
	return c.Status(http.StatusInternalServerError).SendString("Failed to create blog post")
}

func GetBlogPostHandler(c *fiber.Ctx, db *sql.DB) error {
	id := c.Params("id") // Get id from URL parameter
	if id == "" {
		return c.Status(http.StatusNotFound).SendString("Post ID Not Found")
	}

	var post models.BlogPost
	sqlStatementGet := fmt.Sprintf(`SELECT id, title, content, created_at, updated_at, author FROM %s WHERE id = $1`, blogPostTable)
	rowResult := db.QueryRow(sqlStatementGet, id)

	err := rowResult.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt, &post.Author)

	if err == nil {
		return c.Status(http.StatusOK).JSON(post)
	}
	if err == sql.ErrNoRows {
		return c.Status(http.StatusNotFound).SendString("Post Not Found")
	}
	return c.Status(http.StatusInternalServerError).SendString("Cannot Retrieve Post")

}

func UpdateBlogPostHandler(c *fiber.Ctx, db *sql.DB) error {
	id := c.Params("id") // Get id from URL parameter
	if id == "" {
		return c.Status(http.StatusBadRequest).SendString("Missing Post ID Parameter")
	}

	var updatedPost models.BlogPost
	err := c.BodyParser(&updatedPost)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}
	location := time.FixedZone("GMT+7", 7*60*60)
	updatedTime := time.Now().In(location)
	sqlStatementUpdated := fmt.Sprintf(`UPDATE %s SET title=$1, content=$2, updated_at=$3, author=$4 WHERE id=$5`, blogPostTable)
	result, err := db.Exec(sqlStatementUpdated, updatedPost.Title, updatedPost.Content, updatedTime, updatedPost.Author, id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	if rowsAffected, err := result.RowsAffected(); err == nil && rowsAffected > 0 {
		return c.SendStatus(http.StatusOK)
	}
	return c.Status(http.StatusInternalServerError).SendString("Failed to update blog post")
}

func DeleteBlogPostHandler(c *fiber.Ctx, db *sql.DB) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).SendString("Missing Post ID Parameter")
	}

	sqlStatementDelete := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, blogPostTable)
	result, err := db.Exec(sqlStatementDelete, id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	if rowsAffected, err := result.RowsAffected(); err == nil && rowsAffected > 0 {
		return c.SendStatus(http.StatusOK)
	}
	return c.Status(http.StatusInternalServerError).SendString("Failed to delete blog post")
}

func GetAllBlogPostHandler(c *fiber.Ctx, db *sql.DB) error {
	sqlStatementGetAll := fmt.Sprintf(`SELECT id, title, content, created_at, updated_at, author FROM %s`, blogPostTable)
	rows, err := db.Query(sqlStatementGetAll)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	defer rows.Close()
	blogPosts := []models.BlogPost{}
	for rows.Next() {
		var post models.BlogPost
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt, &post.Author)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
		blogPosts = append(blogPosts, post)
	}
	err = rows.Err()
	if err == nil {
		return c.Status(http.StatusOK).JSON(blogPosts)
	}
	return c.Status(http.StatusInternalServerError).SendString(err.Error())
}
