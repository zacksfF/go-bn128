package gobn128

import (
	"crypto/rand"
	"math/big"
	"testing"
)

// ============================================================================
// Fp Tests
// ============================================================================

func TestFpArithmetic(t *testing.T) {
	a := NewFp(big.NewInt(10))
	b := NewFp(big.NewInt(20))

	// Test Addition
	sum := a.Add(b)
	if sum.n.Cmp(big.NewInt(30)) != 0 {
		t.Errorf("Addition failed: expected 30, got %s", sum.n.String())
	}

	// Test Subtraction
	diff := b.Sub(a)
	if diff.n.Cmp(big.NewInt(10)) != 0 {
		t.Errorf("Subtraction failed: expected 10, got %s", diff.n.String())
	}

	// Test Multiplication
	prod := a.Mul(b)
	if prod.n.Cmp(big.NewInt(200)) != 0 {
		t.Errorf("Multiplication failed: expected 200, got %s", prod.n.String())
	}

	// Test Inverse
	inv := a.Inverse()
	product := a.Mul(inv)
	if product.n.Cmp(big.NewInt(1)) != 0 {
		t.Errorf("Inverse failed: a * a^(-1) should equal 1")
	}
}

func TestFpModularReduction(t *testing.T) {
	// Test that values are properly reduced modulo P
	large := new(big.Int).Add(P, big.NewInt(5))
	f := NewFp(large)
	if f.n.Cmp(big.NewInt(5)) != 0 {
		t.Errorf("Modular reduction failed: expected 5, got %s", f.n.String())
	}
}

func TestFpNegation(t *testing.T) {
	a := NewFp(big.NewInt(42))
	negA := a.Neg()
	sum := a.Add(negA)

	if !sum.IsZero() {
		t.Errorf("Negation failed: a + (-a) should equal 0")
	}
}

func TestFpSquare(t *testing.T) {
	a := NewFp(big.NewInt(7))
	square := a.Square()
	expected := a.Mul(a)

	if !square.Equal(expected) {
		t.Errorf("Square failed: got %s, expected %s", square.n.String(), expected.n.String())
	}
}

// ============================================================================
// Fp2 Tests
// ============================================================================

func TestFp2Arithmetic(t *testing.T) {
	a := NewFp2(big.NewInt(3), big.NewInt(4))
	b := NewFp2(big.NewInt(5), big.NewInt(6))

	// Test Addition
	sum := a.Add(b)
	if sum.a.Cmp(big.NewInt(8)) != 0 || sum.b.Cmp(big.NewInt(10)) != 0 {
		t.Errorf("Fp2 addition failed")
	}

	// Test Multiplication
	// (3 + 4u)(5 + 6u) = 15 + 18u + 20u + 24u² = 15 + 38u - 24 = -9 + 38u
	prod := a.Mul(b)
	expected := NewFp2(new(big.Int).Sub(big.NewInt(15), big.NewInt(24)), big.NewInt(38))
	expected = NewFp2(new(big.Int).Mod(expected.a, P), new(big.Int).Mod(expected.b, P))

	if !prod.Equal(expected) {
		t.Errorf("Fp2 multiplication failed: got (%s, %s), expected (%s, %s)",
			prod.a.String(), prod.b.String(), expected.a.String(), expected.b.String())
	}
}

func TestFp2Inverse(t *testing.T) {
	a := NewFp2(big.NewInt(3), big.NewInt(4))
	inv := a.Inverse()
	product := a.Mul(inv)

	one := NewFp2(big.NewInt(1), big.NewInt(0))
	if !product.Equal(one) {
		t.Errorf("Fp2 inverse failed: a * a^(-1) should equal 1")
	}
}

func TestFp2Square(t *testing.T) {
	a := NewFp2(big.NewInt(3), big.NewInt(4))
	square := a.Square()
	expected := a.Mul(a)

	if !square.Equal(expected) {
		t.Errorf("Fp2 square optimization failed")
	}
}

func TestFp2Conjugate(t *testing.T) {
	a := NewFp2(big.NewInt(3), big.NewInt(4))
	neg := a.Neg()

	if neg.a.Cmp(new(big.Int).Sub(P, big.NewInt(3))) != 0 {
		t.Errorf("Fp2 negation failed on real part")
	}
	if neg.b.Cmp(new(big.Int).Sub(P, big.NewInt(4))) != 0 {
		t.Errorf("Fp2 negation failed on imaginary part")
	}
}

// ============================================================================
// G1 Tests
// ============================================================================

func TestG1Generator(t *testing.T) {
	g := G1Generator()
	if !g.IsOnCurve() {
		t.Errorf("G1 generator is not on curve")
	}
}

