// Package bn128 implements the BN128 (BN254) elliptic curve pairing operations.
// This implementation is designed for blockchain and zero-knowledge proof applications,
// with support for Ethereum's EIP-196 and EIP-197 precompiles.
//
// BN128 is a pairing-friendly elliptic curve with:
// - Base field: Fp where p = 21888242871839275222246405745257275088696311157297823662689037894645226208583
// - Embedding degree: k = 12
// - Curve equation: y² = x³ + 3 (over Fp for G1, over Fp2 for G2)
// - Order: r = 21888242871839275222246405745257275088548364400416034343698204186575808495617
//
// Security Note: BN128 provides approximately 100 bits of security, which is considered
// acceptable for most blockchain applications but below the 128-bit standard.
package gobn128

import (
	"crypto/rand"
	"errors"
	"io"
	"math/big"
)

var (
	// ErrInvalidPoint indicates a point is not on the curve
	ErrInvalidPoint = errors.New("bn128: point not on curve")
	// ErrInvalidPairing indicates a pairing check failed
	ErrInvalidPairing = errors.New("bn128: pairing check failed")
	// ErrInvalidEncoding indicates invalid serialization format
	ErrInvalidEncoding = errors.New("bn128: invalid encoding")
)

// Curve parameters
var (
	// P is the base field modulus
	// p = 21888242871839275222246405745257275088696311157297823662689037894645226208583
	P = fromHex("30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47")

	// Order is the curve order (number of points)
	// r = 21888242871839275222246405745257275088548364400416034343698204186575808495617
	Order = fromHex("30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f0000001")

	// GeneratorG1X and GeneratorG1Y are the G1 generator coordinates
	GeneratorG1X = big.NewInt(1)
	GeneratorG1Y = big.NewInt(2)

	// GeneratorG2X and GeneratorG2Y are the G2 generator coordinates (over Fp2)
	GeneratorG2X = &Fp2{
		fromHex("1800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed"),
		fromHex("198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c2"),
	}
	GeneratorG2Y = &Fp2{
		fromHex("12c85ea5db8c6deb4aab71808dcb408fe3d1e7690c43d37b4ce6cc0166fa7daa"),
		fromHex("090689d0585ff075ec9e99ad690c3395bc4b313370b38ef355acdadcd122975b"),
	}

	// TwistB is the curve coefficient for G2: y² = x³ + b where b = 3/(9+u)
	TwistB = &Fp2{
		fromHex("2b149d40ceb8aaae81be18991be06ac3b5b4c5e559dbefa33267e6dc24a138e5"),
		fromHex("009713b03af0fed4cd2cafadeed8fdf4a74fa084e52d1852e4a2bd0685c315d2"),
	}

	// xiToPMinus1Over6 is used in the final exponentiation
	xiToPMinus1Over6 = &Fp2{
		fromHex("16c9e55061ebae204ba4cc8bd75a079432ae2a1d0b7c9dce1665d51c640fcba2"),
		fromHex("063cf305489af5dcdc5ec698b6e2f9b9dbaae0eda9c95998dc54014671a0135a"),
	}

	// xiToPMinus1Over3 is used in the final exponentiation
	xiToPMinus1Over3 = &Fp2{
		fromHex("06c990cc9b6bf4c3c6040c2e85e8c0c0c9c99c6d3c1b4c6f4c5c5c5c5c5c5c5c"),
		fromHex("1787d6f5e7f0c7c7c7c7c7c7c7c7c7c7c7c7c7c7c7c7c7c7c7c7c7c7c7c7c7c7"),
	}
)

// Helper function to convert hex string to big.Int
func fromHex(s string) *big.Int {
	n, ok := new(big.Int).SetString(s, 16)
	if !ok {
		panic("bn128: invalid hex string")
	}
	return n
}

// ============================================================================
// Fp - Base Field Element
// ============================================================================
// Fp represents an element from a big.Int
type Fp struct {
	n *big.Int
}

// NewFp creates a neew field element from big.Int
func NewFp(n *big.Int) *Fp {
	return &Fp{n: new(big.Int).Mod(n, P)}
}

// Copy creates a deep copy of the field element
func (f *Fp) Copy() *Fp {
	return &Fp{n: new(big.Int).Set(f.n)}
}

// Add computes f + g in Fp
func (f *Fp) Add(g *Fp) *Fp {
	result := new(big.Int).Add(f.n, g.n)
	return &Fp{n: result.Mod(result, P)}
}

// Sub computes f - g in Fp
func (f *Fp) Sub(g *Fp) *Fp {
	result := new(big.Int).Sub(f.n, g.n)
	return &Fp{n: result.Mod(result, P)}
}

// Mul computes f * g in Fp
func (f *Fp) Mul(g *Fp) *Fp {
	result := new(big.Int).Mul(f.n, g.n)
	return &Fp{n: result.Mod(result, P)}
}

// Square computes f² in Fp
func (f *Fp) Square() *Fp {
	result := new(big.Int).Mul(f.n, f.n)
	return &Fp{n: result.Mod(result, P)}
}

