package handlers

import (
	"github.com/Fruitfulfriends-REST-API-server/internal/clients/kafka/producer"
	messageform "github.com/Fruitfulfriends-REST-API-server/internal/rest/forms/message"
	"github.com/Fruitfulfriends-REST-API-server/pkg/rest/response"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Messenger struct {
	kafkaProducer *producer.Producer
	log           *logrus.Logger
}

func NewMessagerHandler(kafkaProducer *producer.Producer, logrus *logrus.Logger) *Messenger {
	return &Messenger{
		kafkaProducer: kafkaProducer,
		log:           logrus,
	}
}

func (h *Messenger) EnrichRoutes(router *gin.Engine) {
	messageRoutes := router.Group("/message")
	messageRoutes.POST("/send", h.sendMessage)
}

func (h *Messenger) sendMessage(c *gin.Context) {
	const op = "handlers.Messenger.sendMessage"
	log := h.log.WithField("operation", op)
	log.Info("send message")

	form, verr := messageform.NewMessageForm().ParseAndValidate(c)
	if verr != nil {
		response.HandleError(verr, c)
		return
	}

	err := h.kafkaProducer.SendMessage(form.(*messageform.MessageForm).Message)
	if err != nil {
		log.WithError(err).Errorf("%s: failed to send message", op)
		response.HandleError(response.ResolveError(err), c)
		return
	}
}
