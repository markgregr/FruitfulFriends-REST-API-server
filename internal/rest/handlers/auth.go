package handlers

import (
	"context"
	"fmt"
	grpccli "github.com/Fruitfulfriends-REST-API-server/internal/clients/grpc"
	authform "github.com/Fruitfulfriends-REST-API-server/internal/rest/forms/auth"
	"github.com/Fruitfulfriends-REST-API-server/internal/rest/models"
	"github.com/Fruitfulfriends-REST-API-server/pkg/rest/helper"
	"github.com/Fruitfulfriends-REST-API-server/pkg/rest/response"
	"github.com/gin-gonic/gin"
	ssov1 "github.com/markgregr/FruitfulFriends-protos/gen/go/sso"
	logrus "github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	"net/http"
)

type Auth struct {
	log   *logrus.Logger
	api   *grpccli.Client
	appID int32
}

func NewAuthHandler(api *grpccli.Client, log *logrus.Logger, appID int32) *Auth {
	return &Auth{
		log:   log,
		api:   api,
		appID: appID,
	}
}

func (h *Auth) EnrichRoutes(router *gin.Engine) {
	authRoutes := router.Group("/auth")
	authRoutes.POST("/register", h.registerAction)
	authRoutes.POST("/login", h.loginAction)
	authRoutes.POST("/logout", h.logoutAction)
}

func (h *Auth) registerAction(c *gin.Context) {
	const op = "handlers.Auth.registerAction"
	log := h.log.WithField("operation", op)
	log.Info("register user")

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
		response.HandleError(response.ResolveError(err), c)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"userId": resp.UserId,
	})
}

func (h *Auth) loginAction(c *gin.Context) {
	const op = "handlers.Auth.loginAction"
	log := h.log.WithField("operation", op)
	log.Info("login user")

	form, verr := authform.NewLoginForm().ParseAndValidate(c)
	if verr != nil {
		response.HandleError(verr, c)
		return
	}

	token, err := h.api.AuthService.Login(c, &ssov1.LoginRequest{
		Email:    form.(*authform.LoginForm).Email,
		Password: form.(*authform.LoginForm).Password,
		AppId:    h.appID,
	})
	if err != nil {
		log.WithError(err).Error("failed to login user")
		response.HandleError(response.ResolveError(err), c)
		return
	}

	c.JSON(http.StatusOK, models.AuthToken{
		AccessToken: token.GetToken(),
	})
}

func (h *Auth) logoutAction(c *gin.Context) {
	accessToken := helper.ExtractTokenFromHeaders(c)
	if accessToken == "" {
		c.Status(http.StatusUnauthorized)
		return
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "app_id", fmt.Sprintf("%d", h.appID))

	_, err := h.api.AuthService.Logout(metadata.AppendToOutgoingContext(ctx, "access_token", accessToken), &emptypb.Empty{})
	if err != nil {
		response.HandleError(response.ResolveError(err), c)
		return
	}

	c.Status(http.StatusNoContent)
}