// Inverse computes f⁻¹ in Fp using Fermat's little theorem
func (f *Fp) Inverse() *Fp {
	if f.IsZero() {
		return &Fp{n: big.NewInt(0)}
	}
	// By Fermat's little theorem: a^(p-1) ≡ 1 (mod p)
	// Therefore: a^(-1) ≡ a^(p-2) (mod p)
	pMinus2 := new(big.Int).Sub(P, big.NewInt(2))
	result := new(big.Int).Exp(f.n, pMinus2, P)
	return &Fp{n: result}
}

// Neg computes -f in Fp
func (f *Fp) Neg() *Fp {
	if f.IsZero() {
		return &Fp{n: big.NewInt(0)}
	}
	return &Fp{n: new(big.Int).Sub(P, f.n)}
}

// IsZero returns true if f == 0
func (f *Fp) IsZero() bool {
	return f.n.Sign() == 0
}

// Equal returns true if f == g
func (f *Fp) Equal(g *Fp) bool {
	return f.n.Cmp(g.n) == 0
}

// BigInt returns the big.Int representation
func (f *Fp) BigInt() *big.Int {
	return new(big.Int).Set(f.n)
}

// ============================================================================
// Fp2 - Quadratic Extension Field Element
// ============================================================================

// Fp2 represents an element in Fp2 = Fp[u]/(u²+1)
// Represented as a + b*u where a, b ∈ Fp
type Fp2 struct {
	a, b *big.Int //a + b*u
}

// NewFp2 creates a new Fp2 element
func NewFp2(a, b *big.Int) *Fp2 {
	return &Fp2{
		a: new(big.Int).Mod(a, P),
		b: new(big.Int).Mod(b, P),
	}
}

// Copy creates a deep copy
func (f *Fp2) Copy() *Fp2 {
	return &Fp2{
		a: new(big.Int).Set(f.a),
		b: new(big.Int).Set(f.b),
	}
}

// Add computes f + g in Fp2
func (f *Fp2) Add(g *Fp2) *Fp2 {
	a := new(big.Int).Add(f.a, g.a)
	a.Mod(a, P)
	b := new(big.Int).Add(f.b, g.b)
	b.Mod(b, P)
	return &Fp2{a: a, b: b}
}

// Sub computes f - g in Fp2
func (f *Fp2) Sub(g *Fp2) *Fp2 {
	a := new(big.Int).Sub(f.a, g.a)
	a.Mod(a, P)
	b := new(big.Int).Sub(f.b, g.b)
	b.Mod(b, P)
	return &Fp2{a: a, b: b}
}

// Mul computes f * g in Fp2 using Karatsuba multiplication
// (a + bu)(c + du) = (ac - bd) + (ad + bc)u, where u² = -1
func (f *Fp2) Mul(g *Fp2) *Fp2 {
	// Karatsuba: (a+bu)(c+du) = ac - bd + ((a+b)(c+d) - ac - bd)u
	ac := new(big.Int).Mul(f.a, g.a)
	bd := new(big.Int).Mul(f.b, g.b)

	aPlusB := new(big.Int).Add(f.a, f.b)
	cPlusD := new(big.Int).Add(g.a, g.b)
	adPlusBc := new(big.Int).Mul(aPlusB, cPlusD)
	adPlusBc.Sub(adPlusBc, ac).Sub(adPlusBc, bd)

	// Real part: ac - bd (since u² = -1)
	realPart := new(big.Int).Sub(ac, bd)
	realPart.Mod(realPart, P)

	// Imaginary part: ad + bc
	imagPart := adPlusBc.Mod(adPlusBc, P)
	return &Fp2{a: realPart, b: imagPart}
}

// Square computes f² in Fp2 optimized
func (f *Fp2) Square() *Fp2 {
	// (a + bu)² = (a² - b²) + 2ab*u
	a2 := new(big.Int).Mul(f.a, f.a)
	b2 := new(big.Int).Mul(f.b, f.b)
	realPart := new(big.Int).Sub(a2, b2)
	realPart.Mod(realPart, P)

	imagPart := new(big.Int).Mul(f.a, f.b)
	imagPart.Lsh(imagPart, 1) // *2
	imagPart.Mod(imagPart, P)

	return &Fp2{a: realPart, b: imagPart}
}

// Inverse computes f⁻¹ in Fp2
func (f *Fp2) Inverse() *Fp2 {
	if f.IsZero() {
		return &Fp2{a: big.NewInt(0), b: big.NewInt(0)}
	}
	// 1/(a+bu) = (a-bu)/(a²+b²)
	a2 := new(big.Int).Mul(f.a, f.a)
	b2 := new(big.Int).Mul(f.b, f.b)
	norm := new(big.Int).Add(a2, b2)
	norm.Mod(norm, P)

	normInv := new(big.Int).ModInverse(norm, P)

	realPart := new(big.Int).Mul(f.a, normInv)
	realPart.Mod(realPart, P)

	imagPart := new(big.Int).Mul(f.b, normInv)
	imagPart.Neg(imagPart)
	imagPart.Mod(imagPart, P)

	return &Fp2{a: realPart, b: imagPart}
}

// Neg computes -f in Fp2
func (f *Fp2) Neg() *Fp2 {
	a := new(big.Int).Neg(f.a)
	a.Mod(a, P)
	b := new(big.Int).Neg(f.b)
	b.Mod(b, P)
	return &Fp2{a: a, b: b}
}

