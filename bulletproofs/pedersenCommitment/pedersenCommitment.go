package main

import (
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
)

type ECPoint struct {
	X, Y *big.Int
}

func main() {
	seed := make([]byte, 32)
	_, err := rand.Read(seed)
	if err != nil {
		fmt.Println("Error generating seed:", err)
		return
	}

	curve := elliptic.P256()

	points := generateECPoints(seed, 5, curve)

	for i, point := range points {
		fmt.Printf("Point %d: (%s, %s)\n", i+1, point.X.Text(10), point.Y.Text(10))
	}

	publicG := points[len(points)-2]
	publicB := points[len(points)-1]

	// Providing EC Points, public values G and B to prover and verifier
	prover(points[:3], publicG, publicB, curve)
}

func generateECPoints(seed []byte, n int, curve elliptic.Curve) []ECPoint {
	var points []ECPoint

	for i := 0; len(points) < n; i++ {
		// Create a new hash for each x value
		hashInput := append(seed, byte(i))
		// P is the order of the curve
		x := hashToBigInt(hashInput, curve.Params().P)

		y, yNeg := findYForX(curve, x)
		if y != nil {
			if randBit(seed) == 0 {
				points = append(points, ECPoint{X: new(big.Int).Set(x), Y: y})
			} else {
				points = append(points, ECPoint{X: new(big.Int).Set(x), Y: yNeg})
			}
		}
	}
	return points
}

func findYForX(curve elliptic.Curve, x *big.Int) (*big.Int, *big.Int) {
	xCubed := new(big.Int).Exp(x, big.NewInt(3), nil)
	a := big.NewInt(-3)
	aX := new(big.Int).Mul(a, x)
	rightSide := new(big.Int).Add(xCubed, aX)
	rightSide.Add(rightSide, curve.Params().B)
	rightSide.Mod(rightSide, curve.Params().P)

	y := new(big.Int).ModSqrt(rightSide, curve.Params().P)
	if y == nil {
		return nil, nil
	}
	yNeg := new(big.Int).Neg(y)
	yNeg.Mod(yNeg, curve.Params().P)
	return y, yNeg
}

func randBit(seed []byte) int {
	hash := sha256.Sum256(seed)
	return int(hash[0] & 1)
}

func hashToBigInt(data []byte, mod *big.Int) *big.Int {
	hash := sha256.Sum256(data)
	x := new(big.Int).SetBytes(hash[:])
	x.Mod(x, mod)
	return x
}

func evaluate() {

}

func prove() {

}

func verify() {

}

func cubicPolynomial() []*big.Int {
	return []*big.Int{
		big.NewInt(3),
		big.NewInt(2),
		big.NewInt(6),
		big.NewInt(1),
	}
}

func evaluatePolynomial(curve elliptic.Curve, Point ECPoint, v *big.Int) []ECPoint {
	polynomial := cubicPolynomial()
	// calculate constant term
	constantX, constantY := curve.ScalarMult(Point.X, Point.Y, polynomial[0].Bytes())
	// calculate linear term
	linearX, linearY := curve.ScalarMult(Point.X, Point.Y, new(big.Int).Mul(polynomial[1], v).Bytes())
	// calculate quadratic term
	quadraticX, quadraticY := curve.ScalarMult(Point.X, Point.Y, new(big.Int).Mul(polynomial[2], new(big.Int).Mul(v, v)).Bytes())
	// calculate cubic term
	v2Sq := new(big.Int).Mul(v, v)
	v2Cubed := new(big.Int).Mul(v, v2Sq)
	cubicX, cubicY := curve.ScalarMult(Point.X, Point.Y, v2Cubed.Bytes())

	return []ECPoint{
		{X: constantX, Y: constantY},
		{X: linearX, Y: linearY},
		{X: quadraticX, Y: quadraticY},
		{X: cubicX, Y: cubicY},
	}
}