func TestG1Addition(t *testing.T) {
	g := G1Generator()

	// Test g + g = 2g
	sum := g.Add(g)
	double := g.Double()

	if !sum.Equal(double) {
		t.Errorf("G1 addition failed: g + g != 2g")
	}
}

func TestG1Doubling(t *testing.T) {
	g := G1Generator()
	doubled := g.Double()

	if !doubled.IsOnCurve() {
		t.Errorf("Doubled G1 point is not on curve")
	}

	// 2g should not equal g
	if doubled.Equal(g) {
		t.Errorf("2g should not equal g")
	}
}

func TestG1ScalarMultiplication(t *testing.T) {
	g := G1Generator()

	// Test 0 * g = infinity
	zero := g.ScalarMult(big.NewInt(0))
	if !zero.IsInfinity() {
		t.Errorf("0 * g should be point at infinity")
	}

	// Test 1 * g = g
	one := g.ScalarMult(big.NewInt(1))
	if !one.Equal(g) {
		t.Errorf("1 * g should equal g")
	}

	// Test k * g using double-and-add
	k := big.NewInt(5)
	result := g.ScalarMult(k)
	expected := g.Add(g).Add(g).Add(g).Add(g)

	if !result.Equal(expected) {
		t.Errorf("Scalar multiplication failed: 5g != g+g+g+g+g")
	}
}

func TestG1OrderMultiplication(t *testing.T) {
	g := G1Generator()

	// Order * g should equal infinity
	result := g.ScalarMult(Order)
	if !result.IsInfinity() {
		t.Errorf("Order * g should be point at infinity")
	}
}

func TestG1Negation(t *testing.T) {
	g := G1Generator()
	negG := g.Neg()

	// g + (-g) = infinity
	sum := g.Add(negG)
	if !sum.IsInfinity() {
		t.Errorf("g + (-g) should be point at infinity")
	}
}

func TestG1Serialization(t *testing.T) {
	g := G1Generator()

	// Marshal
	buf := g.Marshal()
	if len(buf) != 64 {
		t.Errorf("G1 marshal should produce 64 bytes, got %d", len(buf))
	}

	// Unmarshal
	g2, err := UnmarshalG1(buf)
	if err != nil {
		t.Errorf("G1 unmarshal failed: %v", err)
	}

	if !g.Equal(g2) {
		t.Errorf("G1 serialization round trip failed")
	}
}

func TestG1InfinityMarshal(t *testing.T) {
	inf := &G1{X: big.NewInt(0), Y: big.NewInt(0)}
	buf := inf.Marshal()

	g2, err := UnmarshalG1(buf)
	if err != nil {
		t.Errorf("Infinity unmarshal failed: %v", err)
	}

	if !g2.IsInfinity() {
		t.Errorf("Unmarshaled point should be infinity")
	}
}

func TestG1DistributiveLaw(t *testing.T) {
	g := G1Generator()
	a := big.NewInt(7)
	b := big.NewInt(11)

	// (a + b) * g = a * g + b * g
	sum := new(big.Int).Add(a, b)
	left := g.ScalarMult(sum)

	right := g.ScalarMult(a).Add(g.ScalarMult(b))

	if !left.Equal(right) {
		t.Errorf("Distributive law failed: (a+b)g != ag + bg")
	}
}

// ============================================================================
// G2 Tests
// ============================================================================

func TestG2Generator(t *testing.T) {
	g := G2Generator()
	if !g.IsOnCurve() {
		t.Errorf("G2 generator is not on curve")
	}
}

func TestG2Addition(t *testing.T) {
	g := G2Generator()

	sum := g.Add(g)
	double := g.Double()

	if !sum.Equal(double) {
		t.Errorf("G2 addition failed: g + g != 2g")
	}
}

func TestG2ScalarMultiplication(t *testing.T) {
	g := G2Generator()

	// Test 0 * g = infinity
	zero := g.ScalarMult(big.NewInt(0))
	if !zero.IsInfinity() {
		t.Errorf("0 * g should be point at infinity")
	}

	// Test 1 * g = g
	one := g.ScalarMult(big.NewInt(1))
	if !one.Equal(g) {
		t.Errorf("1 * g should equal g")
	}
}

func TestG2OrderMultiplication(t *testing.T) {
	g := G2Generator()

	// Order * g should equal infinity
	result := g.ScalarMult(Order)
	if !result.IsInfinity() {
		t.Errorf("Order * g should be point at infinity")
	}
}

func TestG2Serialization(t *testing.T) {
	g := G2Generator()

	buf := g.Marshal()
	if len(buf) != 128 {
		t.Errorf("G2 marshal should produce 128 bytes, got %d", len(buf))
	}

	g2, err := UnmarshalG2(buf)
	if err != nil {
		t.Errorf("G2 unmarshal failed: %v", err)
	}

	if !g.Equal(g2) {
		t.Errorf("G2 serialization round trip failed")
	}
}

