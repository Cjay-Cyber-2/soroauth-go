package soroauth_test

import (
	"fmt"
	"github.com/soroauth/soroauth-go"
)

// ExampleNonceTracker demonstrates recording a nonce and refusing its reuse.
// Note that the default nonce tracker is in-memory and process-local unless
// backed by durable storage.
func ExampleNonceTracker() {
	tracker := soroauth.NewInMemoryNonceTracker()

	account := "GAB7BCQ72A23Z3BXTJVRW66JXYW25J6Y3D7V3X6X44J6H2H2H2H2H2"
	var nonce uint64 = 100

	firstUse := tracker.CheckAndRecord(account, nonce)
	fmt.Println("first use accepted:", firstUse)

	secondUse := tracker.CheckAndRecord(account, nonce)
	fmt.Println("second use accepted:", secondUse)

	// Output:
	// first use accepted: true
	// second use accepted: false
}
