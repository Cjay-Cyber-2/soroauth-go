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
	nonce := uint64(12345)

	first := tracker.CheckAndRecord(nonce)
	second := tracker.CheckAndRecord(nonce)

	fmt.Println(first, second)
	// Output:
	// true false
}
