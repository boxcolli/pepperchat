package socketboard

import (
	"testing"
	"time"

	"github.com/boxcolli/pepperchat/service/internal/stream"
	"github.com/stretchr/testify/assert"
)

const (
	chatId = "mytest"
	userId = "markov"
	content = "hello"
	pubaddr = ":50050"
	subaddr = ":50051"
)

func TestBoard(t *testing.T) {
	sc, err := stream.NewStreamClient(pubaddr, subaddr)
	assert.NoError(t, err)

	// Read channel
	ch := make(chan []byte, 10)
	go func() {
		for {
			b, ok := <- ch
			if !ok {
				t.Log("channel close")
				return
			}

			t.Logf("channel read: %s\n", string(b))
		}
	} ()

	// Add to board
	sb := NewSocketBoard(sc.GetSubscribeStream())
	sb.Add(chatId, ch)

	// Publish
	err = sc.PublishMessage(chatId, userId, content)
	assert.NoError(t, err)

	time.Sleep(1 * time.Second)
}
