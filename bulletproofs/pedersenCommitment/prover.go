package main

import (
	"crypto/elliptic"
	"crypto/sha256"
	"math/big"
)

func prover(curvePoints []ECPoint, G, B ECPoint, curve elliptic.Curve) {
	// setup secret values
	v0 := sha256.Sum256([]byte("secret 1"))

	v := make([]*big.Int, len(curvePoints))
	v[0] = new(big.Int).SetBytes(v0[:])

	commitments := commit(G, B, curve, v[0])
}

func commit(G, B ECPoint, curve elliptic.Curve, x *big.Int) []ECPoint {
	var commitments []ECPoint
	points := evaluatePolynomial(curve, G, x)
	for _, point := range points {
		cX, cY := curve.Add(point.X, point.Y, B.X, B.Y)
		commitments = append(commitments, ECPoint{X: cX, Y: cY})
	}
	return commitments
}
