package resturant

import (
	"food-siam-si/food-siam-si-menu/internal/handlers/proto"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var RestaurantClient proto.RestaurantServiceClient

func Init() {
	conn, err := grpc.Dial(os.Getenv("RestaurantServiceUrl"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect restaurant service %v", err)
		os.Exit(1)
	}

	RestaurantClient = proto.NewRestaurantServiceClient(conn)
}
