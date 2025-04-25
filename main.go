package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
)

var ( //+
	store    = sessions.NewCookieStore([]byte("your_secret_key"))
	db       *sql.DB
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	clients = make(map[*websocket.Conn]string)
	mu      sync.Mutex
)

func main() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/chat") //Можно прям на хост бд залить порт тот же
	if err != nil {
		panic(err)
	}
	defer db.Close()
	r := gin.Default()
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")
	r.GET("/", indexHandler)                // main.go //
	r.GET("/regauth", regAuthHandler)       // regauth.go //
	r.GET("/logout", logout)                // logoout.go //
	r.POST("/regauth", handleRegAuth)       // regauth.go //
	r.GET("/search", searchHandler)         // serch.go //
	r.POST("/add_friend", addFriendHandler) // addfriend.go //
	r.GET("/chat", chatHandler)             // while main.go //
	r.GET("/ws", handleWebSocket)
	r.Run(":8080")
}

func indexHandler(c *gin.Context) {
	//Срабатывает при открытии главной страницы
	session, _ := store.Get(c.Request, "session") //В переменную записывает значение сессии
	if session.Values["user_id"] == nil {         //Проверка если нету сессии user_id то кидает на страницу регистрации
		c.Redirect(http.StatusFound, "/regauth")
		return
	}

	ID := session.Values["user_id"].(int) //Вписывает в переменную ID значение сесии
	type Friend struct {
		ID    int
		Login string
	}
	var friends []Friend //создает массив для сохранения найденных друзей

	rows, err := db.Query("SELECT u.id, u.login FROM Friends f JOIN Users u ON f.friend_id = u.id WHERE f.user_id = ?", ID) //Выводит всех друзей которые есть у пользователя user_id = ID
	if err != nil {                                                                                                         //Ошбка
		c.String(http.StatusInternalServerError, "Ошибка получения друзей")
		return
	}
	defer rows.Close() //Закрывает rows после выполнения

	for rows.Next() { //Проходится циклом по rows
		var friend Friend                                            //создает переменную куда будут вписываться полученные данные
		if err := rows.Scan(&friend.ID, &friend.Login); err != nil { //вписываем login в переменную login
			c.String(http.StatusInternalServerError, "Ошибка получения друзей") //Ошибка
			return
		}
		friends = append(friends, friend) //Записываение в массив полученных данных
	}

	c.HTML(http.StatusOK, "index.html", gin.H{ //Передача данных сессии и массива друзей на главную страницу
		"session": session.Values, //Нынешняя сессии пользователя
		"friends": friends,        //Массив его друзей
	})
}

func chatHandler(c *gin.Context) {
	type Message struct {
		SenderID   int    `json:"sender_id"`
		ReceiverID int    `json:"receiver_id"`
		Message    string `json:"content"`
		Timestamp  string `json:"timestamp"`
	}

	type Receiver struct {
		ID    int
		Login string
	}

	session, _ := store.Get(c.Request, "session")
	if session.Values["user_id"] == nil {
		c.Redirect(http.StatusFound, "/regauth")
		return
	}

	ID := session.Values["user_id"].(int)
	FriendID := c.Query("friendID")

	free, err := db.Query("SELECT id, login FROM Users WHERE id = ?", FriendID)
	if err != nil {
		c.String(http.StatusInternalServerError, "Ошибка при получении логина друга")
	}
	defer free.Close()

	var receivers []Receiver
	for free.Next() {
		var receiver Receiver
		if err := free.Scan(&receiver.ID, &receiver.Login); err != nil {
			c.String(http.StatusInternalServerError, "Ошибка при получении друга")
			return
		}
		receivers = append(receivers, receiver)
	}

	rows, err := db.Query("SELECT sender_id, receiver_id, message, timestamp FROM Messages WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?) ORDER BY timestamp ASC", ID, FriendID, FriendID, ID)
	if err != nil {
		c.String(http.StatusInternalServerError, "Ошибка при получении сообщений")
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.SenderID, &msg.ReceiverID, &msg.Message, &msg.Timestamp); err != nil {
			c.String(http.StatusInternalServerError, "Ошибка при считывании данных")
			return
		}
		messages = append(messages, msg)
	}

	c.HTML(http.StatusOK, "chat.html", gin.H{
		"session":   session.Values,
		"receivers": receivers,
		"messages":  messages,
	})
}

func handleWebSocket(c *gin.Context) {
	type Message struct {
		SenderID   string `json:"SenderID"`
		ReceiverID string `json:"ReceiverID"`
		Content    string `json:"Content"`
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Error while upgrading connection:", err)
		return
	}
	defer conn.Close()
	userID := c.Query("sender_id")
	clients[conn] = userID

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			delete(clients, conn)
			break
		}

		mu.Lock()
		_, err = db.Exec("INSERT INTO Messages (sender_id, receiver_id, message) VALUES (?, ?, ?)", msg.SenderID, msg.ReceiverID, msg.Content)
		mu.Unlock()

		if err != nil {
			fmt.Println("Error saving message to DB:", err)
			continue
		}

		for client, id := range clients {
			if id == msg.ReceiverID {
				err := client.WriteJSON(msg)
				if err != nil {
					client.Close()
					delete(clients, client)
				}
			}
			if id == msg.SenderID {
				err := client.WriteJSON(msg)
				if err != nil {
					client.Close()
					delete(clients, client)
				}
			}
		}
	}
}

//накатить стилей
