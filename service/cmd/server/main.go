package main

import (
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
