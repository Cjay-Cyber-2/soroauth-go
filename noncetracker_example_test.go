package soroauth_test

import (
	"fmt"
	"github.com/soroauth/soroauth-go"
)


const (
	nonceTrackerExampleKey = "GAB7BCQ72A23Z3BXTJVRW66JXYW25J6Y3D7V3X6X44J6H2H2H2H2H2"
)
// ExampleNonceTracker demonstrates recording a nonce and refusing a reuse.
// Note that the default in-memory nonce tracker is process-local unless backed by durable storage.
func ExampleNonceTracker() {
	tracker := soroauth.NewInMemoryNonceTracker()
	address := "GAB7BCQ72A23Z3BXTJVRW66JXYW25J6Y3D7V3X6X44J6H2H2H2H2H2"
	nonce := uint64(100)

	// First use of the nonce should succeed.
	err := tracker.CheckAndRecord(address, nonce)
	fmt.Println("first use err:", err == nil)

	// Second use of the same nonce should be refused.
	err = tracker.CheckAndRecord(address, nonce)
	fmt.Println("reuse refused:", err != nil)

	// Output:
	// first use err: true
	// reuse refused: true
}
