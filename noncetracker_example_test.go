
package soroauth_test

import (
	"fmt"

	"github.com/soroauth/soroauth-go"
)


const (
	nonceTrackerExampleKey = "GAB7BCQ72A23Z3BXTJVRW66JXYW25J6Y3D7V3X6X44J6H2H2H2H2"
)

func ExampleNonceTracker() {
	tracker := soroauth.NewMemoryNonceTracker()
	var nonce uint64 = 42

	fmt.Println(tracker.CheckAndRecord(nonce))
	fmt.Println(tracker.CheckAndRecord(nonce))
	// Output:
	// true
	// false
}

func ExampleNewMemoryNonceTracker() {
	tracker := soroauth.NewMemoryNonceTracker()
	var nonce uint64 = 100

	firstUse := tracker.CheckAndRecord(nonce)
	secondUse := tracker.CheckAndRecord(nonce)

	fmt.Println("first use allowed:", firstUse)
	fmt.Println("second use allowed:", secondUse)
	// Output:
	// first use allowed: true
	// second use allowed: false
}