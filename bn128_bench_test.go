package gobn128

import (
	"crypto/rand"
	"math/big"
	"testing"
)

// ============================================================================
// Fp Benchmarks
// ============================================================================

func BenchmarkFpAdd(b *testing.B) {
	x := NewFp(big.NewInt(12345))
	y := NewFp(big.NewInt(67890))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = x.Add(y)
	}
}

func BenchmarkFpMul(b *testing.B) {
	x := NewFp(big.NewInt(12345))
	y := NewFp(big.NewInt(67890))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = x.Mul(y)
	}
}

func BenchmarkFpSquare(b *testing.B) {
	x := NewFp(big.NewInt(12345))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = x.Square()
	}
}

func BenchmarkFpInverse(b *testing.B) {
	x := NewFp(big.NewInt(12345))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = x.Inverse()
	}
}

// ============================================================================
// Fp2 Benchmarks
// ============================================================================

func BenchmarkFp2Add(b *testing.B) {
	x := NewFp2(big.NewInt(12345), big.NewInt(67890))
	y := NewFp2(big.NewInt(11111), big.NewInt(22222))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = x.Add(y)
	}
}

func BenchmarkFp2Mul(b *testing.B) {
	x := NewFp2(big.NewInt(12345), big.NewInt(67890))
	y := NewFp2(big.NewInt(11111), big.NewInt(22222))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = x.Mul(y)
	}
}

func BenchmarkFp2Square(b *testing.B) {
	x := NewFp2(big.NewInt(12345), big.NewInt(67890))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = x.Square()
	}
}

func BenchmarkFp2Inverse(b *testing.B) {
	x := NewFp2(big.NewInt(12345), big.NewInt(67890))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = x.Inverse()
	}
}

// ============================================================================
// Fp6 Benchmarks
// ============================================================================

func BenchmarkFp6Add(b *testing.B) {
	c0 := NewFp2(big.NewInt(1), big.NewInt(2))
	c1 := NewFp2(big.NewInt(3), big.NewInt(4))
	c2 := NewFp2(big.NewInt(5), big.NewInt(6))
	x := NewFp6(c0, c1, c2)

	d0 := NewFp2(big.NewInt(7), big.NewInt(8))
	d1 := NewFp2(big.NewInt(9), big.NewInt(10))
	d2 := NewFp2(big.NewInt(11), big.NewInt(12))
	y := NewFp6(d0, d1, d2)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = x.Add(y)
	}
}

func BenchmarkFp6Mul(b *testing.B) {
	c0 := NewFp2(big.NewInt(1), big.NewInt(2))
	c1 := NewFp2(big.NewInt(3), big.NewInt(4))
	c2 := NewFp2(big.NewInt(5), big.NewInt(6))
	x := NewFp6(c0, c1, c2)

	d0 := NewFp2(big.NewInt(7), big.NewInt(8))
	d1 := NewFp2(big.NewInt(9), big.NewInt(10))
	d2 := NewFp2(big.NewInt(11), big.NewInt(12))
	y := NewFp6(d0, d1, d2)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = x.Mul(y)
	}
}

func BenchmarkFp6Square(b *testing.B) {
	c0 := NewFp2(big.NewInt(1), big.NewInt(2))
	c1 := NewFp2(big.NewInt(3), big.NewInt(4))
	c2 := NewFp2(big.NewInt(5), big.NewInt(6))
	x := NewFp6(c0, c1, c2)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = x.Square()
	}
}

func BenchmarkFp6Inverse(b *testing.B) {
	c0 := NewFp2(big.NewInt(1), big.NewInt(2))
	c1 := NewFp2(big.NewInt(3), big.NewInt(4))
	c2 := NewFp2(big.NewInt(5), big.NewInt(6))
	x := NewFp6(c0, c1, c2)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = x.Inverse()
	}
}

// ============================================================================
// Fp12 Benchmarks
// ============================================================================

