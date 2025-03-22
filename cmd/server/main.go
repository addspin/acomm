package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/addspin/acomm/pkg/models"
)

const (
	dbPath = "./db/comm.db"
	addr   = ":8082"
)

type CommentRequest struct {
	Text string `json:"text"`
}

// Получение всех комментариев для новости по ID
func handleGetComments(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID новости из параметра запроса
	newsIDStr := r.URL.Query().Get("id")
	if newsIDStr == "" {
		http.Error(w, "Отсутствует ID новости", http.StatusBadRequest)
		return
	}

	newsID, err := strconv.ParseInt(newsIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Некорректный ID новости", http.StatusBadRequest)
		return
	}

	comments, err := models.GetCommentsByNewsID(newsID)
	if err != nil {
		log.Printf("Ошибка при получении комментариев: %v", err)
		http.Error(w, "Не удалось получить комментарии", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

// Добавление нового комментария к новости
func handleAddComment(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID новости из параметра запроса
	newsIDStr := r.URL.Query().Get("id")
	if newsIDStr == "" {
		http.Error(w, "Отсутствует ID новости", http.StatusBadRequest)
		return
	}

	newsID, err := strconv.ParseInt(newsIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Некорректный ID новости", http.StatusBadRequest)
		return
	}

	var req CommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Некорректный формат запроса", http.StatusBadRequest)
		return
	}

	if req.Text == "" {
		http.Error(w, "Текст комментария не может быть пустым", http.StatusBadRequest)
		return
	}

	commentID, err := models.AddComment(newsID, req.Text)
	if err != nil {
		log.Printf("Ошибка при добавлении комментария: %v", err)
		http.Error(w, "Не удалось добавить комментарий", http.StatusInternalServerError)
		return
	}

	response := struct {
		ID int64 `json:"id"`
	}{
		ID: commentID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Настройка маршрутов HTTP
func setupRoutes() {
	// Эндпоинт для получения комментариев для конкретной новости
	http.HandleFunc("/api/comm_news", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
			return
		}
		handleGetComments(w, r)
	})

	// Эндпоинт для добавления комментариев к конкретной новости
	http.HandleFunc("/api/comm_add_news", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
			return
		}
		handleAddComment(w, r)
	})
}

func main() {
	// Инициализация базы данных
	if err := models.InitDB(dbPath); err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}

	// Настройка маршрутов
	setupRoutes()

	// Обработка корректного завершения
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Запуск сервера на %s", addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatalf("Ошибка сервера: %v", err)
		}
	}()

	// Ожидание сигнала прерывания
	<-c
	log.Println("Завершение работы сервера...")
}
