package handler

import (
	"math"
	"net/http"
	"strconv"

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
	domain := c.Query("domain")
	limitStr := c.DefaultQuery("limit", "10")
	pageStr := c.DefaultQuery("page", "1")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	offset := (page - 1) * limit

	// получаем отфильтрованные записи и общее количество
	data, total, err := h.svc.GetPaginated(domain, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch"})
		return
	}

	pages := int(math.Ceil(float64(total) / float64(limit)))

	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"pages": pages,
		"page":  page,
		"limit": limit,
		"data":  data,
	})
}

func (h *DomainHandler) ParseAndSave(c *gin.Context) {
	domain := c.Query("domain")
	if domain == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing domain"})
		return
	}

	entity, err := h.svc.ParseAndSave(domain)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, entity)
}


func (h *DomainHandler) DeleteDomain(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	err = h.svc.DeleteByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *DomainHandler) CleanData(c *gin.Context) {
	if err := h.svc.CleanData(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "db cleaned"})
}

func (h *DomainHandler) ParseBatch(c *gin.Context) {
	var req struct {
		Domains []string `json:"domains"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || len(req.Domains) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "domains[] required"})
		return
	}

	saved, err := h.svc.ParseBatchAndSave(req.Domains)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, saved)
}

func (h *DomainHandler) StreamParseBatch(c *gin.Context) {
	var req struct {
		Domains []string `json:"domains"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || len(req.Domains) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "domains[] required"})
		return
	}

	saved, err := h.svc.StreamParseBatchAndSave(req.Domains)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, saved)
}