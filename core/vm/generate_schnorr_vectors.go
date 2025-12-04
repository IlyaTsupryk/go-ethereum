//go:build ignore
// +build ignore

package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
)

type testVector struct {
	Input       string `json:"Input"`
	Expected    string `json:"Expected"`
	Gas         uint64 `json:"Gas"`
	Name        string `json:"Name"`
	NoBenchmark bool   `json:"NoBenchmark"`
}

func main() {
	vectors := []testVector{}

	// Generate a few valid signatures
	for i := 0; i < 3; i++ {
		// Generate a random private key
		privKey, err := btcec.NewPrivateKey()
		if err != nil {
			panic(err)
		}

		// Get the public key (x-only)
		pubKey := privKey.PubKey()
		pubKeyBytes := pubKey.SerializeCompressed()[1:] // Remove the prefix byte for x-only

		// Create a message hash
		message := sha256.Sum256([]byte(fmt.Sprintf("test message %d", i)))

		// Sign the message
		sig, err := schnorr.Sign(privKey, message[:])
		if err != nil {
			panic(err)
		}

		// Construct the input: hash (32) + signature (64) + pubkey (32) = 128 bytes
		input := hex.EncodeToString(message[:]) + hex.EncodeToString(sig.Serialize()) + hex.EncodeToString(pubKeyBytes)

		vectors = append(vectors, testVector{
			Input:       input,
			Expected:    "0000000000000000000000000000000000000000000000000000000000000001",
			Gas:         3000,
			Name:        fmt.Sprintf("Valid signature test %d", i),
			NoBenchmark: false,
		})
	}

	// Add some invalid signature tests
	// Test with wrong signature
	{
		privKey, _ := btcec.NewPrivateKey()
		pubKey := privKey.PubKey()
		pubKeyBytes := pubKey.SerializeCompressed()[1:]
		message := sha256.Sum256([]byte("test message"))
		sig, _ := schnorr.Sign(privKey, message[:])

		// Modify the signature to make it invalid
		sigBytes := sig.Serialize()
		sigBytes[0] ^= 0xFF
		invalidSig := hex.EncodeToString(sigBytes)

		input := hex.EncodeToString(message[:]) + invalidSig + hex.EncodeToString(pubKeyBytes)
		vectors = append(vectors, testVector{
			Input:       input,
			Expected:    "",
			Gas:         3000,
			Name:        "Invalid signature - modified signature",
			NoBenchmark: false,
		})
	}

	// Test with wrong public key
	{
		privKey, _ := btcec.NewPrivateKey()
		message := sha256.Sum256([]byte("test message"))
		sig, _ := schnorr.Sign(privKey, message[:])

		// Use a different public key
		wrongPrivKey, _ := btcec.NewPrivateKey()
		wrongPubKey := wrongPrivKey.PubKey()
		wrongPubKeyBytes := wrongPubKey.SerializeCompressed()[1:]

		input := hex.EncodeToString(message[:]) + hex.EncodeToString(sig.Serialize()) + hex.EncodeToString(wrongPubKeyBytes)
		vectors = append(vectors, testVector{
			Input:       input,
			Expected:    "",
			Gas:         3000,
			Name:        "Invalid signature - wrong public key",
			NoBenchmark: false,
		})
	}

	// Test with wrong message
	{
		privKey, _ := btcec.NewPrivateKey()
		pubKey := privKey.PubKey()
		pubKeyBytes := pubKey.SerializeCompressed()[1:]
		message := sha256.Sum256([]byte("test message"))
		sig, _ := schnorr.Sign(privKey, message[:])

		// Use a different message
		wrongMessage := sha256.Sum256([]byte("wrong message"))

		input := hex.EncodeToString(wrongMessage[:]) + hex.EncodeToString(sig.Serialize()) + hex.EncodeToString(pubKeyBytes)
		vectors = append(vectors, testVector{
			Input:       input,
			Expected:    "",
			Gas:         3000,
			Name:        "Invalid signature - wrong message",
			NoBenchmark: false,
		})
	}

	// Test with invalid input lengths
	vectors = append(vectors, testVector{
		Input:       "",
		Expected:    "",
		Gas:         3000,
		Name:        "Empty input",
		NoBenchmark: false,
	})

	vectors = append(vectors, testVector{
		Input:       "00",
		Expected:    "",
		Gas:         3000,
		Name:        "Too short input (1 byte)",
		NoBenchmark: false,
	})

	// Create 127 bytes (one byte short)
	shortInput := make([]byte, 127)
	rand.Read(shortInput)
	vectors = append(vectors, testVector{
		Input:       hex.EncodeToString(shortInput),
		Expected:    "",
		Gas:         3000,
		Name:        "Too short input (127 bytes)",
		NoBenchmark: false,
	})

	// Create 129 bytes (one byte too long)
	longInput := make([]byte, 129)
	rand.Read(longInput)
	vectors = append(vectors, testVector{
		Input:       hex.EncodeToString(longInput),
		Expected:    "",
		Gas:         3000,
		Name:        "Too long input (129 bytes)",
		NoBenchmark: false,
	})

	// Marshal to JSON
	jsonData, err := json.MarshalIndent(vectors, "", "  ")
	if err != nil {
		panic(err)
	}

	// Write to file
	err = os.WriteFile("testdata/precompiles/schnorrVerify.json", jsonData, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println("Generated schnorrVerify.json test vectors successfully!")
}
