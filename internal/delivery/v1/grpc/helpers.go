package grpc

import (
	"errors"

	"github.com/DRSN-tech/go-backend/pkg/e"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GRPCErrorResponse(err error) error {
	switch {
	case errors.Is(err, e.ErrNoProducts):
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.Internal, e.ErrInternalServerError.Error())
	}
}
