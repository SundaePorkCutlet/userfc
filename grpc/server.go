package grpc

import (
	"context"
	"fmt"
	"net"
	"userfc/cmd/user/service"
	"userfc/infrastructure/log"
	"userfc/pb"
	"userfc/utils"

	"google.golang.org/grpc"
)

type UserGRPCServer struct {
	pb.UnimplementedUserServiceServer
	userService *service.UserService
}

func NewUserGRPCServer(userService *service.UserService) *UserGRPCServer {
	return &UserGRPCServer{
		userService: userService,
	}
}

func (s *UserGRPCServer) GetUserInfoByUserId(ctx context.Context, req *pb.GetUserInfoByUserIdRequest) (*pb.GetUserInfoByUserIdResponse, error) {
	user, err := s.userService.GetUserByUserId(ctx, req.UserId)
	if err != nil {
		log.Logger.Error().Err(err).Int64("user_id", req.UserId).Msg("Failed to get user by ID")
		return nil, err
	}

	return &pb.GetUserInfoByUserIdResponse{
		Id:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		Role:  user.Role,
	}, nil
}

func (s *UserGRPCServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	claims, err := utils.ValidateToken(req.Token)
	if err != nil {
		return &pb.ValidateTokenResponse{
			Valid: false,
			Error: err.Error(),
		}, nil
	}

	return &pb.ValidateTokenResponse{
		Valid:  true,
		UserId: claims.UserID,
		Email:  claims.Email,
	}, nil
}

func StartGRPCServer(port string, userService *service.UserService) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Logger.Fatal().Err(err).Msg("Failed to listen for gRPC")
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, NewUserGRPCServer(userService))

	log.Logger.Info().Str("port", port).Msg("gRPC server is running")

	if err := grpcServer.Serve(lis); err != nil {
		log.Logger.Fatal().Err(err).Msg("Failed to serve gRPC")
	}
}
