package database

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite" // драйвер для работы с файлом бд
)

// Глобальная переменная бд, которую будем дергать из обработчиков
var DB *sql.DB

// Функция инициализации базы данных
func InitDB() {
	var err error
	// Открываем (или создаем) локальный файл базы данных
	DB, err = sql.Open("sqlite", "expenses.db")
	if err != nil {
		fmt.Println("Ошибка открытия файла базы данных:", err)
		return
	}

	// Автоматически создаем таблицу при запуске, если её ещё нет
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS expenses (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			amount INTEGER,
			description TEXT,
			date TEXT
		);
	`)
	if err != nil {
		fmt.Println("Не удалось автоматически создать таблицу:", err)
	}
}
