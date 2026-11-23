package soroauth_test

import (
	"fmt"

	"github.com/soroauth/soroauth-go"
)

// ExampleNonceTracker shows a nonce being recorded and a reuse being refused.
// Note: the default nonce tracker is process-local and in-memory; it does not
// persist across process restarts unless backed by durable storage.
func ExampleNonceTracker() {
	tracker := soroauth.NewInMemoryNonceTracker()

	// Address and nonce to record.
	address := "GAB7BCQ72A23Z3BXTJVRW66JXYW25J6Y3D7V3X6X44J6H2H2H2H2H2H2"
	var nonce uint64 = 42

	// First use records successfully.
	err := tracker.CheckAndRecord(address, nonce)
	fmt.Println("first use err:", err)

	// Second use of the same nonce refuses reuse.
	err = tracker.CheckAndRecord(address, nonce)
	fmt.Println("second use err:", err != nil)

	// Output:
	// first use err: <nil>
	// second use err: true
}
