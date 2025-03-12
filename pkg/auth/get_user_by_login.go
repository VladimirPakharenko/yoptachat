package auth

import "yoptachat/pkg/models"

// getUserByLogin получает пользователя по логину.
func (a *AuthService) getUserByLogin(login string) (*models.User, error) {
	user := &models.User{}
	err := a.db.QueryRow("SELECT id, login, phone, password FROM Users WHERE login = ?", login).Scan(&user.ID, &user.Login, &user.Phone, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}
