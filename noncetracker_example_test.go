package soroauth_test

import (
	"fmt"

	"github.com/soroauth/soroauth-go"
)


const (
	nonceTrackerExampleKey = "GAB7BCQ72A23Z3BXTJVRW66JXYW25J6Y3D7V3X6X44J6H2H2H2H2"
)
// ExampleNonceTracker demonstrates recording a nonce and refusing its reuse.
// Note that the NonceTracker is process-local unless backed by durable storage.
func ExampleNonceTracker() {
	tracker := soroauth.NewInMemoryNonceTracker()

	address := "GAB7BCQ72A23Z3BXTJVRW66JXYW25J6Y3D7V3X6X44J6H2H2H2H2H2"
	var nonce uint64 = 100

	firstCheck := tracker.CheckAndRecord(address, nonce)
	fmt.Println("first use allowed:", firstCheck == nil)

	secondCheck := tracker.CheckAndRecord(address, nonce)
	fmt.Println("second use allowed:", secondCheck == nil)

	// Output:
	// first use allowed: true
	// second use allowed: false
}
