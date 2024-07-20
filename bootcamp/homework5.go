package main

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/cloudflare/bn256"
)

func main() {
	homework5()
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
	fmt.Println("curve order: ", bn256.Order)
	p := bigFromBase10("65000549695646603732796438742359905742825358107623003571877145026864184071783")
	fmt.Println("prime p: ", p)

	beta2 := new(bn256.G2).ScalarBaseMult(big.NewInt(1))
	alfa1 := new(bn256.G1).ScalarBaseMult(big.NewInt(2))
	gamma2 := new(bn256.G2).ScalarBaseMult(big.NewInt(2))
	delta2 := new(bn256.G2).ScalarBaseMult(big.NewInt(2))

	labels := []string{"beta2", "gamma2", "delta2"}
	verifyAndReduceConversion(alfa1, []*bn256.G2{beta2, gamma2, delta2}, labels, p)
}

// pointToBigInt converts a *bn256.G1 point to big.Int representations
func pointToBigInt(g1 *bn256.G1) (*big.Int, *big.Int) {
	x, y := new(big.Int), new(big.Int)
	bytes := g1.Marshal()
	x.SetBytes(bytes[:32])
	y.SetBytes(bytes[32:])
	return x, y
}

// pointToBigIntG2 converts a *bn256.G2 point to big.Int representations
func pointToBigIntG2(g2 *bn256.G2) (*big.Int, *big.Int, *big.Int, *big.Int) {
	x0, x1, y0, y1 := new(big.Int), new(big.Int), new(big.Int), new(big.Int)
	bytes := g2.Marshal()
	x0.SetBytes(bytes[:32])
	x1.SetBytes(bytes[32:64])
	y0.SetBytes(bytes[64:96])
	y1.SetBytes(bytes[96:])
	return x0, x1, y0, y1
}

// modP reduces a big.Int value to fit within the prime p using modulo operation
func modP(value *big.Int, p *big.Int) *big.Int {
	return new(big.Int).Mod(value, p)
}

// verifyAndReduceConversion checks if the conversion from hex to uint format is correct and applies reduction if necessary
func verifyAndReduceConversion(g1 *bn256.G1, g2 []*bn256.G2, labels []string, p *big.Int) {
	x1, y1 := pointToBigInt(g1)

	// Reduce and print G1 point
	x1 = modP(x1, p)
	y1 = modP(y1, p)
	fmt.Printf("G1 Point: (x: %s, y: %s)\n", x1.String(), y1.String())

	// Loop through the array of G2 points
	for i, point := range g2 {
		x0, x1G2, y0, y1G2 := pointToBigIntG2(point)

		// Reduce each coordinate
		x0 = modP(x0, p)
		x1G2 = modP(x1G2, p)
		y0 = modP(y0, p)
		y1G2 = modP(y1G2, p)

		// Print reduced coordinates with label
		fmt.Printf("%s G2 Point: ((x0: %s, y0: %s), (x1: %s, y1: %s))\n", labels[i], x0.String(), y0.String(), x1G2.String(), y1G2.String())
	}
}
