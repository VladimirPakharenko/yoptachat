package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func regAuthHandler(c *gin.Context) { //Заход на страницу регистрации
	session, _ := store.Get(c.Request, "session") //В переменную записывает значение сессии
	if session.Values["user_id"] != nil {         //Если сессия user_id есть то редирект на главную
		c.Redirect(http.StatusFound, "/")
		return
	}
	c.HTML(http.StatusOK, "regauth.html", nil) //Переход на страницу регистрации
}

func handleRegAuth(c *gin.Context) { //Регистрация и Авторизация
	login := c.PostForm("login")       //Получаем логин
	password := c.PostForm("password") //Получаем пароль
	action := c.PostForm("action")     //Получаем метод отправки
	if action == "register" {          //Если метод регистрации то ...
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) //Шифруем пароль
		if err != nil {                                                                          //Ошибка
			c.String(http.StatusInternalServerError, "Ошибка регистрации: не удалось хешировать пароль")
			return
		}
		_, err = db.Exec("INSERT INTO Users (login, password) VALUES (?, ?)", login, hashedPassword) //Добавляем пользователя в БД
		if err != nil {                                                                              //Ошибка
			fmt.Println("Ошибка при вставке в БД:", err)
			c.String(http.StatusInternalServerError, "Ошибка регистрации: не удалось сохранить пользователя")
			return
		}
		handleLogin(c, login, password) //Перенаправляем запрос на авторизацию и передаем ему логин и пароль
	} else if action == "login" { //Если метод авторизации то ...
		handleLogin(c, login, password) //Перенаправляем запрос на авторизацию и передаем ему логин и пароль
	} else {
		c.String(http.StatusBadRequest, "Неверное действие") //Ошибка
	}
}

func handleLogin(c *gin.Context, login, password string) { //Авторизация и создание сесиий
	var storedPassword string                                                                              //Переменная получения пароля
	var ID int                                                                                             //Переменная получения id
	err := db.QueryRow("SELECT id, password FROM Users WHERE login = ?", login).Scan(&ID, &storedPassword) //Записываем в ID и storedPassword данные из запроса
	if err != nil || bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password)) != nil {      //Ошибка
		c.String(http.StatusUnauthorized, "Неверные учетные данные")
		return
	}
	session, _ := store.Get(c.Request, "session") //В переменную записывает значение сессии

	session.Values["login"] = login   //Значение сессии пол логину гл пользователя
	session.Values["user_id"] = ID    //Значение сессии пол id гл пользователя
	session.Save(c.Request, c.Writer) //Сохранение в сессии данных
	c.Redirect(http.StatusFound, "/") //Переход на гл страницу
}
