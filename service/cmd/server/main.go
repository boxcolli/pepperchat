package main

import (
	// "encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	// "time"

	// "github.com/boxcolli/pepperchat/service/internal/docrepo"
	"github.com/boxcolli/pepperchat/service/internal/socketboard"
	"github.com/boxcolli/pepperchat/service/internal/stream"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/peterbourgon/ff/v4"
)

const (
	// projectId = "-"
	// databaseId = "pepperdb-1"
	// pubaddr = ":50050"
	// subaddr = ":50051"
	timeLayout = "2006-01-02 15:04:05"
	chSize = 100
)

func main() {
	fs := flag.NewFlagSet("service", flag.ContinueOnError)
	var (
		paddr = fs.String("paddr", ":50050", "addrerss of the publisher server")
		saddr = fs.String("saddr", ":50050", "addrerss of the subscriber server")
		port = fs.String("port", "8080", "port number to listen to")
	)
	ff.Parse(fs, os.Args[1:], ff.WithEnvVars())

	// DB client
	// var dr docrepo.DocRepo
	// {
	// 	var err error
	// 	dr, err = docrepo.NewDocRepo(projectId, databaseId)
	// 	if err != nil {
	// 		log.Fatalf("Failed to create client with database: %v", err)
	// 	}
	// }

	// Stream client
	var st stream.StreamClient
	{
		var err error
		st, err = stream.NewStreamClient(*paddr, *saddr)
		if err != nil {
			log.Fatalf("Failed to create client with stream: %v", err)
		}
		log.Println("StreamClient success")
	}

	// SocketBoard
	var sb = socketboard.NewSocketBoard(st.GetSubscribeStream())

    app := fiber.New()

	// POST /chat/{chat_id}
	// app.Post("/chat/:chat_id", func(c *fiber.Ctx) error {
	// 	err := dr.CreateChat(c.Context(), c.Params("chat_id"))

	// 	switch err {
	// 	case docrepo.ErrUnknown:
	// 		return c.SendStatus(fiber.StatusInternalServerError)
	// 	case docrepo.ErrNotFound:
	// 		return c.SendStatus(fiber.StatusNotFound)
	// 	default:
	// 		return c.SendString("success")
	// 	}
	// })

	// GET /chat/{chat_id}/ws => WebSocket
	app.Get("/chat/:chat_id/ws", func(c *fiber.Ctx) error {
		log.Println("Received: ", c.OriginalURL())

		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired

	}, websocket.New(func(c *websocket.Conn) {
		// (websocket is open now)

		// Get parameter
		chatId := c.Params("chat_id")
		if chatId == "" {
			if err := c.WriteJSON("Bad Request"); err != nil {
				log.Printf("GetWebSocket: failed to marshal message: %v\n", err)
				return
			}
			return
		}

		// Make channels
		mch := make(chan socketboard.OutType, chSize)
		sb.Add(chatId, mch)
		stop := make(chan bool)

		// Defer clean up functions
		defer sb.Del(chatId, mch)
		defer close(mch)

		// Listen for client side closure
		c.SetCloseHandler(func(code int, text string) error {
			close(stop)
			return nil
		})

		// Transfer message
		for {
			select {
			case <- stop:
				log.Printf("GetWebSocket(%s): connection closed\n", chatId)
				return

			case m := <- mch:
				log.Printf("GetWebSocket(%s): transfer message: %+v\n", chatId, string(m))
				err := c.WriteMessage(websocket.BinaryMessage, m)
				// err := c.WriteJSON(m)
				if err != nil {
					// Something went wrong
					log.Printf("GetWebSocket(%s): failed to write json: %v\n", chatId, err)
					return 
				}

				// Successfully sent message
			}
		}
	}))

	// POST /chat/{chat_id}/message
	app.Post("/chat/:chat_id/message", func(c *fiber.Ctx) error {
		log.Println("Received: ", c.OriginalURL())
		
		chatId := c.Params("chat_id")
		userId := c.Get("x-user-id")
		body := make(map[string]interface{})
		if err := c.BodyParser(&body); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		// Publish message
		err := st.PublishMessage(chatId, userId, body["content"].(string))
		
		switch err {
		case stream.ErrUnknown:
			return c.SendStatus(fiber.StatusInternalServerError)
		default:
			return c.SendString("success")
		}
	})

	// // GET /chat/{chat_id}/message/{offset}/{limit}
	// app.Get("/chat/:chat_id/message/:offset/:limit", func(c *fiber.Ctx) error {
	// 	// Extract offset
	// 	offset, err := time.Parse(timeLayout, c.Params("offset"))
	// 	if err != nil { return c.SendStatus(fiber.StatusBadRequest) }

	// 	// Extract limit
	// 	limit, err := c.ParamsInt("limit")
	// 	if err != nil { return c.SendStatus(fiber.StatusBadRequest)}

	// 	// Query messages
	// 	messages, err := dr.GetMessage(c.Context(), c.Params("chat_id"), offset, limit)
	// 	if err != nil { return c.SendStatus(fiber.StatusInternalServerError)}

	// 	// Convert into JSON
	// 	jsonByte, err := json.Marshal(messages)
	// 	if err != nil { return c.SendStatus(fiber.StatusInternalServerError)}

	// 	return c.Send(jsonByte)
	// })
	log.Println("Now listening on", *port)
    log.Fatal(app.Listen(fmt.Sprintf(":%s", *port)))
}
