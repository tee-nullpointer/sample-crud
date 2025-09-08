package handler

import (
	"errors"
	"io"
	"net/http"
	"sample-crud/internal/domain"
	customerrors "sample-crud/pkg/errors"
	"sample-crud/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tee-nullpointer/go-common-kit/pkg/logger"
	"go.uber.org/zap"
)

type ProductHandler struct {
	productService domain.ProductService
}

func NewProductHandler(productService domain.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

func (h *ProductHandler) Create(c *gin.Context) {
	var request domain.ProductCreation
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Warn("Invalid product creation request", zap.Error(err))
		if errors.Is(err, io.EOF) || err.Error() == "EOF" {
			c.Error(customerrors.NewBadRequestError("Request body empty", err.Error()))
			return
		}
		c.Error(customerrors.NewBadRequestError(err.Error(), err.Error()))
		return
	}
	id, err := h.productService.CreateProduct(c.Request.Context(), request.Name)
	if err != nil {
		c.Error(err)
		return
	}
	logger.SInfo("Product created with id %v", id)
	c.JSON(http.StatusCreated, response.Created(id))
}

func (h *ProductHandler) FindByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Warn("Invalid product id", zap.String("id", c.Param("id")), zap.Error(err))
		c.Error(customerrors.NewBadRequestError("Invalid product id", err.Error()))
		return
	}
	productInfo, err2 := h.productService.FindByID(c.Request.Context(), uint(id))
	if err2 != nil {
		c.Error(err2)
		return
	}
	logger.SInfo("Product found with data : %+v", productInfo)
	c.JSON(http.StatusOK, response.Success(productInfo))
}

func (h *ProductHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Warn("Invalid product id", zap.String("id", c.Param("id")), zap.Error(err))
		c.Error(customerrors.NewBadRequestError("Invalid product id", err.Error()))
		return
	}

	var request domain.ProductUpdate
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Warn("Invalid product update request", zap.Error(err))
		if errors.Is(err, io.EOF) || err.Error() == "EOF" {
			c.Error(customerrors.NewBadRequestError("Request body empty", err.Error()))
			return
		}
		c.Error(customerrors.NewBadRequestError(err.Error(), err.Error()))
		return
	}

	err = h.productService.UpdateProduct(c.Request.Context(), uint(id), request.Name)
	if err != nil {
		c.Error(err)
		return
	}

	logger.SInfo("Product updated with id %v", id)
	c.JSON(http.StatusOK, response.Success(nil))
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Warn("Invalid product id", zap.String("id", c.Param("id")), zap.Error(err))
		c.Error(customerrors.NewBadRequestError("Invalid product id", err.Error()))
		return
	}
	err = h.productService.DeleteProduct(c.Request.Context(), uint(id))
	if err != nil {
		c.Error(err)
		return
	}
	logger.SInfo("Product deleted with id %v", id)
	c.JSON(http.StatusOK, response.Success(nil))
}
