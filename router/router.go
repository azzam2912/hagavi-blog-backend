package router

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"hagavi-blog-go/controller"
)

func NewConnection(db *sql.DB) *fiber.App {
	app := fiber.New()
	app.Use(cors.New())
	api := app.Group("/api")
	post := api.Group("/post")
	post.Post("/", func(c *fiber.Ctx) error {
		return contoller.CreateBlogPostHandler(c, db)
	})
	post.Get("/:id", func(c *fiber.Ctx) error {
		return contoller.GetBlogPostHandler(c, db)
	})
	post.Get("/", func(c *fiber.Ctx) error {
		return contoller.GetAllBlogPostHandler(c, db)
	})
	post.Put("/:id", func(c *fiber.Ctx) error {
		return contoller.UpdateBlogPostHandler(c, db)
	})
	post.Delete("/:id", func(c *fiber.Ctx) error {
		return contoller.DeleteBlogPostHandler(c, db)
	})
	return app
}
