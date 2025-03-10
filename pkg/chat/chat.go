package chat

import (
    "log"
    "net/http"
    "yoptachat/pkg/db"
    "html/template"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("session")
    if err != nil {
        http.Redirect(w, r, "/regauth.html", http.StatusSeeOther)
        return
    }

    login := cookie.Value
    var userInfo string
    err = db.DB.QueryRow("SELECT login FROM users WHERE login = ?", login).Scan(&userInfo)
    if err != nil {
        log.Println("Ошибка получения данных пользователя:", err)
        http.Redirect(w, r, "/regauth.html", http.StatusSeeOther)
        return
    }

    renderTemplate(w, "templates/index.html", userInfo)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data ...interface{}) {
    t, err := template.ParseFiles(tmpl)
    if err != nil {
        log.Println("Ошибка при рендеринге шаблона:", err)
        http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
        return
    }
    t.Execute(w, data)
}