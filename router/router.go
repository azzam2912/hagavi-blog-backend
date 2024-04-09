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
	apiBlog := app.Group("/api/blog")
	post := apiBlog.Group("/post")
	post.Post("/", func(c *fiber.Ctx) error {
		return controller.CreateBlogPostHandler(c, db)
	})
	post.Get("/:id", func(c *fiber.Ctx) error {
		return controller.GetBlogPostHandler(c, db)
	})
	post.Get("/", func(c *fiber.Ctx) error {
		return controller.GetAllBlogPostHandler(c, db)
	})
	post.Put("/:id", func(c *fiber.Ctx) error {
		return controller.UpdateBlogPostHandler(c, db)
	})
	post.Delete("/:id", func(c *fiber.Ctx) error {
		return controller.DeleteBlogPostHandler(c, db)
	})
	return app
}
