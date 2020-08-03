package main // import "github.com/adisbladis/trustix"

import (
	"fmt"
)

func main() {

	fmt.Println("HELLO")

	// contractAddr := common.HexToAddress("0x21e6fc92f93c8a1bb41e2be64b4e1f88a54d3576")
	// // Fake derivation hash
	// drvHash := sha256.Sum256([]byte("fakeDerivationHash"))
	// // Fake signer address
	// signer := common.HexToAddress("0x21e6fc92f93c8a1bb41e2be64b4e1f88a54d3576")

	// conn, err := ethclient.Dial("./geth.sock")
	// if err != nil {
	// 	log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	// }

	// // Instantiate the contract at known address
	// registry, err := registry.NewNarRegistry(contractAddr, conn)
	// if err != nil {
	// 	log.Fatalf("Failed to instantiate a NarRegistry contract: %v", err)
	// }

	// narInfoHash, err := registry.LookupNarInfoHash(nil, signer, drvHash)
	// if err != nil {
	// 	log.Fatalf("Failed to retrieve narinfo hash: %v", err)
	// }
	// fmt.Println("Hash: ", narInfoHash)

}
