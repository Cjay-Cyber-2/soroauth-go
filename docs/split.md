# How soroauth and soroauth-explain Divide the Problem

`soroauth` and `soroauth-explain` solve two complementary halves of Soroban authorization security. This document explains the exact boundary between them, why they are separate repositories/libraries, and what each is responsible for.

## The Split

- **`soroauth` (this library): Bytes and Structure.**
  It parses, validates, signs, diffs, and inspects XDR authorization entries at the byte and structure level. It guarantees that an entry is well-formed, that its signature is mathematically valid for the given network passphrase and signer address, and that its delegate trees and expiration ledgers follow the rules of CAP-71-01.

- **`soroauth-explain`: Interpretation and Meaning.**
  It answers *what an authorization entry actually does in human terms*—translating function arguments, invoking contract methods, inspecting token transfers, and explaining the contract business logic of an authorization payload.

## Why They Are Separate

1. **Different Honesty Boundaries:**
   - `soroauth` is a security-critical signing library. Its correctness relies on minimal dependencies, strict adherence to XDR specs, cryptographic correctness, and zero ambiguity. It must be auditable and reproducible down to the exact byte.
   - `soroauth-explain` deals with human-readable interpretation, contract ABIs, heuristic formatting, and rich explanations. These change rapidly as ecosystem tooling evolves and contract standards shift.

2. **Different Failure Consequences:**
   - If `soroauth` misreads structure or bytes, transactions fail to sign, signatures are rejected, or funds are locked. Safety invariants are absolute.
   - If `soroauth-explain` misinterprets a human-readable label or formats an argument unclearly, the user sees an imperfect summary, but the cryptographic authorization remains intact and secure.

## Concrete Example: One Entry, Two Lenses

Consider an incoming Soroban authorization entry invoking `transfer` on a token contract.

1. **Through `soroauth inspect` (this library):**
   Reports the structure:
   - **Type:** `SorobanAuthorizationEntry` (v6)
   - **Credentials:** Address `GABC...`, nonce `42`, signature valid.
   - **Root Invocation:** Contract ID `C123...`, function `transfer`, argument count: 3.
   - **Signatures / Tree:** Valid signature over the hash preimage, valid ledger bounds (`validUntilLedger: 5000000`).

2. **Through `soroauth-explain`:**
   Interprets the meaning:
   - "This entry authorizes account `GABC...` to transfer **100 USDC** to account `GXYZ...` via contract `C123...` (USDC Token Contract), expiring at ledger ~5,000,000."

## Summary

`soroauth` tells you **that** an entry is valid and signed correctly. `soroauth-explain` tells you **what** it asks you to approve.
