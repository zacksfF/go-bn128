# Technical Deep Dive: BN128 Curve Arithmetic

## What is BN128?

BN128 (also called BN254 or alt_bn128) is a pairing-friendly elliptic curve used in zero-knowledge proofs. It's the backbone of zkSNARKs in Ethereum, Zcash, and many privacy-focused blockchains.

**Why BN128?** It lets us do fancy math operations called "pairings" that make zkSNARKs possible.

## The Math Stack

Think of BN128 as a tower of mathematical structures, each built on top of the previous one:

```
Fp (base field)
  ↓
Fp2 (quadratic extension)
  ↓
Fp6 (sextic extension)
  ↓
Fp12 (dodecic extension)
```

Let's break down each layer.

---

## 1. Fp - The Base Field

**What it is**: Numbers modulo a big prime p.

**The prime**:
```
p = 21888242871839275222246405745257275088696311157297823662689037894645226208583
```

This is a 254-bit prime. All our arithmetic happens modulo this number.

**Operations**:
- Addition: $(a + b) \bmod p$
- Multiplication: $(a \times b) \bmod p$
- Inversion: $a^{-1}$ such that $a \times a^{-1} \equiv 1 \pmod{p}$

**Why this prime?** It's specially chosen to make pairing computations efficient.

**Code example**:
```go
type Fp struct {
    n *big.Int  // value mod p
}

func (f *Fp) Add(g *Fp) *Fp {
    result := new(big.Int).Add(f.n, g.n)
    return &Fp{n: result.Mod(result, P)}
}
```

**Complexity**: $O(n^2)$ for multiplication using big.Int (where n = bit length)

---

## 2. Fp2 - Quadratic Extension

**What it is**: Complex numbers over Fp, where $i^2 = -1$.

**Representation**: $a + b \cdot u$ where $u^2 = -1$

**Why?** We need Fp2 to define the twisted curve for G2.

**Math**:
- Addition: $(a_1 + b_1u) + (a_2 + b_2u) = (a_1 + a_2) + (b_1 + b_2)u$
- Multiplication: $(a_1 + b_1u)(a_2 + b_2u) = (a_1a_2 - b_1b_2) + (a_1b_2 + b_1a_2)u$
  - We use **Karatsuba** to save multiplications

**Karatsuba optimization**:
```
Instead of 4 multiplications: a₁a₂, a₁b₂, b₁a₂, b₁b₂
We do 3: a₁a₂, b₁b₂, (a₁+b₁)(a₂+b₂)
```

**Code snippet**:
```go
type Fp2 struct {
    a, b *big.Int  // represents a + b*u
}

func (f *Fp2) Mul(g *Fp2) *Fp2 {
    // Karatsuba: (a+bu)(c+du) = (ac-bd) + ((a+b)(c+d)-ac-bd)u
    ac := new(big.Int).Mul(f.a, g.a)
    bd := new(big.Int).Mul(f.b, g.b)
    
    aPlusB := new(big.Int).Add(f.a, f.b)
    cPlusD := new(big.Int).Add(g.a, g.b)
    adPlusBc := new(big.Int).Mul(aPlusB, cPlusD)
    adPlusBc.Sub(adPlusBc, ac).Sub(adPlusBc, bd)
    
    realPart := new(big.Int).Sub(ac, bd).Mod(realPart, P)
    imagPart := adPlusBc.Mod(adPlusBc, P)
    
    return &Fp2{a: realPart, b: imagPart}
}
```

**Complexity**: $O(n^2)$ for field operations

---

## 3. Fp6 - Sextic Extension

**What it is**: Fp6 = Fp2[v] where $v^3 = \xi$ and $\xi = u + 9$.

**Representation**: $c_0 + c_1v + c_2v^2$ where each $c_i \in \text{Fp2}$

**Why?** Fp6 is a stepping stone to build Fp12, which is where pairing outputs live.

