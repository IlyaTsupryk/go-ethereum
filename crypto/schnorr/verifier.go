// Package schnorr implements BIP-340 Schnorr signature verification using the btcsuite library.
package schnorr

import (
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
)

// Verify verifies a BIP-340 Schnorr signature using the btcsuite library.
// It takes a 32-byte message hash, 64-byte signature, and 32-byte x-only public key.
// Returns true if the signature is valid, false otherwise.
func Verify(hash, sig, pubkey []byte) bool {
	signature, err := schnorr.ParseSignature(sig)
	if err != nil {
		return false
	}

	publicKey, err := schnorr.ParsePubKey(pubkey)
	if err != nil {
		return false
	}

	return signature.Verify(hash, publicKey)
}
