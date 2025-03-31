package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func searchHandler(c *gin.Context) { //Поиск пользователей по тексту из инпута
	letter := c.Query("letter") //Получем значение записанное в инпут
	var results []struct {      //Создаем структуру для записи полученных пользователей со значнеиями ID и Login
		ID    int
		Login string
	}

	rows, err := db.Query("SELECT id, login FROM Users WHERE login LIKE concat('%', ?, '%')", letter) //запрос на нахождение подходящих значений из бд по данным из инпута
	if err != nil {                                                                                   //Ошибка
		c.String(http.StatusInternalServerError, "Ошибка поиска")
		return
	}
	defer rows.Close() //Закрывает rows после выполнения

	for rows.Next() { //Проходится циклом по rows
		var user struct { //создаем значение структуры в которую будут вписываться данные
			ID    int
			Login string
		}
		if err := rows.Scan(&user.ID, &user.Login); err != nil { //Вписываемданные полученные из запроса
			c.String(http.StatusInternalServerError, "Ошибка получения результатов") //Ошибка
			return
		}
		results = append(results, user) //Вписываем данные в массив
	}

	c.JSON(http.StatusOK, results) //передаем данные в виде JSON значения
}
