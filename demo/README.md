# Schorr signature verification demo

A simple Ethereum smart contract demo for storing and retrieving strings per user.

## Prerequisites

- Node.js and npm
- Local Ethereum node running at `http://127.0.0.1:8545`

## Setup

```bash
npm install
```

## Usage

**Compile contracts:**
```bash
npx hardhat compile
```

**Deploy contract:**
```bash
npm run deploy
```

**A demo script showing how strings can be added and retrieved:**  
Set CONTRACT_ADDRESS and run:
```bash
npm run add-strings
```

## Contract

`StringStorage.sol` - Stores string arrays per user address.

- `addString(string)` - Add a string for msg.sender
- `getString(address, uint256)` - Retrieve string by user and index