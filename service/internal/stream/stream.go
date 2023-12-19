package stream

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	pb "github.com/boxcolli/go-transistor/api/gen/transistor/v1"
	"github.com/boxcolli/go-transistor/types"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

var (
	ErrUnknown = errors.New("unknown error")
)

type StreamClient interface {
	Close()
	PublishMessage(chatId, userId, content string) error
	GetSubscribeStream() pb.TransistorService_SubscribeClient
}

type streamClient struct {
	pclient	pb.TransistorServiceClient
	pconn	*grpc.ClientConn
	pctx	context.Context
	pstream	pb.TransistorService_PublishClient

	sclient	pb.TransistorServiceClient
	sconn	*grpc.ClientConn
	sctx	context.Context
	sstream	pb.TransistorService_SubscribeClient
}

func NewStreamClient(pubaddr, subaddr string) (StreamClient, error) {
	c := streamClient{}
	var err error
	// Open clients
	for {
		c.pconn, err = grpc.Dial(pubaddr, dialOpts...)
		if err != nil { continue }
		break
	}
	for {
		c.sconn, err = grpc.Dial(subaddr, dialOpts...)
		if err != nil { continue }
		break
	}
	c.pclient = pb.NewTransistorServiceClient(c.pconn)
	c.sclient = pb.NewTransistorServiceClient(c.sconn)

	// Call RPCs
	c.pctx = context.Background()
	c.sctx = context.Background()
	for {
		c.pstream, err = c.pclient.Publish(c.pctx)
		if err != nil { continue }
		break
	}
	for {
		c.sstream, err = c.sclient.Subscribe(c.sctx)
		if err != nil { continue }
		break
	}

	// Initiate subscription
	for {
		err = c.sstream.Send(&pb.SubscribeRequest{
			Change: &pb.Change{
				Op: pb.Operation_OPERATION_ADD,
				Topic: &pb.Topic{
					Tokens: []string{"chat"},
				},
			},
		})
		if err != nil { continue }
		break
	}

	return &c, nil
}

func (c *streamClient) Close() {
	c.pconn.Close()
	c.sconn.Close()
}

func (c *streamClient) PublishMessage(chatId, userId, content string) error {
	// Prepare data
	timestamp := time.Now().UTC()
	messageId, err := uuid.NewV7()
	if err != nil { return ErrUnknown }
	data := map[string]interface{}{
        "message_id": messageId.String(),
		"username": userId,
		"content": content,
		"created_time": timestamp.String(),	
    }

    // Marshal the JSON with data
    jsonData, err := json.Marshal(data)
    if err != nil {
        return ErrUnknown
    }
	
	msg := types.Message{
		Topic: types.Topic{"chat", chatId},
		Data: jsonData,
		TP: timestamp,
	}
	mar := msg.Marshal()

	// Send message
	err = c.pstream.Send(&pb.PublishRequest{
		Msg: mar,
	})
	if err != nil {
		// Unsuccessful
		return ErrUnknown
	}

	// Success
	return nil
}

func (c *streamClient) GetSubscribeStream() pb.TransistorService_SubscribeClient {
	return c.sstream
}