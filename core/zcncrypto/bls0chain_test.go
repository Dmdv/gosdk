package zcncrypto

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/0chain/gosdk/bls"
	"github.com/0chain/gosdk/miracl"
	"github.com/0chain/gosdk/core/encryption"
)

var verifyPublickey = `041eeb1b4eb9b2456799d8e2a566877e83bc5d76ff38b964bd4b7796f6a6ccae6f1966a4d91d362669fafa3d95526b132a6341e3dfff6447e0e76a07b3a7cfa6e8034574266b382b8e5174477ab8a32a49a57eda74895578031cd2d41fd0aef446046d6e633f5eb68a93013dfac1420bf7a1e1bf7a87476024478e97a1cc115de9`
var signPrivatekey = `18c09c2639d7c8b3f26b273cdbfddf330c4f86c2ac3030a6b9a8533dc0c91f5e`

var data = `TEST`
var blsWallet *Wallet

// This is a basic unit test to print out the generator of MIRACL. We used this
// to compare against generator of herumi/bls, to make library compatible.
func TestGenerator(t *testing.T) {
	base := BN254.ECP2_generator()
	fmt.Println("base?", base.ToString())
}

// This is a basic unit test to check that MIRACL generates correct public key.
func TestHerumiPKcompatibility(t *testing.T) {
	var skStr = signPrivatekey
	var sk bls.SecretKey
	sk.DeserializeHexStr(skStr)
	pk := sk.GetPublicKey()

	skStr2 := sk.SerializeToHexStr()
	if skStr2 != skStr {
		panic("Secret Key deserialize failed: [skStr, skStr2]: " + skStr + " " + skStr2)
	}

	// Expect 'pk' to be: ([1bdfed3a85690775ee35c61678957aaba7b1a1899438829f1dc94248d87ed368,18a02c6bd223ae0dfda1d2f9a3c81726ab436ce5e9d17c531ff0a385a13a0b49],[039ac7dfc3364e851ebd2631ea6f1685609fc66d50223cc696cb59ff2fee47ac,17f6dfafec19bfa87bf791a4d694f43fec227ae6f5a867490e30328cac05eaff])
	fmt.Println("pk", pk.ToString())
	fmt.Println("pk", pk.SerializeToHexStr())
}

func TestSetHexString(t *testing.T) {
	testSetHexStringCase("11")
	testSetHexStringCase("1e5f0362da9a74615cb5e3013bab322f")
}

func testSetHexStringCase(n string) {
	var id bls.ID
	err := id.SetHexString(n)
	if err != nil {
		panic(err)
	}
	a := id.GetHexString()
	fmt.Println("Got", a)
}

// Basic unit test to test that secret key Set works.
func TestSecretKeySet(t *testing.T) {
	var sk bls.SecretKey
	var id bls.ID

	sk.SetByCSPRNG()
	Msk := sk.GetMasterSecretKey(2)
	mpk := bls.GetMasterPublicKey(Msk)

	// These other SetHexString calls are additional cases that can be used to test.
	// err := id.SetHexString("000000000" + "000000000" + "00000000" + "00000000" + "00000000000" + "0001e5f0362da" + "9a74615cb5e3013bab322f")
	// err := id.SetHexString("00")
	// err := id.SetHexString("01")
	//
	// The ideal unit test should be testing way more 'id' than just these few.
	// Perhaps a for-loop or random sampling would be a good approach.
	err := id.SetHexString("0000000000000000000000000000000000000000000000000000000000000005")
	if err != nil {
		panic(err)
	}

	var sij bls.SecretKey
	sij.Set(Msk, &id)

	var expectedSijPK bls.PublicKey
	expectedSijPK.Set(mpk, &id)

	if !expectedSijPK.IsEqual( sij.GetPublicKey() ) {
		t.Fatalf("Should've been a valid share.")
	}
}

// A simple unit test to test serialization and deserialization of a private key.
// It's simple, but necessary because did a big port replacing herumi/bls with
// miracl/core, and it's easy to make simple mistakes like this (we did).
func TestSerialization(t *testing.T) {
	privateKey := `212ba4f27ffffff5a2c62effffffffcdb939ffffffffff8a15ffffffffffff8d`
	var primarySk bls.SecretKey
	primarySk.DeserializeHexStr(privateKey)
	d := primarySk.SerializeToHexStr()
	if privateKey != d {
		fmt.Println("before:", privateKey)
		fmt.Println("after:", d)
		t.Fatalf("Basic de/serialization test failed.")
	}

	var pk bls.PublicKey
	err := pk.DeserializeHexStr(`04106806dfd2410c9072daed4892280a944dce4c81da48f854c59a6c1e4d4e2725206048b53a71242dcf370baf15cce63532dbb50e6646c803fb6609063140e134097635737e1c9dd8c6caaa7f375a72dddbfd6c2a21557f37d73938aed76cbb2416082a343a30f0621b308b01cd019bcb8795652e018d61d4afa1159b76df0aac`)
	if err != nil {
		fmt.Println("Got err:", err)
		t.Fatalf("Couldn't deserialize public key.")
	}
}

