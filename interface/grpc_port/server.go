package grpcport

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/google/uuid"
	"github.com/guilhermealvess/guicpay/domain/usecase"
	"github.com/guilhermealvess/guicpay/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type app struct {
	accountService     *accountServer
	authService        *authService
	transactionService *transactionService
}

func NewApp(u usecase.AccountUseCase) *app {
	return &app{
		accountService:     &accountServer{usecase: u},
		authService:        &authService{usecase: u},
		transactionService: &transactionService{usecase: u},
	}
}

func (a *app) Start(port int) {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			RequestIDInterceptor,
		),
	)

	pb.RegisterAccountsServer(s, a.accountService)
	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func RequestIDInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	requestID := uuid.NewString()
	ctx = context.WithValue(ctx, "requestID", requestID)
	grpc.SetHeader(ctx, metadata.Pairs("request-id", requestID))
	log.Printf("Nova solicitação - ID: %s", requestID)
	resp, err := handler(ctx, req)
	return resp, err
}