**Key operation - Multiplication by non-residue**:
```
mulByNonResidue(a + bu) = (9a - b) + (a + 9b)u
```

This trick makes Fp6 multiplication faster.

**Math**:
For multiplication $(c_0 + c_1v + c_2v^2)(d_0 + d_1v + d_2v^2)$, we use Karatsuba again:

$$
\text{result} = \begin{cases}
r_0 = c_0d_0 + \xi(c_1d_2 + c_2d_1) \\
r_1 = c_0d_1 + c_1d_0 + \xi(c_2d_2) \\
r_2 = c_0d_2 + c_1d_1 + c_2d_0
\end{cases}
$$

**Complexity**: $O(n^2)$ with Karatsuba optimizations

---

## 4. Fp12 - The Target Group

**What it is**: Fp12 = Fp6[w] where $w^2 = v$.

**Representation**: $c_0 + c_1w$ where $c_0, c_1 \in \text{Fp6}$

**Why?** Pairings output elements in Fp12. This is where GT lives.

**Math**:
Multiplication: $(a + bw)(c + dw) = (ac + bdv) + (ad + bc)w$

Note: $w^2 = v$, so $bdw^2 = bdv$

**Exponentiation** is crucial here (used in final exponentiation):
```
f^e using square-and-multiply
```

**Complexity**: 
- Multiplication: $O(n^2)$
- Exponentiation: $O(n^3)$ where n = bit length of exponent

---

## 5. G1 - The First Curve Group

**Curve equation**: $y^2 = x^3 + 3$ over Fp

**Generator point**:
```
G1 = (1, 2)
```

**Order** (number of points):
```
r = 21888242871839275222246405745257275088548364400416034343698204186575808495617
```

**Point operations**:

1. **Point Addition** (affine coordinates):
```
Given P = (x₁, y₁) and Q = (x₂, y₂)
If x₁ ≠ x₂:
  λ = (y₂ - y₁) / (x₂ - x₁)
  x₃ = λ² - x₁ - x₂
  y₃ = λ(x₁ - x₃) - y₁
Result: R = (x₃, y₃)
```

2. **Point Doubling**:
```
λ = (3x₁²) / (2y₁)
x₃ = λ² - 2x₁
y₃ = λ(x₁ - x₃) - y₁
```

3. **Scalar Multiplication** (double-and-add):
```
To compute k·P:
  result = O (point at infinity)
  for each bit i in k (from LSB to MSB):
    if bit i is 1:
      result = result + P
    P = 2·P
  return result
```

**Complexity**:
- Addition: $O(n^2)$ (due to field operations)
- Doubling: $O(n^2)$
- Scalar multiplication: $O(n^3)$ where n = scalar bit length (typically 256)

**Why affine?** Simpler formulas, easier to verify. Projective coordinates are faster but more complex.

---

## 6. G2 - The Twisted Curve Group

**Curve equation**: $y^2 = x^3 + b'$ over Fp2

Where $b' = 3/(9+u) = (19485874751759354771024239261021720505790618469301721065564631296452457478373 + 266929791119991161246907387137283842545076965332900288569378510910307636690u)$

**Generator point**: Lives in Fp2 × Fp2 (each coordinate is a pair of Fp elements)

**Why a twisted curve?** 
- We can't use the same curve as G1 over Fp2
- The twist gives us the right structure for pairings
- Points in G2 are larger (128 bytes vs 64 bytes for G1)

**Operations**: Same formulas as G1, but arithmetic is in Fp2

**Complexity**: ~4x slower than G1 (because Fp2 operations are more expensive)

---

## 7. Pairing - The Magic Operation

**What is a pairing?**

A pairing is a bilinear map:
$$
e: G1 \times G2 \to GT
$$

**Bilinearity** means:
$$
e(aP, bQ) = e(P, Q)^{ab}
$$

This property is what makes zkSNARKs work!

**The Optimal Ate Pairing**:

