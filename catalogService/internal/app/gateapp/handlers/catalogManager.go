package handlers

import (
	"catalogServiceGit/internal/core/entities"
	"catalogServiceGit/internal/core/entities/api"
	"catalogServiceGit/internal/core/services"
	"catalogServiceGit/internal/core/services/DTO"
	"catalogServiceGit/internal/log"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
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

// CreateProduct godoc
// @Summary Creates new product
// @Schemes
// @Description Creates new product
// @Tags Product
// @Param Authorization header string true "Used to pass access token"
// @Param request body api.CreateProductRequest true "Create product request"
// @Accept json
// @Produce json
// @Success 200 {object} entities.ApiReply
// @Failure 400 {object} entities.ApiReply{error=entities.Error}
// @Failure 500 {object} entities.ApiReply{error=entities.Error}
// @Router /catalog/v1/products [post]
func CreateProductHandler(ctx *gin.Context, svc services.GateService) {
	var request api.CreateProductRequest
	err := bindAndValidateJSON(&request, ctx)
	err = svc.CreateProduct(request.ProductDTO)
	if err != nil {
		log.Error("Internal error: ", err.Error())
		ctx.Error(entities.NewError(http.StatusInternalServerError, "internal error"))
		return
	}
	log.Debug("Successful creating product")
	ctx.JSON(200, entities.ApiReply{
		Error:   nil,
		Message: "OK",
	})
}

// ChangeProductHandler godoc
// @Summary Change product information
// @Schemes
// @Description Change product information
// @Tags Product
// @Param Authorization header string true "Used to pass access token"
// @Param request body api.ChangeProductRequest true "Create product request"
// @Produce json
// @Success 200 {object} entities.ApiReply
// @Failure 400 {object} entities.ApiReply{error=entities.Error}
// @Failure 500 {object} entities.ApiReply{error=entities.Error}
// @Router /catalog/v1/products/{id} [post]
func ChangeProductHandler(ctx *gin.Context, svc services.GateService) {
	var request api.ChangeProductRequest
	err := bindAndValidateJSON(&request, ctx)
	err = svc.ChangeProduct(request.ProductDTO)
	if err != nil {
		log.Error("Internal error: ", err.Error())
		ctx.Error(entities.NewError(http.StatusInternalServerError, "internal error"))
		return
	}
	log.Debug("Successful changing products")
	ctx.JSON(200, entities.ApiReply{
		Error:   nil,
		Message: "OK",
	})
}

// DeleteProduct godoc
// @Summary Deletes product
// @Schemes
// @Description Deletes product
// @Tags Product
// @Param Authorization header string true "Used to pass access token"
// @Param request body api.DeleteProductRequest true "Delete product request"
// @Accept json
// @Produce json
// @Success 200 {object} entities.ApiReply
// @Failure 400 {object} entities.ApiReply{error=entities.Error}
// @Failure 500 {object} entities.ApiReply{error=entities.Error}
// @Router /catalog/v1/products/{id} [delete]
func DeleteProductHandler(ctx *gin.Context, svc services.GateService) {
	var request api.DeleteProductRequest
	err := bindAndValidateJSON(&request, ctx)
	err = svc.DeleteProduct(DTO.ProductDTO{ID: request.ProductID})
	if err != nil {
		log.Error("internal error: ", err.Error())
		ctx.Error(entities.NewError(http.StatusInternalServerError, "internal error"))
		return
	}
	log.Debug("Successful deleteng product")
	ctx.JSON(200, entities.ApiReply{
		Error:   nil,
		Message: "OK",
	})
}

// OrderProduct godoc
// @Summary Order products
// @Schemes
// @Description Order products
// @Param Authorization header string true "Used to pass access token"
// @Param request body api.OrderProductRequest true "Order product request"
// @Tags Product
// @Produce json
// @Success 200 {object} entities.ApiReply
// @Failure 500 {object} entities.ApiReply{error=entities.Error}
// @Router /catalog/v1/products/{id}/order [post]
func OrderProductHandler(ctx *gin.Context, svc services.GateService) {
	var request api.OrderProductRequest
	err := bindAndValidateJSON(&request, ctx)
	err = svc.OrderProduct(DTO.ProductDTO{ID: request.ProductID})
	if err != nil {
		log.Error("Internal error: ", err.Error())
		ctx.Error(entities.NewError(http.StatusInternalServerError, "internal error"))
		return
	}
	log.Debug("Successful ordered product")
	ctx.JSON(200, entities.ApiReply{
		Error:   nil,
		Message: "OK",
	})
}

// RateProduct godoc
// @Summary Rates some product
// @Schemes
// @Description Rates some product
// @Param Authorization header string true "Used to pass access token"
// @Param request body api.RateProductRequest true "Order product request"
// @Tags Product
// @Produce json
// @Success 200 {object} entities.ApiReply
// @Failure 500 {object} entities.ApiReply{error=entities.Error}
// @Router /catalog/v1/{id}/rate [post]
func RateProductHandler(ctx *gin.Context, svc services.GateService) {
	var request api.RateProductRequest
	err := bindAndValidateJSON(&request, ctx)
	err = svc.RateProduct(DTO.ProductDTO{ID: request.ProductID, Ratting: request.NewRatting})
	if err != nil {
		log.Error("internal error: ", err.Error())
		log.Error("internal error: ", err.Error())
		ctx.Error(entities.NewError(http.StatusInternalServerError, "internal error"))
		return
	}
	ctx.JSON(200, entities.ApiReply{
		Error:   nil,
		Message: "OK",
	})
}

// GetSeveralProducts godoc
// @Summary Get Several Products
// @Schemes
// @Description Get Several Products
// @Param Authorization header string true "Used to pass access token"
// @Param request body api.GetAllProductsRequest true "Get some products"
// @Tags Product
// @Produce json
// @Success 200 {object} entities.ApiReply
// @Failure 500 {object} entities.ApiReply{error=entities.Error}
// @Router /catalog/v1/products [get]
func GetSeveralProductsHandler(ctx *gin.Context, svc services.GateService) {
	var request api.GetAllProductsRequest
	v := validator.New()
	if err := ctx.ShouldBindUri(&request); err != nil {
		log.Info("Can't load body: ", err.Error())
		err = ctx.Error(entities.NewError(http.StatusBadRequest, fmt.Sprintf("invalid body: %s", err.Error())))
		return
	}

	if err := v.Struct(request); err != nil {
		log.Error("Internal handler error: ", err.Error())
		err = ctx.Error(entities.NewError(http.StatusBadRequest, fmt.Sprintf("invalid body: %s", err.Error())))
		return
	}
	products, err := svc.GetSeveralProducts(request.SelectionDTO)
	if err != nil {
		log.Error("Internal handler error: ", err.Error())
		ctx.Error(entities.NewError(http.StatusInternalServerError, "internal error"))
		return
	}
	log.Debug("Successful getting several products")
	ctx.JSON(200, entities.ApiReply{
		Data:    products,
		Error:   nil,
		Message: "OK",
	})
}
