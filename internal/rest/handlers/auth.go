package handlers

import (
	"errors"
	grpccli "github.com/Fruitfulfriends-REST-API-server/internal/clients/grpc"
	authform "github.com/Fruitfulfriends-REST-API-server/internal/rest/forms/auth"
	"github.com/Fruitfulfriends-REST-API-server/internal/rest/models"
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

func NewAuthHandler(api *grpccli.Client, log *logrus.Entry) *Auth {
	return &Auth{
		log: log.WithField("rest", "handlers"),
		api: api,
	}
}

func (h *Auth) EnrichRoutes(router *gin.Engine) {
	authRoutes := router.Group("/auth")
	authRoutes.POST("/regiser", h.registerAction)
	authRoutes.POST("/login", h.loginAction)
}

func (h *Auth) registerAction(c *gin.Context) {
	const op = "handlers.Auth.signupAction"
	log := h.log.WithField("operation", op)
	log.Info("signing up user")

	form, verr := authform.NewRegisterForm().ParseAndValidate(c)
	if verr != nil {
		response.HandleError(verr, c)
		return
	}

	resp, err := h.api.AuthService.Register(c, &ssov1.RegisterRequest{
		Email:    form.(*authform.RegisterForm).Email,
		Password: form.(*authform.RegisterForm).Password,
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

func (h *Auth) loginAction(c *gin.Context) {
	form, verr := authform.NewLoginForm().ParseAndValidate(c)
	if verr != nil {
		response.HandleError(verr, c)
		return
	}

	token, err := h.api.AuthService.Login(c, &ssov1.LoginRequest{
		Email:    form.(*authform.LoginForm).Email,
		Password: form.(*authform.LoginForm).Password,
	})
	if err != nil {
		response.HandleError(response.ResolveError(errors.New("invalid username or password")), c)
		return
	}

	c.JSON(http.StatusOK, models.AuthToken{
		AccessToken: token.GetToken(),
	})
}
