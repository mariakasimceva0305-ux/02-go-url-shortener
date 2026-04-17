package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"go-url-shortener/internal/service"
)

type Handler struct {
	service *service.LinkService
}

func NewHandler(service *service.LinkService) *Handler {
	return &Handler{service: service}
}

type createLinkRequest struct {
	OriginalURL string `json:"original_url" binding:"required,url"`
	TTLHours    int    `json:"ttl_hours"`
}

func (h *Handler) Register(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/swagger", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"title": "Go URL Shortener API",
			"hint":  "Use endpoints under /api/v1 and /r/{code}",
		})
	})

	v1 := r.Group("/api/v1")
	{
		v1.POST("/links", h.createLink)
		v1.GET("/links", h.listLinks)
		v1.GET("/links/:id", h.getLink)
		v1.GET("/links/:id/analytics", h.analytics)
		v1.DELETE("/links/:id", h.deactivate)
	}
	r.GET("/r/:code", h.redirect)
}

func (h *Handler) createLink(c *gin.Context) {
	var req createLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	link, err := h.service.Create(req.OriginalURL, req.TTLHours)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"id":           link.ID,
		"original_url": link.OriginalURL,
		"short_code":   link.ShortCode,
		"short_url":    h.service.ShortURL(link.ShortCode),
		"created_at":   link.CreatedAt,
		"expires_at":   link.ExpiresAt,
	})
}

func (h *Handler) listLinks(c *gin.Context) {
	items, err := h.service.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *Handler) getLink(c *gin.Context) {
	item, err := h.service.Get(c.Param("id"))
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *Handler) analytics(c *gin.Context) {
	out, err := h.service.Analytics(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *Handler) deactivate(c *gin.Context) {
	if err := h.service.Deactivate(c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) redirect(c *gin.Context) {
	link, err := h.service.ResolveByCode(c.Param("code"), c.Request.UserAgent(), c.ClientIP(), c.GetHeader("Referer"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, link.OriginalURL)
}
