package soroauth_test

import (
	"fmt"
	"github.com/soroauth/soroauth-go"
)

// ExampleNonceTracker demonstrates recording a nonce and refusing reuse.
//
// Note: the default nonce tracker is in-memory and process-local. It does not
// persist across process restarts unless backed by durable storage.
func ExampleNonceTracker() {
	tracker := soroauth.NewInMemoryNonceTracker()
	address := "GCEZWKCA5VJHQN7WJWBR7LHXUXVPUN2KY7GBK4AEI5YXR5D2K67X4JOK"
	var nonce uint64 = 42

	// First use records the nonce successfully.
	firstUse, err := tracker.CheckAndRecord(address, nonce)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("first use recorded:", firstUse)

	// Second use of the same nonce is refused.
	reused, err := tracker.CheckAndRecord(address, nonce)
	fmt.Println("reuse allowed:", reused, "error:", err != nil)

	// Output:
	// first use recorded: true
	// reuse allowed: false error: true
}
