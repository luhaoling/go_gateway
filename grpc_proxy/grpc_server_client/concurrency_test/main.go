package main

import (
	"context"
	"fmt"
	pb "github.com/e421083458/go_gateway/grpc_proxy/grpc_server_client/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	addr := "127.0.0.1:8012"
	processTime := int64(20)

	ctx, _ := context.WithTimeout(context.Background(), time.Duration(processTime)*time.Second)

	wg := sync.WaitGroup{}
	var totalCount int64
	var successCount int64
	var failCount int64
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(ctx context.Context) {
			defer wg.Done()
			conn, err := grpc.Dial(addr, grpc.WithInsecure())
			if err != nil {
				log.Fatalf("did not connect: %v", err)
			}
			defer conn.Close()
			c := pb.NewEchoClient(conn)

			for {
				select {
				case <-ctx.Done():
					fmt.Println("ctx.Done")
					return
				default:
				}
				atomic.AddInt64(&totalCount, 1)
				if err := unaryCallWithMetadata(c, "this is examples/metadata"); err != nil {
					atomic.AddInt64(&failCount, 1)
				} else {
					atomic.AddInt64(&successCount, 1)
				}
			}
		}(ctx)
	}
	wg.Wait()
	fmt.Printf("result qps:%v succ:%v fail:%v", totalCount/processTime, successCount, failCount)
}

var AccessToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTA5NDg4ODIsImlzcyI6ImFwcF9pZF9iIn0.RgxX45dqhbijoSoRNkipS-6_m1ylyupwkt76dKs0Y1E"

func unaryCallWithMetadata(c pb.EchoClient, message string) error {

	md := metadata.Pairs("timestamp", time.Now().Format(time.StampNano))
	md.Append("authorization", "Bearer "+AccessToken)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	_, err := c.UnaryEcho(ctx, &pb.EchoRequest{Message: message})
	if err != nil {
		return err
	}
	return nil
}
