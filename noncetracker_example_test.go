package soroauth_test

import (
	"fmt"
	"github.com/soroauth/soroauth-go"
	"context"
)


const (
	nonceTrackerExampleKey = "GAB7BCQ72A23Z3BXTJVRW66JXYW25J6Y3D7V3X6X44J6H2H2H2H2"
)

func ExampleNewMemoryNonceTracker() {
	ctx := context.Background()
	tracker := soroauth.NewMemoryNonceTracker()
	var nonce int64 = 100

	err1 := tracker.Reserve(ctx, nonceTrackerExampleKey, nonce)
	err2 := tracker.Reserve(ctx, nonceTrackerExampleKey, nonce)

	fmt.Println("first use allowed:", err1 == nil)
	fmt.Println("second use allowed:", err2 == nil)
	// Output:
	// first use allowed: true
	// second use allowed: false
}

func ExampleNewInMemoryNonceTracker() {
	ctx := context.Background()
	tracker := soroauth.NewInMemoryNonceTracker()
	var nonce int64 = 100

	err1 := tracker.Reserve(ctx, nonceTrackerExampleKey, nonce)
	err2 := tracker.Reserve(ctx, nonceTrackerExampleKey, nonce)

	fmt.Println("first reserve error:", err1)
	fmt.Println("second reserve error:", err2 != nil)
	// Output:
	// first reserve error: <nil>
	// second reserve error: true
}
func ExampleNonceTracker() {
	// Note: The default nonce tracker is per-process in-memory unless durable storage is provided.
	tracker := soroauth.NewNonceTracker()
	account := "GAB7BCQ72A23Z3BXTJVRW66JXYW25J6Y3D7V3X6X44J6H2H2H2H2"
	var nonce int64 = 100

	firstUse := tracker.CheckAndRecord(account, nonce)
	fmt.Println("first use allowed:", firstUse)

	secondUse := tracker.CheckAndRecord(account, nonce)
	fmt.Println("second use allowed:", secondUse)

	// Output:
	// first use allowed: true
	// second use allowed: false
}