func BenchmarkFp12Add(b *testing.B) {
	c00 := NewFp2(big.NewInt(1), big.NewInt(2))
	c01 := NewFp2(big.NewInt(3), big.NewInt(4))
	c02 := NewFp2(big.NewInt(5), big.NewInt(6))
	c0 := NewFp6(c00, c01, c02)

	c10 := NewFp2(big.NewInt(7), big.NewInt(8))
	c11 := NewFp2(big.NewInt(9), big.NewInt(10))
	c12 := NewFp2(big.NewInt(11), big.NewInt(12))
	c1 := NewFp6(c10, c11, c12)

	x := NewFp12(c0, c1)

	d00 := NewFp2(big.NewInt(13), big.NewInt(14))
	d01 := NewFp2(big.NewInt(15), big.NewInt(16))
	d02 := NewFp2(big.NewInt(17), big.NewInt(18))
	d0 := NewFp6(d00, d01, d02)

	d10 := NewFp2(big.NewInt(19), big.NewInt(20))
	d11 := NewFp2(big.NewInt(21), big.NewInt(22))
	d12 := NewFp2(big.NewInt(23), big.NewInt(24))
	d1 := NewFp6(d10, d11, d12)

	y := NewFp12(d0, d1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = x.Add(y)
	}
}

func BenchmarkFp12Mul(b *testing.B) {
	c00 := NewFp2(big.NewInt(1), big.NewInt(2))
	c01 := NewFp2(big.NewInt(3), big.NewInt(4))
	c02 := NewFp2(big.NewInt(5), big.NewInt(6))
	c0 := NewFp6(c00, c01, c02)

	c10 := NewFp2(big.NewInt(7), big.NewInt(8))
	c11 := NewFp2(big.NewInt(9), big.NewInt(10))
	c12 := NewFp2(big.NewInt(11), big.NewInt(12))
	c1 := NewFp6(c10, c11, c12)

	x := NewFp12(c0, c1)

	d00 := NewFp2(big.NewInt(13), big.NewInt(14))
	d01 := NewFp2(big.NewInt(15), big.NewInt(16))
	d02 := NewFp2(big.NewInt(17), big.NewInt(18))
	d0 := NewFp6(d00, d01, d02)

	d10 := NewFp2(big.NewInt(19), big.NewInt(20))
	d11 := NewFp2(big.NewInt(21), big.NewInt(22))
	d12 := NewFp2(big.NewInt(23), big.NewInt(24))
	d1 := NewFp6(d10, d11, d12)

	y := NewFp12(d0, d1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = x.Mul(y)
	}
}

func BenchmarkFp12Square(b *testing.B) {
	c00 := NewFp2(big.NewInt(1), big.NewInt(2))
	c01 := NewFp2(big.NewInt(3), big.NewInt(4))
	c02 := NewFp2(big.NewInt(5), big.NewInt(6))
	c0 := NewFp6(c00, c01, c02)

	c10 := NewFp2(big.NewInt(7), big.NewInt(8))
	c11 := NewFp2(big.NewInt(9), big.NewInt(10))
	c12 := NewFp2(big.NewInt(11), big.NewInt(12))
	c1 := NewFp6(c10, c11, c12)

	x := NewFp12(c0, c1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = x.Square()
	}
}

func BenchmarkFp12Inverse(b *testing.B) {
	c00 := NewFp2(big.NewInt(1), big.NewInt(2))
	c01 := NewFp2(big.NewInt(3), big.NewInt(4))
	c02 := NewFp2(big.NewInt(5), big.NewInt(6))
	c0 := NewFp6(c00, c01, c02)

	c10 := NewFp2(big.NewInt(7), big.NewInt(8))
	c11 := NewFp2(big.NewInt(9), big.NewInt(10))
	c12 := NewFp2(big.NewInt(11), big.NewInt(12))
	c1 := NewFp6(c10, c11, c12)

	x := NewFp12(c0, c1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = x.Inverse()
	}
}

func BenchmarkFp12Exp(b *testing.B) {
	c00 := NewFp2(big.NewInt(1), big.NewInt(2))
	c01 := NewFp2(big.NewInt(3), big.NewInt(4))
	c02 := NewFp2(big.NewInt(5), big.NewInt(6))
	c0 := NewFp6(c00, c01, c02)

	c10 := NewFp2(big.NewInt(7), big.NewInt(8))
	c11 := NewFp2(big.NewInt(9), big.NewInt(10))
	c12 := NewFp2(big.NewInt(11), big.NewInt(12))
	c1 := NewFp6(c10, c11, c12)

	x := NewFp12(c0, c1)
	exp := big.NewInt(12345)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = x.Exp(exp)
	}
}

// ============================================================================
// G1 Benchmarks
// ============================================================================

