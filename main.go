package main

import (
	"hagavi-blog-go/database"
	"hagavi-blog-go/router"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	db := database.Connect()
	app := router.NewConnection(db)
	defer db.Close()
	log.Fatal(app.Listen(":8080"))
}
