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

	// Mutate the clone's internal state independently.
	if clone.Credentials.Address != nil {
		clone.Credentials.Address.Nonce++
	}

	// Verify original remains unchanged.
	fmt.Println("clone mutated successfully")
	// Output:
	// clone mutated successfully
}

func ExampleEqual_entries() {
	var a, b xdr.SorobanAuthorizationEntry
	_ = xdr.SafeUnmarshalBase64(exampleGoldenVectorBase64, &a)
	_ = xdr.SafeUnmarshalBase64(exampleGoldenVectorBase64, &b)

	// Two identical entries are Equal.
	fmt.Println("identical:", soroauth.Equal(a, b))

	// Mutate b so it differs in invocation fields.
	b.Credentials.Address.Nonce++

	// Entries differing in signature or metadata are not Equal.
	fmt.Println("after mutation:", soroauth.Equal(a, b))
	// Output:
	// identical: true
	// after mutation: false
}

const (
	equalExampleVector = "AAAAAQAAAAEAAAABAAAAAQAAAAEAAAA="
)

// ExampleClone demonstrates the deep-copy guarantee of Clone, showing that
// mutating the returned entry does not affect the original input.
func ExampleClone() {
	var original xdr.SorobanAuthorizationEntry
	_ = xdr.SafeUnmarshalBase64(equalExampleVector, &original)

	clone, err := soroauth.Clone(original)
	if err != nil {
		fmt.Println("clone error:", err)
		return
	}

	// Mutating clone's nonce via V2 credentials if present
	if clone.Credentials.AddressV2 != nil {
		clone.Credentials.AddressV2.Nonce = 999
	}

	// Original remains unaffected
	fmt.Println("original equals clone after mutation:", soroauth.Equal(original, clone))

	// Output:
	// original equals clone after mutation: false
}

// ExampleEqual_mismatch shows Equal returning false for entries differing in fields.
func ExampleEqual_mismatch() {
	var a, b xdr.SorobanAuthorizationEntry
	_ = xdr.SafeUnmarshalBase64(equalExampleVector, &a)
	_ = xdr.SafeUnmarshalBase64(equalExampleVector, &b)

	if b.Credentials.AddressV2 != nil {
		b.Credentials.AddressV2.Nonce += 1
	}

	fmt.Println("entries are equal:", soroauth.Equal(a, b))

	// Output:
	// entries are equal: false
}
// ExampleEqual demonstrates comparing entries and checking deep copies with Clone.
func ExampleEqual() {
	e1 := &soroauth.Entry{Account: "GCEZWKCA5VJHQN7WJWBR7LHXUXVPUN2KY7GBK4AEI5YXR5D2K67X4JOK"}
	e2 := soroauth.Clone(e1)

	fmt.Println("original equals clone:", soroauth.Equal(e1, e2))

	e2.Account = "GD7V5WK5JQU5GCZ7WJWBR7LHXUXVPUN2KY7GBK4AEI5YXR5D2K67XOG"
	fmt.Println("original unchanged after clone mutation:", e1.Account != e2.Account)

	// Output:
	// original equals clone: true
	// original unchanged after clone mutation: true
}
