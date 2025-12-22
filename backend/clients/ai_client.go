package clients

import (
	"context"
	"log"
	"os"
	"time"

	pb "github.com/wesuuu/helpnow/backend/gen/ai_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AIClient struct {
	conn   *grpc.ClientConn
	client pb.AIServiceClient
}

var GlobalAIClient *AIClient

func InitAIClient() {
	addr := os.Getenv("AI_SERVICE_ADDR")
	if addr == "" {
		addr = "localhost:50051"
	}

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to AI Service: %v", err)
	}

	c := pb.NewAIServiceClient(conn)
	GlobalAIClient = &AIClient{
		conn:   conn,
		client: c,
	}
	log.Println("Connected to AI Service at", addr)
}

func (c *AIClient) Close() {
	c.conn.Close()
}

func (c *AIClient) ExecuteRoutine(ctx context.Context, routineID string, agentID string, inputParams map[string]string, userID string) (*pb.ExecuteRoutineResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	return c.client.ExecuteRoutine(ctx, &pb.ExecuteRoutineRequest{
		RoutineId:   routineID,
		AgentId:     agentID,
		InputParams: inputParams,
		UserId:      userID,
	})
}
