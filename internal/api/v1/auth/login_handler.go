package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type LoginHandler struct {
	AuthService *services.AuthService
}

// Login form sent from frontend.
type loginForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"` // Plaintext password
}

type LoginResponse struct {
	Success  bool   `json:"success"`
	ErrorMsg string `json:"error_msg"`
	Jwt      string `json:"jwt"`
	Role     string `json:"role"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

func (h *LoginHandler) HandleLogin(c *gin.Context) {
	utils.Logger.Trace().Msg("Login request received")
	var form loginForm

	// Bind with form
	if err := c.ShouldBind(&form); err != nil {
		utils.Logger.Error().Err(err).Msg("Error processing login request")
		response := LoginResponse{
			Success:  false,
			ErrorMsg: "Malformed request",
		}

		c.JSON(200, &response) // TODO: Change the error code in the future
		return
	}

	// Authenticate user
	isAuthenticated, user, role := h.AuthService.AuthenticateUser(form.Username, form.Password)
	if !isAuthenticated {
		response := LoginResponse{
			Success:  false,
			ErrorMsg: "Unauthorized: Incorrect username/password",
		}

		c.JSON(200, &response)
		return
	}

	// Generate jwt
	claim := models.NewClaim(user, role)
	jwt, err := services.JwtHandler.GenerateJwt(claim)
	if err != nil {
		utils.Logger.Error().Err(err)
		c.JSON(500, "Internal server error")
	}

	response := LoginResponse{
		Success:  true,
		Jwt:      jwt,
		Role:     role,
		Name:     user.Name,
		Username: user.Username,
	}
	c.JSON(200, &response)
}
