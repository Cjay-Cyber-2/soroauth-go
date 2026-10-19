# How soroauth and soroauth-explain Divide the Problem

This document outlines the architectural boundary between `soroauth` and `soroauth-explain`. It explains which library answers which question, why they are deliberately kept separate, and walks through a concrete authorization entry inspected by both.

## The Split: Bytes and Structure Here, Interpretation There

- **`soroauth` (This Library)**: Answers questions about **bytes, structure, cryptographic preimages, and protocol validity**. It builds, signs, inspects, and verifies XDR-encoded `SorobanAuthorizationEntry` instances against consensus rules, credential arms, and cryptographic signatures. It deals with exact wire formats, hash id preimages, nonces, and expiration ledgers.
- **`soroauth-explain`**: Answers questions about **meaning, intent, and human-readable semantic interpretation**. It takes a structurally valid authorization entry and translates what the contract invocation tree actually means for the resource owners — what assets are transferred, what privileges are granted, and what policy limits apply.

### Why Interpretation is Not Done Here

`soroauth` deliberately refuses to interpret the semantic meaning of contract arguments, ScVal payloads, or custom contract methods. This is driven by starkly different honesty boundaries and failure consequences:

1. **Different Honesty Boundaries**: A wire and cryptographic library must be strictly deterministic, protocol-faithful, and zero-trust regarding transaction semantics. Its correctness depends solely on matching Stellar Core and protocol specifications byte-for-byte. Semantic interpretation, on the other hand, depends on contract ABIs, off-chain registries, heuristic analysis, and business logic that can evolve independently of consensus rules.
2. **Different Failure Consequences**: If `soroauth` fails, a signature is invalid or rejected by the network (a clear cryptographic rejection). If a semantic interpreter misinterprets contract intent, a user might be misled into approving a malicious transaction thinking it does something benign. Conflating the two risks introducing semantic heuristics into cryptographic validation or vice-versa.

---

## Concrete Example: One Entry Through Both

Consider a Soroban authorization entry where a passkey smart wallet approves a token transfer contract call.

### 1. What `soroauth` Reports (`soroauth inspect`)

`soroauth` decodes the XDR structure, validates credential arms, and reports the raw shape without guessing what the contract function does:

```json
{
  "credential_type": "address_v2",
  "address_bound": true,
  "address": "GA...",
  "nonce": 42,
  "valid_until_ledger": 1298400,
  "top_level_signed": true,
  "root_contract": "CD...",
  "root_function": "transfer",
  "sub_invocations": 0
}
```

It knows that address `GA...` authorized function `transfer` on contract `CD...` with nonce `42` expiring at ledger `1298400`, carrying a 64-byte passkey signature. It does **not** know what token is being transferred, what amount `42` means, or whether the recipient is trusted.

### 2. What `soroauth-explain` Reports (`soroauth-explain`)

`soroauth-explain` takes that same entry, inspects contract metadata/ABIs, and produces a semantic breakdown:

```json
{
  "summary": "Approve transfer of 100.0 USDC from GA... to GB...",
  "intent": "token_transfer",
  "risk_level": "low",
  "details": {
    "asset": "USDC (CC...)",
    "amount": "100000000",
    "recipient": "GB..."
  }
}
```

---

## References

- [`soroauth` Architecture](../ARCHITECTURE.md)
- [`soroauth` Glossary](glossary.md)
