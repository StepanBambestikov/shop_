package handlers

import (
	"net/http"

	"gitea.teneshag.ru/gigabit/goauth/internal/core/entities"
	"gitea.teneshag.ru/gigabit/goauth/internal/core/entities/api"
	"gitea.teneshag.ru/gigabit/goauth/internal/core/services"
	"gitea.teneshag.ru/gigabit/goauth/internal/log"

	"github.com/gin-gonic/gin"
)

// MeHandler godoc
// @Summary Returns user info
// @Schemes
// @Description Returns information about user
// @Tags User
// @Param Authorization header string true "Used to pass access token"
// @Accept json
// @Produce json
// @Success 200 {object} entities.ApiReply{data=entities.User}
// @Failure 400 {object} entities.ApiReply{error=entities.Error}
// @Failure 500 {object} entities.ApiReply{error=entities.Error}
// @Router /api/v1/users/me [get]
func MeHandler(ctx *gin.Context) {
	_user, ok := ctx.Get("user")
	if !ok {
		return
	}
	user, ok := _user.(*entities.User)
	if !ok {
		log.Error("Can't get user from context: ", _user)
		ctx.Error(entities.NewError(http.StatusInternalServerError, "internal error"))
		return
	}

	ctx.JSON(200, entities.ApiReply{
		Data:    user,
		Error:   nil,
		Message: "OK",
	})
}

// DeleteUserHandler godoc
// @Summary Deletes user
// @Schemes
// @Description Deletes user
// @Tags User
// @Param Authorization header string true "Used to pass access token"
// @Accept json
// @Produce json
// @Success 200 {object} entities.ApiReply
// @Failure 400 {object} entities.ApiReply{error=entities.Error}
// @Failure 500 {object} entities.ApiReply{error=entities.Error}
// @Router /api/v1/users/me [delete]
func DeleteUserHandler(ctx *gin.Context, svc services.UserService) {
	_user, ok := ctx.Get("user")
	if !ok {
		return
	}
	user, ok := _user.(*entities.User)
	if !ok {
		log.Error("Can't get user from context: ", _user)
		ctx.Error(entities.NewError(http.StatusInternalServerError, "internal error"))
		return
	}

	err := svc.DeleteUser(user)
	if err != nil {
		log.Error("Can't delete user: ", err.Error())
	}

	ctx.JSON(200, entities.ApiReply{
		Error:   nil,
		Message: "OK",
	})
}

// DeleteUserByIdHandler godoc
// @Summary Deletes user by id
// @Schemes
// @Description Deletes user by id
// @Tags User
// @Param userId path string true "User Id"
// @Produce json
// @Success 200 {object} entities.ApiReply
// @Failure 400 {object} entities.ApiReply{error=entities.Error}
// @Failure 500 {object} entities.ApiReply{error=entities.Error}
// @Router /internal/api/v1/users/{userId} [delete]
func DeleteUserByIdHandler(ctx *gin.Context, svc services.UserService) {
	var form struct {
		UserId string `uri:"userId" binding:"required" validate:"required,uuid4"`
	}
	if err := ctx.ShouldBindUri(&form); err != nil { // request code
		ctx.Error(entities.NewError(http.StatusBadRequest, "user id not provided"))
		return
	}
	err := svc.DeleteUserById(form.UserId)
	if err != nil {
		log.Error("Can't delete user: ", err.Error())
		ctx.Error(entities.NewError(http.StatusBadRequest, "user does not exist"))
		return
	}

	ctx.JSON(200, entities.ApiReply{
		Error:   nil,
		Message: "OK",
	})
}

// DeleteUserByIdHandler godoc
// @Summary Deletes user by id
// @Schemes
// @Description Gives user one of theese roles: notverified, verified, seller
// @Tags User
// @Param userId path string true "User Id"
// @Param request body api.GiveRoleRequest true "Login request"
// @Accept json
// @Produce json
// @Success 200 {object} entities.ApiReply
// @Failure 400 {object} entities.ApiReply{error=entities.Error}
// @Failure 500 {object} entities.ApiReply{error=entities.Error}
// @Router /internal/api/v1/users/{userId}/giveRole [post]
func GiveUserRoleHandler(ctx *gin.Context, svc services.UserService) {
	var request api.GiveRoleRequest
	var form struct {
		UserId string `uri:"userId" binding:"required" validate:"required,uuid4"`
	}

	if err := ctx.ShouldBindUri(&form); err != nil { // request code
		ctx.Error(entities.NewError(http.StatusBadRequest, "user id not provided"))
		return
	}
	if err := ctx.ShouldBindJSON(&request); err != nil { // request code
		ctx.Error(entities.NewError(http.StatusBadRequest, "bad request"))
		return
	}
	role, ok := entities.UserRoleFromString(request.Role)
	if !ok {
		ctx.Error(entities.NewError(http.StatusBadRequest, "unknown role"))
		return
	}
	err := svc.SetRoleById(form.UserId, role)
	if err != nil {
		log.Error("Can't set role: ", err.Error())
		ctx.Error(entities.NewError(http.StatusInternalServerError, "internal error"))
		return
	}

	ctx.JSON(200, entities.ApiReply{
		Error:   nil,
		Message: "OK",
	})
}
