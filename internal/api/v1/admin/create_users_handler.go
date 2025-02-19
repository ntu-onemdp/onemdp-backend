package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type CreateUserHandler struct {
	UserService *services.UserService
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
