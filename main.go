package main

import (
	"flag"
	"fmt"
	"github.com/google/trillian"
	"google.golang.org/grpc"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	tLogEndpoint = flag.String("tlog_endpoint", "", "The gRPC endpoint of the Trillian Log Server.")
	tLogID       = flag.Int64("tlog_id", 0, "Trillian Log ID")
)

func main() {
	flag.Parse()

	logEndpoint := *tLogEndpoint
	if *tLogEndpoint == "" {
		logEndpoint = ":8090"
	}

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

	// Establish gRPC connection w/ Trillian Log Server
	log.Printf("[main] Establishing connection w/ Trillian Log Server [%s]", *tLogEndpoint)
	conn, err := grpc.Dial(logEndpoint, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Create a Trillian Log Server client
	log.Println("[main] Creating new Trillian Log Client")
	tLogClient := trillian.NewTrillianLogClient(conn)

	// Eventually this personality will be a server
	log.Printf("[main] Creating Server using LogID [%d]", logID)
	server := newServer(tLogClient, logID)

	// Leaves comprise a primary LeafValue (thing) and may have associated ExtraData(extra)
	// The LeafValue will become the hashed value for a node in the Merkle Tree
	log.Println("[main] Creating a 'Thing' and something 'Extra'")
	thing := newThing(fmt.Sprintf("[%s] Thing", time.Now().Format(time.RFC3339)))
	extra := newExtra("Extra")
	fmt.Println(thing)

	// Eventually it will be convenient to explicit Request and Response types
	resp := &Response{}

	// Try to put this Request (Thing+Extra) in the Log
	log.Println("[main] Submitting it for inclusion in the Trillian Log")
	resp, err = server.put(&Request{
		thing: *thing,
		extra: *extra,
	})
	log.Printf("[main] put: %s", resp.status)

	// Try to get this Request (Thing+Extra) from the Log
	log.Println("[main] Retrieving it from the Trillian Log")
	resp, err = server.get(&Request{
		thing: *thing,
		extra: *extra,
	})
	log.Printf("[main] get: %s", resp.status)
}
