# Schnorr Signature Verification Demo

Demonstrates Schnorr signature verification using a custom precompiled contract in go-ethereum.

## Features

- **Custom Precompile**: Schnorr signature verification at address `0x0101`
- **Smart Contract**: StringStorage with signature-verified string addition
- **Demo Scripts**: Examples showing both standard and signature-verified operations

## Prerequisites

- Node.js and npm
- Custom go-ethereum node with Schnorr precompile running at `http://127.0.0.1:8545`

## Setup

```bash
npm install
```

## Usage

**1. Compile contracts:**
```bash
npx hardhat compile
```

**2. Deploy contract:**
```bash
npm run deploy
```

**3. Run demos:**

Standard string addition (no signature verification):
```bash
npm run add-strings
```

Secure string addition with Schnorr signature verification:
```bash
npm run add-string-safe
```

## Contract

**StringStorage.sol** - Stores string arrays per user address with optional Schnorr verification.

### Methods

- `addString(string)` - Add a string for msg.sender (no verification)
- `addStringSecure(string, bytes signature, bytes pubkey)` - Add a string with Schnorr signature verification
- `getString(address, uint256)` - Retrieve string by user and index
- `verifySchnorrSignature(bytes32, bytes, bytes)` - Internal Schnorr verification using precompile at `0x0101`