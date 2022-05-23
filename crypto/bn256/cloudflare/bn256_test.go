package bn256

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/stretchr/testify/assert"
	"math/big"
	"strconv"
	"testing"
)

func TestG1Marshal(t *testing.T) {
	gorigin, Ga, err := RandomG1(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	ma := Ga.Marshal()

	Gb := new(G1)
	_, err = Gb.Unmarshal(ma)
	if err != nil {
		t.Fatal(err)
	}
	mb := Gb.Marshal()

	if !bytes.Equal(ma, mb) {
		t.Fatal("bytes are different")
	}
	hash1 := crypto.Keccak256Hash([]byte("1"))
	selfG1 := new(G1).ScalarBaseMult(new(big.Int).SetBytes(hash1.Bytes()))
	fmt.Println(hex.EncodeToString(selfG1.Marshal()))
	fmt.Println(selfG1.String())
	dfG2 := &G2{twistGen}
	fmt.Println("generatorG2", dfG2.String(), gorigin)

	//W := Ga
	PK := new(G2).ScalarBaseMult(gorigin)

	gorigin2, _, err := RandomG1(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	PK2 := new(G2).ScalarBaseMult(gorigin2)

	msg := "1"
	signature := sign(gorigin, []byte(msg))

	signature2 := sign(gorigin2, []byte(msg))

	mergeSign := new(G1).Add(signature, signature2)
	mergePK := new(G2).Add(PK, PK2)

	G1List := make([]*G1, 0)
	G2List := make([]*G2, 0)
	G1List = append(G1List, new(G1).Neg(mergeSign))
	G2List = append(G2List, dfG2)

	G1List = append(G1List, hashToG1([]byte(msg)))
	G2List = append(G2List, mergePK)
	pr1 := Pair(mergeSign, dfG2)
	pr2 := Pair(hashToG1([]byte(msg)), mergePK)
	fmt.Println("pr1==pr2", pr1.String() == pr2.String())
	fmt.Println("pr1:", pr1.String())
	fmt.Println("pr2:", pr2.String())
	assert.True(t, PairingCheck(G1List, G2List))
	fmt.Println("----------")
	fmt.Println("signature", signature.String())
	fmt.Println("pk", PK.String())
	fmt.Println("msgG1", hashToG1([]byte(msg)).String(), hashToG1([]byte(msg)).p.y.String())

	G21, _ := new(big.Int).SetString("1800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed", 16)
	fmt.Println("g2dec:", G21)

	// The prime q in the base field F_q for G1
	q, _ := new(big.Int).SetString("21888242871839275222246405745257275088696311157297823662689037894645226208583", 10)
	signBytes, _ := hexutil.Decode("23f5d1a185178b2e02f5f20c032d5ddda8ff05382350c05698aec4c1491d6553")
	signatureBigInt := new(big.Int).SetBytes(signBytes)
	fmt.Println("negate sign:", hex.EncodeToString(new(big.Int).Sub(q, new(big.Int).Mod(signatureBigInt, q)).Bytes()))

	fmt.Println(len(P.Bytes()))

	fmt.Println("-------------------C++")
	sh256Out := sha256.Sum256([]byte(strconv.Itoa(1)))
	x := sh256Out[:]
	L := len(P.Bytes())
	maxL := new(big.Int).SetUint64((1 << L) - 1)
	fmt.Println(x, (1<<L)-1, len(maxL.Bytes()), maxL)

	fmt.Println(sha256.Sum256([]byte("h2c")))

}

func hashToG1(msg []byte) *G1 {
	hash1 := crypto.Keccak256Hash(msg)
	selfG1 := new(G1).ScalarBaseMult(new(big.Int).SetBytes(hash1.Bytes()))
	return selfG1
}

func sign(prik *big.Int, msg []byte) *G1 {
	msgG1 := hashToG1(msg)
	return new(G1).ScalarMult(msgG1, prik)
}

func TestG2Marshal(t *testing.T) {
	_, Ga, err := RandomG2(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	ma := Ga.Marshal()

	Gb := new(G2)
	_, err = Gb.Unmarshal(ma)
	if err != nil {
		t.Fatal(err)
	}
	mb := Gb.Marshal()

	if !bytes.Equal(ma, mb) {
		t.Fatal("bytes are different")
	}
}

func TestBilinearity(t *testing.T) {
	for i := 0; i < 2; i++ {
		a, p1, _ := RandomG1(rand.Reader)
		b, p2, _ := RandomG2(rand.Reader)
		e1 := Pair(p1, p2)

		e2 := Pair(&G1{curveGen}, &G2{twistGen})
		e2.ScalarMult(e2, a)
		e2.ScalarMult(e2, b)

		if *e1.p != *e2.p {
			t.Fatalf("bad pairing result: %s", e1)
		}
	}
}

func TestTripartiteDiffieHellman(t *testing.T) {
	a, _ := rand.Int(rand.Reader, Order)
	b, _ := rand.Int(rand.Reader, Order)
	c, _ := rand.Int(rand.Reader, Order)

	pa, pb, pc := new(G1), new(G1), new(G1)
	qa, qb, qc := new(G2), new(G2), new(G2)

	pa.Unmarshal(new(G1).ScalarBaseMult(a).Marshal())
	qa.Unmarshal(new(G2).ScalarBaseMult(a).Marshal())
	pb.Unmarshal(new(G1).ScalarBaseMult(b).Marshal())
	qb.Unmarshal(new(G2).ScalarBaseMult(b).Marshal())
	pc.Unmarshal(new(G1).ScalarBaseMult(c).Marshal())
	qc.Unmarshal(new(G2).ScalarBaseMult(c).Marshal())

	k1 := Pair(pb, qc)
	k1.ScalarMult(k1, a)
	k1Bytes := k1.Marshal()

	k2 := Pair(pc, qa)
	k2.ScalarMult(k2, b)
	k2Bytes := k2.Marshal()

	k3 := Pair(pa, qb)
	k3.ScalarMult(k3, c)
	k3Bytes := k3.Marshal()

	if !bytes.Equal(k1Bytes, k2Bytes) || !bytes.Equal(k2Bytes, k3Bytes) {
		t.Errorf("keys didn't agree")
	}
}

func BenchmarkG1(b *testing.B) {
	x, _ := rand.Int(rand.Reader, Order)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		new(G1).ScalarBaseMult(x)
	}
}

func BenchmarkG2(b *testing.B) {
	x, _ := rand.Int(rand.Reader, Order)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		new(G2).ScalarBaseMult(x)
	}
}
func BenchmarkPairing(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Pair(&G1{curveGen}, &G2{twistGen})
	}
}
