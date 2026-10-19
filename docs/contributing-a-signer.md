# Contributing a new signer

This guide describes the complete process for adding a new signer adapter to
`soroauth`. Whether you are integrating a hardware wallet, a cloud key
management service, or an external signing proxy, all adapters must follow a
consistent design, contract, and verification standard.

## 1. Module layout and `go.mod` placement

Signer adapters that pull in heavy third-party SDKs (such as AWS SDK, GCP SDK,
or wallet SDKs) **must not** bloat the root `github.com/soroauth/soroauth-go`
module. Instead, they reside in isolated subdirectories under the `adapters/`
folder:

- Directory: `adapters/<name>/` (e.g., `adapters/awskms`, `adapters/gcpkms`,
  `adapters/walletsdk`).
- Module: Each adapter is its own Go module with its own `go.mod` file (e.g.,
  `module github.com/soroauth/soroauth-go/adapters/gcpkms`).
- CI: Every adapter must be included in the adapters CI workflow
  (`.github/workflows/ci-adapters.yml`), which validates formatting (`gofmt`),
  runs `go vet ./...`, and executes tests with the race detector
  (`go test -race ./...`).

## 2. The `Signer` contract and conformance suite

Every adapter must implement the `soroauth.Signer` interface defined in the
root package (`signer.go`):

```go
type Signer interface {
    Address() string
    Sign(ctx context.Context, preimage xdr.HashIdPreimage, payload [32]byte) (xdr.ScVal, error)
}
```

### Contract Rules:

1. **Address()**: Returns the classic Stellar account address (`G…`) or target
   identifier that the signer manages.
2. **Sign(...)**: Takes a context, the hash ID preimage (for inspection/logging
   if required, though typically unused by KMS signers), and the exact 32-byte
   SHA-256 payload. It returns a signed `xdr.ScVal` (for classic accounts,
   typically created via `soroauth.Ed25519SignatureScVal`).

### Conformance Suite:

Adapters should be tested against the standard conformance expectations used
across existing signers:

- Rejection of empty/invalid inputs and mismatched addresses.
- Local signature verification before returning the `ScVal` (ensuring the
  private key actually signed the expected payload).
- Propagating upstream errors correctly.
- Respecting context cancellation and timeouts.

## 3. Recurring traps and how to avoid them

When writing a signer adapter, several subtle bugs are common. Study the
existing adapters for exact implementations:

### Trap 1: Typed-nil clients

A caller might pass a nil pointer wrapped in an interface (e.g.,
`var client *kms.KeyManagementClient` passed as an interface). A simple
`client == nil` check evaluates to **false** because the interface contains a
type descriptor.

- **The Fix:** Use a reflection helper to check for typed nils, as demonstrated
  in `adapters/gcpkms` and `adapters/awskms` (`isNilClient`). Refuse typed-nil
  clients at construction time rather than panicking later during `Sign()`.

### Trap 2: Context handling

Always check `ctx.Err()` at the beginning of `Sign()` before dispatching network
calls or hardware operations, and pass the context through to the underlying
client.

- **Example:** See `adapters/gcpkms/signer.go` and `adapters/awskms/signer.go` for
  proper context checks
  (`if err := ctx.Err(); err != nil { return ..., err }`).

### Trap 3: Never logging secrets

Private keys, seeds, raw private key bytes, or authorization tokens must never
be logged or included in error messages. Error messages should wrap sentinels
like `soroauth.ErrSignatureMismatch` or `soroauth.ErrMissingSigner` without
exposing sensitive material.

## 4. Documentation and evidence in `docs/signers.md`

Adding a new signer requires updating `docs/signers.md`:

- **Table Entry:** Add the new signer to the shipped signer set table with its
  threat model, evidence level, and recommendation.
- **Threat Model & Evidence:** State clearly whether the signer is proven by
  live e2e testnet scenarios or only unit-tested against mocks/fakes. Unproven
  hardware or cloud claims must explicitly note the lack of live e2e
  verification in this repository.
