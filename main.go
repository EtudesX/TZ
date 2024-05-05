package main

import (
	"encoding/json" // Пакет для работы с JSON
	"fmt"           // Пакет для форматированного ввода и вывода
	"github.com/gorilla/mux" // Пакет для маршрутизации веб-запросов
	"log"           // Пакет для логирования
	"net/http"      // Пакет для работы с HTTP
)

func main() {
	initDB() // Инициализация базы данных
	defer db.Close() // Закрытие соединения с базой данных при завершении работы

	r := mux.NewRouter() // Создание нового роутера

	// Определение маршрутов
	r.HandleFunc("/books", getBooksHandler).Methods("GET")
	r.HandleFunc("/book", addBookHandler).Methods("POST")
	r.HandleFunc("/book/{id}", getBookHandler).Methods("GET")
	r.HandleFunc("/book/{id}", deleteBookHandler).Methods("DELETE")
	r.HandleFunc("/genres", getGenresHandler).Methods("GET")
	r.HandleFunc("/books/genre/{genre}", getBooksByGenreHandler).Methods("GET")
	r.HandleFunc("/books/search", searchBooksHandler).Methods("GET")

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r)) // Запуск сервера на порту 8080
}

// Обработчик для получения всех книг
func getBooksHandler(w http.ResponseWriter, r *http.Request) {
	books, err := getAllBooks() // Получение списка всех книг
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(books) // Вывод списка книг в формате JSON
}

// Обработчик для добавления книги
func addBookHandler(w http.ResponseWriter, r *http.Request) {
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := addBook(book); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(book) // Вывод добавленной книги в формате JSON
}

// Обработчик для получения книги по ID
func getBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID := vars["id"]
	book, err := getBookByID(bookID)
	if err != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(book) // Вывод книги в формате JSON
}

// Обработчик для удаления книги по ID
func deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID := vars["id"]
	if err := deleteBook(bookID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK) // Отправка статуса успешного удаления
}

// Обработчик для получения списка жанров
func getGenresHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK) // Здесь должна быть логика получения жанров
}

// Обработчик для получения списка книг по жанру
func getBooksByGenreHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	genre := vars["genre"]
	books, err := getBooksByGenre(genre)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(books) // Вывод списка книг в формате JSON
}

// Обработчик для поиска книг по названию или автору
func searchBooksHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q") // Получение поискового запроса из URL
	if query == "" {
		http.Error(w, "Query parameter 'q' is missing", http.StatusBadRequest)
		return
	}
	books, err := searchBooks(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(books) // Вывод результатов поиска в формате JSON
}
