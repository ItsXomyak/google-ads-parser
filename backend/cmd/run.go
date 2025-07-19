package cmd

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"parser/internal/adapter/repository"
	"parser/internal/handler"
	"parser/internal/model"
	usecase "parser/internal/usecase"
	"parser/pkg/postgres"
)
func Run() {
	_ = godotenv.Load()

		db, err := postgres.NewPostgresConnection()
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    if err := postgres.TestConnection(db); err != nil {
        log.Fatal("Database connection test failed:", err)
    }

		log.Println("Connected to database successfully")


	_ = db.AutoMigrate(&model.Domain{})

	repo := repository.NewDomainRepository(db)
	svc := usecase.NewDomainService(repo)
	h := handler.NewDomainHandler(svc)

	r := gin.Default()

	r.POST("/domains", h.PostDomains)
	r.GET("/domains", h.GetDomains)
	r.DELETE("/domains/:id", h.DeleteDomain)
	r.POST("/domains/clear", h.CleanData)
	r.POST("/parse/batch", h.ParseBatch)
	r.POST("/stream/parse/batch", h.StreamParseBatch) // пока не рабочая хрень

	r.Run(":8080")
}
