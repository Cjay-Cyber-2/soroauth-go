# Contributing a new signer adapter

This guide explains how to add a new signer adapter to `soroauth` (such as a hardware security module, remote signing service, or cloud KMS). There are currently a dozen signers across several modules, and following established patterns ensures consistency, safety, and correctness.

## 1. Module Layout & `go.mod` Placement

Signer adapters that pull in heavy third-party SDK dependencies (like AWS KMS, Google Cloud KMS, or Vault) must live in their own out-of-tree directories under the `adapters/` top-level directory (e.g., `adapters/awskms`, `adapters/gcpkms`).

- **Module Isolation:** Each adapter is its own Go module with its own `go.mod` and `go.sum`. This keeps the root `soroauth` package lightweight and free of heavy cloud or vendor SDKs.
- **Import Paths:** The module path should follow the convention `github.com/soroauth/soroauth-go/adapters/<name>`.
- **CI Integration:** Every new adapter must be wired into the CI adapters workflow (`.github/workflows/ci-adapters.yml`) so that its `gofmt`, `go vet`, and `go test -race` checks run automatically.

## 2. The `Signer` Contract & Conformance Suite

Every signer must implement the `soroauth.Signer` interface:

```go
type Signer interface {
    Address() string
    Sign(ctx context.Context, preimage xdr.HashIdPreimage, payload [32]byte) (xdr.ScVal, error)
}
```

### Conformance
Do not write an ad-hoc test suite for adapter mechanics. Use the repository's dedicated conformance suite located in `signertest` to validate that your signer meets all contract expectations, error behavior, and signature wrapping rules.

## 3. Recurring Traps and How to Avoid Them

### A. Typed-Nil Clients
**The Trap:** A caller passes a nil pointer wrapped in an interface (e.g., `var client *kms.KeyManagementClient` passed as an interface parameter). In Go, `client == nil` evaluates to `false` because the interface holds a non-nil type descriptor with a nil value. This leads to runtime panics later during `Sign`.

**The Fix:** Use reflection-based checks (as demonstrated in `adapters/gcpkms` and `adapters/awskms`) to detect typed-nil interface values during construction and refuse them immediately with `soroauth.ErrMissingSigner`.

### B. Never Logging Secrets
**The Trap:** Debug statements, error paths, or telemetry logging might accidentally record raw private keys, seed bytes, or payload signatures.

**The Fix:** Never log, print, or expose private key material, raw seeds, or sensitive signing inputs. Rely exclusively on public keys, addresses, and opaque error wrappers.

### C. Context Handling
**The Trap:** Initiating network or hardware calls without honoring context cancellation or timeouts.

**The Fix:** Always check `ctx.Err()` at the beginning of `Sign()` before dispatching external client requests.

## 4. Evidence and Documentation Requirements in `docs/signers.md`

When you introduce a new signer adapter, you must update `docs/signers.md`:

- Add the signer to the summary table specifying its threat model, evidence level, and caveats.
- Be explicit about live e2e test status: an untested hardware or remote claim is not acceptable. If there is no live testnet scenario proving the path in this repo, state clearly that operational review and contract-specific verification are required.
