package soroauth_test

import (
	"fmt"
	"github.com/soroauth/soroauth-go"
)

// ExampleNonceTracker demonstrates recording a nonce and refusing a reuse.
// Note: The default NonceTracker is process-local and in-memory unless backed by durable storage.
func ExampleNonceTracker() {
	tracker := soroauth.NewInMemoryNonceTracker()
	address := "GCEZWKCA5VLDNRLN3RPRJMRZOX3Z6G5CHCGSNFHEYVXM3XOJMDS674JZ"
	var nonce uint64 = 100

	// First use records successfully
	err := tracker.CheckAndRecord(address, nonce)
	fmt.Println("first use error:", err == nil)

	// Second use with the same nonce is refused
	err = tracker.CheckAndRecord(address, nonce)
	fmt.Println("reuse refused:", err != nil)

	// Output:
	// first use error: true
	// reuse refused: true
}