// Test the following code we ported from herumi.
// ```
// var sk bls.SecretKey
// sk.SetByCSPRNG()
// pk := sk.SerializeToHexStr()
// ```
func TestSetByCSPRNG(t *testing.T) {
	// Not interesting test case anymore.
	return
	testSetByCSPRNGCase(t, []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, "212ba4f27ffffff5a2c62effffffffcdb939ffffffffff8a15ffffffffffff8d")
	testSetByCSPRNGCase(t, []byte{178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178, 178}, "1e2520a9b2b2b2abc9e17cb2b2b2b2912e2eb2b2b2b2b26416b2b2b2b2b2b266")
	testSetByCSPRNGCase(t, []byte{5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5}, "0505050505050505050505050505050505050505050505050505050505050505")
}

// Test the following code we ported from herumi.
// ```
// var sk bls.SecretKey
// sk.SetByCSPRNG()
// pk := sk.SerializeToHexStr()
// ```
func testSetByCSPRNGCase(t *testing.T, seed []byte, expected_pk string) {
	var sk bls.SecretKey
	r := bytes.NewReader(seed)
	bls.SetRandFunc(r)
	sk.SetByCSPRNG()
	pk := sk.SerializeToHexStr()
	if pk != expected_pk {
		fmt.Println("pk:", pk)
		fmt.Println("expected_pk:", expected_pk)
		t.Fatalf("Did not get right secret key.")
	}

	// Do a basic sanity test that Serialize/Deserialize is working.
	sk.DeserializeHexStr(pk)
	pk2 := sk.SerializeToHexStr()
	if pk != pk2 {
		fmt.Println("before ser :", pk)
		fmt.Println("after deser:", pk2)
		t.Fatalf("Basic de/serialization test failed.")
	}
}

func TestSignatureScheme(t *testing.T) {
	sigScheme := NewSignatureScheme("bls0chain")
	switch sigScheme.(type) {
	case SignatureScheme:
		// pass
	default:
		t.Fatalf("Signature scheme invalid")
	}
	w, err := sigScheme.GenerateKeys()
	if err != nil {
		t.Fatalf("Generate Key failed %s", err.Error())
	}
	if w.ClientID == "" || w.ClientKey == "" || len(w.Keys) != 1 || w.Mnemonic == "" {
		t.Fatalf("Invalid keys generated")
	}
	blsWallet = w
}

func TestSSSignAndVerify(t *testing.T) {
	signScheme := NewSignatureScheme("bls0chain")
	signScheme.SetPrivateKey(signPrivatekey)
	hash := Sha3Sum256(data)

	fmt.Println("hash", hash)
	fmt.Println("privkey", signScheme.GetPrivateKey())

	var sk bls.SecretKey
	sk.DeserializeHexStr(signScheme.GetPrivateKey())
	pk := sk.GetPublicKey()
	fmt.Println("pubkey", pk.ToString())

	signature, err := signScheme.Sign(hash)
	if err != nil {
		t.Fatalf("BLS signing failed")
	}
	verifyScheme := NewSignatureScheme("bls0chain")
	verifyScheme.SetPublicKey(verifyPublickey)
	if ok, err := verifyScheme.Verify(signature, hash); err != nil || !ok {
		t.Fatalf("Verification failed\n")
	}
}

func BenchmarkBLSSign(b *testing.B) {
	sigScheme := NewSignatureScheme("bls0chain")
	sigScheme.SetPrivateKey(signPrivatekey)
	for i := 0; i < b.N; i++ {
		_, err := sigScheme.Sign(encryption.Hash(data))
		if err != nil {
			b.Fatalf("BLS signing failed")
		}
	}
}

func TestRecoveryKeys(t *testing.T) {
	sigScheme := NewSignatureScheme("bls0chain")
	TestSignatureScheme(t)
	w, err := sigScheme.RecoverKeys(blsWallet.Mnemonic)
	if err != nil {
		t.Fatalf("set Recover Keys failed")
	}
	if w.ClientID != blsWallet.ClientID || w.ClientKey != blsWallet.ClientKey {
		t.Fatalf("Recover key didn't match with generated keys")
	}
}