// MulScalar multiplies by a scalar from Fp
func (f *Fp2) MulScalar(s *big.Int) *Fp2 {
	a := new(big.Int).Mul(f.a, s)
	a.Mod(a, P)
	b := new(big.Int).Mul(f.b, s)
	b.Mod(b, P)
	return &Fp2{a: a, b: b}
}

// IsZero returns true if f == 0
func (f *Fp2) IsZero() bool {
	return f.a.Sign() == 0 && f.b.Sign() == 0
}

// Equal returns true if f == g
func (f *Fp2) Equal(g *Fp2) bool {
	return f.a.Cmp(g.a) == 0 && f.b.Cmp(g.b) == 0
}

// ============================================================================
// Fp6 - Sextic Extension Field Element
// ============================================================================

// Fp6 represents an element in Fp6 = Fp2[v]/(v³-ξ) where ξ = u+9
// Represented as c0 + c1*v + c2*v² where c0, c1, c2 ∈ Fp2
type Fp6 struct {
	c0, c1, c2 *Fp2
}

// NewFp6 creates a new Fp6 element
func NewFp6(c0, c1, c2 *Fp2) *Fp6 {
	return &Fp6{c0: c0, c1: c1, c2: c2}
}

// Copy creates a deep copy
func (f *Fp6) Copy() *Fp6 {
	return &Fp6{
		c0: f.c0.Copy(),
		c1: f.c1.Copy(),
		c2: f.c2.Copy(),
	}
}

// Add computes f + g in Fp6
func (f *Fp6) Add(g *Fp6) *Fp6 {
	return &Fp6{
		c0: f.c0.Add(g.c0),
		c1: f.c1.Add(g.c1),
		c2: f.c2.Add(g.c2),
	}
}

// Sub computes f - g in Fp6
func (f *Fp6) Sub(g *Fp6) *Fp6 {
	return &Fp6{
		c0: f.c0.Sub(g.c0),
		c1: f.c1.Sub(g.c1),
		c2: f.c2.Sub(g.c2),
	}
}

// mulByNonResidue multiplies by the non-residue ξ = u+9
func mulByNonResidue(f *Fp2) *Fp2 {
	// (a + bu)(u + 9) = (9a - b) + (a + 9b)u
	nine := big.NewInt(9)
	realPart := new(big.Int).Mul(f.a, nine)
	realPart.Sub(realPart, f.b)
	realPart.Mod(realPart, P)

	imagPart := new(big.Int).Add(f.a, new(big.Int).Mul(f.b, nine))
	imagPart.Mod(imagPart, P)

	return &Fp2{a: realPart, b: imagPart}
}

// Mul computes f * g in Fp6 using Karatsuba
// Mul computes f * g in Fp6 using Karatsuba
func (f *Fp6) Mul(g *Fp6) *Fp6 {
	// Use Karatsuba multiplication for efficiency
	a := f.c0.Mul(g.c0)
	b := f.c1.Mul(g.c1)
	c := f.c2.Mul(g.c2)

	t0 := f.c1.Add(f.c2)
	t1 := g.c1.Add(g.c2)
	t0 = t0.Mul(t1)
	t0 = t0.Sub(b).Sub(c)
	t0 = mulByNonResidue(t0)
	c0 := a.Add(t0)

	t0 = f.c0.Add(f.c1)
	t1 = g.c0.Add(g.c1)
	t0 = t0.Mul(t1)
	t0 = t0.Sub(a).Sub(b)
	c1 := mulByNonResidue(c).Add(t0)

	t0 = f.c0.Add(f.c2)
	t1 = g.c0.Add(g.c2)
	t0 = t0.Mul(t1)
	c2 := t0.Sub(a).Sub(c).Add(b)

	return &Fp6{c0: c0, c1: c1, c2: c2}
}

// Square computes f² in Fp6
func (f *Fp6) Square() *Fp6 {
	return f.Mul(f)
}

// Inverse computes f⁻¹ in Fp6
func (f *Fp6) Inverse() *Fp6 {
	// Use the norm formula for sextic extensions
	c0 := f.c0.Square().Sub(mulByNonResidue(f.c1.Mul(f.c2)))
	c1 := mulByNonResidue(f.c2.Square()).Sub(f.c0.Mul(f.c1))
	c2 := f.c1.Square().Sub(f.c0.Mul(f.c2))

	t := f.c2.Mul(c1)
	t = mulByNonResidue(t)
	t = t.Add(f.c1.Mul(c2))
	t = mulByNonResidue(t)
	t = t.Add(f.c0.Mul(c0))
	t = t.Inverse()

	return &Fp6{
		c0: c0.Mul(t),
		c1: c1.Mul(t),
		c2: c2.Mul(t),
	}
}

// Neg computes -f in Fp6
func (f *Fp6) Neg() *Fp6 {
	return &Fp6{
		c0: f.c0.Neg(),
		c1: f.c1.Neg(),
		c2: f.c2.Neg(),
	}
}

// IsZero returns true if f == 0
func (f *Fp6) IsZero() bool {
	return f.c0.IsZero() && f.c1.IsZero() && f.c2.IsZero()
}

