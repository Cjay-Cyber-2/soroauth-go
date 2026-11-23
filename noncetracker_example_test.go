package soroauth_test

import (
	"context"
	"fmt"
	"github.com/soroauth/soroauth-go"
)


const (
	nonceTrackerExampleKey = "GAB7BCQ72A23Z3BXTJVRW66JXYW25J6Y3D7V3X6X44J6H2H2H2H2"
)
// ExampleNonceTracker demonstrates recording a nonce and refusing a reuse.
func ExampleNonceTracker() {
	ctx := context.Background()
	tracker := soroauth.NewInMemoryNonceTracker()

	address := "GAB7BCQ72A23Z3BXTJVRW66JXYW25J6Y3D7V3X6X44J6H2H2H2H2H2"
	var nonce int64 = 100

	err := tracker.Reserve(ctx, address, nonce)
	fmt.Println("first reservation err:", err == nil)

	err = tracker.Reserve(ctx, address, nonce)
	fmt.Println("second reservation err:", err != nil)

	// Output:
	// first reservation err: true
	// second reservation err: true
}
