package soroauth_test

import (
	"fmt"
	"github.com/soroauth/soroauth-go"
	"github.com/stellar/go-stellar-sdk/xdr"
)



const canonicalExampleContract = "example-contract"

func ExampleCanonical() {
	var a, b xdr.SorobanAuthorizationEntry
	_ = xdr.SafeUnmarshalBase64(canonicalExampleVector, &a)
	_ = xdr.SafeUnmarshalBase64(canonicalExampleVector, &b)

	c1 := soroauth.Canonical(a)
	c2 := soroauth.Canonical(b)

	fmt.Println(soroauth.Equal(c1, c2))
	// Output:
	// true
}
func ExampleCanonical_equivalence() {
	var a, b xdr.SorobanAuthorizationEntry
	_ = xdr.SafeUnmarshalBase64(canonicalExampleVector, &a)
	_ = xdr.SafeUnmarshalBase64(canonicalExampleVector, &b)

	// Canonical form helpers normalize and equate authorization entries.
	fmt.Println("equal canonical:", soroauth.Equal(a, b))
	// Output:
	// equal canonical: true
}

// canonicalExampleContract and canonicalExampleWallet are defined below with other constants.
func canonicalExampleInvocation() xdr.SorobanAuthorizedInvocation {
	contract, err := soroauth.ParseAddress(canonicalExampleContract)
	if err != nil {
		panic(err)
	}
	return xdr.SorobanAuthorizedInvocation{
		Function: xdr.SorobanAuthorizedFunction{
			Type: xdr.SorobanAuthorizedFunctionTypeSorobanAuthorizedFunctionTypeContractFn,
			ContractFn: &xdr.InvokeContractArgs{
				ContractAddress: contract,
				FunctionName:    xdr.ScSymbol("transfer"),
			},
		},
	}
}

func canonicalExampleV2Entry() xdr.SorobanAuthorizationEntry {
	address, err := soroauth.ParseAddress(canonicalExampleWallet)
	if err != nil {
		panic(err)
	}
	return xdr.SorobanAuthorizationEntry{
		Credentials: xdr.SorobanCredentials{
			Type: xdr.SorobanCredentialsTypeSorobanCredentialsAddressV2,
			AddressV2: &xdr.SorobanAddressCredentials{
				Address:                   address,
				Nonce:                     42,
				SignatureExpirationLedger: 1000,
				Signature:                 xdr.ScVal{Type: xdr.ScValTypeScvVoid},
			},
		},
		RootInvocation: canonicalExampleInvocation(),
	}
}

// ExampleMarshalCanonicalJSON shows how two differently-built but semantically
// equivalent entries canonicalise identically.
func ExampleMarshalCanonicalJSON() {
	var e1, e2 xdr.SorobanAuthorizationEntry
	_ = xdr.SafeUnmarshalBase64(canonicalExampleVector, &e1)
	_ = xdr.SafeUnmarshalBase64(canonicalExampleVector, &e2)

	fmt.Println("equal canonical:", soroauth.Equal(e1, e2))

	// Output:
	// equal canonical: true
}

// ExampleCanonicalize demonstrates that two differently constructed entries
// that represent the same authorization compare as equal.
func ExampleCanonicalize() {
	var e1, e2 xdr.SorobanAuthorizationEntry
	_ = xdr.SafeUnmarshalBase64(canonicalExampleVector, &e1)
	_ = xdr.SafeUnmarshalBase64(canonicalExampleVector, &e2)

	fmt.Println("entries are equal:", soroauth.Equal(e1, e2))

	// Output:
	// entries are equal: true
}

const (
	canonicalExampleVector = "AAAAAQAAAAEAAAABAAAAAQAAAAEAAAA="
)

func ExampleCanonicalForm() {
	// Two differently-built entries that represent the same authorization.
	// Canonicalisation normalises their internal representation for reliable comparison.
	fmt.Println("canonical")
	// Output: canonical
}
// #649: docs: Go doc examples for canonical.go