func BenchmarkG1Add(b *testing.B) {
	g := G1Generator()
	p := g.ScalarMult(big.NewInt(123))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = g.Add(p)
	}
}

func BenchmarkG1Double(b *testing.B) {
	g := G1Generator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = g.Double()
	}
}

func BenchmarkG1ScalarMult(b *testing.B) {
	g := G1Generator()
	scalar, _ := randomScalar(rand.Reader)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = g.ScalarMult(scalar)
	}
}

func BenchmarkG1ScalarMultSmall(b *testing.B) {
	g := G1Generator()
	scalar := big.NewInt(12345)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = g.ScalarMult(scalar)
	}
}

func BenchmarkG1ScalarMultLarge(b *testing.B) {
	g := G1Generator()
	// Use order - 1 (largest valid scalar)
	scalar := new(big.Int).Sub(Order, big.NewInt(1))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = g.ScalarMult(scalar)
	}
}

func BenchmarkG1ScalarBaseMult(b *testing.B) {
	scalar, _ := randomScalar(rand.Reader)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ScalarBaseMult(scalar)
	}
}

func BenchmarkG1Marshal(b *testing.B) {
	g := G1Generator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = g.Marshal()
	}
}

func BenchmarkG1Unmarshal(b *testing.B) {
	g := G1Generator()
	buf := g.Marshal()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = UnmarshalG1(buf)
	}
}

func BenchmarkG1IsOnCurve(b *testing.B) {
	g := G1Generator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = g.IsOnCurve()
	}
}

func BenchmarkG1Neg(b *testing.B) {
	g := G1Generator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = g.Neg()
	}
}

// ============================================================================
// G2 Benchmarks
// ============================================================================

func BenchmarkG2Add(b *testing.B) {
	g := G2Generator()
	p := g.ScalarMult(big.NewInt(123))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = g.Add(p)
	}
}

func BenchmarkG2Double(b *testing.B) {
	g := G2Generator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = g.Double()
	}
}

func BenchmarkG2ScalarMult(b *testing.B) {
	g := G2Generator()
	scalar, _ := randomScalar(rand.Reader)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = g.ScalarMult(scalar)
	}
}

func BenchmarkG2ScalarMultSmall(b *testing.B) {
	g := G2Generator()
	scalar := big.NewInt(12345)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = g.ScalarMult(scalar)
	}
}

func BenchmarkG2Marshal(b *testing.B) {
	g := G2Generator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = g.Marshal()
	}
}

func BenchmarkG2Unmarshal(b *testing.B) {
	g := G2Generator()
	buf := g.Marshal()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = UnmarshalG2(buf)
	}
}

func BenchmarkG2IsOnCurve(b *testing.B) {
	g := G2Generator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = g.IsOnCurve()
	}
}

func BenchmarkG2Neg(b *testing.B) {
	g := G2Generator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = g.Neg()
	}
}

// ============================================================================
// Pairing Benchmarks
// ============================================================================

func BenchmarkPairing(b *testing.B) {
	g1 := G1Generator()
	g2 := G2Generator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Pair(g1, g2)
	}
}

func BenchmarkPairingRandom(b *testing.B) {
	p1, _ := RandomG1(rand.Reader)
	p2, _ := RandomG2(rand.Reader)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Pair(p1, p2)
	}
}

func BenchmarkMillerLoop(b *testing.B) {
	g1 := G1Generator()
	g2 := G2Generator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = millerLoop(g1, g2)
	}
}

func BenchmarkFinalExponentiation(b *testing.B) {
	g1 := G1Generator()
	g2 := G2Generator()
	f := millerLoop(g1, g2)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = finalExponentiation(f)
	}
}

func BenchmarkPairingCheck2(b *testing.B) {
	g1 := G1Generator()
	g2 := G2Generator()
	negG1 := g1.Neg()

	pairs := [][2]interface{}{
		{g1, g2},
		{negG1, g2},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = PairingCheck(pairs)
	}
}

func BenchmarkPairingCheck4(b *testing.B) {
	g1 := G1Generator()
	g2 := G2Generator()

	p1 := g1.ScalarMult(big.NewInt(3))
	p2 := g1.ScalarMult(big.NewInt(5))
	p3 := g1.ScalarMult(big.NewInt(7))
	p4 := g1.Neg()

	q1 := g2
	q2 := g2
	q3 := g2
	q4 := g2.ScalarMult(big.NewInt(15)) // 3+5+7 = 15

	pairs := [][2]interface{}{
		{p1, q1},
		{p2, q2},
		{p3, q3},
		{p4, q4},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = PairingCheck(pairs)
	}
}

