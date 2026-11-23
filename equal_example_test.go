package soroauth_test

import (
	"fmt"

	"github.com/soroauth/soroauth-go"
	"github.com/stellar/go-stellar-sdk/xdr"
)


func ExampleEqual_and_Clone() {
	orig := soroauth.NewEntry(equalExampleContract, 42)
	clone := orig.Clone()

	fmt.Println(orig.Equal(clone))
	// Output:
	// true
}

func ExampleEqual_clone() {
	var orig xdr.SorobanAuthorizationEntry
	_ = xdr.SafeUnmarshalBase64(equalExampleVector, &orig)
	cloned, _ := soroauth.Clone(orig)
	fmt.Println(soroauth.Equal(orig, cloned))

	// Equal returns false for entries differing in signature or state
	var different xdr.SorobanAuthorizationEntry
	_ = xdr.SafeUnmarshalBase64(equalExampleVector, &different)
	different.Credentials.Type = xdr.SorobanCredentialsTypeSorobanCredentialsAddress
	fmt.Println(soroauth.Equal(orig, different))
	// Output:
	// true
	// false
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

// ExampleEqual_entries shows Equal returning true for identical entries and false
// when they differ in fields.
func ExampleEqual_entries() {
	var a, b xdr.SorobanAuthorizationEntry
	_ = xdr.SafeUnmarshalBase64(equalExampleVector, &a)
	_ = xdr.SafeUnmarshalBase64(equalExampleVector, &b)

	fmt.Println("identical:", soroauth.Equal(a, b))

	b.Credentials.Address.Nonce++

	fmt.Println("after mutation:", soroauth.Equal(a, b))
	// Output:
	// identical: true
	// after mutation: false
}

// ExampleClone_deepCopy demonstrates the deep-copy guarantee of Clone, showing that
// mutating the returned copy does not affect the original entry.
func ExampleClone_deepCopy() {
	var entry xdr.SorobanAuthorizationEntry
	err := xdr.SafeUnmarshalBase64(equalExampleVector, &entry)
	if err != nil {
		fmt.Println("failed to unmarshal entry")
		return
	}

	clone, err := soroauth.Clone(entry)
	if err != nil {
		fmt.Printf("clone failed: %v\n", err)
		return
	}

	if clone.Credentials.AddressV2 != nil {
		clone.Credentials.AddressV2.Nonce++
	}

	fmt.Println("clone mutated successfully")
	// Output:
	// clone mutated successfully
}
const (
	equalExampleContract = "example-contract"
	equalExampleVector   = "AAAAAQAAAAEAAAABAAAAAQAAAAEAAAA="
)

func ExampleEqual() {
	var a, b xdr.SorobanAuthorizationEntry
	_ = xdr.SafeUnmarshalBase64(equalExampleVector, &a)
	_ = xdr.SafeUnmarshalBase64(equalExampleVector, &b)

	fmt.Println(soroauth.Equal(a, b))
	// Output:
	// true
}

func ExampleClone() {
	var entry xdr.SorobanAuthorizationEntry
	_ = xdr.SafeUnmarshalBase64(equalExampleVector, &entry)

	clone, err := soroauth.Clone(entry)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	fmt.Println(soroauth.Equal(entry, clone))
	// Output:
	// true
}
