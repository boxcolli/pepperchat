package test

import (
	// "encoding/base64"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

const (
	addr = ":8080"
)

/*

Result:
=== RUN   TestWebsocket
    websocket_test.go:23: connecting to ws://:8080/chat/mychat/ws
    websocket_test.go:33: waiting for message...
    websocket_test.go:76: Response Status Code: 200
    websocket_test.go:36: recv: {"content":"hello from markov","created_time":"2023-12-20 02:02:23.066977921 +0000 UTC","message_id":"018c84f5-afda-7e8a-a996-6f6a3f0cb958","username":"markov"}
--- PASS: TestWebsocket (0.01s)
PASS
ok      github.com/boxcolli/pepperchat/service/test     0.397s

*/

func TestWebsocket(t *testing.T) {
	// URL
	u := url.URL{Scheme: "ws", Host: addr, Path: "/chat/mychat/ws"}
	t.Logf("connecting to %s", u.String())

	// Dial
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	assert.NoError(t, err)
	defer c.Close()

	go postMessage(t)

	// Read 1 message
	t.Log("waiting for message...")
	_, msg, err := c.ReadMessage()
	assert.NoError(t, err)
	t.Logf("recv: %s", string(msg))

	// decoded := make([]byte, 0)
	// _, err = base64.StdEncoding.Decode(decoded, msg)
	// assert.NoError(t, err)

	// t.Logf("decoded: %s\n", string(decoded))
}

func postMessage(t *testing.T) {
    // URL of the server
    url := "http://localhost:8080/chat/mychat/message"

    // JSON payload
    payload := map[string]string{"content": "hello from markov"}
    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        log.Fatal(err)
    }
    body := bytes.NewReader(payloadBytes)

    // Create a new request
    req, err := http.NewRequest("POST", url, body)
    if err != nil {
        log.Fatal(err)
    }

    // Set headers
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("x-user-id", "markov")

    // Send the request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    // Print the response status code
    t.Log("Response Status Code:", resp.StatusCode)
}