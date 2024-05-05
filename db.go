package main

import (
	"fmt"         // Пакет для форматированного ввода и вывода
	"log"         // Пакет для логирования
	"strings"     // Пакет для работы со строками

	"gorm.io/driver/postgres" // Драйвер PostgreSQL для GORM
	"gorm.io/gorm"            // ORM библиотека
)

var db *gorm.DB
var err error

// Структура, описывающая модель книги в базе данных
type Book struct {
	gorm.Model  // Встраивание поля gorm.Model для ID, CreatedAt и т.д.
	Title       string
	Author      string
	Description string
	Genre       string
}

// Функция инициализации базы данных
func initDB() {
	// Строка подключения к базе данных
	dsn := "host=localhost user=youruser dbname=yourdb password=yourpassword port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err) // Логирование ошибки, если подключение не удалось
	}

	db.AutoMigrate(&Book{}) // Автоматическое создание и миграция таблиц
}

// Функция добавления книги в базу данных
func addBook(book Book) error {
	return db.Create(&book).Error // Создание записи книги в базе данных
}

// Функция получения всех книг из базы данных
func getAllBooks() ([]Book, error) {
	var books []Book
	result := db.Find(&books) // Поиск всех книг
	return books, result.Error
}

// Функция поиска книг по названию или автору
func searchBooks(query string) ([]Book, error) {
	var books []Book
	result := db.Where("lower(title) LIKE ? OR lower(author) LIKE ?", "%"+strings.ToLower(query)+"%", "%"+strings.ToLower(query)+"%").Find(&books)
	return books, result.Error
}

// Функция удаления книги по идентификатору
func deleteBook(bookID uint) error {
	return db.Delete(&Book{}, bookID).Error // Удаление книги по ID
}

// Функция получения книги по идентификатору
func getBookByID(bookID uint) (*Book, error) {
	var book Book
	result := db.First(&book, bookID) // Поиск первой записи по ID
	return &book, result.Error
}

// Функция получения списка книг по жанру
func getBooksByGenre(genre string) ([]Book, error) {
	var books []Book
	result := db.Where("lower(genre) = ?", strings.ToLower(genre)).Find(&books) // Поиск книг по жанру
	return books, result.Error
}
