package auth

import (
	"fmt"
	"yoptachat/models"
)

// Register регистрирует нового пользователя.
func (a *AuthService) Register(login, phone, password string) (int, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		fmt.Println("пароль не зашифровался")
		return 0, err
	}

	user := models.User{Login: login, Phone: phone, Password: hashedPassword}                                                      //создание user'а подобного структуре User и присвоение ему регистрационных данных
	result, err := a.db.Exec("INSERT INTO Users (login, phone, password) VALUES (?, ?, ?)", user.Login, user.Phone, user.Password) //выполнение запроса по добавлению пользователя в бд->таблицу Users
	if err != nil {
		fmt.Println("пользователь не добавлен в бд")
		return 0, err
	}

	id, err := result.LastInsertId()
	fmt.Println(id)
	return int(id), err
}
