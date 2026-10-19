# How soroauth and soroauth-explain divide the problem

This document describes the architectural boundary between **`soroauth`** (this
library and its tooling) and **`soroauth-explain`**, and why interpreting what
an authorization entry means is deliberately not done here.

## The Split: Bytes and Structure vs. Interpretation

- **`soroauth` (Bytes and Structure):** This library builds, signs, inspects,
  and verifies Soroban authorization entries at the XDR and cryptographic
  level. It handles key material, credentials arms (legacy, V2, and delegates),
  preimage computation, payload hashing, and structural inspection
  (`Inspect`, `DescribeSignature`). It knows *what* signed and *what* cryptographic
  shape the signature takes, but it does not know what the contract code
  actually does.
- **`soroauth-explain` (Interpretation and Meaning):** That companion library
  takes a structurally valid authorization entry or invocation tree and interprets
  its meaning — mapping contract IDs to known interfaces, decoding contract
  argument types, explaining token transfers or function semantics, and
  evaluating authorization policies against a human-readable threat or intent
  model.

## Why They Are Separate

1. **Different Honesty Boundaries:** Code that handles private keys, signing
   prehashes, and XDR marshalling must maintain a minimal, auditable attack
   surface with zero external runtime dependencies for contract metadata, ABI
   parsing, or network RPC lookups. Mixing cryptographic signing with dynamic
   contract interpretation conflates low-level security-critical code with
   heuristic protocol semantics.
2. **Different Failure Consequences:** A failure in `soroauth` results in an
   invalid signature or malformed XDR bytes that fail on-chain or get rejected
   by the host. A failure or heuristic error in explanation/interpretation
   could mislead a human user into approving a malicious transaction. Keeping
   interpretation separate ensures that the core signing library never
   over-promises semantic understanding.

## Concrete Comparison: One Entry Through Both

Consider an authorization entry for a token transfer contract call (`transfer(from, to, amount)`):

- **`soroauth` (`Inspect`) reports:**

  - Credential type: `address_v2`
  - Address bound: `true`
  - Address: `GA3...`
  - Nonce: `42`
  - Valid until ledger: `1234567`
  - Root contract: `CC4...`
  - Root function: `transfer`
  - Sub-invocations: `0`
  - Top-level signed: `true` (64-byte binary passkey signature)

- **`soroauth-explain` reports:**

  - Intent: "Transfer 100.00 USDC from `GA3...` to `GB7...` via Stellar Asset Contract `CC4...`"
  - Security advisory: "This contract is verified on Stellar Expert; recipient `GB7...` is a known exchange address."

Neither library tries to do the other's job. `soroauth` guarantees the bytes
and signature are correct; `soroauth-explain` tells you what those bytes mean for
a human reviewer.
