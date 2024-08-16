package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
)

func printCurveParams(curve elliptic.Curve) {
	params := curve.Params()
	fmt.Printf("Curve Name: %s\n", params.Name)
	fmt.Printf("P (Prime Modulus): %s\n", params.P.Text(10))
	fmt.Printf("N (Order): %s\n", params.N.Text(10))
	fmt.Printf("B (Coefficient): %s\n", params.B.Text(10))
	fmt.Printf("Gx (Base Point X): %s\n", params.Gx.Text(10))
	fmt.Printf("Gy (Base Point Y): %s\n", params.Gy.Text(10))
	fmt.Printf("Bit Size: %d\n", params.BitSize)
	fmt.Println(curve.Params().IsOnCurve(curve.Params().Gx, curve.Params().Gy))
}

func gen() {
	// Example private key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Println("Error generating private key:", err)
		return
	}

	// Example message
	message := "Hello, World!"

	// Step 1: Hash the message
	hash := sha256.New()
	hash.Write([]byte(message))
	digest := hash.Sum(nil)

	// Step 2: Generate a random integer k
	k, err := rand.Int(rand.Reader, privateKey.Params().N)
	if err != nil {
		fmt.Println("Error generating random k:", err)
		return
	}

	// Step 3: Calculate the elliptic curve point (x1, y1)
	x1, _ := privateKey.Curve.ScalarBaseMult(k.Bytes())

	// Step 4: Compute r
	r := new(big.Int).Mod(x1, privateKey.Params().N)

	// Ensure r is not zero
	if r.Sign() == 0 {
		fmt.Println("r is zero, regenerate k")
		return
	}

	// Step 5: Compute s
	e := new(big.Int).SetBytes(digest)
	kInv := new(big.Int).ModInverse(k, privateKey.Params().N)
	s := new(big.Int).Mul(r, privateKey.D)
	s.Add(s, e)
	s.Mul(s, kInv)
	s.Mod(s, privateKey.Params().N)

	// Ensure s is not zero
	if s.Sign() == 0 {
		fmt.Println("s is zero, regenerate k")
		return
	}

	// Output the signature
	fmt.Printf("Signature (r, s): (%s, %s)\n", r.Text(10), s.Text(10))

	pubKey := privateKey.Public()

	ecdsaPubKey, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		fmt.Errorf("could not convert to ecdsa.PublicKey")
		return
	}

	fmt.Println(ecdsa.Verify(ecdsaPubKey, digest, r, s))
}

func main() {
	curve := elliptic.P256()
	printCurveParams(curve)
	gen()
}
