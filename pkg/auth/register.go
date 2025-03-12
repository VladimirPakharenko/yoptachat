package auth

import "yoptachat/pkg/models"

// Register регистрирует нового пользователя.
func (a *AuthService) Register(login, phone, password string) (int, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return 0, err
	}

	user := models.User{Login: login, Phone: phone, Password: hashedPassword}                                                      //создание user'а подобного структуре User и присвоение ему регистрационных данных
	result, err := a.db.Exec("INSERT INTO Users (login, phone, password) VALUES (?, ?, ?)", user.Login, user.Phone, user.Password) //выполнение запроса по добавлению пользователя в бд->таблицу Users
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return int(id), err
}