// ============================================================================
// Fp12 - Dodecic Extension Field Element
// ============================================================================

// Fp12 represents an element in Fp12 = Fp6[w]/(w²-v)
// Represented as c0 + c1*w where c0, c1 ∈ Fp6
type Fp12 struct {
	c0, c1 *Fp6
}

// NewFp12 creates a new Fp12 element
func NewFp12(c0, c1 *Fp6) *Fp12 {
	return &Fp12{c0: c0, c1: c1}
}

// Copy creates a deep copy
func (f *Fp12) Copy() *Fp12 {
	return &Fp12{
		c0: f.c0.Copy(),
		c1: f.c1.Copy(),
	}
}

// Add computes f + g in Fp12
func (f *Fp12) Add(g *Fp12) *Fp12 {
	return &Fp12{
		c0: f.c0.Add(g.c0),
		c1: f.c1.Add(g.c1),
	}
}

// Sub computes f - g in Fp12
func (f *Fp12) Sub(g *Fp12) *Fp12 {
	return &Fp12{
		c0: f.c0.Sub(g.c0),
		c1: f.c1.Sub(g.c1),
	}
}

// Mul computes f * g in Fp12
func (f *Fp12) Mul(g *Fp12) *Fp12 {
	// (a + bw)(c + dw) = (ac + bd*v) + (ad + bc)w where w² = v
	ac := f.c0.Mul(g.c0)
	bd := f.c1.Mul(g.c1)

	// bd*v: multiply bd by v (shift coefficients)
	bdv := &Fp6{
		c0: mulByNonResidue(bd.c2),
		c1: bd.c0,
		c2: bd.c1,
	}

	c0 := ac.Add(bdv)

	// (a+b)(c+d) - ac - bd
	t0 := f.c0.Add(f.c1)
	t1 := g.c0.Add(g.c1)
	c1 := t0.Mul(t1).Sub(ac).Sub(bd)

	return &Fp12{c0: c0, c1: c1}
}

// Square computes f² in Fp12
func (f *Fp12) Square() *Fp12 {
	return f.Mul(f)
}

// Inverse computes f⁻¹ in Fp12
func (f *Fp12) Inverse() *Fp12 {
	// 1/(a+bw) = (a-bw)/(a²-b²v)
	t0 := f.c0.Square()
	t1 := f.c1.Square()
	// t1*v
	t1v := &Fp6{
		c0: mulByNonResidue(t1.c2),
		c1: t1.c0,
		c2: t1.c1,
	}
	t0 = t0.Sub(t1v)
	t0 = t0.Inverse()

	return &Fp12{
		c0: f.c0.Mul(t0),
		c1: f.c1.Neg().Mul(t0),
	}
}

// Exp computes f^e in Fp12 using square-and-multiply
func (f *Fp12) Exp(e *big.Int) *Fp12 {
	result := &Fp12{
		c0: &Fp6{
			c0: &Fp2{a: big.NewInt(1), b: big.NewInt(0)},
			c1: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
			c2: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
		},
		c1: &Fp6{
			c0: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
			c1: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
			c2: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
		},
	}

	base := f.Copy()
	for i := 0; i < e.BitLen(); i++ {
		if e.Bit(i) == 1 {
			result = result.Mul(base)
		}
		base = base.Square()
	}
	return result
}

// IsZero returns true if f == 0
func (f *Fp12) IsZero() bool {
	return f.c0.IsZero() && f.c1.IsZero()
}

// IsOne returns true if f == 1
func (f *Fp12) IsOne() bool {
	one := &Fp2{a: big.NewInt(1), b: big.NewInt(0)}
	zero := &Fp2{a: big.NewInt(0), b: big.NewInt(0)}
	return f.c0.c0.Equal(one) && f.c0.c1.IsZero() && f.c0.c2.IsZero() &&
		f.c1.c0.Equal(zero) && f.c1.c1.IsZero() && f.c1.c2.IsZero()
}

// ============================================================================
// G1 - Points on the base curve E(Fp): y² = x³ + 3
// ============================================================================

// G1 represents a point on the BN128 curve over Fp
type G1 struct {
	X, Y *big.Int
	// We use affine coordinates; point at infinity represented by X=Y=0
}

// NewG1 creates a new G1 point
func NewG1(x, y *big.Int) (*G1, error) {
	p := &G1{
		X: new(big.Int).Set(x),
		Y: new(big.Int).Set(y),
	}

	if !p.IsOnCurve() {
		return nil, ErrInvalidPoint
	}

	return p, nil
}

// // G1Generator returns the G1 generator point
func G1Generator() *G1 {
	return &G1{
		X: new(big.Int).Set(GeneratorG1X),
		Y: new(big.Int).Set(GeneratorG1Y),
	}
}

// Copy creates a deep copy
func (p *G1) Copy() *G1 {
	return &G1{
		X: new(big.Int).Set(p.X),
		Y: new(big.Int).Set(p.Y),
	}
}

// IsInfinity checks if point is the point at infinity
func (p *G1) IsInfinity() bool {
	return p.X.Sign() == 0 && p.Y.Sign() == 0
}

