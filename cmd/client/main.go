package main

import (
	"context"
	"log"
	"time"

	"github.com/ebobo/learn_buf_grpc/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	user_server_address = "localhost:9092"
)

func main() {
	conn, err := grpc.Dial(user_server_address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect %v", err)
	}
	defer conn.Close()
	c := proto.NewUserManagementClient(conn)

	// https://go.dev/blog/context
	// In Go servers, each incoming request is handled in its own goroutine.
	// Request handlers often start additional goroutines to access backends such as databases
	// and RPC services. The set of goroutines working on a request typically needs access to
	// request-specific values such as the identity of the end user, authorization tokens,
	// and the request’s deadline. When a request is canceled or times out, all the goroutines
	// working on that request should exit quickly so the system can reclaim any resources they are using.
	// At Google, we developed a context package that makes it easy to pass request-scoped values,
	// cancelation signals, and deadlines across API boundaries to all the goroutines involved
	// in handling a request. The package is publicly available as context.

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	var new_users = make(map[string]int32)

	new_users["Qi"] = 39
	new_users["Ellen"] = 40

	for name, age := range new_users {
		r, err := c.CreateUser(ctx, &proto.NewUser{Name: name, Age: age})
		if err != nil {
			log.Fatalf("could not create user %v", err)
		}
		log.Printf(`User Details: Name: %s Age: %d Id: %d`, r.GetName(), r.GetAge(), r.GetId())
	}

	params := &proto.GetUsersParams{}
	r, err := c.GetUser(ctx, params)
	if err != nil {
		log.Fatalf("could not get user list %v", err)
	}
	log.Print("\n User List: \n")
	log.Printf(" %v\n", r.GetUsers())

}

// use go mod tidy to download all the pakages we imported
