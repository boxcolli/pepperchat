package stream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	pubaddr = ":50050"
	subaddr = ":50051"
)

func TestSubscribe(t *testing.T) {
	sc, err := NewStreamClient(pubaddr, subaddr)
	assert.NoError(t, err)

	stream := sc.GetSubscribeStream()
	res, err := stream.Recv()
	assert.NoError(t, err)
	t.Logf("res.Msg: %+v\n", res.GetMsg())
}

func TestPublish(t *testing.T) {
	sc, err := NewStreamClient(pubaddr, subaddr)
	assert.NoError(t, err)

	err = sc.PublishMessage("mychat", "markov", "no cache")
	assert.NoError(t, err)
}