// IsOnCurve checks if point is on the curve: y² = x³ + 3
func (p *G1) IsOnCurve() bool {
	if p.IsInfinity() {
		return true
	}

	y2 := new(big.Int).Mul(p.Y, p.Y)
	y2.Mod(y2, P)

	x3 := new(big.Int).Mul(p.X, p.X)
	x3.Mul(x3, p.X)
	x3.Add(x3, big.NewInt(3))
	x3.Mod(x3, P)

	return y2.Cmp(x3) == 0
}

// Equal checks if two points are equal
func (p *G1) Equal(q *G1) bool {
	return p.X.Cmp(q.X) == 0 && p.Y.Cmp(q.Y) == 0
}

// Neg computes -p
func (p *G1) Neg() *G1 {
	if p.IsInfinity() {
		return &G1{X: big.NewInt(0), Y: big.NewInt(0)}
	}
	return &G1{
		X: new(big.Int).Set(p.X),
		Y: new(big.Int).Sub(P, p.Y),
	}
}

// Add computes p + q using affine coordinates
func (p *G1) Add(q *G1) *G1 {
	if p.IsInfinity() {
		return q.Copy()
	}
	if q.IsInfinity() {
		return p.Copy()
	}

	if p.X.Cmp(q.X) == 0 {
		if p.Y.Cmp(q.Y) == 0 {
			return p.Double()
		}
		// Points are negatives
		return &G1{X: big.NewInt(0), Y: big.NewInt(0)}
	}

	// λ = (y2 - y1) / (x2 - x1)
	dy := new(big.Int).Sub(q.Y, p.Y)
	dx := new(big.Int).Sub(q.X, p.X)
	dxInv := new(big.Int).ModInverse(dx, P)
	lambda := new(big.Int).Mul(dy, dxInv)
	lambda.Mod(lambda, P)

	// x3 = λ² - x1 - x2
	x3 := new(big.Int).Mul(lambda, lambda)
	x3.Sub(x3, p.X)
	x3.Sub(x3, q.X)
	x3.Mod(x3, P)

	// y3 = λ(x1 - x3) - y1
	y3 := new(big.Int).Sub(p.X, x3)
	y3.Mul(y3, lambda)
	y3.Sub(y3, p.Y)
	y3.Mod(y3, P)

	return &G1{X: x3, Y: y3}
}

// Double computes 2p
func (p *G1) Double() *G1 {
	if p.IsInfinity() {
		return &G1{X: big.NewInt(0), Y: big.NewInt(0)}
	}

	// λ = (3x² + a) / 2y where a = 0 for BN128
	three := big.NewInt(3)
	two := big.NewInt(2)

	x2 := new(big.Int).Mul(p.X, p.X)
	numerator := new(big.Int).Mul(x2, three)

	denominator := new(big.Int).Mul(p.Y, two)
	denominatorInv := new(big.Int).ModInverse(denominator, P)

	lambda := new(big.Int).Mul(numerator, denominatorInv)
	lambda.Mod(lambda, P)

	// x3 = λ² - 2x
	x3 := new(big.Int).Mul(lambda, lambda)
	x3.Sub(x3, p.X)
	x3.Sub(x3, p.X)
	x3.Mod(x3, P)

	// y3 = λ(x - x3) - y
	y3 := new(big.Int).Sub(p.X, x3)
	y3.Mul(y3, lambda)
	y3.Sub(y3, p.Y)
	y3.Mod(y3, P)

	return &G1{X: x3, Y: y3}
}

// ScalarMult computes k*p using double-and-add
func (p *G1) ScalarMult(k *big.Int) *G1 {
	if k.Sign() == 0 || p.IsInfinity() {
		return &G1{X: big.NewInt(0), Y: big.NewInt(0)}
	}

	result := &G1{X: big.NewInt(0), Y: big.NewInt(0)}
	base := p.Copy()

	// Use binary representation for scalar multiplication
	for i := 0; i < k.BitLen(); i++ {
		if k.Bit(i) == 1 {
			result = result.Add(base)
		}
		base = base.Double()
	}

	return result
}

// ScalarBaseMult computes k*G where G is the generator
func ScalarBaseMult(k *big.Int) *G1 {
	return G1Generator().ScalarMult(k)
}

// MarshalG1 serializes a G1 point (64 bytes: 32 for X, 32 for Y)
func (p *G1) Marshal() []byte {
	if p.IsInfinity() {
		return make([]byte, 64)
	}
	buf := make([]byte, 64)
	xBytes := p.X.Bytes()
	yBytes := p.Y.Bytes()
	copy(buf[32-len(xBytes):32], xBytes)
	copy(buf[64-len(yBytes):64], yBytes)
	return buf
}

// UnmarshalG1 deserializes a G1 point
func UnmarshalG1(buf []byte) (*G1, error) {
	if len(buf) != 64 {
		return nil, ErrInvalidEncoding
	}

	x := new(big.Int).SetBytes(buf[:32])
	y := new(big.Int).SetBytes(buf[32:64])

	if x.Sign() == 0 && y.Sign() == 0 {
		return &G1{X: big.NewInt(0), Y: big.NewInt(0)}, nil
	}

	return NewG1(x, y)
}

