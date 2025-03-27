package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
)

var ( //+
	store = sessions.NewCookieStore([]byte("your_secret_key"))
	db    *sql.DB
)

func main() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/chat")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	r := gin.Default()
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")
	r.GET("/", indexHandler)
	r.GET("/regauth", regAuthHandler)
	r.GET("/logout", logout)
	r.POST("/regauth", handleRegAuth)
	r.GET("/search", searchHandler)
	r.POST("/add_friend", addFriendHandler)
	r.GET("/chat", chatHandler)
	r.Run(":8080")
}

func indexHandler(c *gin.Context) { //Срабатывает при открытии главной страницы
	session, _ := store.Get(c.Request, "session") //В переменную записывает значение сессии
	if session.Values["user_id"] == nil {         //Проверка если нету сессии user_id то кидает на страницу регистрации
		c.Redirect(http.StatusFound, "/regauth")
		return
	}

	ID := session.Values["user_id"].(int) //Вписывает в переменную ID значение сесии
	var friends []string                  //создает массив для сохранения найденных друзей

	rows, err := db.Query("SELECT u.login FROM Friends f JOIN Users u ON f.friend_id = u.id WHERE f.user_id = ?", ID) //Выводит всех друзей которые есть у пользователя user_id = ID
	if err != nil {                                                                                                   //Ошбка
		c.String(http.StatusInternalServerError, "Ошибка получения друзей")
		return
	}
	defer rows.Close() //Закрывает rows после выполнения

	for rows.Next() { //Проходится циклом по rows
		var login string                          //создает переменную куда будут вписываться полученные данные
		if err := rows.Scan(&login); err != nil { //вписываем login в переменную login
			c.String(http.StatusInternalServerError, "Ошибка получения друзей") //Ошибка
			return
		}
		friends = append(friends, fmt.Sprint(login)) //Записываение в массив полученных данных
	}

	c.HTML(http.StatusOK, "index.html", gin.H{ //Передача данных сессии и массива друзей на главную страницу
		"session": session.Values, //Нынешняя сессии пользователя
		"friends": friends,        //Массив его друзей
	})
}

func logout(c *gin.Context) { //Выход из сессии
	session, err := store.Get(c.Request, "session") //В переменную записывает значение сессии
	if err != nil {                                 //Ошибка
		c.String(http.StatusInternalServerError, "Ошибка выхода из системы")
		return
	}
	delete(session.Values, "user_id")       //Удаление сессии user_id
	err = session.Save(c.Request, c.Writer) //сохранение изменений
	if err != nil {                         //Ошибка
		c.String(http.StatusInternalServerError, "Ошибка выхода из системы")
		return
	}
	c.Redirect(http.StatusFound, "/") //Редирект обратно на главную
}

func regAuthHandler(c *gin.Context) { //Заход на страницу регистрации
	session, _ := store.Get(c.Request, "session")//В переменную записывает значение сессии
	if session.Values["user_id"] != nil { //Если сессия user_id есть то редирект на главную
		c.Redirect(http.StatusFound, "/")
		return
	}
	c.HTML(http.StatusOK, "regauth.html", nil)//Переход на страницу регистрации
}

func handleRegAuth(c *gin.Context) { //Регистрация и Авторизация
	login := c.PostForm("login")//Получаем логин
	password := c.PostForm("password")//Получаем пароль
	action := c.PostForm("action")//Получаем метод отправки
	if action == "register" {//Если метод регистрации то ...
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)//Шифруем пароль
		if err != nil {//Ошибка
			c.String(http.StatusInternalServerError, "Ошибка регистрации: не удалось хешировать пароль")
			return
		}
		_, err = db.Exec("INSERT INTO Users (login, password) VALUES (?, ?)", login, hashedPassword)//Добавляем пользователя в БД
		if err != nil {//Ошибка
			fmt.Println("Ошибка при вставке в БД:", err)
			c.String(http.StatusInternalServerError, "Ошибка регистрации: не удалось сохранить пользователя")
			return
		}
		handleLogin(c, login, password)//Перенаправляем запрос на авторизацию и передаем ему логин и пароль
	} else if action == "login" {//Если метод авторизации то ...
		handleLogin(c, login, password)//Перенаправляем запрос на авторизацию и передаем ему логин и пароль
	} else {
		c.String(http.StatusBadRequest, "Неверное действие")//Ошибка
	}
}

func handleLogin(c *gin.Context, login, password string) { //Авторизация и создание сесиий
	var storedPassword string//Переменная получения пароля
	var ID int//Переменная получения id
	err := db.QueryRow("SELECT id, password FROM Users WHERE login = ?", login).Scan(&ID, &storedPassword)//Записываем в ID и storedPassword данные из запроса
	if err != nil || bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password)) != nil {//Ошибка
		c.String(http.StatusUnauthorized, "Неверные учетные данные")
		return
	}
	session, _ := store.Get(c.Request, "session")//В переменную записывает значение сессии

	session.Values["login"] = login//Значение сессии пол логину гл пользователя
	session.Values["user_id"] = ID//Значение сессии пол id гл пользователя
	session.Save(c.Request, c.Writer)//Сохранение в сессии данных
	c.Redirect(http.StatusFound, "/")//Переход на гл страницу
}

