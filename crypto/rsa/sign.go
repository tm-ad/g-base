package rsa

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
)

func SignSha1WithRsa(data []byte, privateKey *rsa.PrivateKey) (sign []byte, err error) {
	hashed := sha1.Sum(data)
	return rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA1, hashed[:])
}

func SignSha256WithRsa(data []byte, privateKey *rsa.PrivateKey) (sign []byte, err error) {
	hashed := sha256.Sum256(data)
	return rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
}

func VerifySignSha1WithRsa(data, signData []byte, publicKey *rsa.PublicKey) error {
	hashed := sha1.Sum(data)
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA1, hashed[:], signData)
}

func VerifySignSha256WithRsa(data, signData []byte, publicKey *rsa.PublicKey) error {
	hashed := sha256.Sum256(data)
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signData)
}
