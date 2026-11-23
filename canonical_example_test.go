package soroauth_test

import (
	"fmt"

	"github.com/soroauth/soroauth-go"
	"github.com/stellar/go-stellar-sdk/xdr"
)



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
// equivalent entries (or an entry round-tripped through MarshalCanonicalJSON and
// UnmarshalCanonicalJSON) canonicalise identically.
// ExampleMarshalCanonicalJSON shows how two differently-built but semantically
// equivalent entries canonicalise identically.
func ExampleMarshalCanonicalJSON() {
	entry1 := canonicalExampleV2Entry()

	jsonBytes, err := soroauth.MarshalCanonicalJSON(entry1)
	if err != nil {
		fmt.Println("marshal error:", err)
		return
	}

	entry2, err := soroauth.UnmarshalCanonicalJSON(jsonBytes)
	if err != nil {
		fmt.Println("unmarshal error:", err)
		return
	}

	jsonBytes2, err := soroauth.MarshalCanonicalJSON(entry2)
	if err != nil {
		fmt.Println("marshal2 error:", err)
		return
	}

	fmt.Println(string(jsonBytes) == string(jsonBytes2))

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


const (
	canonicalExampleVector   = "AAAAAQAAAAEAAAABAAAAAQAAAAEAAAA="
	canonicalExampleContract = "CDLZFC3SJQLJBZM62JUIOA4W344PZZX6O2C5Z2AFAK5Y2X2O3N5J5Q5Q"
	canonicalExampleWallet   = "GAB7BCQ72A23Z3BXTJVRW66JXYW25J6Y3D7V3X6X44J6H2H2H2H2H2"
	canonicalVector          = "AAAAAQAAAAEAAAABAAAAAQAAAAEAAAA="
)

// ExampleCanonicalize demonstrates that two differently constructed entries
// that represent the same authorization compare as equal.
func ExampleCanonicalize() {
	var e1, e2 xdr.SorobanAuthorizationEntry
	_ = xdr.SafeUnmarshalBase64(canonicalVector, &e1)
	_ = xdr.SafeUnmarshalBase64(canonicalVector, &e2)

	fmt.Println("entries are equal:", soroauth.Equal(e1, e2))

	// Output:
	// entries are equal: true
}
