package socketboard

import (
	"encoding/json"
	"log"
	"sync"

	pb "github.com/boxcolli/go-transistor/api/gen/transistor/v1"
	"github.com/boxcolli/pepperchat/service/internal/types"
)

type SocketBoard interface {
	Add(chatId string, ch chan<- []byte)
	Del(chatId string, ch chan<- []byte)
}

type socketBoard struct {
	board map[string](map[chan<- []byte]bool)
	mx    sync.RWMutex
}

func NewSocketBoard(stream pb.TransistorService_SubscribeClient) SocketBoard {
	b := &socketBoard{
		board: make(map[string]map[chan<- []byte]bool),
	}
 
	go func() {
		for {
			// Read message
			res, err := stream.Recv()	// Block
			if err != nil {
				log.Fatalf("SocketBoard: stream end: %v\n", err)
				break
			}

			// Convert transistor message
			trmsg := res.GetMsg()
			tokens := trmsg.GetTopic().GetTokens()
			if len(tokens) < 2 || tokens[0] != "chat" {
				log.Printf("SocketBoard: wrong message topic: %+v\n", trmsg)
				continue
			}
		
			// Extract chat message
			var msg types.Message
			msgByte := trmsg.GetData().GetValue()
			err = json.Unmarshal(msgByte, &msg)
			if err != nil {
				log.Printf("SocketBoard: unmarshal failed: %v\n", err)
				continue
			}
			
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
	} ()

	return b
}

// Add implements SocketBoard.
func (b *socketBoard) Add(chatId string, ch chan<- []byte) {
	b.mx.Lock()
	defer b.mx.Unlock()

	if _, ok := b.board[chatId]; !ok {
		b.board[chatId] = make(map[chan<- []byte]bool)
	}
	b.board[chatId][ch] = true
}

// Del implements SocketBoard.
func (b *socketBoard) Del(chatId string, ch chan<- []byte) {
	b.mx.Unlock()
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
