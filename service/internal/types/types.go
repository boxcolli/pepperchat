package types

import "time"

type Message struct {
	MessageId string    `json:"message_idd" firestore:"message_id"`
	UserId    string    `json:"user_id" firestore:"user_id"`
	Content   string    `json:"content" firestore:"content"`
	Timestamp time.Time `json:"timestamp" firestore:"timestamp"`
}