package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
)


func main () {
	env:= godotenv.Load()
	
	if env != nil {
		fmt.Println("Error loading .env file")
	}


	port := os.Getenv("PORT")

	if port == "" {
		port = ":3000"
	} else {
		port = ":" + port
	}

    app:= fiber.New(fiber.Config{
		IdleTimeout: 5  * time.Second,
	})

	app.Use(compress.New())

	
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!, first fiber app")
	})

	go func () {
		if err := app.Listen(port); err != nil {
			log.Panic(err)
		}

	}()

	c:= make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

   	<- c // Block main thread


	app.Shutdown()

	fmt.Println("Server is shutting down...")
}