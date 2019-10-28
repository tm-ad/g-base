package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func GetPublicKey(key []byte) (*rsa.PublicKey, error) {
	// decode public key
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.New(`decode public key failed`)
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pub.(*rsa.PublicKey), nil
}

func GetPrivateKey(key []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.New("decode private key failed")
	}
	pri, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err == nil {
		return pri, nil
	}
	pri2, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pri2.(*rsa.PrivateKey), nil
}

func GenRsaKey(bits int) (publicPem, privateKey []byte, err error) {
	// 生成私钥文件
	key, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	derStream := x509.MarshalPKCS1PrivateKey(key)
	privateBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	// 获取私钥字符串
	privateKey = pem.EncodeToMemory(privateBlock)

	// 生成公钥
	publicKey := &key.PublicKey
	derPikx, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, nil, err
	}
	publicBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPikx,
	}

	publicPem = pem.EncodeToMemory(publicBlock)

	return publicPem, privateKey, nil
}
