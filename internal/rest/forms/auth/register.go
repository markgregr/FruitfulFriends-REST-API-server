package auth

import (
	"encoding/json"
	"github.com/Fruitfulfriends-REST-API-server/internal/rest/forms"
	"github.com/Fruitfulfriends-REST-API-server/pkg/rest/response"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterForm struct {
	Email    string
	Password string
}

func NewRegisterForm() *RegisterForm {
	return &RegisterForm{}
}

func (f *RegisterForm) ParseAndValidate(c *gin.Context) (forms.Former, response.Error) {
	const op = "forms.auth.RegisterForm.ParseAndValidate"

	body, err := io.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()

	if err != nil {
		log.WithError(err).Error("unable to read body")
		return nil, response.NewInternalError()
	}

	var request *RegisterRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		ve := response.NewValidationError()
		ve.SetError(response.GeneralErrorKey, response.InvalidRequestStructure, "invalid request structure")

		return nil, ve
	}

	errors := make(map[string]response.ErrorMessage)
	f.validateAndSetEmail(request, errors)
	f.validateAndSetPassword(request, errors)

	if len(errors) > 0 {
		return nil, response.NewValidationError(errors)
	}

	return f, nil
}

func (f *RegisterForm) ConvertToMap() map[string]interface{} {
	return map[string]interface{}{
		"email": f.Email,
	}
}

func (f *RegisterForm) validateAndSetEmail(request *RegisterRequest, errors map[string]response.ErrorMessage) {
	if request.Email == "" {
		errors["email"] = response.ErrorMessage{
			Code:    response.MissedValue,
			Message: "missed value",
		}
		return
	}

	f.Email = request.Email
}

func (f *RegisterForm) validateAndSetPassword(request *RegisterRequest, errors map[string]response.ErrorMessage) {
	if request.Password == "" {
		errors["password"] = response.ErrorMessage{
			Code:    response.MissedValue,
			Message: "missed value",
		}
		return
	}

	f.Password = request.Password
}
