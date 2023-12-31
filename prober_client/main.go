// Package main implements a client for Prober service.
package main

import (
	"context"
	"flag"
	"log"

	pb "github.com/CodeYourFuture/immersive-go-course/grpc-client-server/prober"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	url  = flag.String("endpoint", "https://google.com", "the url to make request to")
	n    = flag.Int("n", 500, "the number of requests to be made")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewProberClient(conn)

	// Contact the server and print out its response.
	ctx := context.Background() // TODO: add a timeout

	// TODO: endpoint should be a flag
	// TODO: add number of times to probe
	r, err := c.DoProbes(ctx, &pb.ProbeRequest{Endpoint: *url, RequestNum: int32(*n)})
	if err != nil {
		log.Fatalf("could not probe: %v", err)
	}
	log.Printf("Response Time: %.2f", r.GetAvgLatencyMsecs())
}
