# How soroauth and soroauth-explain Divide the Problem

This document describes the architectural split between **soroauth** (this repository) and **soroauth-explain**, defining which library answers which question and why interpreting an authorization entry's policy or semantic meaning is deliberately kept out of this library.

## The Division of Labor

| Concern | `soroauth` | `soroauth-explain` |
| :--- | :--- | :--- |
| **Domain** | Bytes, XDR, cryptographic preimages, structural inspection, and signing | Semantic interpretation, authorization policy, and contract intent |
| **Core Question** | *Are these bytes well-formed, cryptographically sound, and correctly signed?* | *What does this authorization actually mean in plain English or business logic?* |
| **Inputs** | `SorobanAuthorizationEntry`, network passphrases, raw keys, and signer adapters | Decoded authorization structures and contract invocation trees |
| **Outputs** | Signed XDR entries, hex payloads, cryptographic error verdicts, and structural `EntryInfo` summaries | Human-readable policy summaries, parameter interpretations, and capability declarations |

### Why Bytes and Structure Live Here

`soroauth` produces and verifies signatures that move real assets on the Stellar network. Its responsibility is strictly cryptographic and structural honesty: adhering precisely to CAP-46-11 and CAP-71 specifications, maintaining binary reproducibility via immutable golden vectors, and ensuring that no signature can be forged, replayed, or misbound.

### Why Interpretation Lives Elsewhere

Interpreting what a contract function argument means (e.g., whether `ScVal` index `0` represents a token amount, an expiration timestamp, or a recipient address) requires domain-specific contract knowledge, ABI definitions, or heuristic analyzers. Mixing interpretation logic into `soroauth` would:

1. **Blur Honesty Boundaries:** Cryptographic verification must be mechanical and unambiguous. Semantic heuristics are subjective and prone to false interpretations.
2. **Inflate Failure Consequences:** A bug in cryptographic preimage generation loses funds. A bug in semantic explanation should never corrupt transaction assembly or signing flows.

## Concrete Walkthrough: An Entry Through Both

Consider a Soroban authorization entry where an account authorizes a transfer of 100 USDC.

1. **Through `soroauth` (`Inspect` / CLI):**
   - Reports structural summary: Credential type `address_v2`, address `GCEZ...`, nonce `42`, valid until ledger `86400`, top-level signed `true`, and root invocation function `transfer`. It exposes raw `ScVal` vectors and positional argument types without guessing their business meaning.

2. **Through `soroauth-explain`:**
   - Consumes the structural output or raw entry and maps contract identifiers or argument types to known interfaces (e.g., Soroban Token standard), interpreting argument `0` as recipient `GB7...` and argument `1` as amount `100,000,000` (100 USDC), rendering a clear human-readable policy summary.
