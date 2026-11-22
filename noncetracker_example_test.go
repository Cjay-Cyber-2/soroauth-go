package soroauth_test

import (
	"fmt"
	"github.com/soroauth/soroauth-go"
)

// ExampleNonceTracker demonstrates recording a nonce and refusing reuse.
// Note that the nonce tracker is per-process unless backed by durable storage.
func ExampleNonceTracker() {
	tracker := soroauth.NewInMemoryNonceTracker()
	account := "GCEZWKCA5VJHQN7WJWBR7LHXUXVPUN2KY7GBK4AEI5YXR5D2K67X4JOK"

	// First use of a nonce is accepted.
	firstUse := tracker.CheckAndRecord(account, 100)
	fmt.Println("first use accepted:", firstUse)

	// Reusing the same nonce for the same account is refused.
	secondUse := tracker.CheckAndRecord(account, 100)
	fmt.Println("reuse accepted:", secondUse)

	// Output:
	// first use accepted: true
	// reuse accepted: false
}