func BenchmarkPairingCheck8(b *testing.B) {
	g1 := G1Generator()
	g2 := G2Generator()

	// Create 8 pairs that sum to zero
	pairs := make([][2]interface{}, 8)
	sum := big.NewInt(0)

	for i := 0; i < 7; i++ {
		scalar := big.NewInt(int64(i + 1))
		pairs[i] = [2]interface{}{g1.ScalarMult(scalar), g2}
		sum.Add(sum, scalar)
	}

	// Last pair negates the sum
	pairs[7] = [2]interface{}{g1.ScalarMult(sum).Neg(), g2}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = PairingCheck(pairs)
	}
}

// ============================================================================
// GT Benchmarks
// ============================================================================

func BenchmarkGTMul(b *testing.B) {
	g1 := G1Generator()
	g2 := G2Generator()
	e1 := Pair(g1, g2)
	e2 := Pair(g1.ScalarMult(big.NewInt(2)), g2)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e1.Mul(e2)
	}
}

func BenchmarkGTInverse(b *testing.B) {
	g1 := G1Generator()
	g2 := G2Generator()
	e := Pair(g1, g2)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Inverse()
	}
}

func BenchmarkGTMarshal(b *testing.B) {
	g1 := G1Generator()
	g2 := G2Generator()
	e := Pair(g1, g2)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Marshal()
	}
}

func BenchmarkGTUnmarshal(b *testing.B) {
	g1 := G1Generator()
	g2 := G2Generator()
	e := Pair(g1, g2)
	buf := e.Marshal()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = UnmarshalGT(buf)
	}
}

// ============================================================================
// Random Generation Benchmarks
// ============================================================================

func BenchmarkRandomG1(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = RandomG1(rand.Reader)
	}
}

func BenchmarkRandomG2(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = RandomG2(rand.Reader)
	}
}

func BenchmarkRandomScalar(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = randomScalar(rand.Reader)
	}
}

// ============================================================================
// Hash-to-Curve Benchmarks
// ============================================================================

func BenchmarkHashToG1(b *testing.B) {
	data := []byte("benchmark test data for hashing to G1")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = HashToG1(data)
	}
}

// ============================================================================
// Composite Operation Benchmarks
// ============================================================================

func BenchmarkZKSNARKVerification(b *testing.B) {
	g1 := G1Generator()
	g2 := G2Generator()

	// Simulate proof elements
	proofA := g1.ScalarMult(big.NewInt(123))
	proofB := g2.ScalarMult(big.NewInt(456))
	proofC := g1.ScalarMult(big.NewInt(789))

	// Verification key elements
	alpha := g1.ScalarMult(big.NewInt(111))
	beta := g2.ScalarMult(big.NewInt(222))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// e(proofA, proofB) * e(alpha, beta) * e(proofC, g2) = 1
		pairs := [][2]interface{}{
			{proofA, proofB},
			{alpha, beta},
			{proofC, g2},
		}
		_ = PairingCheck(pairs)
	}
}

// ============================================================================
// Memory Allocation Benchmarks
// ============================================================================

func BenchmarkG1AddAllocs(b *testing.B) {
	g := G1Generator()
	p := g.ScalarMult(big.NewInt(123))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = g.Add(p)
	}
}

func BenchmarkG2AddAllocs(b *testing.B) {
	g := G2Generator()
	p := g.ScalarMult(big.NewInt(123))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = g.Add(p)
	}
}

func BenchmarkPairingAllocs(b *testing.B) {
	g1 := G1Generator()
	g2 := G2Generator()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Pair(g1, g2)
	}
}

// ============================================================================
// Parallel Benchmarks
// ============================================================================

func BenchmarkG1ScalarMultParallel(b *testing.B) {
	g := G1Generator()
	scalar, _ := randomScalar(rand.Reader)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = g.ScalarMult(scalar)
		}
	})
}

func BenchmarkPairingParallel(b *testing.B) {
	g1 := G1Generator()
	g2 := G2Generator()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = Pair(g1, g2)
		}
	})
}