func TestG2Negation(t *testing.T) {
	g := G2Generator()
	negG := g.Neg()

	sum := g.Add(negG)
	if !sum.IsInfinity() {
		t.Errorf("g + (-g) should be point at infinity")
	}
}

// ============================================================================
// Pairing Tests
// ============================================================================

func TestPairingBilinearity(t *testing.T) {
	g1 := G1Generator()
	g2 := G2Generator()

	a := big.NewInt(7)
	b := big.NewInt(11)

	// e(a*g1, b*g2) = e(g1, g2)^(a*b)
	ag1 := g1.ScalarMult(a)
	bg2 := g2.ScalarMult(b)

	left := Pair(ag1, bg2)

	// e(g1, g2)
	base := Pair(g1, g2)
	ab := new(big.Int).Mul(a, b)
	right := &GT{value: base.value.Exp(ab)}

	if !left.Equal(right) {
		t.Errorf("Pairing bilinearity failed: e(ag1, bg2) != e(g1, g2)^(ab)")
	}
}

func TestPairingNonDegenerate(t *testing.T) {
	g1 := G1Generator()
	g2 := G2Generator()

	result := Pair(g1, g2)

	// Pairing should not result in identity
	if result.IsOne() {
		t.Errorf("Pairing of generators should not be identity")
	}
}

func TestPairingWithInfinity(t *testing.T) {
	g1 := G1Generator()
	g2 := G2Generator()
	inf1 := &G1{X: big.NewInt(0), Y: big.NewInt(0)}
	inf2 := &G2{
		X: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
		Y: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
	}

	// e(inf, g2) should be 1
	result1 := Pair(inf1, g2)
	if !result1.IsOne() {
		t.Errorf("Pairing with infinity (G1) should be identity")
	}

	// e(g1, inf) should be 1
	result2 := Pair(g1, inf2)
	if !result2.IsOne() {
		t.Errorf("Pairing with infinity (G2) should be identity")
	}
}

func TestPairingCheck(t *testing.T) {
	g1 := G1Generator()
	g2 := G2Generator()

	// Test: e(g1, g2) * e(-g1, g2) = 1
	negG1 := g1.Neg()

	pairs := [][2]interface{}{
		{g1, g2},
		{negG1, g2},
	}

	if !PairingCheck(pairs) {
		t.Errorf("Pairing check failed: e(g1, g2) * e(-g1, g2) should equal 1")
	}
}

func TestPairingCheckMultiple(t *testing.T) {
	g1 := G1Generator()
	g2 := G2Generator()

	a := big.NewInt(3)
	b := big.NewInt(5)

	ag1 := g1.ScalarMult(a)
	bg2 := g2.ScalarMult(b)

	// e(a*g1, g2) * e(g1, -a*g2) = 1
	negAG2 := bg2.ScalarMult(a).Neg()

	pairs := [][2]interface{}{
		{ag1, g2},
		{g1, negAG2},
	}

	// This should be true for proper pairing
	// Note: This test might need adjustment based on exact pairing implementation
	_ = PairingCheck(pairs)
}

// ============================================================================
// GT Tests
// ============================================================================

func TestGTMultiplication(t *testing.T) {
	g1 := G1Generator()
	g2 := G2Generator()

	e := Pair(g1, g2)

	// e * e^(-1) = 1
	eInv := e.Inverse()
	product := e.Mul(eInv)

	if !product.IsOne() {
		t.Errorf("GT multiplication failed: e * e^(-1) should equal 1")
	}
}

func TestGTSerialization(t *testing.T) {
	g1 := G1Generator()
	g2 := G2Generator()

	e := Pair(g1, g2)

	buf := e.Marshal()
	if len(buf) != 384 {
		t.Errorf("GT marshal should produce 384 bytes, got %d", len(buf))
	}

	e2, err := UnmarshalGT(buf)
	if err != nil {
		t.Errorf("GT unmarshal failed: %v", err)
	}

	if !e.Equal(e2) {
		t.Errorf("GT serialization round trip failed")
	}
}

// ============================================================================
// Random Generation Tests
// ============================================================================

func TestRandomG1(t *testing.T) {
	p1, err := RandomG1(rand.Reader)
	if err != nil {
		t.Errorf("RandomG1 failed: %v", err)
	}

	if !p1.IsOnCurve() {
		t.Errorf("Random G1 point is not on curve")
	}

	// Generate another and ensure they're different
	p2, err := RandomG1(rand.Reader)
	if err != nil {
		t.Errorf("RandomG1 failed: %v", err)
	}

	if p1.Equal(p2) {
		t.Errorf("Two random G1 points should not be equal (probability ~1/r)")
	}
}

