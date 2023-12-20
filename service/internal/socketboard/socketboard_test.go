package socketboard

import (
	"testing"

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
	ch := make(chan OutType, 10)

	// Add to board
	sb := NewSocketBoard(sc.GetSubscribeStream())
	sb.Add(chatId, ch)

	// Publish
	err = sc.PublishMessage(chatId, userId, content)
	assert.NoError(t, err)

	b, ok := <- ch
	if !ok {
		t.Log("channel close")
		return
	}

	t.Logf("channel read: %+v\n", b)
	t.Logf("channel read: %s\n", string(b))
}
