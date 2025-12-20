package grpc

import (
	"github.com/DRSN-tech/go-backend/pkg/e"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GRPCErrorResponse(err error) error {
	switch {
	default:
		return status.Error(codes.Internal, e.ErrInternalServerError.Error())
	}
}
