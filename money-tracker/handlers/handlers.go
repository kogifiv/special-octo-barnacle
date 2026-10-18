package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	// ЗАМЕНИ "myproject" НА НАЗВАНИЕ СВОЕГО МОДУЛЯ ИЗ GO.MOD:
	"money-tracker/database"
	"money-tracker/models"
)

// Главная страница: читает данные из базы и рендерит HTML
func HomePage(w http.ResponseWriter, r *http.Request) {
	// Если ввели левый путь, отдаем 404
	if r.URL.Path != "/" {
		http.Error(w, "Страница не найдена", 404)
		return
	}

	// Делаем обычный SELECT запрос к базе данных
	rows, err := database.DB.Query("SELECT id, amount, description, date FROM expenses ORDER BY date DESC")
	if err != nil {
		http.Error(w, "Ошибка запроса к БД: "+err.Error(), 500)
		return
	}
	defer rows.Close()

	var list []models.Expense

	// Перебираем строки из ответа базы
	for rows.Next() {
		var exp models.Expense
		err := rows.Scan(&exp.ID, &exp.Amount, &exp.Description, &exp.Date)
		if err == nil {
			list = append(list, exp)
		}
	}

	// Проверка для линтера (rows.Err)
	if dbErr := rows.Err(); dbErr != nil {
		http.Error(w, "Ошибка чтения строк бд: "+dbErr.Error(), 500)
		return
	}

	// Собираем шаблоны вместе
	t, err := template.ParseFiles("templates/layout.html", "templates/index.html")
	if err != nil {
		http.Error(w, "Ошибка парсинга шаблонов: "+err.Error(), 500)
		return
	}

	// Запускаем рендер и передаем массив list
	err = t.ExecuteTemplate(w, "layout", list)
	if err != nil {
		http.Error(w, "Ошибка вывода HTML: "+err.Error(), 500)
		return
	}
}

// Обработчик добавления новой траты из формы
func AddExpense(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Разрешен только метод POST", 405)
		return
	}

	// Достаем данные из полей HTML-формы
	r.ParseForm()
	sumStr := r.FormValue("amount")
	info := r.FormValue("description")
	dt := r.FormValue("date")

	// Парсим сумму в число
	sum, _ := strconv.Atoi(sumStr)

	// Записываем в базу. Плейсхолдеры '?' защищают от SQL-инъекций
	_, err := database.DB.Exec("INSERT INTO expenses (amount, description, date) VALUES (?, ?, ?)", sum, info, dt)
	if err != nil {
		http.Error(w, "Ошибка вставки строки в БД", 500)
		return
	}

	// Паттерн PRG: делаем редирект, чтобы форма не отправлялась повторно при F5
	http.Redirect(w, r, "/", 302)
}
