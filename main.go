package main // import "github.com/tweag/trustix"

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/google/trillian"
	pb "github.com/tweag/trustix/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

var (
	tLogEndpoint = flag.String("tlog_endpoint", ":8090", "The gRPC endpoint of the Trillian Log Server.")
	tLogID       = flag.Int64("tlog_id", 0, "Trillian Log ID")
	tInputHash   = flag.String("input_hash", "", "The input hash (derivation hash)")
	tOutputHash  = flag.String("output_hash", "", "The output hash (output hash)")
	action       = flag.String("action", "submit", "Operation mode (submit/daemon)")
)

type pbServer struct {
	pb.UnimplementedTrustixServer
	trillianServer *server
}

func (s *pbServer) SubmitMapping(ctx context.Context, in *pb.SubmitRequest) (*pb.SubmitReply, error) {
	inputHash := hex.EncodeToString(in.InputHash)
	outputHash := hex.EncodeToString(in.OutputHash)

	fmt.Println(inputHash)
	fmt.Println(outputHash)

	fmt.Println(in)

	value := newInput(inputHash, outputHash)
	resp := &Response{}

	log.Println("[main] Submitting it for inclusion in the Trillian Log")
	resp, err := s.trillianServer.put(&Request{
		input: *value,
	})
	if err != nil {
		return nil, err
	}

	log.Printf("[main] put: %s", resp.status)

	log.Println("[main] Retrieving it from the Trillian Log")
	resp, err = s.trillianServer.get(&Request{
		input: *value,
	})
	log.Printf("[main] get: %s", resp.status)

	return &pb.SubmitReply{
		Status: pb.SubmitReply_OK,
	}, nil
}

func main() {
	flag.Parse()

	port := ":8081"

	logID := *tLogID
	if logID == 0 {
		envLogId := os.Getenv("LOG_ID")
		if envLogId == "" {
			envLogId = "0"
		}
		iEnvLogID, _ := strconv.Atoi(envLogId)
		logID = int64(iEnvLogID)

		if logID == 0 {
			panic("logID 0 is invalid")
		}
	}

	// if *action == "daemon"
	switch *action {
	case "daemon":

		log.Printf("[main] Establishing connection w/ Trillian Log Server [%s]", *tLogEndpoint)
		conn, err := grpc.Dial(*tLogEndpoint, grpc.WithInsecure())
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		// Create a Trillian Log Server client
		log.Println("[main] Creating new Trillian Log Client")
		tLogClient := trillian.NewTrillianLogClient(conn)

		log.Printf("[main] Creating Server using LogID [%d]", logID)
		server := newServer(tLogClient, logID)

		lis, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()

		pb.RegisterTrustixServer(s, &pbServer{
			trillianServer: server,
		})
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	case "submit":

		address := "localhost:8081"

		inputHash := *tInputHash
		outputHash := *tOutputHash

		if inputHash == "" {
			inputHash = os.Getenv("INPUT_HASH")
			if inputHash == "" {
				panic("Input hash cannot be empty")
			}
		}
		if outputHash == "" {
			outputHash = os.Getenv("OUTPUT_HASH")
			if outputHash == "" {
				panic("Output hash cannot be empty")
			}
		}

		conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		c := pb.NewTrustixClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		inputHashBytes, err := hex.DecodeString(inputHash)
		if err != nil {
			panic(err)
		}

		outputHashBytes, err := hex.DecodeString(outputHash)
		if err != nil {
			panic(err)
		}

		r, err := c.SubmitMapping(ctx, &pb.SubmitRequest{
			OutputHash: outputHashBytes,
			InputHash:  inputHashBytes,
		})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}

		fmt.Println(r.Status)

	default:
		panic(fmt.Errorf("Action '%s' not supported", *action))
	}
}
