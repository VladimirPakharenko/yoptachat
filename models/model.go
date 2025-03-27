package models

type User struct {
    ID       int
    Login    string
    Phone    string
    Password string
}

// Message представляет сообщение между пользователями.
type Message struct {
    ID         int    `json:"id"`
    SenderID   int    `json:"sender_id"`
    ReceiverID int    `json:"receiver_id"`
    Content    string `json:"content"`
}