package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"parser/internal/model"
	service "parser/internal/usecase"
)

type DomainHandler struct {
	svc *service.DomainService
}

func NewDomainHandler(svc *service.DomainService) *DomainHandler {
	return &DomainHandler{svc: svc}
}

func (h *DomainHandler) PostDomains(c *gin.Context) {
	var domains []model.Domain
	if err := c.ShouldBindJSON(&domains); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, d := range domains {
		err := h.svc.Save(&d) 
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Saved"})
}

func (h *DomainHandler) GetDomains(c *gin.Context) {
	domains, err := h.svc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch"})
		return
	}
	c.JSON(http.StatusOK, domains)
}