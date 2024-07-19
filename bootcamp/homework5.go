package main

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/cloudflare/bn256"
)

func main() {
	fmt.Println("curve order: ", bn256.Order)
	verifyOrderAndMod()
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

func bigFromBase10(s string) *big.Int {
	n, _ := new(big.Int).SetString(s, 10)
	return n
}

func verifyOrderAndMod() {
	verifyParamU()
	// u is the BN parameter that determines the prime: 1868033³.
	var u = bigFromBase10("6518589491078791937")

	// Compute p = 36u⁴ + 36u³ + 24u² + 6u + 1
	u2 := new(big.Int).Mul(u, u)
	u3 := new(big.Int).Mul(u2, u)
	u4 := new(big.Int).Mul(u3, u)

	p := new(big.Int).Mul(u4, big.NewInt(36))
	temp := new(big.Int).Mul(u3, big.NewInt(36))
	p.Add(p, temp)
	temp.Mul(u2, big.NewInt(24))
	p.Add(p, temp)
	temp.Mul(u, big.NewInt(6))
	p.Add(p, temp)
	p.Add(p, big.NewInt(1))

	// Print computed modulus p
	fmt.Println("Computed modulus p:", p)

	// Given p value
	givenP := bigFromBase10("65000549695646603732796438742359905742825358107623003571877145026864184071783")
	fmt.Println("Given modulus p:", givenP)
	fmt.Println("Modulus p matches:", p.Cmp(givenP) == 0)

	// Compute order = 36u⁴ + 36u³ + 18u² + 6u + 1
	order := new(big.Int).Mul(u4, big.NewInt(36))
	temp.Mul(u3, big.NewInt(36))
	order.Add(order, temp)
	temp.Mul(u2, big.NewInt(18))
	order.Add(order, temp)
	temp.Mul(u, big.NewInt(6))
	order.Add(order, temp)
	order.Add(order, big.NewInt(1))

	// Print computed order
	fmt.Println("Computed order:", order)

	// Given order value
	givenOrder := bigFromBase10("65000549695646603732796438742359905742570406053903786389881062969044166799969")
	fmt.Println("Given order:", givenOrder)
	fmt.Println("Order matches:", order.Cmp(givenOrder) == 0)
}

func verifyParamU() {
	// Define the known parameter u
	var u = bigFromBase10("6518589491078791937")

	// Compute 1868033³
	base := bigFromBase10("1868033")
	cubed := new(big.Int).Exp(base, big.NewInt(3), nil)

	// Print the results
	fmt.Println("Computed 1868033³:", cubed)
	fmt.Println("Given u:", u)
	fmt.Println("Match:", cubed.Cmp(u) == 0)
}

func homework5() {
	// alpha1 := new(bn256.G1).ScalarBaseMult(big.NewInt(5))
	// beta2 := new(bn256.G2).ScalarBaseMult(big.NewInt(6))

}
