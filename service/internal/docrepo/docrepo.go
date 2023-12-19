package docrepo

import (
	"context"
	"errors"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrUnknown = errors.New("unknown error")
	ErrConflict = errors.New("already exists")
	ErrNotFound = errors.New("not found")
)

type DocRepo interface {
	SignUp(ctx context.Context, name string) (string, error)
	SignIn(ctx context.Context, name string) (string, error)
	GetUser(ctx context.Context, id string) (string, error)
	CreateChat(ctx context.Context, name string) (string, error)
	CreateMessage(ctx context.Context, chatId, userId, content string) error
	Close()

}

type docRepo struct {
	client *firestore.Client
}

func NewDocRepo(projectId string) (DocRepo, error) {
	ctx := context.Background()
    client, err := firestore.NewClient(ctx, projectId)
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

func(r docRepo) SignUp(ctx context.Context, name string) (id string, err error) {
	// Run a transaction
	err = r.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		
		// Check if the name is already taken
		nameRef := r.client.Collection("usernames").Doc(name)
		nameDoc, err := tx.Get(nameRef) // Attempt to retrieve the document
		if err != nil {

			if status.Code(err) == codes.NotFound {
				
				// Create new user
				userRef, _, err := r.client.Collection("users").Add(ctx, map[string]interface{}{
					"name": name,
				})
				if err != nil { return ErrUnknown }

				id = userRef.ID

				// Create new username
				return tx.Set(nameRef, map[string]interface{}{
					"user_id": userRef.ID,
				})
				
			} else {
				// Return any other error
				return ErrUnknown
			}
		}

		if nameDoc.Exists() {
			// The name already exists
			return ErrConflict
		}

		// Success
		return nil
	})

	return
}

func (r docRepo) SignIn(ctx context.Context, name string) (string, error) {
	nameRef := r.client.Collection("usernames").Doc(name)
	snapshot, err := nameRef.Get(ctx)
	if err != nil {
		return "", ErrUnknown
	}

	if snapshot.Exists() {
		data := snapshot.Data()
		return data["user_id"].(string), nil
	} else {
		return "", ErrNotFound
	}
}

func (r docRepo) GetUser(ctx context.Context, id string) (string, error) {
	userRef := r.client.Collection("users").Doc(id)
	snapshot, err := userRef.Get(ctx)
	if err != nil {
		return "", ErrUnknown
	}

	if snapshot.Exists() {
		data := snapshot.Data()
		return data["username"].(string), nil
	} else {
		return "", ErrNotFound
	}
}

func (r docRepo) CreateChat(ctx context.Context, name string) (id string, err error) {
	// Run a transaction
	err = r.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		
		// Check if the name is already taken
		nameRef := r.client.Collection("chatnames").Doc(name)
		nameDoc, err := tx.Get(nameRef) // Attempt to retrieve the document
		if err != nil {

			if status.Code(err) == codes.NotFound {
				
				// Create new chat
				chatRef, _, err := r.client.Collection("chats").Add(ctx, map[string]interface{}{
					"name": name,
				})
				if err != nil { return ErrUnknown }

				id = chatRef.ID

				// Create new chatname
				return tx.Set(nameRef, map[string]interface{}{
					"chat_id": chatRef.ID,
				})
				
			} else {
				// Return any other error
				return ErrUnknown
			}
		}

		if nameDoc.Exists() {
			// The name already exists
			return ErrConflict
		}

		// Success
		return nil
	})

	return
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
