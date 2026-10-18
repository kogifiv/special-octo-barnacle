package main

import (
	"fmt"
	"net/http"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Страница не найдена", 404)
		return
	}
	fmt.Fprint(w, "Добро пожаловать в учет личных трат")
}

func aboutPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Семестровое практическое задание: Веб-разработка на Go: Учёт личных трат")
}

func pingPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Метод не поддерживается", 405)
		return
	}

	w.WriteHeader(200)
	fmt.Fprint(w, "pong")
}

func main() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/about", aboutPage)
	http.HandleFunc("/ping", pingPage)

	fmt.Println("Сервер стартовал на http://localhost:8090")

	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		fmt.Println("Ошибка при запуске:", err)
	}
}
