package services

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
	constants "github.com/ntu-onemdp/onemdp-backend/config"
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/repositories"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type EduvisorService struct {
	EduvisorModel *models.User
	apiKey        string
	endpoint      string
}

var Eduvisor *EduvisorService

func NewEduvisorService() *EduvisorService {
	// Retrieve eduvisor user object
	eduvisor := repositories.Users.GetEduvisor()

	// Eduvisor does not exist: create user and add
	if eduvisor == nil {
		utils.Logger.Trace().Msg("eduvisor not found in system, creating from scratch. ")

		uid, err := uuid.NewRandom()
		if err != nil {
			utils.Logger.Error().Err(err).Msg("Error generating UUID")
			utils.Logger.Warn().Msg("Eduvisor is not initialized")
			return nil
		}

		// Create eduvisor user
		eduvisor = models.CreateUser(uid.String(), constants.EDUVISOR_NAME, constants.EDUVISOR_EMAIL, "N.A.", "bot")

		// Insert into users table directly. (Do not use services as this is not a normal student user.)
		if err = repositories.Users.DangerouslyInsertUser(eduvisor); err != nil {
			utils.Logger.Error().Err(err).Msg("Error inserting eduvisor into users")
			utils.Logger.Warn().Msg("Eduvisor is not initialized")
			return nil
		}
	}

	// // Generate JWT token for Eduvisor
	// claim := models.NewClaim(eduvisor.Uid)
	// jwt, err := services.JwtHandler.GenerateJwt(claim)
	// if err != nil {
	// 	utils.Logger.Warn().Msg("Error generating jwt for eduvisor")
	// }

	// // Save jwt token to /eduvisor
	// err = os.WriteFile("config/eduvisor-jwt.txt", []byte(jwt), 0644)
	// if err != nil {
	// 	utils.Logger.Warn().Msg("Error writing jwt key to file")
	// }

	// // Show eduvisor jwt token
	// utils.Logger.Info().Str("eduvisor token", jwt).Msg("")

	// Retrieve eduvisor api key from environment variables.
	apiKey := os.Getenv("EDUVISOR_API_KEY")

	// Retrieve eduvisor endpoint from environment variables
	endpoint := os.Getenv("EDUVISOR_API_ENDPOINT")

	utils.Logger.Info().Str("endpoint", endpoint).Bool("Is API key set", len(apiKey) > 0).Msg("Eduvisor service initialized.")

	return &EduvisorService{
		EduvisorModel: eduvisor,
		apiKey:        apiKey,
		endpoint:      endpoint,
	}
}

type postSummary struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type sendThreadResponse struct {
	Title    string `json:"post_title"`
	Content  string `json:"post_content"`
	Response string `json:"response"`
}

func (s *EduvisorService) SendThread(post *models.DbPost) {
	postSummary := postSummary{
		Title:   post.Title,
		Content: post.PostContent,
	}

	// Serialize post content and title.
	data, err := json.Marshal(postSummary)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error encountered when serializing post")
		return
	}

	endpoint := s.endpoint + "/response"

	req, _ := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(data))
	req.Header.Add("x-api-key", s.apiKey)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error sending to eduvisor")
		return
	}

	defer res.Body.Close()

	utils.Logger.Debug().Msg("response code: " + res.Status)

	rawBody, _ := io.ReadAll(res.Body)

	response := &sendThreadResponse{}
	err = json.Unmarshal(rawBody, &response)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error unmarshalling response")
		return
	}

	utils.Logger.Info().Bytes("raw response body", rawBody).Interface("response", response).Msg("Response received from eduvisor")
}
