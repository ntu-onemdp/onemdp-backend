package services

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
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
		eduvisor = models.CreateSpecialUser(uid.String(), constants.EDUVISOR_NAME, constants.EDUVISOR_EMAIL, models.Bot.String())

		// Insert into users table directly. (Do not use services as this is not a normal student user.)
		if err = repositories.Users.RegisterUser(eduvisor); err != nil {
			utils.Logger.Error().Err(err).Msg("Error inserting eduvisor into users")
			utils.Logger.Warn().Msg("Eduvisor is not initialized")
			return nil
		}
	}

	// Retrieve eduvisor api key from environment variables.
	apiKey, found := os.LookupEnv("EDUVISOR_API_KEY")
	if !found {
		utils.Logger.Warn().Msg("EDUVISOR_API_KEY is not set in .env")
	}

	// Retrieve eduvisor endpoint from environment variables
	endpoint, found := os.LookupEnv("EDUVISOR_API_ENDPOINT")
	if !found {
		utils.Logger.Warn().Msg("EDUVISOR_API_ENDPOINT is not set in .env")
	}

	utils.Logger.Info().Str("endpoint", endpoint).Msg("Eduvisor service initialized.")

	return &EduvisorService{
		EduvisorModel: eduvisor,
		apiKey:        apiKey,
		endpoint:      endpoint,
	}
}

// Automatically attach API key to all requests made to eduvisor.
// This function automatically prepends eduvisor endpoint to all paths as well.
func (s *EduvisorService) req(method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, s.endpoint+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-api-key", s.apiKey)
	return req, nil
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

	req, _ := s.req(http.MethodPost, "/response", bytes.NewBuffer(data))

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

	// No valid response from eduvisor
	if response.Response == "I don't know." {
		return
	}

	title := "Eduvisor response"
	if err = Posts.CreateNewPost(Eduvisor.EduvisorModel.Uid, nil, post.ThreadId, title, response.Response, false); err != nil {
		utils.Logger.Error().Err(err).Msg("Error inserting eduvisor's response to posts db")
		return
	}

	utils.Logger.Info().Msg("Eduvisor's response created")
}

// Upload file to Eduvisor. Eduvisor will update the vectorstore on its end.
// Current implementation uploads only 1 file at a time.
// TODO: Handle multiple file uploads (eduvisor accepts multiple file uploads already)
func (s *EduvisorService) Upload(file *multipart.FileHeader) error {
	// Open file
	f, err := file.Open()
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error opening file")
		return err
	}
	defer f.Close()

	// Prepare buffer and multipart writer - write bytes from file into buffer
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Create form file field
	// Eduvisor reads from the field 'files'
	fileWriter, err := w.CreateFormFile("files", file.Filename)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error creating file writer")
		return err
	}

	// Copy file data into multipart field
	if _, err := io.Copy(fileWriter, f); err != nil {
		utils.Logger.Error().Err(err).Msg("Error writing file data into multipart field")
		return err
	}

	// Close multipart writer
	if err := w.Close(); err != nil {
		utils.Logger.Error().Err(err).Msg("Error closing multipart writer")
		return err
	}

	// Create request and set headers
	req, err := s.req(http.MethodPost, "/upload", &b)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error creating request")
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Send the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error sending the request to eduvisor")
		return err
	}
	defer res.Body.Close()

	utils.Logger.Info().Str("res", res.Status).Msg("Successfully uploaded file to Eduvisor")
	return nil
}
