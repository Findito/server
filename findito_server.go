package main

import (
	"context"
	"io"
	"log"
	"net"
	"time"

	//"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/Findito/server/proto"
)

type finditoServer struct {
	pb.UnimplementedFinditoServiceServer
	// TODO: Implement db conn
}

func (s *finditoServer) ReportFoundItem(ctx context.Context, in *pb.ReportFoundItemRequest) (*pb.ReportFoundItemResponse, error) {
	log.Printf("Received: %v", in.GetName())

	if in.GetName() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Name is required")
	}

	return &pb.ReportFoundItemResponse{
		ItemId:  "12345",
		Success: true,
		Message: "Item reported successfully",
	}, nil

}

func (s *finditoServer) SearchItemsAlongRoute(ctx context.Context, in *pb.SearchItemsAlongRouteRequest) (*pb.SearchItemsAlongRouteResponse, error) {
	log.Printf("Recieved along the route item search request")

	return &pb.SearchItemsAlongRouteResponse{
		Items: nil,
	}, nil
}

func (s *finditoServer) TrackNearbyItems(stream pb.FinditoService_TrackNearbyItemsServer) error {
	log.Printf("Recieved track nearby items request")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("Client closed the stream")
			return nil
		}

		if err != nil {
			log.Printf("Error receiving location update: %v", err)
			return err
		}

		items := fetchNearbyItems(req.GetLocation(), req.RadiusKm)

		for _, item := range items {
			log.Printf("Sending item: %v", item)
			if err := stream.Send(item); err != nil {
				log.Printf("Error sending item: %v", err)
			}
		}

	}
}

func fetchNearbyItems(location *pb.Location, radiusKm float64) []*pb.FoundItem {
	log.Printf("Fetching nearby items for location: %v", location)

	// TODO: Implement search algorithm

	return []*pb.FoundItem{
		{
			Id:          "abc123",
			Name:        "Lost Keys",
			Description: "Set of car keys on a red keychain",
			Location: &pb.Location{
				Latitude:  50.049683,
				Longitude: 19.944544,
			},
			FoundTimestamp: time.Now().Unix(),
		},
	}
}

func RunServer(port string) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterFinditoServiceServer(s, &finditoServer{})

	log.Printf("server listening at %v", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	port := ":50051" // lub inny port, np. ":8080"
	RunServer(port)
}
