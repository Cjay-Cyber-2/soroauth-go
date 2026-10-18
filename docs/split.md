# How soroauth and soroauth-explain divide the problem

This library (`soroauth` / `soroauth-go`) and its companion library `soroauth-explain` occupy distinct roles in handling Soroban authorization entries. Understanding this boundary is essential for working with either tool effectively.

## The Split: Bytes and Structure Here, Interpretation There

- **`soroauth` (this library)**: Deals entirely with **bytes, structure, preimages, and signatures**. It constructs, signs, verifies, inspects, and decodes `SorobanAuthorizationEntry` XDR structures and their credential arms. It tells you *who* approved *what* call tree, what cryptographic preimage was signed, and whether the signature is structurally valid.
- **`soroauth-explain`**: Deals entirely with **meaning and interpretation**. It takes a decoded authorization entry and interprets what the contract call actually does in human-readable terms (e.g., explaining that a specific invocation transfers 100 native XLM tokens, changes a policy parameter, or revokes a role).

## Why They Are Separate

1. **Different Honesty Boundaries**: 
   `soroauth` operates under strict cryptographic and protocol honesty boundaries. It performs no network calls, maintains no contract ABI registries, and makes no assumptions about what a contract function means. If a contract function is named `execute_malicious_drain`, `soroauth` reports that exact function name and leaves the judgment to the caller. Mixing interpretation logic into a signature-generation core would expand the trust and attack surface unnecessarily.

2. **Different Failure Consequences**:
   A bug in `soroauth` can lead to invalid signatures, failed transactions, or lost funds. A bug or limitation in an explainer library might produce an unclear or incomplete human description, which is a documentation or UX issue rather than a cryptographic failure.

## Concrete Example: One Entry Through Both

Consider an authorization entry for a token transfer:

### 1. What `soroauth` (`Inspect`) Sees

`soroauth` reports the structural facts:
```json
{
  "credential_type": "address_v2",
  "address_bound": true,
  "address": "GCEZ…",
  "nonce": 42,
  "valid_until_ledger": 1234570,
  "top_level_signed": true,
  "root_contract": "CC…",
  "root_function": "transfer",
  "sub_invocations": 0
}
```
It knows the credential arm, the signer address, the nonce, the ledger expiration, and that the root contract function is `transfer` with opaque XDR arguments.

### 2. What `soroauth-explain` Sees

Using contract metadata, symbol registries, or ABIs, `soroauth-explain` interprets the entry:
> **Authorized Action**: `GCEZ…` authorizes contract `CC…` to execute `transfer(to: GBUY…, amount: 10000000)`. This moves `1.0000000 XLM` from the signer's account.

## Linked References

- See the primary documentation index in [`docs/README.md`](README.md).
- See [`soroauth-explain`](https://github.com/soroauth/soroauth-explain) for interpretation features.
