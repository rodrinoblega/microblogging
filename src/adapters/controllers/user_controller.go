package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rodrinoblega/microblogging/src/entities"
	"net/http"
)

type CreateUserUseCaseInterface interface {
	Execute(username string) (*entities.User, error)
}

type UserController struct {
	createUserUseCase CreateUserUseCaseInterface
}

func NewUserController(createUserUseCase CreateUserUseCaseInterface) *UserController {
	return &UserController{createUserUseCase: createUserUseCase}
}

func (uc *UserController) CreateUser(c *gin.Context) {
	var request struct {
		Username string `json:"username"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := uc.createUserUseCase.Execute(request.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": fmt.Sprintf("User created successfully. ID %v", user.ID)})
}
