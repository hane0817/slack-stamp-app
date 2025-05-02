package controllers

import (
	"backend/usecase"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	//dbを内側に閉じる　Dbとの違い
	db *sql.DB
}

func NewUserController(db *sql.DB) *UserController {
	return &UserController{db: db}
}

func (uc *UserController) RegisterUser(c *gin.Context) {
	var req struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// コントローラーがアプリケーションに依存するように入れなおす
	input := usecase.RegisterUserInput{
		Name:     req.Name,
		Password: req.Password,
	}

	result, err := usecase.RegisterUser(Db, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": result.ID, "name": result.Name})
}
