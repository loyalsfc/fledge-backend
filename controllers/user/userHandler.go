package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/loyalsfc/fledge-backend/internal/database"
	"github.com/loyalsfc/fledge-backend/utils"
	"golang.org/x/crypto/bcrypt"
)

type parameters struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type signin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ApiCfg struct {
	DB *database.Queries
}

func (apiCfg *ApiCfg) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := apiCfg.DB.GetUsers(r.Context())

	if err != nil {
		utils.ErrResponse(403, w, fmt.Sprintf("Error getting users %v", err))
	}

	utils.JsonResponse(200, w, utils.HandleUsersToUsers(users))

}

func (apiCfg *ApiCfg) CreateUser(w http.ResponseWriter, r *http.Request) {
	decorder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decorder.Decode(&params)

	if err != nil {
		utils.ErrResponse(500, w, "invalid parameters")
		return
	}

	username := generateUsername(params.Name)

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:             uuid.New(),
		Name:           params.Name,
		Email:          params.Email,
		Username:       username,
		Password:       passwordEncryption(params.Password),
		CoverPicture:   sql.NullString{Valid: true, String: "https://res.cloudinary.com/dplpf3g05/image/upload/fl_preserve_transparency/v1714867594/pexels-photo-5109665_hzr15h.jpg"},
		ProfilePicture: sql.NullString{Valid: true, String: "/dummy.jpg"},
	})

	if err != nil {
		utils.ErrResponse(400, w, fmt.Sprintf("Error parsing json %v", err))
		return
	}

	utils.JsonResponse(200, w, utils.HandleUserToUser(user))
}

func (apiCfg ApiCfg) Signin(w http.ResponseWriter, r *http.Request) {
	decorder := json.NewDecoder(r.Body)

	params := signin{}

	err := decorder.Decode(&params)

	if err != nil {
		utils.ErrResponse(400, w, fmt.Sprintf("Error paring json %v", err))
	}

	user, err := apiCfg.DB.SignIn(r.Context(), params.Email)

	if err != nil {
		utils.ErrResponse(401, w, "invalid email or password")
		return
	}

	isPasswordMatched := utils.ComparePassword(user.Password, params.Password)

	if !isPasswordMatched {
		utils.ErrResponse(401, w, "invalid email or password")
		return
	}

	jwtString, err := utils.CreateToken(user.Username)

	if err != nil {
		utils.ErrResponse(500, w, "Internal error occured")
		return
	}

	payload := utils.SignInPayload{
		User:        user,
		AccessToken: jwtString,
	}

	utils.JsonResponse(200, w, utils.HandleLoginToLogin(payload))
}

func generateUsername(name string) string {
	// Split the name into first and last name
	parts := strings.Fields(name)
	if len(parts) < 2 {
		// Handle cases where there's no space between the first and last name
		// You can implement your own logic here
		return name
	}
	firstName := strings.ToLower(parts[0])
	lastName := strings.ToLower(parts[1])

	// Generate the username by concatenating the first letter of the first name with the last name
	username := string(firstName[0]) + lastName

	// You might want to ensure the username is unique in your system
	// You can append a unique identifier if necessary
	// For example:
	// username += "_" + generateUniqueIdentifier()

	return username
}

func passwordEncryption(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		fmt.Println("Error generating hashped password", err)
	}

	return string(hashedPassword)
}

func (apiCfg ApiCfg) GetUser(w http.ResponseWriter, r *http.Request) {
	profileUsername := chi.URLParam(r, "username")
	user, err := apiCfg.DB.GetUser(r.Context(), profileUsername)

	if err != nil {
		utils.ErrResponse(404, w, fmt.Sprintf("Error occured: %v", err))
		return
	}

	utils.JsonResponse(200, w, utils.HandleUserToUser(user))
}

type ProfileImageBody struct {
	ProfileImage string `json:"profile_image"`
}

func (apiCfg ApiCfg) ChangeUserProfile(w http.ResponseWriter, r *http.Request, username string) {
	decorder := json.NewDecoder(r.Body)

	params := ProfileImageBody{}

	decorder.Decode(&params)

	var profileImage sql.NullString

	if params.ProfileImage != "" {
		profileImage = sql.NullString{String: params.ProfileImage, Valid: true}
	} else {
		profileImage = sql.NullString{Valid: false}
	}

	user, err := apiCfg.DB.ChangeProfilePicture(r.Context(), database.ChangeProfilePictureParams{
		Username:       username,
		ProfilePicture: profileImage,
	})

	if err != nil {
		utils.ErrResponse(401, w, fmt.Sprintf("An error occured: %v", err))
	}

	utils.JsonResponse(200, w, utils.HandleUserToUser(user))
}

type ProfileCoverBody struct {
	CoverImage string `json:"cover_image"`
}

func (apiCfg ApiCfg) ChangeUserCoverImage(w http.ResponseWriter, r *http.Request, username string) {
	decorder := json.NewDecoder(r.Body)

	params := ProfileCoverBody{}

	decorder.Decode(&params)

	var CoverImage sql.NullString

	if params.CoverImage != "" {
		CoverImage = sql.NullString{String: params.CoverImage, Valid: true}
	} else {
		CoverImage = sql.NullString{Valid: false}
	}

	user, err := apiCfg.DB.ChangeCoverPicture(r.Context(), database.ChangeCoverPictureParams{
		Username:     username,
		CoverPicture: CoverImage,
	})

	if err != nil {
		utils.ErrResponse(401, w, fmt.Sprintf("An error occured: %v", err))
	}

	utils.JsonResponse(200, w, utils.HandleUserToUser(user))
}

type UpdateUserParams struct {
	Name       string `json:"name"`
	Bio        string `json:"bio"`
	Profession string `json:"profession"`
}

func (apiCfg ApiCfg) UpdateUserProfile(w http.ResponseWriter, r *http.Request, username string) {
	decoder := json.NewDecoder(r.Body)

	params := UpdateUserParams{}

	decoder.Decode(&params)
	fmt.Println(username)
	user, err := apiCfg.DB.UpdateUserProfile(r.Context(), database.UpdateUserProfileParams{
		Username:   username,
		Name:       params.Name,
		Bio:        sql.NullString{Valid: true, String: params.Bio},
		Profession: sql.NullString{Valid: true, String: params.Profession},
	})

	if err != nil {
		utils.ErrResponse(400, w, fmt.Sprintf("An error occured %v ", err))
		return
	}

	utils.JsonResponse(200, w, utils.HandleUserToUser(user))
}

func (apiCfg ApiCfg) SuggestedUsers(w http.ResponseWriter, r *http.Request, username string) {
	users, err := apiCfg.DB.GetSuggestedUsers(r.Context(), username)

	if err != nil {
		utils.ErrResponse(200, w, fmt.Sprintf("error %v", err))
		return
	}

	utils.JsonResponse(200, w, utils.HandleUsersToUsers(users))
}
