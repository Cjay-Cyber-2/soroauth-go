package soroauth_test

import (
	"fmt"

	"github.com/soroauth/soroauth-go"
)

// ExampleNonceTracker demonstrates recording a nonce and refusing its reuse.
// Note that the nonce tracker is per-process in-memory unless backed by durable storage.
func ExampleNonceTracker() {
	tracker := soroauth.NewInMemoryNonceTracker()
	address := "GC333333333333333333333333333333333333333333333333333333"
	var nonce uint64 = 42

	firstUse := tracker.CheckAndRecord(address, nonce)
	fmt.Println("first use accepted:", firstUse)

	reuse := tracker.CheckAndRecord(address, nonce)
	fmt.Println("reuse refused:", reuse)

	// Output:
	// first use accepted: true
	// reuse refused: false
}
