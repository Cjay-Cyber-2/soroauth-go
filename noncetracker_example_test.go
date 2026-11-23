package soroauth_test

import (
	"fmt"

	"github.com/soroauth/soroauth-go"
)

const (
	nonceTrackerExampleKey = "GAB7BCQ72A23Z3BXTJVRW66JXYW25J6Y3D7V3X6X44J6H2H2H2H2"
)

// ExampleNonceTracker demonstrates recording a nonce and refusing a reuse.
// Note that the default in-memory nonce tracker is process-local unless backed by durable storage.
func ExampleNonceTracker() {
	tracker := soroauth.NewInMemoryNonceTracker()
	address := "GAB7BCQ72A23Z3BXTJVRW66JXYW25J6Y3D7V3X6X44J6H2H2H2H2"
	nonce := uint64(100)

	err := tracker.Record(address, nonce)
	fmt.Println("first use err:", err == nil)

	err = tracker.Record(address, nonce)

	// Output:
	// first use err: true
	// reuse refused: true
}
