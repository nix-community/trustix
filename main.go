package main

import (
	"flag"
	"github.com/google/trillian"
	"google.golang.org/grpc"
	"log"
	"os"
	"strconv"
)

var (
	tLogEndpoint = flag.String("tlog_endpoint", ":8090", "The gRPC endpoint of the Trillian Log Server.")
	tLogID       = flag.Int64("tlog_id", 0, "Trillian Log ID")
	tInputHash   = flag.String("input_hash", "", "The input hash (derivation hash)")
	tOutputHash  = flag.String("output_hash", "", "The output hash (output hash)")
)

func main() {
	flag.Parse()

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

	value := newInput(inputHash, outputHash)
	resp := &Response{}

	log.Println("[main] Submitting it for inclusion in the Trillian Log")
	resp, err = server.put(&Request{
		input: *value,
	})
	log.Printf("[main] put: %s", resp.status)

	log.Println("[main] Retrieving it from the Trillian Log")
	resp, err = server.get(&Request{
		input: *value,
	})
	log.Printf("[main] get: %s", resp.status)
}
