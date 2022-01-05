package main

import (
	"context"
	"log"
	"math/rand"
	"net"

	"github.com/ebobo/learn_buf_grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":9092"
)

type UserManagementServer struct {
	user_list *proto.UserList
	proto.UnimplementedUserManagementServer
}

func NewUserManagementServer() *UserManagementServer {
	return &UserManagementServer{user_list: &proto.UserList{}}
}

func (server *UserManagementServer) Run() error {
	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("unable to listen %v", err)
	}
	gs := grpc.NewServer()
	reflection.Register(gs)

	proto.RegisterUserManagementServer(gs, server)

	log.Printf("server listening at %v", lis.Addr())
	return gs.Serve(lis)
}

func (s *UserManagementServer) CreateUser(ctx context.Context, in *proto.NewUser) (*proto.User, error) {
	log.Printf("Handle CreateUser %v", in.GetName())
	var user_id int32 = int32(rand.Intn(1000))

	created_user := &proto.User{Id: user_id, Name: in.GetName(), Age: in.GetAge()}
	s.user_list.Users = append(s.user_list.Users, created_user)

	return created_user, nil
}

func (s *UserManagementServer) GetUser(ctx context.Context, in *proto.GetUsersParams) (*proto.UserList, error) {
	return s.user_list, nil
}

func main() {
	var user_management_server *UserManagementServer = NewUserManagementServer()

	if err := user_management_server.Run(); err != nil {
		log.Fatalf("failed to serve %v", err)
	}
}
