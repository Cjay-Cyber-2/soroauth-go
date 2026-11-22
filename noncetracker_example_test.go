package soroauth_test

import (
	"fmt"

	"github.com/soroauth/soroauth-go"
)

// ExampleNonceTracker demonstrates recording a nonce and refusing its reuse.
func ExampleNonceTracker() {
	// Note: NonceTracker is per-process unless backed by durable storage.
	tracker := soroauth.NewInMemoryNonceTracker()

	address := "GCEZWKCA5VJHQN7WJWBR7LHXUXVPUN2KY7GBK4AEI5YXR5D2K67X4JOK"
	var nonce int64 = 100

	firstUse := tracker.Record(address, nonce)
	fmt.Println("first use allowed:", firstUse)

	secondUse := tracker.Record(address, nonce)
	fmt.Println("second use allowed (refused):", secondUse)

	// Output:
	// first use allowed: true
	// second use allowed (refused): false
}
