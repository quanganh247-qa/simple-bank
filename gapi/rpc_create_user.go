package gapi

import (
	"context"
	"time"

	"github.com/hibiken/asynq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	db "tutorial.sqlc.dev/app/db/sqlc"
	"tutorial.sqlc.dev/app/db/util"
	"tutorial.sqlc.dev/app/pb"
	"tutorial.sqlc.dev/app/val"
	"tutorial.sqlc.dev/app/worker"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	hashedPwd, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}
	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			Username:       req.GetUsername(),
			HashedPassword: hashedPwd,
			FullName:       req.GetFullName(),
			Email:          req.GetEmail(),
		},
		AfterCreate: func(user db.User) error {
			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}

			payload := &worker.PayloadSendVerifyEmail{
				Username: user.Username,
			}
			return server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, payload, opts...)
		},
	}
	txResult, err := server.store.CreateUserTx(ctx, arg)
	if err != nil {

		return nil, status.Errorf(codes.AlreadyExists, "username already exists: %s", err)

	}

	rsp := &pb.CreateUserResponse{
		User: convertUser(txResult.User),
	}

	return rsp, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (validations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		validations = append(validations, fieldValidation("username", err))
	}
	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		validations = append(validations, fieldValidation("password", err))
	}
	if err := val.ValidateFullname(req.GetFullName()); err != nil {
		validations = append(validations, fieldValidation("full_name", err))
	}
	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		validations = append(validations, fieldValidation("email", err))
	}
	return validations
}
