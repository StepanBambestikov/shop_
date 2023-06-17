package middleware

import (
	"catalogServiceGit/internal/core/entities"
	"catalogServiceGit/internal/integrations"
	"catalogServiceGit/internal/log"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Authorization(ctx *gin.Context, userStorage integrations.UserStorage) {
	var authHeader = ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.Error(entities.NewError(http.StatusUnauthorized, "authorization header needed"))
		return
	}

	var parts = strings.Split(authHeader, " ")
	if len(parts) != 2 {
		ctx.Error(entities.NewError(http.StatusBadRequest, "wrong auth token format"))
		return
	}

	user, err := userStorage.GetFromToken(parts[1])
	if err != nil && errors.Is(err, &integrations.ErrUnauthorized{}) {
		ctx.Error(entities.NewError(http.StatusUnauthorized, "unauthorized: token is expired or revoked"))
		return
	} else if err != nil {
		log.Error("Auth error: ", err.Error())
		ctx.Error(entities.NewError(http.StatusInternalServerError, "internal error: can't authorize user"))
		return
	}

	// log.Info("User: ", user, err)

	ctx.Set("user", user)
}
