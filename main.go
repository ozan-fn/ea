package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/mongodb/v2"
	"github.com/joho/godotenv"
)

var store = mongodb.New(mongodb.Config{
	ConnectionURI: os.Getenv("MONGODB_URI"),
	Database:      "fiber",
	Collection:    "fiber_storage",
})

func main() {
	godotenv.Load()

	app := fiber.New()

	// 2. Hubungkan storage ke Session
	store := session.New(session.Config{
		Storage:        store,
		CookieSecure:   true,
		CookieHTTPOnly: true,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			return err
		}

		// 3. Gunakan sess.Get dan sess.Set seperti biasa!
		// Jauh lebih ringkas daripada manual JWT
		count := sess.Get("count")
		currCount := 0
		if count != nil {
			currCount = count.(int)
		}

		currCount++
		sess.Set("count", currCount)

		if err := sess.Save(); err != nil {
			return err
		}

		return c.SendString(fmt.Sprintf("Counter (Stateless & Simple): %d", currCount))
	})

	log.Fatal(app.Listen(":8080"))
}
