package main

import (
	"context"
	"flag"
	"fmt"

	// "log"
	"math"
	"net"
	"net/http"
	"os"
	"time"

	pb "github.com/CodeYourFuture/immersive-go-course/grpc-client-server/prober"
	"github.com/charmbracelet/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement prober.ProberServer.
type server struct {
	pb.UnimplementedProberServer
}

func (s *server) DoProbes(ctx context.Context, in *pb.ProbeRequest) (*pb.ProbeReply, error) {
	var i int32
	var sum, n float32

	client := http.Client{Timeout: 1 * time.Second}
	for i = 0; i < in.GetRequestNum(); i++ {
		log.Info(fmt.Sprintf("Making request to %s", in.GetEndpoint()))
		start := time.Now()
		resp, err := client.Get(in.GetEndpoint())
		if os.IsTimeout(err) {
			log.Errorf("Request to %s timed out: %s", in.GetEndpoint(), err.Error())
			continue
		}
		if err != nil {
			log.Errorf("Error reaching url. Check if url provided is correct. Error: %s", err.Error())
			continue
		}
		if resp.StatusCode == 200 {
			n += 1.0
			elapsed := time.Since(start)
			elapsedMsecs := float32(elapsed / time.Millisecond)
			sum += elapsedMsecs
			log.Info(fmt.Sprintf("Response time: %vms", elapsedMsecs))

		}
	}

	averageElapsedMsecs := sum / n
	if math.IsNaN(float64(averageElapsedMsecs)) {
		averageElapsedMsecs = 0
	}

	recordLatency(in.GetEndpoint(), float64(averageElapsedMsecs))

	return &pb.ProbeReply{AvgLatencyMsecs: averageElapsedMsecs}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterProberServer(s, &server{})
	log.Info(fmt.Sprintf("Server listening at %v", lis.Addr()))
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Prometheus
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func recordLatency(label string, value float64) {
	latency.WithLabelValues(label).Set(value)

}

var (
	latency = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "latency",
		Name:      "client_requests_ms",
		Help:      "The average latency of client's requests",
	},
		[]string{"request_url"})
)