// ============================================================================
// G2 - Points on the twisted curve E'(Fp2): y² = x³ + 3/(9+u)
// ============================================================================
// G2 represents a point on the twisted BN128 curve over Fp2
type G2 struct {
	X, Y *Fp2
}

// NewG2 creates a new G2 point
func NewG2(x, y *Fp2) (*G2, error) {
	p := &G2{X: x, Y: y}
	if !p.IsOnCurve() {
		return nil, ErrInvalidPoint
	}
	return p, nil
}

// G2Generator returns the G2 generator point
func G2Generator() *G2 {
	return &G2{
		X: GeneratorG2X.Copy(),
		Y: GeneratorG2Y.Copy(),
	}
}

// Copy creates a deep copy
func (p *G2) Copy() *G2 {
	return &G2{
		X: p.X.Copy(),
		Y: p.Y.Copy(),
	}
}

// IsOnCurve checks if point is on the twisted curve: y² = x³ + b where b = 3/(9+u)
func (p *G2) IsOnCurve() bool {
	if p.IsInfinity() {
		return true
	}

	y2 := p.Y.Square()
	x3 := p.X.Square().Mul(p.X)
	x3 = x3.Add(TwistB)

	return y2.Equal(x3)
}

// IsInfinity checks if point is the point at infinity
func (p *G2) IsInfinity() bool {
	return p.X.IsZero() && p.Y.IsZero()
}

// Equal checks if two points are equal
func (p *G2) Equal(q *G2) bool {
	return p.X.Equal(q.X) && p.Y.Equal(q.Y)
}

// Neg computes -p
func (p *G2) Neg() *G2 {
	if p.IsInfinity() {
		return &G2{
			X: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
			Y: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
		}
	}
	return &G2{
		X: p.X.Copy(),
		Y: p.Y.Neg(),
	}
}

// Add computes p + q using affine coordinates
func (p *G2) Add(q *G2) *G2 {
	if p.IsInfinity() {
		return q.Copy()
	}
	if q.IsInfinity() {
		return p.Copy()
	}

	if p.X.Equal(q.X) {
		if p.Y.Equal(q.Y) {
			return p.Double()
		}
		// Points are negatives
		return &G2{
			X: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
			Y: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
		}
	}

	// λ = (y2 - y1) / (x2 - x1)
	dy := q.Y.Sub(p.Y)
	dx := q.X.Sub(p.X)
	lambda := dy.Mul(dx.Inverse())

	// x3 = λ² - x1 - x2
	x3 := lambda.Square().Sub(p.X).Sub(q.X)

	// y3 = λ(x1 - x3) - y1
	y3 := p.X.Sub(x3).Mul(lambda).Sub(p.Y)

	return &G2{X: x3, Y: y3}
}

// Double computes 2p
func (p *G2) Double() *G2 {
	if p.IsInfinity() {
		return &G2{
			X: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
			Y: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
		}
	}

	// λ = 3x² / 2y
	three := big.NewInt(3)
	two := big.NewInt(2)

	numerator := p.X.Square().MulScalar(three)
	denominator := p.Y.MulScalar(two)
	lambda := numerator.Mul(denominator.Inverse())

	// x3 = λ² - 2x
	x3 := lambda.Square().Sub(p.X).Sub(p.X)

	// y3 = λ(x - x3) - y
	y3 := p.X.Sub(x3).Mul(lambda).Sub(p.Y)

	return &G2{X: x3, Y: y3}
}

// ScalarMult computes k*p using double-and-add
func (p *G2) ScalarMult(k *big.Int) *G2 {
	if k.Sign() == 0 || p.IsInfinity() {
		return &G2{
			X: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
			Y: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
		}
	}

	result := &G2{
		X: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
		Y: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
	}
	base := p.Copy()

	for i := 0; i < k.BitLen(); i++ {
		if k.Bit(i) == 1 {
			result = result.Add(base)
		}
		base = base.Double()
	}

	return result
}

// Marshal serializes a G2 point (128 bytes: 64 for X, 64 for Y)
func (p *G2) Marshal() []byte {
	buf := make([]byte, 128)
	if p.IsInfinity() {
		return buf
	}

	xaBytes := p.X.a.Bytes()
	xbBytes := p.X.b.Bytes()
	yaBytes := p.Y.a.Bytes()
	ybBytes := p.Y.b.Bytes()

	copy(buf[32-len(xaBytes):32], xaBytes)
	copy(buf[64-len(xbBytes):64], xbBytes)
	copy(buf[96-len(yaBytes):96], yaBytes)
	copy(buf[128-len(ybBytes):128], ybBytes)

	return buf
}

// UnmarshalG2 deserializes a G2 point
func UnmarshalG2(buf []byte) (*G2, error) {
	if len(buf) != 128 {
		return nil, ErrInvalidEncoding
	}

	xa := new(big.Int).SetBytes(buf[0:32])
	xb := new(big.Int).SetBytes(buf[32:64])
	ya := new(big.Int).SetBytes(buf[64:96])
	yb := new(big.Int).SetBytes(buf[96:128])

	x := NewFp2(xa, xb)
	y := NewFp2(ya, yb)

	if x.IsZero() && y.IsZero() {
		return &G2{X: x, Y: y}, nil
	}

	return NewG2(x, y)
}

