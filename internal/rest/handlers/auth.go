package handlers

import (
	grpccli "github.com/Fruitfulfriends-REST-API-server/internal/clients/grpc"
	"github.com/Fruitfulfriends-REST-API-server/internal/rest/forms/auth"
	"github.com/Fruitfulfriends-REST-API-server/pkg/rest/response"
	"github.com/gin-gonic/gin"
	ssov1 "github.com/markgregr/FruitfulFriends-protos/gen/go/sso"
	logrus "github.com/sirupsen/logrus"
	"net/http"
)

type Auth struct {
	log *logrus.Entry
	api *grpccli.Client
}

func New(api *grpccli.Client, log *logrus.Entry) *Auth {
	return &Auth{
		log: log.WithField("rest", "handlers"),
		api: api,
	}
}

func (h *Auth) EnrichRoutes(router *gin.Engine) {
	authRoutes := router.Group("/auth")
	authRoutes.POST("/sign-up", h.signupAction)
	//authRoutes.POST("/sign-in", h.signinAction)
}

func (h *Auth) signupAction(c *gin.Context) {
	const op = "handlers.Auth.signupAction"
	log := h.log.WithField("operation", op)
	log.Info("signing up user")

	form, verr := auth.NewRegisterForm().ParseAndValidate(c)
	if verr != nil {
		response.HandleError(verr, c)
		return
	}

	resp, err := h.api.AuthService.Register(c, &ssov1.RegisterRequest{
		Email:    form.(*auth.RegisterForm).Email,
		Password: form.(*auth.RegisterForm).Password,
	})
	if err != nil {
		log.WithError(err).Errorf("%s: failed to register user", op)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"userId": resp.UserId,
	})
}
