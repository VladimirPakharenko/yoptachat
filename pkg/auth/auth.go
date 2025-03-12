package auth

import (
    "database/sql"
    "net/http"
    "golang.org/x/crypto/bcrypt"
	"github.com/gorilla/sessions"
    "yoptachat/pkg/models"
)

// AuthService предоставляет методы для аутентификации пользователей.
type AuthService struct {
    db *sql.DB
	store *sessions.CookieStore
}

// NewAuthService создает новый AuthService.
func NewAuthService(db *sql.DB, store *sessions.CookieStore) *AuthService {
    return &AuthService{db, store}
}

// Register регистрирует нового пользователя.
func (a *AuthService) Register(login, phone, password string) (int, error) {
    hashedPassword, err := hashPassword(password)
    if err != nil {
        return 0, err
    }

    user := models.User{Login: login, Phone: phone, Password: hashedPassword} //создание user'а подобного структуре User и присвоение ему регистрационных данных
    result, err := a.db.Exec("INSERT INTO Users (login, phone, password) VALUES (?, ?, ?)", user.Login, user.Phone, user.Password)//выполнение запроса по добавлению пользователя в бд->таблицу Users
    if err != nil {
        return 0, err
    }

    id, err := result.LastInsertId()
    return int(id), err
}

// Login аутентифицирует пользователя и создает сессию.
func (a *AuthService) Login(w http.ResponseWriter, r *http.Request) {
    login := r.FormValue("login")
    password := r.FormValue("password")

    user, err := a.getUserByLogin(login)
    if err != nil || !comparePasswords(user.Password, password) {
        http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
        return
    }

    // Создание сессии
    session, _ := a.store.Get(r, "session-name")
    session.Values["user"] = user.ID
    session.Save(r, w)

    http.Redirect(w, r, "/index.html", http.StatusSeeOther)
}

// getUserByLogin получает пользователя по логину.
func (a *AuthService) getUserByLogin(login string) (*models.User, error) {
    user := &models.User{}
    err := a.db.QueryRow("SELECT id, login, phone, password FROM Users WHERE login = ?", login).Scan(&user.ID, &user.Login, &user.Phone, &user.Password)
    if err != nil {
        return nil, err
    }
    return user, nil
}

// hashPassword хэширует пароль.
func hashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

// comparePasswords сравнивает хэшированный пароль с обычным.
func comparePasswords(hashedPwd, plainPwd string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
    return err == nil
}
