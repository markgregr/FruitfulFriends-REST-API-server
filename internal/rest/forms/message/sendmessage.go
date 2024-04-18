package message

import (
	"github.com/Fruitfulfriends-REST-API-server/internal/rest/forms"
	"github.com/Fruitfulfriends-REST-API-server/pkg/rest/response"
	"github.com/gin-gonic/gin"
)

type MessageRequest struct {
	Message string `json:"message"`
}

type MessageForm struct {
	Message string
}

func NewMessageForm() *MessageForm {
	return &MessageForm{}
}

func (m *MessageForm) ParseAndValidate(c *gin.Context) (forms.Former, response.Error) {
	var request *MessageRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		ve := response.NewValidationError()
		ve.SetError(response.GeneralErrorKey, response.InvalidRequestStructure, "invalid request structure")

		return nil, ve
	}

	errors := make(map[string]response.ErrorMessage)
	m.validateAndSetMessage(request, errors)

	if len(errors) > 0 {
		return nil, response.NewValidationError(errors)
	}

	return m, nil
}

func (m *MessageForm) ConvertToMap() map[string]interface{} {
	return map[string]interface{}{
		"message": m.Message,
	}
}

func (m *MessageForm) validateAndSetMessage(request *MessageRequest, errors map[string]response.ErrorMessage) {
	if request.Message == "" {
		errors["message"] = response.ErrorMessage{
			Code:    response.EmptyField,
			Message: "message is empty",
		}
	} else {
		m.Message = request.Message
	}
}
