package chain

import (
	"math/big"
	"crypto/ecdsa"
	"crypto/rand"

	"github.com/akamensky/base58"
	"encoding/hex"
	"fmt"
	"crypto/sha256"
)

//holds the two bigints required to represent (sign/verify) an ECDSA signature
type ecdsaSignature struct {
	R, S *big.Int
}

func createAddress() (*ecdsa.PrivateKey, *ecdsa.PublicKey, string) {
	privKey, _ := ecdsa.GenerateKey(Curve, rand.Reader)
	pubKey := privKey.PublicKey
	address := encodePublicKey(pubKey)
	return privKey, &pubKey, address
}

func createSignature(msg []byte, priv *ecdsa.PrivateKey) (ecdsaSignature) {
	var esig ecdsaSignature
	digest := sha256.Sum256(msg)
	r, s, _ := ecdsa.Sign(rand.Reader, priv, digest[:])
	esig.R = r
	esig.S = s
	return esig
}

func verifyESig(pub *ecdsa.PublicKey, msg []byte, esig ecdsaSignature) bool {
	digest := sha256.Sum256(msg)
	return ecdsa.Verify(pub, digest[:], esig.R, esig.S)
}

func encodePrivateKey(privKey *ecdsa.PrivateKey) string {
	return base58.Encode(privKey.D.Bytes())
}

func decodePrivateKey(str string) *ecdsa.PrivateKey {
	//TODO: curve asetus jostain vakiosta
	privKey := new(ecdsa.PrivateKey)
	privKey.D = b58ToBigInt(str)
	privKey.Curve = Curve
	return privKey
}

func encodePublicKey(pubKey ecdsa.PublicKey) string {
	xBytes := pubKey.X.Bytes()
	yBytes := pubKey.Y.Bytes()
	whole := mergeTwoByteSlices(xBytes, yBytes)
	return base58.Encode(whole)
}

func mergeTwoByteSlices(slice1 []byte, slice2 []byte) []byte {
	s1Len := len(slice1)
	s2Len := len(slice2)
	finalBytes := make([]byte, 1+s1Len+s2Len)
	//https://stackoverflow.com/questions/37210379/convert-int-to-a-single-byte-in-go#37210523
	finalBytes[0] = byte(s1Len)
	copy(finalBytes[1:1+s1Len], slice1[:])
	copy(finalBytes[1+s2Len:], slice2[:])
	return finalBytes
}

func splitTwoByteSlices(whole []byte) ([]byte, []byte) {
	//int(byte) seems to always produce a positive valued integer (-1 = 255)
	size1 := int(whole[0])
	slice1End := 1+size1
	//if want big.int: https://stackoverflow.com/questions/24757814/golang-convert-byte-array-to-big-int
	slice1 := whole[1:slice1End]
	slice2 := whole[slice1End+1:]
	return slice1, slice2
}

func decodePublicKey(b58 string) *ecdsa.PublicKey {
	pubKey := new(ecdsa.PublicKey)
	data, _ := base58.Decode(b58) //TODO: error handling
	xBytes, yBytes := splitTwoByteSlices(data)
	pubKey.X = bytesToBigInt(xBytes)
	pubKey.Y = bytesToBigInt(yBytes)
	//pubKey.X = b58ToBigInt(x58)
	//pubKey.Y = b58ToBigInt(y58)
	pubKey.Curve = Curve
	return pubKey
}

func b58ToBigInt(str string) *big.Int {
	xBytes, _ := base58.Decode(str)
	return bytesToBigInt(xBytes)
}

func bytesToBigInt(data []byte) *big.Int {
	x := new(big.Int)
	x.SetBytes(data)
	return x
}

func HexToPublicKey(xHex string, yHex string) *ecdsa.PublicKey {
	xBytes, _ := hex.DecodeString(xHex)
	x := new(big.Int)
	x.SetBytes(xBytes)

	yBytes, _ := hex.DecodeString(yHex)
	y := new(big.Int)
	y.SetBytes(yBytes)

	pub := new(ecdsa.PublicKey)
	pub.X = x
	pub.Y = y

	pub.Curve = Curve

	return pub
}

func hexToPrivateKey(hexStr string)  *ecdsa.PrivateKey {
	bytes, err := hex.DecodeString(hexStr)
	print(err)

	k := new(big.Int)
	k.SetBytes(bytes)
	println("k:")
	fmt.Println(k.String())

	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = Curve
	priv.D = k
	priv.PublicKey.X, priv.PublicKey.Y = Curve.ScalarBaseMult(k.Bytes())
	fmt.Printf("X: %d, Y: %d", priv.PublicKey.X, priv.PublicKey.Y)
	println()

	return priv
}