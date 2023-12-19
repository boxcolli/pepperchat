package main

import (
	"fmt"
	"log"

	"github.com/boxcolli/pepperchat/service/internal/docrepo"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

const (
	projectId = "sandboxcolli"
)

func main() {
	// DB client
	var dr docrepo.DocRepo
	{
		var err error
		dr, err = docrepo.NewDocRepo(projectId)
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}
	}
    app := fiber.New()

    // GET /api/register
    app.Get("/api/*", func(c *fiber.Ctx) error {
        msg := fmt.Sprintf("✋ %s", c.Params("*"))
        return c.SendString(msg) // => ✋ register
    })

    // GET /flights/LAX-SFO
    app.Get("/flights/:from-:to", func(c *fiber.Ctx) error {
        msg := fmt.Sprintf("💸 From: %s, To: %s", c.Params("from"), c.Params("to"))
        return c.SendString(msg) // => 💸 From: LAX, To: SFO
    })

    // GET /dictionary.txt
    app.Get("/:file.:ext", func(c *fiber.Ctx) error {
        msg := fmt.Sprintf("📃 %s.%s", c.Params("file"), c.Params("ext"))
        return c.SendString(msg) // => 📃 dictionary.txt
    })

    // GET /john/75
    app.Get("/:name/:age/:gender?", func(c *fiber.Ctx) error {
        msg := fmt.Sprintf("👴 %s is %s years old", c.Params("name"), c.Params("age"))
        return c.SendString(msg) // => 👴 john is 75 years old
    })

    // GET /john
    app.Get("/:name", func(c *fiber.Ctx) error {
        msg := fmt.Sprintf("Hello, %s 👋!", c.Params("name"))
        return c.SendString(msg) // => Hello john 👋!
    })

	// POST /auth/signup/{username}
	app.Post("/auth/signup/:username", func(c *fiber.Ctx) error {
		id, err := dr.SignUp(c.Context(), c.Params("username"))

		switch err {
		case docrepo.ErrConflict:
			return c.SendStatus(fiber.StatusConflict)
		case docrepo.ErrUnknown:
			return c.SendStatus(fiber.StatusInternalServerError)
		default:
			return c.SendString(id)
		}
	})

	// GET /auth/signin/{username} => user_id
	app.Get("/auth/signin/:username", func(c *fiber.Ctx) error {
		id, err := dr.SignIn(c.Context(), c.Params("username"))

		switch err {
		case docrepo.ErrUnknown:
			return c.SendStatus(fiber.StatusInternalServerError)
		case docrepo.ErrNotFound:
			return c.SendStatus(fiber.StatusNotFound)
		default:
			return c.SendString(id)
		}
	})

	// GET /user/{user_id} => username
	app.Get("/user/:user_id", func(c *fiber.Ctx) error {
		name, err := dr.GetUser(c.Context(), c.Params("user_id"))

		switch err {
		case docrepo.ErrUnknown:
			return c.SendStatus(fiber.StatusInternalServerError)
		case docrepo.ErrNotFound:
			return c.SendStatus(fiber.StatusNotFound)
		default:
			return c.SendString(name)
		}
	})

	// POST /chat/{chatname} => chat_id
	app.Post("/chat/:chatname", func(c *fiber.Ctx) error {
		id, err := dr.CreateChat(c.Context(), c.Params("chatname"))

		switch err {
		case docrepo.ErrUnknown:
			return c.SendStatus(fiber.StatusInternalServerError)
		case docrepo.ErrNotFound:
			return c.SendStatus(fiber.StatusNotFound)
		default:
			return c.SendString(id)
		}
	})

	// GET /chat/{chatname}/ws => WebSocket
	app.Get("/chat/:chatname/ws", func(c *fiber.Ctx) error {
		userId := c.Get("x-user-id")
		if userId == "" {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		return c.Next()
	}, websocket.New(func(c *websocket.Conn) {
		
		// (websocket is open now)
		for {
			
		}
	}))

	// POST /chat/{chat_id}/message
	app.Post("chat/:chat_id", func(c *fiber.Ctx) error {
		userId := c.Get("x-user-id")
		body := make(map[string]interface{})
		if err := c.BodyParser(&body); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		err := dr.CreateMessage(c.Context(), c.Params("chat_id"), userId, body["content"].(string))
		
		switch err {
		case docrepo.ErrUnknown:
			return c.SendStatus(fiber.StatusInternalServerError)
		default:
			return c.SendString("success")
		}
	})

    log.Fatal(app.Listen(":3000"))
}