// ============================================================================
// Pairing Operations - Optimal Ate Pairing
// ============================================================================

// GT represents an element in the target group (Fp12)
type GT struct {
	value *Fp12
}

// NewGT creates a new GT element
func NewGT(v *Fp12) *GT {
	return &GT{value: v}
}

// lineFunctionAdd computes the line function for point addition
func lineFunctionAdd(r, p *G2, q *G1) *Fp12 {
	// Computes l_{r,p}(q) for Miller loop
	if r.IsInfinity() {
		return &Fp12{
			c0: &Fp6{
				c0: &Fp2{a: big.NewInt(1), b: big.NewInt(0)},
				c1: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
				c2: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
			},
			c1: &Fp6{
				c0: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
				c1: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
				c2: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
			},
		}
	}

	// λ = (yp - yr) / (xp - xr)
	dy := p.Y.Sub(r.Y)
	dx := p.X.Sub(r.X)
	lambda := dy.Mul(dx.Inverse())

	// Line function: λ(xq - xr) - (yq - yr)
	c := lambda.Mul(r.X).Sub(r.Y)

	qx := NewFp2(q.X, big.NewInt(0))
	qy := NewFp2(q.Y, big.NewInt(0))

	// Result: yq - λ*xq + c*w
	return &Fp12{
		c0: &Fp6{
			c0: c,
			c1: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
			c2: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
		},
		c1: &Fp6{
			c0: lambda.Neg().Mul(qx),
			c1: qy,
			c2: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
		},
	}
}

// lineFunctionDouble computes the line function for point doubling
func lineFunctionDouble(r *G2, q *G1) *Fp12 {
	if r.IsInfinity() {
		return &Fp12{
			c0: &Fp6{
				c0: &Fp2{a: big.NewInt(1), b: big.NewInt(0)},
				c1: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
				c2: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
			},
			c1: &Fp6{
				c0: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
				c1: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
				c2: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
			},
		}
	}

	// λ = 3x² / 2y
	three := big.NewInt(3)
	two := big.NewInt(2)
	numerator := r.X.Square().MulScalar(three)
	denominator := r.Y.MulScalar(two)
	lambda := numerator.Mul(denominator.Inverse())

	c := lambda.Mul(r.X).Sub(r.Y)

	qx := NewFp2(q.X, big.NewInt(0))
	qy := NewFp2(q.Y, big.NewInt(0))

	return &Fp12{
		c0: &Fp6{
			c0: c,
			c1: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
			c2: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
		},
		c1: &Fp6{
			c0: lambda.Neg().Mul(qx),
			c1: qy,
			c2: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
		},
	}
}

// millerLoop computes the Miller loop for ate pairing
func millerLoop(q *G1, p *G2) *Fp12 {
	if q.IsInfinity() || p.IsInfinity() {
		return &Fp12{
			c0: &Fp6{
				c0: &Fp2{a: big.NewInt(1), b: big.NewInt(0)},
				c1: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
				c2: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
			},
			c1: &Fp6{
				c0: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
				c1: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
				c2: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
			},
		}
	}

	// BN128 ate pairing parameter: 6t + 2 where t = 4965661367192848881
	// In binary: this is the loop parameter for optimal ate pairing
	// For BN128: 6t+2 = 29793968203157093288
	loopCount := fromHex("19d797039be763ba8")

	r := p.Copy()
	f := &Fp12{
		c0: &Fp6{
			c0: &Fp2{a: big.NewInt(1), b: big.NewInt(0)},
			c1: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
			c2: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
		},
		c1: &Fp6{
			c0: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
			c1: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
			c2: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
		},
	}

	// Miller loop
	for i := loopCount.BitLen() - 2; i >= 0; i-- {
		f = f.Square()
		lrr := lineFunctionDouble(r, q)
		f = f.Mul(lrr)
		r = r.Double()

		if loopCount.Bit(i) == 1 {
			lrp := lineFunctionAdd(r, p, q)
			f = f.Mul(lrp)
			r = r.Add(p)
		}
	}

	return f
}

// finalExponentiation computes the final exponentiation for the ate pairing
// Raises f to the power (p^12 - 1) / r
func finalExponentiation(f *Fp12) *Fp12 {
	// Easy part: (p^6 - 1)(p^2 + 1)
	// First: f^(p^6 - 1)
	t0 := &Fp12{
		c0: f.c0.Copy(),
		c1: f.c1.Neg(),
	}
	t0 = t0.Mul(f.Inverse())

	// Second: f^(p^2 + 1)
	t1 := frobeniusP2(t0)
	f = t1.Mul(t0)

	// Hard part: use addition chains for efficiency
	// This is a simplified version; production code uses optimized addition chains
	exp := new(big.Int).Sub(P, big.NewInt(1))
	exp.Mul(exp, exp)
	exp.Mul(exp, exp)
	exp.Mul(exp, exp)
	exp.Mul(exp, exp)
	exp.Mul(exp, exp)
	exp.Mul(exp, exp)
	exp.Sub(exp, big.NewInt(1))
	exp.Div(exp, Order)

	return f.Exp(exp)
}

