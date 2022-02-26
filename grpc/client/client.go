package main

import (
	"context"
	"log"
	"time"

	pb "github.com/akshaykhairmode/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	conn, err := grpc.Dial("127.0.0.1:9999", grpc.WithTransportCredentials(insecure.NewCredentials())) //We connect to our server.
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewCalculatorClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	input := []int32{1, 2, 3, 4, 5}

	methods := map[string]func(ctx context.Context, in *pb.Request, opts ...grpc.CallOption) (*pb.Reply, error){
		"Add": c.Add,
		"Sub": c.Sub,
		"Mul": c.Mul,
		"Div": c.Div,
	}

	for operation, function := range methods {
		reply, err := function(ctx, &pb.Request{Nums: input})
		if err != nil {
			log.Printf("Error : %v", err)
		}
		log.Printf("Operation : %v | Result : %v", operation, reply)
	}

}
