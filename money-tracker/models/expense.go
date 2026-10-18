package models

// Структура для одной траты.
// Названия полей с большой буквы, чтобы Go видел их из других пакетов и в HTML.
type Expense struct {
	ID          int
	Amount      int
	Description string
	Date        string
}
