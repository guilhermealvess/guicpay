package grpcport

import (
	"context"
	"database/sql"
	"errors"

	"github.com/guilhermealvess/guicpay/domain/usecase"
	"github.com/guilhermealvess/guicpay/pkg/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type accountServer struct {
	pb.AccountsServer
	usecase usecase.AccountUseCase
}

func (s *accountServer) Create(ctx context.Context, input *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	output, err := s.usecase.ExecuteNewAccount(ctx, usecase.NewAccountInput{
		Name:           input.CustomerName,
		Email:          input.Email,
		Password:       input.Password,
		Type:           input.AccountType,
		DocumentNumber: input.DocumentNumber,
		PhoneNumber:    input.PhoneNumber,
	})

	if err != nil {
		return nil, buildStatusError(err)
	}

	return &pb.CreateAccountResponse{Id: output.String()}, nil
}

func (s *accountServer) Fetch(ctx context.Context, input *pb.FetchAccountRequest) (*pb.FetchAccountResponse, error) {
	accountID, ok := getAccountContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.FailedPrecondition, "TODO: ")
	}

	output, err := s.usecase.FindByID(ctx, accountID)
	if err != nil {
		return nil, buildStatusError(err)
	}

	return &pb.FetchAccountResponse{
		Id:           output.ID.String(),
		AccountType:  output.AccountType,
		CustomerName: output.CustomerName,
		Email:        output.Email,
		Balance:      output.Balance,
		Status:       output.Status,
	}, nil
}

func (s *accountServer) List(ctx context.Context, input *pb.ListRequest) (*pb.ListResponse, error) {
	output, err := s.usecase.FindAll(ctx)
	if err != nil {
		return nil, buildStatusError(err)
	}

	accounts := make([]*pb.FetchAccountResponse, 0)
	for _, it := range output {
		account := pb.FetchAccountResponse{
			Id:           it.ID.String(),
			AccountType:  it.AccountType,
			CustomerName: it.CustomerName,
			Email:        it.Email,
			Balance:      it.Balance,
			Status:       it.Status,
		}
		accounts = append(accounts, &account)
	}

	return &pb.ListResponse{Accounts: accounts}, nil
}

type authService struct {
	pb.AuthServer
	usecase usecase.AccountUseCase
}

func (s *authService) Auth(ctx context.Context, input *pb.AuthRequest) (*pb.AuthResponse, error) {
	return &pb.AuthResponse{}, nil
}

type transactionService struct {
	pb.TransactionsServer
	usecase usecase.AccountUseCase
}

func (s *transactionService) Deposit(ctx context.Context, input *pb.DepositRequest) (*pb.TransactionResponse, error) {
	return &pb.TransactionResponse{}, nil
}

func (s *transactionService) Transfer(ctx context.Context, input *pb.TransferRequest) (*pb.TransactionResponse, error) {
	return &pb.TransactionResponse{}, nil
}

func buildStatusError(err error) error {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return status.Errorf(codes.NotFound, err.Error())

	default:
		return status.Errorf(codes.Internal, err.Error())
	}
}
