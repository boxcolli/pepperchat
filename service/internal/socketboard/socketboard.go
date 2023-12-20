package socketboard

import (
	"encoding/json"
	"log"
	"sync"

	pb "github.com/boxcolli/go-transistor/api/gen/transistor/v1"
	"github.com/boxcolli/pepperchat/service/internal/types"
)

type OutType []byte

type SocketBoard interface {
	Add(chatId string, ch chan<- OutType)
	Del(chatId string, ch chan<- OutType)
}

type socketBoard struct {
	stream 	pb.TransistorService_SubscribeClient
	board 	map[string](map[chan<- OutType]bool)
	mx    	sync.RWMutex
}

func NewSocketBoard(stream pb.TransistorService_SubscribeClient) SocketBoard {
	b := &socketBoard{
		stream: stream,
		board: make(map[string]map[chan<- OutType]bool),
	}
 
	go b.run()

	return b
}

func (b *socketBoard) run() {
	for {
		// Read message
		res, err := b.stream.Recv()	// Block
		if err != nil {
			log.Fatalf("SocketBoard: stream end: %v\n", err)
			break
		}

		// Convert response -> transistor message
		trmsg := res.GetMsg()
		tokens := trmsg.GetTopic().GetTokens()
		if len(tokens) < 2 || tokens[0] != "chat" {
			log.Printf("SocketBoard: wrong message topic: %+v\n", trmsg)
			continue
		}
	
		// Convert transistor message -> chat message
		var msg types.Message
		msgByte := trmsg.GetData().GetValue()
		err = json.Unmarshal(msgByte, &msg)
		if err != nil {
			log.Printf("SocketBoard: unmarshal failed: %v\n", err)
			continue
		}
		log.Printf("SocketBoard: received: %s\n", string(msgByte))
		
		// Push message
		b.mx.RLock()
		chset, ok := b.board[tokens[1]]
		if !ok {
			b.mx.RUnlock()
			continue
		}
		for sub := range chset {
			sub <- msgByte
		}
		b.mx.RUnlock()
	}
}

// Add implements SocketBoard.
func (b *socketBoard) Add(chatId string, ch chan<- OutType) {
	b.mx.Lock()
	defer b.mx.Unlock()

	if _, ok := b.board[chatId]; !ok {
		b.board[chatId] = make(map[chan<- OutType]bool)
	}
	b.board[chatId][ch] = true
}

// Del implements SocketBoard.
func (b *socketBoard) Del(chatId string, ch chan<- OutType) {
	b.mx.Lock()
	defer b.mx.Unlock()

	// Delete ch entry
	chset, ok := b.board[chatId]
	if !ok { return }
	delete(b.board[chatId], ch)

	// Delete chatId entry
	if len(chset) == 0 {
		delete(b.board, chatId)
	}
}
