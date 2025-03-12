package chat

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "yoptachat/pkg/models"
)

// ChatService предоставляет методы для работы с чатами.
type ChatService struct {
    db *sql.DB
}

// NewChatService создает новый ChatService.
func NewChatService(db *sql.DB) *ChatService {
    return &ChatService{db}
}

// SearchUsers ищет пользователей по логину.
func (c *ChatService) SearchUsers(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query().Get("query")
    users, err := c.getUsersByLogin(query)
    if err != nil {
        http.Error(w, "Ошибка при поиске пользователей", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

// SendMessage отправляет сообщение между пользователями.
func (c *ChatService) SendMessage(w http.ResponseWriter, r *http.Request) {
    var msg models.Message
    if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
        http.Error(w, "Неверный формат данных", http.StatusBadRequest)
        return
    }

    if err := c.saveMessage(msg); err != nil {
        http.Error(w, "Ошибка при отправке сообщения", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

// GetMessages получает сообщения между двумя пользователями.
func (c *ChatService) GetMessages(w http.ResponseWriter, r *http.Request) {
    senderID := r.URL.Query().Get("sender_id")
    receiverID := r.URL.Query().Get("receiver_id")

    messages, err := c.getMessages(senderID, receiverID)
    if err != nil {
        http.Error(w, "Ошибка при получении сообщений", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(messages)
}

// getUsersByLogin получает пользователей по логину.
func (c *ChatService) getUsersByLogin(login string) ([]models.User, error) {
    rows, err := c.db.Query("SELECT id, login, phone FROM Users WHERE login LIKE ?", "%"+login+"%")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []models.User
    for rows.Next() {
        var user models.User
        if err := rows.Scan(&user.ID, &user.Login, &user.Phone); err != nil {
            return nil, err
        }
        users = append(users, user)
    }
    return users, nil
}

// saveMessage сохраняет сообщение в базе данных.
func (c *ChatService) saveMessage(msg models.Message) error {
    _, err := c.db.Exec("INSERT INTO Messages (sender_id, receiver_id, content) VALUES (?, ?, ?)", msg.SenderID, msg.ReceiverID, msg.Content)
    return err
}

// getMessages получает сообщения между двумя пользователями.
func (c *ChatService) getMessages(senderID, receiverID string) ([]models.Message, error) {
    rows, err := c.db.Query("SELECT id, sender_id, receiver_id, content FROM Messages WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", senderID, receiverID, receiverID, senderID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var messages []models.Message
    for rows.Next() {
        var message models.Message
        if err := rows.Scan(&message.ID, &message.SenderID, &message.ReceiverID, &message.Content); err != nil {
            return nil, err
        }
        messages = append(messages, message)
    }
    return messages, nil
}