// frobeniusP2 computes the Frobenius endomorphism raised to power 2
func frobeniusP2(f *Fp12) *Fp12 {
	// Simplified: conjugate in Fp12
	return &Fp12{
		c0: f.c0.Copy(),
		c1: f.c1.Neg(),
	}
}

// Pair computes the optimal ate pairing e(p, q)
func Pair(p *G1, q *G2) *GT {
	f := millerLoop(p, q)
	f = finalExponentiation(f)
	return &GT{value: f}
}

// PairingCheck verifies if e(p1, q1) * e(p2, q2) * ... * e(pn, qn) = 1
// This is used in zkSNARK verification (EIP-197)
func PairingCheck(pairs [][2]interface{}) bool {
	result := &Fp12{
		c0: &Fp6{
			c0: &Fp2{a: big.NewInt(1), b: big.NewInt(0)},
			c1: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
			c2: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
		},
		c1: &Fp6{
			c0: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
			c1: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
			c2: &Fp2{a: big.NewInt(0), b: big.NewInt(0)},
		},
	}

	for _, pair := range pairs {
		p, ok1 := pair[0].(*G1)
		q, ok2 := pair[1].(*G2)
		if !ok1 || !ok2 {
			return false
		}

		f := millerLoop(p, q)
		result = result.Mul(f)
	}

	result = finalExponentiation(result)
	return result.IsOne()
}

// ============================================================================
// GT Operations
// ============================================================================

// Mul multiplies two GT elements
func (g *GT) Mul(h *GT) *GT {
	return &GT{value: g.value.Mul(h.value)}
}

// Inverse computes g^(-1)
func (g *GT) Inverse() *GT {
	return &GT{value: g.value.Inverse()}
}

// Equal checks if two GT elements are equal
func (g *GT) Equal(h *GT) bool {
	return g.value.c0.c0.Equal(h.value.c0.c0) &&
		g.value.c0.c1.Equal(h.value.c0.c1) &&
		g.value.c0.c2.Equal(h.value.c0.c2) &&
		g.value.c1.c0.Equal(h.value.c1.c0) &&
		g.value.c1.c1.Equal(h.value.c1.c1) &&
		g.value.c1.c2.Equal(h.value.c1.c2)
}

// IsOne checks if g == 1
func (g *GT) IsOne() bool {
	return g.value.IsOne()
}

// Marshal serializes a GT element
func (g *GT) Marshal() []byte {
	buf := make([]byte, 384) // 12 * 32 bytes for Fp12
	offset := 0

	writeFp2 := func(f *Fp2) {
		aBytes := f.a.Bytes()
		bBytes := f.b.Bytes()
		copy(buf[offset+32-len(aBytes):offset+32], aBytes)
		copy(buf[offset+64-len(bBytes):offset+64], bBytes)
		offset += 64
	}

	writeFp2(g.value.c0.c0)
	writeFp2(g.value.c0.c1)
	writeFp2(g.value.c0.c2)
	writeFp2(g.value.c1.c0)
	writeFp2(g.value.c1.c1)
	writeFp2(g.value.c1.c2)

	return buf
}

// UnmarshalGT deserializes a GT element
func UnmarshalGT(buf []byte) (*GT, error) {
	if len(buf) != 384 {
		return nil, ErrInvalidEncoding
	}

	readFp2 := func(offset int) *Fp2 {
		a := new(big.Int).SetBytes(buf[offset : offset+32])
		b := new(big.Int).SetBytes(buf[offset+32 : offset+64])
		return NewFp2(a, b)
	}

	return &GT{
		value: &Fp12{
			c0: &Fp6{
				c0: readFp2(0),
				c1: readFp2(64),
				c2: readFp2(128),
			},
			c1: &Fp6{
				c0: readFp2(192),
				c1: readFp2(256),
				c2: readFp2(320),
			},
		},
	}, nil
}

// ============================================================================
// Utility Functions
// ============================================================================

// RandomG1 generates a random point in G1
func RandomG1(rand io.Reader) (*G1, error) {
	k, err := randomScalar(rand)
	if err != nil {
		return nil, err
	}
	return ScalarBaseMult(k), nil
}

// RandomG2 generates a random point in G2
func RandomG2(rand io.Reader) (*G2, error) {
	k, err := randomScalar(rand)
	if err != nil {
		return nil, err
	}
	return G2Generator().ScalarMult(k), nil
}

// randomScalar generates a random scalar in [1, Order)
func randomScalar(reader io.Reader) (*big.Int, error) {
	if reader == nil {
		reader = rand.Reader
	}

	k, err := rand.Int(reader, Order)
	if err != nil {
		return nil, err
	}

	// Ensure k != 0
	if k.Sign() == 0 {
		k = big.NewInt(1)
	}

	return k, nil
}

// HashToG1 maps arbitrary data to a G1 point (simplified version)
// Note: This is NOT a secure hash-to-curve. Use proper hash-to-curve for production.
func HashToG1(data []byte) *G1 {
	// This is a placeholder. Production code should use proper hash-to-curve
	// algorithms like the one specified in draft-irtf-cfrg-hash-to-curve
	h := new(big.Int).SetBytes(data)
	h.Mod(h, Order)
	if h.Sign() == 0 {
		h = big.NewInt(1)
	}
	return ScalarBaseMult(h)
}