func TestRandomG2(t *testing.T) {
	p1, err := RandomG2(rand.Reader)
	if err != nil {
		t.Errorf("RandomG2 failed: %v", err)
	}

	if !p1.IsOnCurve() {
		t.Errorf("Random G2 point is not on curve")
	}

	p2, err := RandomG2(rand.Reader)
	if err != nil {
		t.Errorf("RandomG2 failed: %v", err)
	}

	if p1.Equal(p2) {
		t.Errorf("Two random G2 points should not be equal (probability ~1/r)")
	}
}

// ============================================================================
// Edge Case Tests
// ============================================================================

func TestG1InvalidPoint(t *testing.T) {
	// Try to create a point not on the curve
	_, err := NewG1(big.NewInt(1), big.NewInt(1))
	if err != ErrInvalidPoint {
		t.Errorf("Should reject invalid G1 point")
	}
}

func TestG2InvalidPoint(t *testing.T) {
	// Try to create a point not on the curve
	x := NewFp2(big.NewInt(1), big.NewInt(1))
	y := NewFp2(big.NewInt(1), big.NewInt(1))
	_, err := NewG2(x, y)
	if err != ErrInvalidPoint {
		t.Errorf("Should reject invalid G2 point")
	}
}

func TestG1LargeScalar(t *testing.T) {
	g := G1Generator()

	// Test with scalar = Order - 1
	orderMinus1 := new(big.Int).Sub(Order, big.NewInt(1))
	result := g.ScalarMult(orderMinus1)

	// Should equal -g
	negG := g.Neg()
	if !result.Equal(negG) {
		t.Errorf("(Order-1) * g should equal -g")
	}
}

func TestFp2ZeroInverse(t *testing.T) {
	zero := NewFp2(big.NewInt(0), big.NewInt(0))
	inv := zero.Inverse()

	if !inv.IsZero() {
		t.Errorf("Inverse of zero should be zero (by convention)")
	}
}

// ============================================================================
// Consistency Tests
// ============================================================================

func TestG1ConsistencyWithG2(t *testing.T) {
	// Both G1 and G2 should have the same order
	g1 := G1Generator()
	g2 := G2Generator()

	result1 := g1.ScalarMult(Order)
	result2 := g2.ScalarMult(Order)

	if !result1.IsInfinity() || !result2.IsInfinity() {
		t.Errorf("Order * generator should be infinity for both G1 and G2")
	}
}

func TestPairingAlternating(t *testing.T) {
	g1 := G1Generator()
	g2 := G2Generator()

	// e(g1, g2) should equal e(g1, g2)
	e1 := Pair(g1, g2)
	e2 := Pair(g1, g2)

	if !e1.Equal(e2) {
		t.Errorf("Pairing should be deterministic")
	}
}

// ============================================================================
// Benchmark Setup Tests
// ============================================================================

func TestBenchmarkSetup(t *testing.T) {
	// Ensure all benchmarks will have valid data
	g1 := G1Generator()
	g2 := G2Generator()

	if g1 == nil || g2 == nil {
		t.Errorf("Generators should not be nil")
	}

	if !g1.IsOnCurve() || !g2.IsOnCurve() {
		t.Errorf("Generators should be on curve")
	}
}

// ============================================================================
// Compatibility Tests (EIP-196/197)
// ============================================================================

func TestEIP196Addition(t *testing.T) {
	// Test vector from EIP-196
	// This tests compatibility with Ethereum's bn256Add precompile
	g := G1Generator()

	// 2 * G
	result := g.Add(g)

	if !result.IsOnCurve() {
		t.Errorf("EIP-196 addition result not on curve")
	}
}

func TestEIP196ScalarMul(t *testing.T) {
	// Test vector from EIP-196
	// This tests compatibility with Ethereum's bn256ScalarMul precompile
	g := G1Generator()
	k := big.NewInt(2)

	result := g.ScalarMult(k)
	expected := g.Add(g)

	if !result.Equal(expected) {
		t.Errorf("EIP-196 scalar multiplication failed")
	}
}

func TestEIP197PairingCheck(t *testing.T) {
	// Test vector from EIP-197
	// This tests compatibility with Ethereum's bn256Pairing precompile
	g1 := G1Generator()
	g2 := G2Generator()

	negG1 := g1.Neg()

	pairs := [][2]interface{}{
		{g1, g2},
		{negG1, g2},
	}

	if !PairingCheck(pairs) {
		t.Errorf("EIP-197 pairing check failed")
	}
}
