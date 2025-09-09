package handler

import (
	"context"
	"errors"
	"sample-crud/internal/service"
	customerrors "sample-crud/pkg/errors"
	"sample-crud/proto/pb/product"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductGRPCHandler struct {
	productService service.ProductService
	product.UnimplementedProductServiceServer
}

func (p ProductGRPCHandler) GetProduct(ctx context.Context, request *product.GetProductRequest) (*product.GetProductResponse, error) {
	productInfo, err := p.productService.FindByID(ctx, uint(request.GetId()))
	if err != nil {
		var customErr *customerrors.CustomError
		switch {
		case errors.As(err, &customErr):
			return nil, status.Error(customErr.GrpcCode, customErr.Message)
		default:
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}
	resp := &product.GetProductResponse{
		Id:   uint64(productInfo.ID),
		Name: productInfo.Name,
	}
	return resp, nil
}

func NewProductGRPCHandler(productService service.ProductService) *ProductGRPCHandler {
	return &ProductGRPCHandler{productService: productService}
}
