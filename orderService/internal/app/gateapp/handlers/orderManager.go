package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"orderServiceGit/internal/core/entities"
	"orderServiceGit/internal/core/entities/api"
	"orderServiceGit/internal/core/services"
	"orderServiceGit/internal/core/services/DTO"
	"orderServiceGit/internal/log"
)

func bindAndValidateJSON[T any](request *T, ctx *gin.Context) (err error) {
	v := validator.New()
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Info("Can't load body: ", err.Error())
		err = ctx.Error(entities.NewError(http.StatusBadRequest, fmt.Sprintf("invalid body: %s", err.Error())))
		return err
	}

	if err := v.Struct(request); err != nil {
		err = ctx.Error(entities.NewError(http.StatusBadRequest, fmt.Sprintf("invalid body: %s", err.Error())))
		return err
	}
	return nil
}

//// CreateOrder godoc
//// @Summary Creates new order
//// @Schemes
//// @Description Creates new order
//// @Tags Order
//// @Param request body api.CreateOrderRequest true "Create order request"
//// @Accept json
//// @Produce json
//// @Success 200 {object} entities.ApiReply
//// @Failure 400 {object} entities.ApiReply{error=entities.Error}
//// @Failure 500 {object} entities.ApiReply{error=entities.Error}
//// @Router /catalog/v1/orders [post]

func CreateOrderHandler(ctx *gin.Context, svc services.GateService) {
	var request api.CreateOrderRequest
	err := bindAndValidateJSON(&request, ctx)
	err = svc.CreateOrder(request.OrderDTO)
	if err != nil {
		ctx.Error(entities.NewError(http.StatusInternalServerError, "internal error"))
		return
	}

	ctx.JSON(200, entities.ApiReply{
		Error:   nil,
		Message: "OK",
	})
}

// SetOrderStatus godoc
// @Summary Set order status
// @Schemes
// @Description Set order status
// @Tags Order
// @Param label path string true "label"
// @Produce json
// @Success 200 {object} entities.ApiReply
// @Failure 400 {object} entities.ApiReply{error=entities.Error}
// @Failure 500 {object} entities.ApiReply{error=entities.Error}
// @Router /order/v1/internal/orders/setStatus/{id} [post]
func SetOrderStatusHandler(ctx *gin.Context, svc services.GateService) {
	var request api.SetOrderStatusRequest
	err := bindAndValidateJSON(&request, ctx)
	err = svc.SetOrderStatus(DTO.OrderDTO{ID: request.OrderID})
	if err != nil {
		ctx.Error(entities.NewError(http.StatusInternalServerError, "internal error"))
		return
	}

	ctx.JSON(200, entities.ApiReply{
		Error:   nil,
		Message: "OK",
	})
}

// DeleteOrder godoc
// @Summary Deletes order
// @Schemes
// @Description Deletes order
// @Tags Product
// @Param request body api.CreateOrderRequest true "Create order request"
// @Accept json
// @Produce json
// @Success 200 {object} entities.ApiReply
// @Failure 400 {object} entities.ApiReply{error=entities.Error}
// @Failure 500 {object} entities.ApiReply{error=entities.Error}
// @Router /order/v1/orders/{id} [delete]
func DeleteOrderHandler(ctx *gin.Context, svc services.GateService) {
	var request api.DeleteOrderRequest
	err := bindAndValidateJSON(&request, ctx)
	err = svc.DeleteOrder(DTO.OrderDTO{ID: request.OrderID})
	if err != nil {
		ctx.Error(entities.NewError(http.StatusInternalServerError, "internal error"))
		return
	}

	ctx.JSON(200, entities.ApiReply{
		Error:   nil,
		Message: "OK",
	})
}

// GetOrderInfo godoc
// @Summary get order info
// @Schemes
// @Description get order info
// @Tags Order
// @Produce json
// @Success 200 {object} entities.ApiReply
// @Failure 500 {object} entities.ApiReply{error=entities.Error}
// @Router  /order/v1/orders/{id} [get]
func GetOrderInfoHandler(ctx *gin.Context, svc services.GateService) {
	var request api.GetOrderInfoRequest
	v := validator.New()
	if err := ctx.ShouldBindUri(&request); err != nil { //TODO gin.BindSkip
		log.Info("Can't load body: ", err.Error())
		err = ctx.Error(entities.NewError(http.StatusBadRequest, fmt.Sprintf("invalid body: %s", err.Error())))
		return
	}

	if err := v.Struct(request); err != nil {
		err = ctx.Error(entities.NewError(http.StatusBadRequest, fmt.Sprintf("invalid body: %s", err.Error())))
		return
	}
	orders, err := svc.GetOrderByID(DTO.OrderDTO{ID: request.OrderID})
	if err != nil {
		ctx.Error(entities.NewError(http.StatusInternalServerError, "internal error"))
		return
	}

	ctx.JSON(200, entities.ApiReply{
		Data:    orders,
		Error:   nil,
		Message: "OK",
	})
}

// GetUserOrders godoc
// @Summary GetUserOrders
// @Schemes
// @Description GetUserOrders
// @Tags Order
// @Produce json
// @Success 200 {object} entities.ApiReply
// @Failure 500 {object} entities.ApiReply{error=entities.Error}
// @Router /order/v1/orders [get]
func GetUserOrdersHandler(ctx *gin.Context, svc services.GateService) {
	var request api.GetUserOrdersRequest
	v := validator.New()
	if err := ctx.ShouldBindUri(&request); err != nil { //TODO gin.BindSkip
		log.Info("Can't load body: ", err.Error())
		err = ctx.Error(entities.NewError(http.StatusBadRequest, fmt.Sprintf("invalid body: %s", err.Error())))
		return
	}

	if err := v.Struct(request); err != nil {
		err = ctx.Error(entities.NewError(http.StatusBadRequest, fmt.Sprintf("invalid body: %s", err.Error())))
		return
	}
	orders, err := svc.GetUserOrders(&request.SelectionDTO)
	if err != nil {
		ctx.Error(entities.NewError(http.StatusInternalServerError, "internal error"))
		return
	}

	ctx.JSON(200, entities.ApiReply{
		Data:    orders,
		Error:   nil,
		Message: "OK",
	})
}
