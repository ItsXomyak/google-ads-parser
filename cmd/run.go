package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"parser/internal/adapter/repository"
	"parser/internal/handler"
	"parser/internal/model"
	usecase "parser/internal/usecase"
	connect "parser/pkg/postgres"
)
func Run() {
	_ = godotenv.Load()

	db, err := connect.NewPostgresConnection()
	if err != nil {
		panic(err)
	}

	_ = db.AutoMigrate(&model.Domain{})

	repo := repository.NewDomainRepository(db)
	svc := usecase.NewDomainService(repo)
	h := handler.NewDomainHandler(svc)

	r := gin.Default()

	r.POST("/domains", h.PostDomains)
	r.GET("/domains", h.GetDomains)

	r.Run(":8080")
}
