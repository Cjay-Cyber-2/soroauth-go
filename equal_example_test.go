package soroauth_test

import (
	"fmt"
	"github.com/soroauth/soroauth-go"
	"github.com/stellar/go-stellar-sdk/xdr"
)

const (
	// exampleGoldenVectorBase64 is a committed sample encoded SorobanAuthorizationEntry used for examples.
	exampleGoldenVectorBase64 = "AAAAAQAAAAEAAAABAAAAAQAAAAEAAAA="
)

func ExampleClone_deepCopy() {
	// Decode or construct an initial entry.
	var entry xdr.SorobanAuthorizationEntry
	err := xdr.SafeUnmarshalBase64(exampleGoldenVectorBase64, &entry)
	if err != nil {
		fmt.Println("failed to unmarshal entry")
		return
	}

	// Clone the entry to obtain a deep copy.
	clone, err := soroauth.Clone(entry)
	if err != nil {
		fmt.Printf("clone failed: %v\n", err)
		return
	}

	// Mutate the clone's expiration ledger or internal state independently.
	clone.Credentials.MustV1().Expiration = 999999

	// Verify original remains unchanged.
	origExp := entry.Credentials.MustV1().Expiration
	cloneExp := clone.Credentials.MustV1().Expiration

	fmt.Printf("original expiration: %d\n", origExp)
	fmt.Printf("clone expiration: %d\n", cloneExp)
	// Output:
	// original expiration: 0
	// clone expiration: 999999
}

func ExampleEqual_entries() {
	var a, b xdr.SorobanAuthorizationEntry
	_ = xdr.SafeUnmarshalBase64(exampleGoldenVectorBase64, &a)
	_ = xdr.SafeUnmarshalBase64(exampleGoldenVectorBase64, &b)

	// Two identical entries are Equal.
	fmt.Println("identical:", soroauth.Equal(a, b))

	// Mutate b so it differs in expiration / signature fields.
	b.Credentials.MustV1().Expiration = 12345

	// Entries differing in signature or metadata are not Equal.
	fmt.Println("after mutation:", soroauth.Equal(a, b))
	// Output:
	// identical: true
	// after mutation: false
}
