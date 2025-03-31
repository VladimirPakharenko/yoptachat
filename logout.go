package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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
