package admin

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

// Raw request from frontend
type CreateNewUsersRequest struct {
	Semester string    `json:"semester" binding:"required"`
	Users    []NewUser `json:"users" binding:"required"`
}

// Model user
type NewUser struct {
	Email string `json:"email" binding:"required"`
	// Role  string `json:"role" binding:"required"` // 21 July: all new users are granted student roles by default.
}

// Response sent to frontend
type CreateUserResponse struct {
	NumSuccess int                  `json:"num_success"`
	NumFailed  int                  `json:"num_failed"`
	Results    []SingleUserResponse `json:"results"`
}

type SingleUserResponse struct {
	Email  string `json:"email"`
	Result string `json:"result"`
}

func CreateUsersHandler(c *gin.Context) {
	utils.Logger.Info().Msg("Create new user request received")
	var createNewUsersRequest CreateNewUsersRequest
	createUserResponse := CreateUserResponse{
		NumSuccess: 0,
		NumFailed:  0,
		Results:    []SingleUserResponse{},
	}

	if err := c.BindJSON(&createNewUsersRequest); err != nil {
		utils.Logger.Error().Err(err).Msg("Error binding request to NewUsers")
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	for _, newUser := range createNewUsersRequest.Users {
		singleUserResult := SingleUserResponse{
			Email:  newUser.Email,
			Result: "not created", // Defaults to not created
		}

		// Call service to create new user
		if err := services.Users.CreateNewUser(strings.ToUpper(newUser.Email), models.Student.String(), createNewUsersRequest.Semester); err != nil {
			utils.Logger.Error().Err(err).Msg("Error encountered when inserting new user")
			singleUserResult.Result = "failed"
		} else {
			singleUserResult.Result = "success"
		}

		// Append the result to overall response
		createUserResponse.Results = append(createUserResponse.Results, singleUserResult)
		if singleUserResult.Result == "failed" {
			createUserResponse.NumFailed++
		} else {
			createUserResponse.NumSuccess++
		}
	}

	c.JSON(http.StatusCreated, createUserResponse)
}