// Обработчик поиска пользователей
func searchHandler(c *gin.Context) {//Поиск пользователей по тексту из инпута
	letter := c.Query("letter")//Получем значение записанное в инпут
	var results []struct {//Создаем структуру для записи полученных пользователей со значнеиями ID и Login
		ID    int
		Login string
	}

	rows, err := db.Query("SELECT id, login FROM Users WHERE login LIKE concat('%', ?, '%')", letter)//запрос на нахождение подходящих значений из бд по данным из инпута
	if err != nil {//Ошибка
		c.String(http.StatusInternalServerError, "Ошибка поиска")
		return
	}
	defer rows.Close()//Закрывает rows после выполнения

	for rows.Next() {//Проходится циклом по rows
		var user struct {//создаем значение структуры в которую будут вписываться данные
			ID    int
			Login string
		}
		if err := rows.Scan(&user.ID, &user.Login); err != nil {//Вписываемданные полученные из запроса
			c.String(http.StatusInternalServerError, "Ошибка получения результатов")//Ошибка
			return
		}
		results = append(results, user)//Вписываем данные в массив
	}

	c.JSON(http.StatusOK, results)//передаем данные в виде JSON значения 
}

func addFriendHandler(c *gin.Context) {
	var request struct {
		FriendID int `json:"friend_id"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, "Неверный запрос")
		return
	}

	session, _ := store.Get(c.Request, "session")
	ID := session.Values["user_id"].(int)

	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM Friends WHERE user_id = ? AND friend_id = ?)", ID, request.FriendID).Scan(&exists)
	if err != nil {
		c.String(http.StatusInternalServerError, "Ошибка при проверке друзей")
		return
	}

	if exists {
		c.String(http.StatusConflict, "Пользователь уже в друзьях")
		return
	}

	_, err = db.Exec("INSERT INTO Friends (user_id, friend_id) VALUES (?, ?), (?, ?)", ID, request.FriendID, request.FriendID, ID)
	if err != nil {
		c.String(http.StatusInternalServerError, "Ошибка при добавлении друга")
		return
	}
	c.String(http.StatusOK, "Друг добавлен")
}

var upgrader = websocket.Upgrader{}

func chatHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        log.Println(err)
        return
    }
    defer conn.Close()

	session, _ := store.Get(c.Request, "session")
    // Получите ID текущего пользователя из сессии или контекста
    ID := session.Values["user_id"].(int)
    friendLogin := c.Query("friend") // Получаем логин друга из запроса
	row, err := db.Query("SELECT id FROM Users WHERE login = ?", friendLogin)
	if err != nil {
        log.Println("Ошибка загрузки сообщений:", err)
        return
    }
    defer row.Close()
	var friendID int
	for row.Next(){
		var friendIDs int
		if err := row.Scan(&friendIDs); err != nil {
            log.Println("Ошибка при сканировании сообщения:", err)
            continue
        }
		friendID = friendIDs
	}
    // Загрузка предыдущих сообщений
    var messages []struct {
        SenderID  int    `json:"sender_id"`
        Message   string `json:"message"`
        Timestamp string `json:"timestamp"`
    }

    rows, err := db.Query(`
        SELECT sender_id, message, timestamp 
        FROM Messages 
        WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)
        ORDER BY timestamp`, ID, friendID, friendID, ID)
    if err != nil {
        log.Println("Ошибка загрузки сообщений:", err)
        return
    }
    defer rows.Close()

    for rows.Next() {
        var msg struct {
            SenderID  int    `json:"sender_id"`
            Message   string `json:"message"`
            Timestamp string `json:"timestamp"`
        }
        if err := rows.Scan(&msg.SenderID, &msg.Message, &msg.Timestamp); err != nil {
            log.Println("Ошибка при сканировании сообщения:", err)
            continue
        }
        messages = append(messages, msg)
    }

    // Отправка предыдущих сообщений
    conn.WriteJSON(messages)

    for {
        var msg struct {
            To      string `json:"to"`
            Message string `json:"message"`
        }

        err := conn.ReadJSON(&msg)
        if err != nil {
            log.Println(err)
            break
        }

        // Сохранение нового сообщения в базе данных
        _, err = db.Exec("INSERT INTO Messages (sender_id, receiver_id, message) VALUES (?, ?, ?)", ID, friendID, msg.Message)
        if err != nil {
            log.Println("Ошибка при сохранении сообщения:", err)
            continue
        }

        // Отправка нового сообщения обратно
        conn.WriteJSON(msg)
    }
}

//задача 1 сделать референс друзей.+
//создать websocket
//накатить стилей
//друзей через точку поставить
