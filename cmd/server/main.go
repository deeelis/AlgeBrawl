package main

import (
    "log"
    "os"

    "algebrawl/internal/api"
    "algebrawl/internal/database"

    "github.com/gin-gonic/gin"
)

func main() {
    // Подключение к БД
    connStr := "host=postgres port=5432 user=algebrawl password=password dbname=algebrawl sslmode=disable"
    repo, err := database.NewRepository(connStr)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer repo.db.Close()

    // Создание HTTP сервера
    router := gin.Default()
    
    handler := api.NewHandler(repo)
    
    // Регистрация роутов
    router.POST("/register", handler.Register)
    router.GET("/api/new", handler.NewEquationSet)
    router.GET("/api/list", handler.GetEquationSet)
    router.POST("/api/statistics", handler.Statistics)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Server starting on port %s", port)
    if err := router.Run(":" + port); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