The pairing computation has two main phases:

### Phase 1: Miller Loop

The Miller loop computes:
$$
f_{s,Q}(P)
$$

where s is the ate pairing parameter.

**For BN128**: $s = 6t + 2$ where $t = 4965661367192848881$

**Algorithm**:
```
function millerLoop(P ∈ G1, Q ∈ G2):
    R = Q
    f = 1
    
    for each bit i in s (from MSB to LSB-1):
        f = f² · line_{R,R}(P)
        R = 2R
        
        if bit i == 1:
            f = f · line_{R,Q}(P)
            R = R + Q
    
    return f
```

The **line function** evaluates the line through two points at P.

### Phase 2: Final Exponentiation

Take the Miller loop output and raise it to:
$$
e = \frac{p^{12} - 1}{r}
$$

This is a HUGE exponent, so we decompose it:

**Easy part**: $(p^6 - 1)(p^2 + 1)$

**Hard part**: Use efficient addition chains (see Vercauteren's algorithm)

**Why final exponentiation?**
- Ensures the result has order r
- Kills unwanted elements
- Makes the pairing unique

**Complexity**:
- Miller loop: $O(n^4)$ where n = log(r) ≈ 256
- Final exponentiation: $O(n^4)$
- **Total**: ~10-20ms on modern hardware

---

## 8. GT - The Target Group

**What it is**: A subgroup of Fp12 with order r.

**Elements**: Results of pairings live here.

**Operations**:
- Multiplication: Fp12 multiplication
- Inversion: Fp12 inversion
- Exponentiation: Using square-and-multiply

**Why GT?**
- Pairing outputs are in GT
- We can multiply pairing results: $e(P_1,Q_1) \cdot e(P_2,Q_2)$
- zkSNARK verification equations live in GT

**Example pairing check**:
```
Verify: e(A, B) · e(C, D) = 1

Algorithm:
  f1 = millerLoop(A, B)
  f2 = millerLoop(C, D)
  f = f1 · f2
  f = finalExp(f)
  return f == 1
```

---

## Algorithm Design Patterns

### 1. Double-and-Add (Scalar Multiplication)

**Idea**: Process scalar bit-by-bit

```
k·P where k = 1011₂ = 11

Binary:  1  0  1  1
         ↓  ↓  ↓  ↓
Step 0:  P
Step 1:  2P (double)
Step 2:  4P + P = 5P (double + add)
Step 3:  10P + P = 11P (double + add)
```

**Complexity**: $O(\log k)$ point operations

### 2. Karatsuba Multiplication

**Standard**: 4 multiplications for (a+bu)(c+du)

**Karatsuba**: Only 3 multiplications
```
ac, bd, (a+b)(c+d)
Then: (a+b)(c+d) - ac - bd gives us (ad + bc)
```

**Savings**: ~25% fewer multiplications

### 3. Square-and-Multiply (Exponentiation)

**Idea**: Similar to double-and-add but for exponentiation

```
f^e where e = 101₂

Binary:  1  0  1
         ↓  ↓  ↓
Step 0:  f
Step 1:  f² (square)
Step 2:  f⁴ · f = f⁵ (square + multiply)
```

**Complexity**: $O(\log e)$ squarings and multiplications

### 4. Addition Chains (Final Exponentiation)

Instead of computing $x^e$ naively, we find a smart sequence:

```
Example: x^15 = x^(8+4+2+1)
Sequence: x, x², x⁴, x⁸, x⁸·x⁴·x²·x
Only 7 operations vs 14 for naive
```

This is critical for final exponentiation performance.

---

## Complexity Theory Summary

| Operation | Time Complexity | Space Complexity |
|-----------|-----------------|------------------|
| Fp addition | $O(n)$ | $O(n)$ |
| Fp multiplication | $O(n^2)$ | $O(n)$ |
| Fp inversion | $O(n^3)$ | $O(n)$ |
| Fp2 multiplication | $O(n^2)$ | $O(n)$ |
| Fp6 multiplication | $O(n^2)$ | $O(n)$ |
| Fp12 multiplication | $O(n^2)$ | $O(n)$ |
| G1 addition | $O(n^2)$ | $O(n)$ |
| G1 scalar mult | $O(n^3)$ | $O(n)$ |
| G2 addition | $O(n^2)$ | $O(n)$ |
| G2 scalar mult | $O(n^3)$ | $O(n)$ |
| Pairing (Miller) | $O(n^4)$ | $O(n)$ |
| Final exp | $O(n^4)$ | $O(n)$ |
| **Full pairing** | **$O(n^4)$** | **$O(n)$** |

Where n = bit length (typically 256 for BN128)

---

## Security Considerations

**Security level**: ~100 bits (as of 2024)

**Why not 128 bits?**
- BN curves have special structure
- Kim-Barbulescu attack (2016) reduced security
- Still secure for most blockchain applications

**Recommendations**:
- ✅ Good for: Smart contracts, zkSNARKs, short-term secrets
- ⚠️ Consider alternatives for: Long-term encryption (30+ years)
- 🔄 Future-proof: Consider migrating to BLS12-381 (128-bit security)

**Known attacks**:
- Pollard's rho: $O(\sqrt{r})$ for discrete log
- Kim-Barbulescu: Faster attacks on pairing-friendly curves
- No known quantum-safe

---

## Practical Implementation Notes

### Memory Layout

```go
// Fp: 32 bytes (256 bits)
type Fp struct { n *big.Int }

// Fp2: 64 bytes (2 × Fp)
type Fp2 struct { a, b *big.Int }

// G1: 64 bytes (X, Y coordinates)
type G1 struct { X, Y *big.Int }

// G2: 128 bytes (Fp2 coordinates)
type G2 struct { X, Y *Fp2 }

// GT: 384 bytes (Fp12 element)
type GT struct { value *Fp12 }
```

### Serialization

**G1 point**: 64 bytes (32 for X, 32 for Y)
**G2 point**: 128 bytes (64 for X, 64 for Y)
**GT element**: 384 bytes (32 × 12 coefficients)

**Compression**: Not implemented (saves ~50% but slower)

### Performance Tips

1. **Batch verifications**: Verify multiple proofs together
2. **Precomputation**: Cache generator multiples
3. **Lazy reduction**: Delay modular reductions when possible
4. **Assembly**: Critical paths can be 2-3× faster with assembly

---

## References

**Papers**:
1. Barreto-Naehrig (2005): "Pairing-Friendly Elliptic Curves of Prime Order"
2. Vercauteren (2010): "Optimal Pairings"
3. EIP-196/197: Ethereum's BN128 precompiles

**Implementations**:
- Ethereum's go-ethereum: Reference implementation
- gnark-crypto: High-performance Go library
- libsnark: C++ zkSNARK library

**Learn more**:
- "Pairings for Beginners" by Craig Costello
- "Guide to Pairing-Based Cryptography" (book)
- Ben Lynn's PBC library documentation

---

## Quick Mental Model

Think of the math stack as building blocks:

```
🧱 Fp         → Basic numbers (like integers mod p)
🧱🧱 Fp2      → Complex numbers (a + bi)
🧱🧱🧱 Fp6    → Higher-dimensional (like 3D vectors)
🧱🧱🧱🧱 Fp12 → The final structure (4D space)

📍 G1         → Points on a curve (like GPS coordinates)
📍📍 G2       → Points on a fancier curve
🎯 GT         → Magic numbers from pairings

🔗 Pairing    → The bridge connecting G1, G2, GT
```

**The pairing property** (bilinearity):
```
Pairing(k·P, Q) = Pairing(P, k·Q) = Pairing(P, Q)^k
```

This is what makes zkSNARKs work - it lets us verify computations without redoing them!
