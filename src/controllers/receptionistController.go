package controllers

import (
	"clinic/sqlc"
	"clinic/src/dto"
	"clinic/src/security"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReceptionistController struct {
	DB *sqlc.Queries
}

func NewReceptionistController(db *sqlc.Queries) *ReceptionistController {
	return &ReceptionistController{DB: db}
}

func (rc *ReceptionistController) CreateReceptionist(c *gin.Context) {
	var req dto.ReceptionistRegistrationDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	passwordHash, salt, err := security.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot hash password"})
		return
	}
	newRec := sqlc.CreateReceptionistParams{
		Username:     req.Username,
		Email:        req.Email,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		PasswordHash: string(passwordHash),
		Salt:         salt,
	}
	rec, err := rc.DB.CreateReceptionist(c, newRec)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating receptionist" + " " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, rec)
}
