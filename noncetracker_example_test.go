package soroauth_test

import (
	"fmt"

	"github.com/soroauth/soroauth-go"
)

const (
	nonceTrackerExampleKey = "GAB7BCQ72A23Z3BXTJVRW66JXYW25J6Y3D7V3X6X44J6H2H2H2H2H2"
)

// ExampleNonceTracker demonstrates recording a nonce and refusing its reuse.
// Note that the nonce tracker is per-process unless backed by durable storage.
func ExampleNonceTracker() {
	tracker := soroauth.NewInMemoryNonceTracker()

	var account = nonceTrackerExampleKey
	var nonce int64 = 42

	err1 := tracker.Reserve(nil, account, nonce)
	fmt.Println("first use accepted:", err1 == nil)

	err2 := tracker.Reserve(nil, account, nonce)
	fmt.Println("second use accepted:", err2 == nil)

	// Output:
	// first use accepted: true
	// second use accepted: false
}
