package main

import (
	"fmt"
	"io"
	"log"
	"net"

	simonsayspb "github.com/simonsays/simonsayspb"
	"google.golang.org/grpc"
)

type server struct{}

func (s *server) Game(stream simonsayspb.SimonSays_GameServer) error {
	fmt.Printf("Game function was invoked with a streaming request\n")
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}
		var sendErr error
		if join := req.GetJoin(); join != nil {
			sendErr = stream.Send(&simonsayspb.Response{
				Event: &simonsayspb.Response_Turn{
					Turn: simonsayspb.Response_BEGIN,
				},
			})
		} else {
			sendErr = stream.Send(&simonsayspb.Response{
				Event: &simonsayspb.Response_LightUp{
					LightUp: req.GetPress(),
				},
			})
		}
		req.GetPress()
		if sendErr != nil {
			log.Fatalf("Error while sending data to client: %v", sendErr)
			return sendErr
		}
	}

}

func main() {
	fmt.Println("Simon says hello")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatal("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	simonsayspb.RegisterSimonSaysServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatal("failed to serve: %v", err)
	}

}
