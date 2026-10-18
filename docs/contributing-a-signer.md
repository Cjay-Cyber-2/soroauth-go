# Contributing a new signer

This guide describes the complete process for adding a new signer adapter to
`soroauth`: the module layout, the contract, the conformance suite, test
expectations, recurring security/correctness traps, and the documentation
evidence required in `docs/signers.md`.

## 1. Module Layout & Placement

Signer adapters that pull in heavy cloud or third-party SDK dependencies (such
as AWS KMS, GCP KMS, or external wallet SDKs) **must** reside in their own
dedicated subdirectories under the `adapters/` top-level directory (e.g.,
`adapters/awskms`, `adapters/gcpkms`, `adapters/walletsdk`).

- **Separate Go Module:** Each adapter directory must be an independent Go
  module with its own `go.mod` and `go.sum`. This keeps the root `soroauth-go`
  module lightweight and free of heavy cloud SDK transitive dependencies.
- **CI Integration:** Every new adapter must have a corresponding job added to
  `.github/workflows/ci-adapters.yml` that runs `gofmt`, `go vet`, and
  `go test -race ./...` using its own module path and go-version file.

## 2. The `Signer` Contract

Every signer adapter must implement the core `soroauth.Signer` interface defined
in `signer.go`:

```go
type Signer interface {
    Address() string
    Sign(ctx context.Context, preimage xdr.HashIdPreimage, payload [32]byte) (xdr.ScVal, error)
}
```

- **`Address()`**: Returns the string-encoded account address (starting with
  `G…`) or target identifier that this signer represents.
- **`Sign()`**: Takes a context, the signing hash ID preimage, and the exact
  32-byte SHA-256 payload. It must perform the signing operation, verify the
  signature locally before returning, and format it into the appropriate
  classic account signature `ScVal` (or custom `ScVal` where applicable).

## 3. Conformance Suite & Testing Expectations

Adapters must not rely solely on self-signed round-trips. They should be tested
against:

- **Faked/Mock Clients:** Use dependency interfaces (such as `AsymmetricSigner`
  in `gcpkms` or `SignerClient` in `awskms`) so that unit tests exercise the
  entire signing path against test doubles without needing live cloud
  infrastructure.
- **Local Verification:** Before returning an `ScVal`, the adapter must verify
  the signature locally using crypto primitives (e.g., `ed25519.Verify`). If a
  signature does not verify, return `soroauth.ErrSignatureMismatch` immediately
  rather than propagating invalid bytes onto the blockchain.

## 4. Recurring Traps & Defensive Fixes

When writing a signer adapter, watch out for these critical failure modes:

1. **Typed-Nil Client Pointers**
   - *The Trap:* A caller passing a typed-nil pointer (e.g., `var client *kms.KeyManagementClient`)
     results in an interface holding a non-nil type with a nil concrete value
     (`client != nil` evaluates to `true`). Without a deep reflection check,
     the constructor succeeds and panics later during `Sign`.
   - *The Fix:* Use a helper like `isNilClient` (demonstrated in
     `adapters/gcpkms/signer.go` and `adapters/awskms/signer.go`) to reject
     typed-nil clients at construction time with `soroauth.ErrMissingSigner`.

2. **Context Handling**
   - *The Trap:* Forgetting to check `ctx.Err()` before invoking external KMS or
     network APIs, leading to wasted round-trips or blocked cancellations.
   - *The Fix:* Always check `if err := ctx.Err(); err != nil { return ..., err }`
     at the very beginning of `Sign`.

3. **Never Logging Secrets**
   - *The Trap:* Accidentally printing private keys, seed bytes, or raw signing
     material in error messages or logs.
   - *The Fix:* Only log or wrap public metadata such as key versions, resource
     IDs, or addresses. Never log payload bytes or keys.

## 5. Documentation Evidence in `docs/signers.md`

Every new signer adapter must be documented in `docs/signers.md`:

- **Entry in the Shipped Signer Table:** Add your adapter to the summary table
  with its threat model, evidence level, and usage recommendations.
- **Honest Evidence Bar:** If your signer relies on unit tests and fakes without
  a live e2e testnet scenario in `e2e/`, explicitly state that caveat in
  `docs/signers.md`. Untested hardware or cloud claims are not acceptable.
