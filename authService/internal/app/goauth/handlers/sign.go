package handlers

import (
	"fmt"
	"gitea.teneshag.ru/gigabit/goauth/internal/core/entities"
	"gitea.teneshag.ru/gigabit/goauth/internal/core/entities/api"
	"gitea.teneshag.ru/gigabit/goauth/internal/core/services"
	"gitea.teneshag.ru/gigabit/goauth/internal/core/services/impl/signsvc"
	"gitea.teneshag.ru/gigabit/goauth/internal/integrations"
	"gitea.teneshag.ru/gigabit/goauth/internal/log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

// LoginHandler godoc
// @Summary Logins user
// @Schemes
// @Description Logins user
// @Tags Sign
// @Param request body api.LoginRequest true "Login request"
// @Accept json
// @Produce json
// @Success 200 {object} entities.ApiReply{data=api.LoginResponse}
// @Failure 400 {object} entities.ApiReply{error=entities.Error}
// @Failure 500 {object} entities.ApiReply{error=entities.Error}
// @Router /api/v1/sign/in [post]
func LoginHandler(ctx *gin.Context, svc services.UserService) {
	var loginRequest api.LoginRequest
	v := validator.New()

	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
		log.Info("Can't load body: ", err.Error())
		ctx.Error(entities.NewError(http.StatusBadRequest, fmt.Sprintf("invalid body: %s", err.Error())))
		return
	}

	if err := v.Struct(loginRequest); err != nil {
		ctx.Error(entities.NewError(http.StatusBadRequest, fmt.Sprintf("invalid body: %s", err.Error())))
		return
	}

	authStrategy, err := signsvc.NewStrategy(&loginRequest)
	if err != nil {
		ctx.Error(entities.NewError(http.StatusBadRequest, "invalid body: refresh token or login-password pair required"))
		return
	}
	authorizer := signsvc.NewAuthorizer(loginRequest.Login, loginRequest.Password, loginRequest.RefreshToken, authStrategy, svc)
	response, outError, err := authorizer.Authorize()
	if err != nil {
		log.Error("Authorization failed:", err.Error())
		ctx.Error(outError)
		return
	}
	ctx.JSON(200, entities.ApiReply{
		Data:    response,
		Error:   nil,
		Message: "OK",
	})
}

// SignupHandler godoc
// @Summary Signs up user
// @Schemes
// @Description Signs up user
// @Tags Sign
// @Param request body api.SignupRequest true "User sign up request"
// @Accept json
// @Produce json
// @Success 200 {object} entities.ApiReply
// @Failure 400 {object} entities.ApiReply{error=entities.Error}
// @Failure 409 {object} entities.ApiReply{error=entities.Error}
// @Failure 500 {object} entities.ApiReply{error=entities.Error}
// @Router /api/v1/sign/up [post]
func SignupHandler(ctx *gin.Context, svc services.UserService) {
	var signupRequest api.SignupRequest

	v := validator.New()

	if err := ctx.ShouldBindJSON(&signupRequest); err != nil {
		log.Warn("Can't bind json: ", err.Error())
		ctx.Error(entities.NewError(http.StatusBadRequest, fmt.Sprintf("invalid body: %s", err.Error())))
		return
	}

	log.Debug("Sign up request ", signupRequest)

	if err := v.Struct(signupRequest); err != nil {
		ctx.Error(entities.NewError(http.StatusBadRequest, fmt.Sprintf("invalid body: %s", err.Error())))
		return
	}

	_, err := svc.Signup(signupRequest)
	if err != nil && errors.Is(err, &integrations.ErrAlreadyExists{}) {
		log.Error("Error signing up: ", err.Error())
		ctx.Error(entities.NewError(http.StatusConflict, "user already exists"))
		return
	} else if err != nil {
		log.Error("Error signing up: ", err.Error())
		ctx.Error(entities.NewError(http.StatusInternalServerError, "internal error"))
		return
	}
	log.Debug("Successful sign up")
	ctx.JSON(200, entities.ApiReply{
		Error:   nil,
		Message: "OK",
	})
}

// LogoutHandler godoc
// @Summary Logouts user
// @Schemes
// @Description Logouts user
// @Tags Sign
// @Param request body api.LogoutRequest true "Used to pass refresh token"
// @Accept json
// @Produce json
// @Success 200 {object} entities.ApiReply
// @Failure 400 {object} entities.ApiReply{error=entities.Error}
// @Failure 500 {object} entities.ApiReply{error=entities.Error}
// @Router /api/v1/sign/out [post]
func LogoutHandler(ctx *gin.Context, svc services.UserService) {
	var err error
	var logoutRequest api.LogoutRequest

	v := validator.New()

	if err := ctx.ShouldBindJSON(&logoutRequest); err != nil {
		log.Info("Can't load body: ", err.Error())
		ctx.Error(entities.NewError(http.StatusBadRequest, fmt.Sprintf("invalid body: %s", err.Error())))
		return
	}

	if err := v.Struct(logoutRequest); err != nil {
		ctx.Error(entities.NewError(http.StatusBadRequest, fmt.Sprintf("invalid body: %s", err.Error())))
		return
	}

	err = svc.Logout(logoutRequest.RefreshToken)
	if err != nil {
		log.Error("Internal error during logout: ", err.Error())
		ctx.Error(entities.NewError(http.StatusInternalServerError, "internal error"))
		return
	}

	ctx.JSON(200, entities.ApiReply{
		Data:    nil,
		Error:   nil,
		Message: "OK",
	})
}
