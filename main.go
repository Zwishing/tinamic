package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"time"
	. "tinamic/routers"
)

func init()  {
	fmt.Println(time.Time{})
	fmt.Println(time.Now())
}
func main() {

	app := fiber.New()

	SwaggerRoute(app)

	api := app.Group("/api/v1")
	RegisterAPI(api)


	log.Fatal(app.Listen(":3001"))
}