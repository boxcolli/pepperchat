package docrepo

import (
	"context"
	"errors"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/boxcolli/pepperchat/service/internal/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrUnknown = errors.New("unknown error")
	ErrConflict = errors.New("already exists")
	ErrNotFound = errors.New("not found")
)

type DocRepo interface {
	CreateChat(ctx context.Context, chatId string) error
	CreateMessage(ctx context.Context, chatId, userId, content string) error
	GetMessage(ctx context.Context, chatId string, offset time.Time, limit int) ([]types.Message, error)
	Close()
}

type docRepo struct {
	client *firestore.Client
}

func NewDocRepo(projectId, databaseId string) (DocRepo, error) {
	ctx := context.Background()
    client, err := firestore.NewClientWithDatabase(ctx, projectId, databaseId)
    if err != nil {
        return nil, err
    }
	return &docRepo{
		client: client,
	}, nil
}

func(r docRepo) Close() {
	r.client.Close()
}

// func(r docRepo) SignUp(ctx context.Context, name string) (id string, err error) {
// 	// Run a transaction
// 	err = r.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		
// 		// Check if the name is already taken
// 		nameRef := r.client.Collection("usernames").Doc(name)
// 		nameDoc, err := tx.Get(nameRef) // Attempt to retrieve the document
// 		if err != nil {

// 			if status.Code(err) == codes.NotFound {
				
// 				// Create new user
// 				userRef, _, err := r.client.Collection("users").Add(ctx, map[string]interface{}{
// 					"name": name,
// 				})
// 				if err != nil { return ErrUnknown }

// 				id = userRef.ID

// 				// Create new username
// 				return tx.Set(nameRef, map[string]interface{}{
// 					"user_id": userRef.ID,
// 				})
				
// 			} else {
// 				// Return any other error
// 				return ErrUnknown
// 			}
// 		}

// 		if nameDoc.Exists() {
// 			// The name already exists
// 			return ErrConflict
// 		}

// 		// Success
// 		return nil
// 	})

// 	return
// }


func (r docRepo) CreateChat(ctx context.Context, chatId string) error {
	// Run a transaction
	err := r.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		
		// Check if the name is already taken
		chatRef := r.client.Collection("chatnames").Doc(chatId)
		chatDoc, err := tx.Get(chatRef) // Attempt to retrieve the document
		if err != nil {

			if status.Code(err) == codes.NotFound {
				
				// Create new chat
				_, err := r.client.Collection("chats").Doc(chatId).Set(ctx, map[string]interface{}{
					"created_at": time.Now(),
				})
				if err != nil { return ErrUnknown }

				return nil
				
			} else {
				// Return any other error
				return ErrUnknown
			}
		}

		if chatDoc.Exists() {
			// The name already exists
			return ErrConflict
		}

		// Create new chat
		_, err = r.client.Collection("chats").Doc(chatId).Set(ctx, nil)
		if err != nil { return ErrUnknown }
		
		// Success
		return nil
	})

	return err
}

func (r docRepo) CreateMessage(ctx context.Context, chatId, userId, content string) error {
	_, _, err := r.client.Collection("chats").Doc(chatId).Collection("messages").Add(ctx, map[string]interface{}{
		"user_id": userId,
		"content": content,
	})
	if err != nil {
		return ErrUnknown
	}

	// Success
	return nil
}

func (r docRepo) GetMessage(ctx context.Context, chatId string, offset time.Time, limit int) ([]types.Message, error) {
	// Make query
	subcollection := r.client.Collection("chats").Doc(chatId).Collection("messages")
    query := subcollection.OrderBy("created_at", firestore.Desc).Where("created_at", "<", offset).Limit(limit)

	// Get iteration
    iter := query.Documents(ctx)
    defer iter.Stop()

	// Gather messages
	messages := []types.Message{}
    for {
		// Read document
        doc, err := iter.Next()
        if err != nil { break }

		// Convert data type
        var msg types.Message
        err = doc.DataTo(&msg)
		if err != nil { return nil, ErrUnknown }

		// Append
		messages = append(messages, msg)
    }

	return messages, nil
}
