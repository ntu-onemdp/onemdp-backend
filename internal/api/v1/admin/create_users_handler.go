package admin

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type CreateUserHandler struct {
	UserService *services.UserService
	AuthService *services.AuthService
}

// Raw request from frontend
type CreateNewUsersRequest struct {
	Semester int       `json:"semester" binding:"required"`
	Users    []NewUser `json:"users" binding:"required"`
}

// Model user
type NewUser struct {
	Username string `json:"username" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

// Response sent to frontend
type CreateUserResponse struct {
	NumSuccess int                  `json:"num_success"`
	NumFailed  int                  `json:"num_failed"`
	Results    []SingleUserResponse `json:"results"`
}

type SingleUserResponse struct {
	Username string `json:"username"`
	Result   string `json:"result"`
}

func (h *CreateUserHandler) HandleCreateNewUser(c *gin.Context) {
	utils.Logger.Info().Msg("Create new user request received")
	var createNewUsersRequest CreateNewUsersRequest
	createUserResponse := CreateUserResponse{
		NumSuccess: 0,
		NumFailed:  0,
		Results:    []SingleUserResponse{},
	}

	if err := c.BindJSON(&createNewUsersRequest); err != nil {
		utils.Logger.Error().Err(err).Msg("Error binding request to NewUsers")
		c.JSON(400, nil)
	}

	// Create file to store default passwords
	file, err := os.OpenFile("new_users.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error creating file for password storage")
	}
	defer file.Close()

	for _, newUser := range createNewUsersRequest.Users {
		singleUserResult := SingleUserResponse{
			Username: newUser.Username,
			Result:   "",
		}

		// Call service to create new user
		if err := h.UserService.CreateNewUser(newUser.Username, newUser.Name, createNewUsersRequest.Semester); err != nil {
			utils.Logger.Error().Err(err).Msg("Error encountered when inserting new user")
			singleUserResult.Result = "failed"
		} else {
			singleUserResult.Result = "success"

			// Success: Continue to insert to auth table
			default_password := utils.GeneratePassword()                                 // Generate default password
			file.WriteString(fmt.Sprintf("%s,%s\n", newUser.Username, default_password)) // Write username and password to file
			if err := h.AuthService.InsertNewAuth(newUser.Username, default_password); err != nil {
				utils.Logger.Error().Err(err).Msg("Error encountered when inserting new auth")
				singleUserResult.Result = "failed"
			}

		}

		// Append the result to overall response
		createUserResponse.Results = append(createUserResponse.Results, singleUserResult)
		if singleUserResult.Result == "failed" {
			createUserResponse.NumFailed++
		} else {
			createUserResponse.NumSuccess++
		}
	}

	c.JSON(201, createUserResponse)
}
