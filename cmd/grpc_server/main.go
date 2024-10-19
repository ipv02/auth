package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/brianvoe/gofakeit"
	"github.com/ipv02/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const grpcPort = 50051

type server struct {
	user_v1.UnimplementedUserV1Server
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	user_v1.RegisterUserV1Server(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// CreateUser обрабатывает CreateUserRequest для создания нового пользователя.
//
// Логирует информацию о запросе (имя и email) и возвращает CreateUserResponse
// с сгенерированным идентификатором пользователя.
//
// Параметры:
//   - _ctx: Контекст для управления временем жизни запроса и дедлайнами.
//   - req: Запрос CreateUserRequest, содержащий данные пользователя (имя, email).
//
// Возвращает:
//   - *desc.CreateUserResponse: Ответ с идентификатором созданного пользователя.
//   - error: Возвращает ошибку в случае неудачи, или nil при успешном выполнении.
func (s *server) CreateUser(_ context.Context, req *user_v1.CreateUserRequest) (*user_v1.CreateUserResponse, error) {
	log.Printf("CreateRequest - Name: %s, Email: %s", req.GetName(), req.GetEmail())

	return &user_v1.CreateUserResponse{
		Id: 1,
	}, nil
}

// GetUser обрабатывает GetUserRequest для получения информации о пользователе по ID.
//
// Логирует идентификатор пользователя и возвращает GetUserResponse с данными,
// такими как имя, email, роль, а также датами создания и обновления.
//
// Параметры:
//   - _ctx: Контекст для управления временем жизни запроса и дедлайнами.
//   - req: Запрос GetUserRequest, содержащий идентификатор пользователя.
//
// Возвращает:
//   - *desc.GetUserResponse: Ответ с данными пользователя (ID, имя, email, роль, даты создания и обновления).
//   - error: Возвращает ошибку в случае неудачи, или nil при успешном выполнении.
func (s *server) GetUser(_ context.Context, req *user_v1.GetUserRequest) (*user_v1.GetUserResponse, error) {
	log.Printf("Note id: %d", req.GetId())

	return &user_v1.GetUserResponse{
		Id:        req.GetId(),
		Name:      gofakeit.Name(),
		Email:     gofakeit.Email(),
		Role:      user_v1.UserRole_USER,
		CreatedAt: timestamppb.New(gofakeit.Date()),
		UpdatedAt: timestamppb.New(gofakeit.Date()),
	}, nil
}

// UpdateUser обрабатывает UpdateUserRequest для обновления данных пользователя.
//
// Логирует информацию о запросе (ID, имя, email) и возвращает пустой ответ
// при успешном выполнении.
//
// Параметры:
//   - _ctx: Контекст для управления временем жизни запроса и дедлайнами.
//   - req: Запрос UpdateUserRequest, содержащий ID пользователя и обновленные данные (имя, email).
//
// Возвращает:
//   - *emptypb.Empty: Пустой ответ при успешном выполнении операции.
//   - error: Возвращает ошибку в случае неудачи, или nil при успешном выполнении.
func (s *server) UpdateUser(_ context.Context, req *user_v1.UpdateUserRequest) (*emptypb.Empty, error) {
	log.Printf("UpdateRequest - ID: %d, Name: %s, Email: %s", req.GetId(), req.GetName(), req.GetEmail())

	return &emptypb.Empty{}, nil
}

// DeleteUser обрабатывает DeleteUserRequest для удаления пользователя по ID.
//
// Логирует идентификатор удаляемого пользователя и возвращает пустой ответ
// при успешном выполнении.
//
// Параметры:
//   - _ctx: Контекст для управления временем жизни запроса и дедлайнами.
//   - req: Запрос DeleteUserRequest, содержащий ID пользователя для удаления.
//
// Возвращает:
//   - *emptypb.Empty: Пустой ответ при успешном выполнении операции.
//   - error: Возвращает ошибку в случае неудачи, или nil при успешном выполнении.
func (s *server) DeleteUser(_ context.Context, req *user_v1.DeleteUserRequest) (*emptypb.Empty, error) {
	log.Printf("Deleting object with ID: %d", req.GetId())

	return &emptypb.Empty{}, nil
}