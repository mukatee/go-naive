package cryptoff

//https://github.com/akamensky/base58

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

/**

http://codrspace.com/supcik/golang-jwt-ecdsa/

1. luodaan avain javassa, otetaan sieltä se numero
2. ladataan numero tänne, luodaan private key siitä, public key private keystä
3. signature tehdään javassa, tuodaan tänne
4. verrataan signaturea public keyhin täällä
*/

//test to see a signature created in Java can be parsed and verified in Go
func TestJavaVerify(t *testing.T) {
	//this is a hex-string serialized in Java for the private key.
	//so re-build the private key here to get public key for verification.
	//should only pass public key but works for this test
	privKey := HexToPrivateKey("10EDB31521C5ABF4DA520F784F927B390B4A844FCED4BF2639588E9430BDA9D1")
	ePubKey := &privKey.PublicKey

	xHexStr := "4bc55d002653ffdbb53666a2424d0a223117c626b19acef89eefe9b3a6cfd0eb"
	yHexStr := "d8308953748596536b37e4b10ab0d247f6ee50336a1c5f9dc13e3c1bb0435727"
	ePubKey = HexToPublicKey(xHexStr, yHexStr)

	//the signature in ASN.1 format from Java
	sig := "3045022071f06054f450f808aa53294d34f76afd288a23749628cc58add828e8b8f2b742022100f82dcb51cc63b29f4f8b0b838c6546be228ba11a7c23dc102c6d9dcba11a8ff2"
	sigHex, _ := hex.DecodeString(sig)
	ok := verifyMySig(ePubKey, "This is string to sign", sigHex)
	assert.True(t, ok, "Golang should verify Java signature ok")
}

//sign a piece of data in Go and verify it in Go as well. To verify we got the basic signature functionality ok
func TestSelfSignAndVerify(t *testing.T) {
	privKey, _ := ecdsa.GenerateKey(Curve, rand.Reader)
	esig := CreateSignature([]byte("Hello World"), privKey)
	ok := verifyESig(&privKey.PublicKey, []byte("Hello World"), esig)
	assert.True(t, ok, "Golang should create and verify its own signatures OK")
}

//create public and private keys for ECDSA, sign some data with them, serialize the keys, deserialize them,
//verify the signed data is OK for the desearialized keys
//to provide assurance that key serialization works fine in both ways
func TestEncodeDecodeKeys(t *testing.T) {
	privKey, _ := ecdsa.GenerateKey(Curve, rand.Reader)
	pubKey := privKey.PublicKey

	pub58 := EncodePublicKey(&pubKey)
	priv58 := EncodePrivateKey(privKey)

	msg := []byte("Hello ECDSA")
	esig := CreateSignature(msg, privKey)

	pubKey2 := DecodePublicKey(pub58)
	verifyESig(pubKey2, msg, esig)

	privKey2 := DecodePrivateKey(priv58)
	esig2 := CreateSignature(msg, privKey2)
	ok := verifyESig(pubKey2, msg, esig2)
	assert.True(t, ok, "Golang (de)serialized keys should work to sign OK")

}

func TestErroReporting(t *testing.T) {
	//this reports two errors, so all errors are reported
	t.Errorf("Error 1")
	t.Errorf("Error 2")
}

func TestAssertReporting(t *testing.T) {
	//this reports both asserts and the following require. so all asserts are reported, require stops the test
	assert.Equal(t, "wrong", "very wrong", "Failure 1 %s", "a")
	assert.Equal(t, "wrong again", "very wrong again", "Failure 1 %s", "b")
	require.Equal(t, "wrong wrong", "very wrong wrong", "Stop here %s", "c")
	assert.Equal(t, "wrong again 2", "very wrong again 2", "Failure 1 %s", "d")
}

//https://golang.org/src/crypto/ecdsa/ecdsa.go

func verifyMySig(pub *ecdsa.PublicKey, msg string, sig []byte) bool {
	//https://github.com/gtank/cryptopasta/blob/master/sign.go
	digest := sha256.Sum256([]byte(msg))

	var esig ecdsaSignature
	asn1.Unmarshal(sig, &esig)
	//	esig.R.SetString("89498588918986623250776516710529930937349633484023489594523498325650057801271", 0)
	//	esig.S.SetString("67852785826834317523806560409094108489491289922250506276160316152060290646810", 0)
	fmt.Printf("R: %d , S: %d", esig.R, esig.S)
	println()
	return ecdsa.Verify(pub, digest[:], esig.R, esig.S)
}
