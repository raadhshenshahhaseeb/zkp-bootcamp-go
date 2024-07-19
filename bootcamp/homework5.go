package main

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/cloudflare/bn256"
)

func main() {
	basics()
	pairing()
	finalExponentiate()
}

func basics() {
	a, _ := rand.Int(rand.Reader, bn256.Order)

	g2a := new(bn256.G2).ScalarBaseMult(a)

	A := new(bn256.G2).ScalarMult(g2a, big.NewInt(4))
	B := new(bn256.G2).ScalarMult(g2a, big.NewInt(8))

	Y := new(bn256.G2).Add(A, B)
	Z := new(bn256.G2).ScalarMult(g2a, big.NewInt(12))
	X := new(bn256.G2).ScalarMult(A, big.NewInt(3))

	fmt.Println(Y)
	fmt.Println(Z)
	fmt.Println(X)

	if X.String() == Y.String() {
		fmt.Println("equals")
	}

	alpha := new(bn256.G2).ScalarMult(g2a, big.NewInt(0).Add(big.NewInt(12), bn256.Order))

	fmt.Println(alpha)
}

func pairing() {
	three := new(bn256.G2).ScalarBaseMult(big.NewInt(3))
	four := new(bn256.G1).ScalarBaseMult(big.NewInt(4))

	two := new(bn256.G2).ScalarBaseMult(big.NewInt(2))
	six := new(bn256.G1).ScalarBaseMult(big.NewInt(6))

	pairThreeFour := bn256.Pair(four, three)
	pairTwoSix := bn256.Pair(six, two)

	fmt.Println("pairThreeFour with length:")
	fmt.Println(pairThreeFour)
	fmt.Println("pairTwoSix with length:")
	fmt.Println(pairTwoSix)

	if pairThreeFour.String() == pairTwoSix.String() {
		fmt.Println("equals")
	}

	neg2 := new(bn256.G2).Neg(two)
	pairNegTwoSix := bn256.Pair(six, neg2)

	identity := new(bn256.GT).ScalarMult(pairThreeFour, big.NewInt(0))
	result := new(bn256.GT).Add(pairThreeFour, pairNegTwoSix)

	fmt.Println("Result of multiplication (should be identity):")
	fmt.Println(result)
	fmt.Println("Identity element:")
	fmt.Println(identity)

	if result.String() == identity.String() {
		fmt.Println("Result is the identity element")
	} else {
		fmt.Println("Result is not the identity element")
	}

	// Verify that adding the identity element to pairThreeFour results in pairThreeFour
	pairingThreeFourWithIdentity := new(bn256.GT).Add(pairThreeFour, identity)

	if pairingThreeFourWithIdentity.String() == pairThreeFour.String() {
		fmt.Println("pairingThreeFourWithIdentity + pairThreeFour are equal")
		fmt.Println(pairingThreeFourWithIdentity)
		fmt.Println(pairThreeFour)
	}
}

func finalExponentiate() {
	A := new(bn256.G2).ScalarBaseMult(big.NewInt(6))
	B := new(bn256.G1).ScalarBaseMult(big.NewInt(2))
	C := new(bn256.G2).ScalarBaseMult(big.NewInt(4))
	D := new(bn256.G1).ScalarBaseMult(big.NewInt(2))
	E := new(bn256.G2).ScalarBaseMult(big.NewInt(2))
	F := new(bn256.G1).ScalarBaseMult(big.NewInt(2))

	NegB := new(bn256.G1).Neg(B)

	pairingAB := bn256.Pair(NegB, A).Finalize()
	pairingCD := bn256.Pair(D, C).Finalize()
	pairingEF := bn256.Pair(F, E).Finalize()

	result := new(bn256.GT).Add(pairingAB, pairingCD)
	result = new(bn256.GT).Add(result, pairingEF)

	identity := new(bn256.GT).ScalarBaseMult(big.NewInt(0))

	fmt.Println("Result of multiplication (should be identity):")
	fmt.Println(result)
	fmt.Println("Identity element:")
	fmt.Println(identity)

	if result.String() == identity.String() {
		fmt.Println("Result is the identity element")
	} else {
		fmt.Println("Result is not the identity element")
	}
}

func homework5() {
	// alpha1 := new(bn256.G1).ScalarBaseMult(big.NewInt(5))
	// beta2 := new(bn256.G2).ScalarBaseMult(big.NewInt(6))

}