func TestCombinedSignAndVerify(t *testing.T) {
	sk0 := `c36f2f92b673cf057a32e8bd0ca88888e7ace40337b737e9c7459fdc4c521918`
	sk1 := `704b6f489583bf1118432fcfb38e63fc2d4b61e524fb196cbd95413f8eb91c12`

	// Public key comes from sk0.add(sk1), then sk0.GetPublicKey().SerializeToHexStr()
	pk := `04039b5e64e3a88f2282c404978ef4d5f6a2fc0ee343f4bb673587b3b77181a96619ad785b8d6a26787116602a916c701b931ec467127bbf44d05d090a163af3e30ff306356487692a2b2652f631b9ee9742df261d4a25525fe9b6da9e60d25cf21ffc3ad8af10210610cef274d8f0a1af80664c39c43bdf30a1819a6258c878fe`

	hash := Sha3Sum256(data)

	// Create signature for 1st.
	scheme0 := NewSignatureScheme("bls0chain")
	err := scheme0.SetPrivateKey(sk0)
	if err != nil {
		t.Fatalf("Set private key failed - %s", err.Error())
	}
	sig0, err := scheme0.Sign(hash)
	if err != nil {
		t.Fatalf("BLS signing failed")
	}

	// Create signature for 2nd.
	scheme1 := NewSignatureScheme("bls0chain")
	err = scheme1.SetPrivateKey(sk1)
	if err != nil {
		t.Fatalf("Set private key failed - %s", err.Error())
	}
	sig1, err := scheme1.Add(sig0, hash)

	verifyScheme := NewSignatureScheme("bls0chain")
	err = verifyScheme.SetPublicKey(pk)
	if err != nil {
		t.Fatalf("Set public key failed")
	}
	if ok, err := verifyScheme.Verify(sig1, hash); err != nil || !ok {
		fmt.Println("err", err)
		t.Fatalf("Verification failed\n")
	}
}

func TestSplitKey(t *testing.T) {
	// Generate 0th signature based on this primaryKey.

	// TODO: we need to investigate this. Running the unit test with these primary
	// keys fail.
	//
	// My current hypothesis is that this primaryKey is "too small",
	// and could lead to a negative Fr/Fp value (which should NEVER happen).
	// How can a negative value happen? What the SplitKeys operation does, for
	// "splitKeys=2", is essentially the following in pseudocode:
	// 1) generate a random key for i=0, call this k_0
	// 2) for i=1, get the key by doing: primaryKey - k_0
	// So you can see that if primaryKey is "too small", and the random key is
	// likely larger than the primaryKey, then k_1 key is going to be negative,
	// and will break the library.
	//
	// In fact this is directly from MIRACL_Core.pdf: `MIRACL Core has no support
	// for negative numbers. It is the programmers responsibility to make sure
	// that a negative result can never happen. MIRACL Core does not support a
	// general purpose big number library.` (page 5, bottom paragraph)

	// `704...` reduces to `00e...`. These primary keys fail.
	// primaryKeyStr := `704b6f489583bf1118432fcfb38e63fc2d4b61e524fb196cbd95413f8eb91c12`
	// primaryKeyStr := `00e141c1d583bf0be9a6474fb38e63e309e861e524fb1931c895413f8eb91bd9`

	// `c36...` reduces to `09b...`. These primary keys fail.
	// primaryKeyStr := `c36f2f92b673cf057a32e8bd0ca88888e7ace40337b737e9c7459fdc4c521918`
	// primaryKeyStr := `09be39077673cefcd72d653d0ca8885f0207e40337b7378784459fdc4c5218b9`

	// Both of these primary keys pass.
	primaryKeyStr := `212ba4f27ffffff5a2c62effffffffcdb939ffffffffff8a15ffffffffffff8d`
	// primaryKeyStr := `5e1fc9c03d53a8b9a63030acc2864f0c33dffddb3c276bf2b3c8d739269cc018`

	scheme0 := NewBLS0ChainScheme()
	err := scheme0.SetPrivateKey(primaryKeyStr)
	if err != nil {
		t.Fatalf("Set private key failed - %s", err.Error())
	}
	hash := Sha3Sum256(data)
	sig0, err := scheme0.Sign(hash)
	if err != nil {
		t.Fatalf("BLS signing failed")
	}

	// Split keys.
	numSplitKeys := int(2)
	w, err := scheme0.SplitKeys(numSplitKeys)
	if err != nil {
		t.Fatalf("Splitkeys key failed - %s", err.Error())
	}

	// Generate schemes from the split keys.
	sigAggScheme := make([]BLS0ChainScheme, numSplitKeys)
	for i := 0; i < numSplitKeys; i++ {
		sigAggScheme[i].SetPrivateKey(w.Keys[i].PrivateKey)
	}

	// Aggregate the signatures generated by each split key.
	var aggrSig string
	for i := 1; i < numSplitKeys; i++ {
		tmpSig, err := sigAggScheme[i].Sign(hash)
		if err != nil {
			fmt.Println("err", err)
			t.Fatalf("Shouldn't have gotten error with Sign()")
		}

		aggrSig, err = sigAggScheme[0].Add(tmpSig, hash)
		if err != nil {
			fmt.Println("err", err)
			t.Fatalf("Shouldn't have gotten error with Add()")
		}
	}

	if aggrSig != sig0 {
		t.Fatalf("split key signature failed")
	}
}
