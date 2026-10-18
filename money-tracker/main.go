package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	_ "modernc.org/sqlite" // Подключаем драйвер SQLite
)

// Структура для расходов
type Expense struct {
	ID          int
	Amount      int
	Description string
	Date        string
}

// Глобальная переменная для базы данных
var database *sql.DB

func homePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Страница не найдена", 404)
		return
	}

	// Выбираем все расходы из таблицы
	rows, err := database.Query("SELECT id, amount, description, date FROM expenses ORDER BY date DESC")
	if err != nil {
		http.Error(w, "Ошибка чтения из БД: "+err.Error(), 500)
		return
	}
	defer rows.Close()

	var listFromDB []Expense

	for rows.Next() {
		var exp Expense
		err := rows.Scan(&exp.ID, &exp.Amount, &exp.Description, &exp.Date)
		if err == nil {
			listFromDB = append(listFromDB, exp)
		}
	}

	// Проверка для линтера
	if dbErr := rows.Err(); dbErr != nil {
		http.Error(w, "Ошибка базы данных: "+dbErr.Error(), 500)
		return
	}

	// Читаем шаблоны
	t, err := template.ParseFiles("templates/layout.html", "templates/index.html")
	if err != nil {
		http.Error(w, "Ошибка шаблонов: "+err.Error(), 500)
		return
	}

	// Рендерим страницу и сразу проверяем, нет ли там ошибки
	err = t.ExecuteTemplate(w, "layout", listFromDB)
	if err != nil {
		http.Error(w, "Ошибка вывода шаблона: "+err.Error(), 500)
		return
	}
}

func addExpense(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Только POST запросы!", 405)
		return
	}

	r.ParseForm()
	sumStr := r.FormValue("amount")
	info := r.FormValue("description")
	dt := r.FormValue("date")

	sum, _ := strconv.Atoi(sumStr)

	// Запись в базу
	_, err := database.Exec("INSERT INTO expenses (amount, description, date) VALUES (?, ?, ?)", sum, info, dt)
	if err != nil {
		http.Error(w, "Не записалось в базу данных", 500)
		return
	}

	http.Redirect(w, r, "/", 302)
}

func main() {
	var err error
	// Открываем базу
	database, err = sql.Open("sqlite", "my_data.db")
	if err != nil {
		fmt.Println("БД не открылась:", err)
		return
	}
	defer database.Close()

	// Автоматическое создание таблицы
	_, err = database.Exec(`
		CREATE TABLE IF NOT EXISTS expenses (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			amount INTEGER,
			description TEXT,
			date TEXT
		);
	`)
	if err != nil {
		fmt.Println("Ошибка создания таблицы:", err)
		return
	}

	http.HandleFunc("/", homePage)
	http.HandleFunc("/add", addExpense)

	fmt.Println("Сервер с базой данных запущен на http://localhost:9000")
	http.ListenAndServe(":9000", nil)
}
