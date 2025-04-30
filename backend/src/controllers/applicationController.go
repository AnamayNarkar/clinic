package controllers

import (
	"net/http"
	"strconv"

	"clinic/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ApplicationController handles application CRUD operations
type ApplicationController struct {
	DB *sqlc.Queries
}

// NewApplicationController creates a new ApplicationController
func NewApplicationController(db *sqlc.Queries) *ApplicationController {
	return &ApplicationController{DB: db}
}

// CreateApplication creates a new application
func (ac *ApplicationController) CreateApplication(c *gin.Context) {
	var req sqlc.CreateApplicationParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	app, err := ac.DB.CreateApplication(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusCreated, app)
}

// GetApplication retrieves an application by ID
func (ac *ApplicationController) GetApplication(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	app, err := ac.DB.GetApplication(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
		return
	}
	c.JSON(http.StatusOK, app)
}

// ListApplications returns paginated applications
func (ac *ApplicationController) ListApplications(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	apps, err := ac.DB.ListApplications(c, sqlc.ListApplicationsParams{Limit: int32(limit), Offset: int32(offset)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, apps)
}

// UpdateApplication updates an existing application
func (ac *ApplicationController) UpdateApplication(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req sqlc.UpdateApplicationParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	req.ID = id
	app, err := ac.DB.UpdateApplication(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, app)
}
