package cryptoff

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/akamensky/base58"
	"math/big"
)

var Curve = elliptic.P256()

//ecdsaSignature holds the two bigints required to represent (sign/verify) an ECDSA signature
type ecdsaSignature struct {
	R, S *big.Int
}

//createAddress creates a new private and public key, and encodes the public key as string
//together these form an address for the coin
func CreateAddress() (*ecdsa.PrivateKey, *ecdsa.PublicKey, string) {
	privKey, _ := ecdsa.GenerateKey(Curve, rand.Reader)
	pubKey := &privKey.PublicKey
	address := EncodePublicKey(pubKey)
	return privKey, pubKey, address
}

//createSignature hashes the given bytes and signs the resulting hash with the given private key to produce signature
func CreateSignature(msg []byte, priv *ecdsa.PrivateKey) ecdsaSignature {
	var esig ecdsaSignature
	digest := sha256.Sum256(msg)
	r, s, _ := ecdsa.Sign(rand.Reader, priv, digest[:])
	esig.R = r
	esig.S = s
	return esig
}

//verifyESig verifies a given signature matches the given bytes
func verifyESig(pub *ecdsa.PublicKey, msg []byte, esig ecdsaSignature) bool {
	digest := sha256.Sum256(msg)
	return ecdsa.Verify(pub, digest[:], esig.R, esig.S)
}

//encodePrivateKey base58 encodes the private key for storage/presentation
func EncodePrivateKey(privKey *ecdsa.PrivateKey) string {
	return base58.Encode(privKey.D.Bytes())
}

//decodePrivateKey decodes the private key from given base58 encoded string
func DecodePrivateKey(str string) *ecdsa.PrivateKey {
	privKey := new(ecdsa.PrivateKey)
	privKey.D = b58ToBigInt(str)
	privKey.Curve = Curve
	//PublicKey is anonymous nested struct, so its fields can be directly accessed
	//so no need to separately define access
	//https://golangbot.com/structs/
	//	pubKey := &privKey.PublicKey
	//	pubKey.Curve = Curve
	privKey.X, privKey.Y = Curve.ScalarBaseMult(privKey.D.Bytes())
	return privKey
}

//encodePublicKey takes the ecdsa public key, encodes its two big integer parts,
//and merges them to a single base58 encoded string
func EncodePublicKey(pubKey *ecdsa.PublicKey) string {
	xBytes := pubKey.X.Bytes()
	yBytes := pubKey.Y.Bytes()
	whole := MergeTwoByteSlices(xBytes, yBytes)
	return base58.Encode(whole)
}

//mergeTwoByteSlices merges two byte slices into a single slice.
//the new slice will start with single byte for the length of first slice in new slice.
//this is followed by the first slice, which is followed by second slice
//so newslice = [slice1length][slice1][slice2]
func MergeTwoByteSlices(slice1 []byte, slice2 []byte) []byte {
	s1Len := len(slice1)
	s2Len := len(slice2)
	finalBytes := make([]byte, 1+s1Len+s2Len)
	//https://stackoverflow.com/questions/37210379/convert-int-to-a-single-byte-in-go#37210523
	finalBytes[0] = byte(s1Len)
	copy(finalBytes[1:1+s1Len], slice1[:])
	copy(finalBytes[1+s2Len:], slice2[:])
	return finalBytes
}

//splitTwoByteSlices splits a merged byte slice, produced by mergeTwoByteSlices()
func splitTwoByteSlices(whole []byte) ([]byte, []byte) {
	//int(byte) seems to always produce a positive valued integer (-1 = 255). so I just trust this is ok for 1 byte length
	size1 := int(whole[0])
	slice1End := 1 + size1
	//if want big.int: https://stackoverflow.com/questions/24757814/golang-convert-byte-array-to-big-int
	slice1 := whole[1:slice1End]
	slice2 := whole[slice1End+1:]
	return slice1, slice2
}

//decodePublicKey takes as parameter a base58 encoded public key string, decodes it into golang structure
func DecodePublicKey(b58 string) *ecdsa.PublicKey {
	pubKey := new(ecdsa.PublicKey)
	data, _ := base58.Decode(b58) //TODO: error handling
	//the public key is built from to big integers, so get those first, then construct the key
	xBytes, yBytes := splitTwoByteSlices(data)
	pubKey.X = bytesToBigInt(xBytes)
	pubKey.Y = bytesToBigInt(yBytes)
	//pubKey.X = b58ToBigInt(x58)
	//pubKey.Y = b58ToBigInt(y58)
	pubKey.Curve = Curve
	return pubKey
}

//b58ToBigInt parses the bytes encoded as base58 in given string, and converts into big integer
func b58ToBigInt(str string) *big.Int {
	xBytes, _ := base58.Decode(str)
	return bytesToBigInt(xBytes)
}

//bytesToBigInt takes the given slice of bytes and converts them into a big integer. TODO: endian?
func bytesToBigInt(data []byte) *big.Int {
	x := new(big.Int)
	x.SetBytes(data)
	return x
}

//HexToPublicKey converts a hex-encoded string into a golang public key structure
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

//hexToPrivateKey converts a hex-encoded string into a golang private key structure
func HexToPrivateKey(hexStr string) *ecdsa.PrivateKey {
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
