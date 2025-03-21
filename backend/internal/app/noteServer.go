// @title DropNote API
// @version 1.0
// @description API для сервиса заметок с хранением в PostgreSQL и MinIO
// @host localhost:8080
// @BasePath /
package app

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/m3owmurrr/DropNote/backend/docs"
	"github.com/m3owmurrr/DropNote/backend/internal/handler"
	"github.com/m3owmurrr/DropNote/backend/internal/repository/miniorepository"
	"github.com/m3owmurrr/DropNote/backend/internal/repository/pgrepository"
	"github.com/m3owmurrr/DropNote/backend/internal/service"
	"github.com/m3owmurrr/DropNote/backend/internal/utils/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/cors"
)

// RunServer инициализирует и запускает HTTP-сервер DropNote.
//
// 1. Подключается к PostgreSQL для хранения метаданных заметок.
// 2. Создаёт MinIO клиент для хранения содержимого заметок.
// 3. Регистрирует маршруты API через Gorilla Mux.
// 4. Запускает сервер на порту, указанном в `config.Cfg.Server`.
//
// Сервер обрабатывает два основных маршрута API:
// - `POST /api/notes` — создание новой заметки.
// - `GET /api/notes` — получение заметки по ID.
//
// Сервер логирует события через `slog.Logger`.
//
// При ошибках подключения к базе данных или MinIO — завершает работу.
func RunServer() {
	opts := &slog.HandlerOptions{
		AddSource: true,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))

	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/DropNote?sslmode=disable")
	if err != nil {
		logger.Error("database conection error:", "error", err)
		os.Exit(1)
	}

	minioClient, err := minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false,
	})

	noteMetaRepository := pgrepository.NewNoteRepository(db)
	noteS3Repository := miniorepository.NewNoteRepository(minioClient)
	noteService := service.NewNoteSevice(noteMetaRepository, noteS3Repository)
	noteHandler := handler.NewNoteHandler(noteService, logger)
	healthHandler := handler.NewHealthHandler()

	router := mux.NewRouter()
	router.HandleFunc("/health", healthHandler.CheckHealth).Methods(http.MethodGet)

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Разрешить все домены
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	router.Use(mux.CORSMethodMiddleware(router)) // Включаем поддержку CORS в mux

	subRouter := router.PathPrefix("/api").Subrouter()
	subRouter.HandleFunc("/notes", noteHandler.CreateNote).Methods(http.MethodPost, http.MethodOptions) // Добавлен OPTIONS
	// subRouter.HandleFunc("/notes/public", noteHandler.GetPublicNotes).Methods(http.MethodGet)
	subRouter.HandleFunc("/notes/{note_id}", noteHandler.GetNote).Methods(http.MethodGet, http.MethodOptions)

	routerWithCORS := c.Handler(router) // Оборачиваем роутер в CORS middleware

	address := fmt.Sprintf("%s:%s", config.Cfg.Server.Host, config.Cfg.Server.Port)
	server := http.Server{
		Addr:         address,
		Handler:      routerWithCORS,
		ReadTimeout:  config.Cfg.Server.Timeout,
		WriteTimeout: config.Cfg.Server.Timeout,
		IdleTimeout:  config.Cfg.Server.IdleTimeout,
	}

	logger.Info("Server running...", "host", address)
	if err := server.ListenAndServe(); err != nil {
		logger.Error("Server cannot run: ", "error", err)
		os.Exit(1)
	}
}
