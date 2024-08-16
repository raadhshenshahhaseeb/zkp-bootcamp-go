package main

import "crypto/sha256"

func verifier(commitments []ECPoint) {
	u := sha256.Sum256([]byte("secret 1"))
}
