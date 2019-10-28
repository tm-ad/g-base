package rsa

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"io"
	"io/ioutil"
	"math/big"
)

func publicKeyIO(pub *rsa.PublicKey, in io.Reader, out io.Writer, encrypt bool) (err error) {
	k := (pub.N.BitLen() + 7) / 8
	if encrypt {
		k = k - 11
	}
	buf := make([]byte, k)
	var b []byte
	size := 0
	for {
		size, err = in.Read(buf)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if size < k {
			b = buf[:size]
		} else {
			b = buf
		}
		if encrypt {
			b, err = rsa.EncryptPKCS1v15(rand.Reader, pub, b)
		} else {
			b, err = publicKeyDecrypt(pub, b)
		}
		if err != nil {
			return err
		}
		if _, err = out.Write(b); err != nil {
			return err
		}
	}
}

// privateKeyIO 私钥加密或解密Reader
func privateKeyIO(pri *rsa.PrivateKey, r io.Reader, w io.Writer, encrypt bool) (err error) {
	k := (pri.N.BitLen() + 7) / 8
	if encrypt {
		k = k - 11
	}
	buf := make([]byte, k)
	var b []byte
	size := 0
	for {
		size, err = r.Read(buf)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if size < k {
			b = buf[:size]
		} else {
			b = buf
		}
		if encrypt {
			b, err = privateKeyEncrypt(rand.Reader, pri, b)
		} else {
			b, err = rsa.DecryptPKCS1v15(rand.Reader, pri, b)
		}
		if err != nil {
			return err
		}
		if _, err = w.Write(b); err != nil {
			return err
		}
	}
}

// 公钥解密
func publicKeyDecrypt(pub *rsa.PublicKey, data []byte) ([]byte, error) {
	k := (pub.N.BitLen() + 7) / 8
	if k != len(data) {
		return nil, ErrDataLen
	}
	m := new(big.Int).SetBytes(data)
	if m.Cmp(pub.N) > 0 {
		return nil, ErrDataToLarge
	}
	m.Exp(m, big.NewInt(int64(pub.E)), pub.N)
	d := leftPad(m.Bytes(), k)
	if d[0] != 0 {
		return nil, ErrDataBroken
	}
	if d[1] != 0 && d[1] != 1 {
		return nil, ErrKeyPairDismatch
	}
	var i = 2
	for ; i < len(d); i++ {
		if d[i] == 0 {
			break
		}
	}
	i++
	if i == len(d) {
		return nil, nil
	}
	return d[i:], nil
}

// 私钥加密
func privateKeyEncrypt(rand io.Reader, priv *rsa.PrivateKey, hashed []byte) ([]byte, error) {
	tLen := len(hashed)
	k := (priv.N.BitLen() + 7) / 8
	if k < tLen+11 {
		return nil, ErrDataLen
	}
	em := make([]byte, k)
	em[1] = 1
	for i := 2; i < k-tLen-1; i++ {
		em[i] = 0xff
	}
	copy(em[k-tLen:k], hashed)
	m := new(big.Int).SetBytes(em)
	c, err := decrypt(rand, priv, m)
	if err != nil {
		return nil, err
	}
	copyWithLeftPad(em, c.Bytes())
	return em, nil
}

// 私钥加密，公钥解密
func PrivateEncrypt(input []byte, privateKey *rsa.PrivateKey) (output []byte, err error) {
	if input == nil {
		return []byte(""), errors.New(`empty input, encrypt failed`)
	}
	if privateKey == nil {
		return []byte(""), errors.New(`invalid key, encrypt failed`)
	}
	cipher := bytes.NewBuffer(nil)
	err = privateKeyIO(privateKey, bytes.NewBuffer(input), cipher, true)
	if err != nil {
		return []byte(""), err
	}

	return ioutil.ReadAll(cipher)
}

func PublicDecrypt(input []byte, publicKey *rsa.PublicKey) (output []byte, err error) {
	if input == nil {
		return []byte(""), errors.New(`empty input, decrypt failed`)
	}
	if publicKey == nil {
		return []byte(""), errors.New(`invalid key, decrypt failed`)
	}
	plain := bytes.NewBuffer(nil)
	err = publicKeyIO(publicKey, bytes.NewBuffer(input), plain, false)
	if err != nil {
		return []byte(""), err
	}
	return ioutil.ReadAll(plain)
}

// 公钥加密，私钥解密
func PublicEncrypt(input []byte, publicKey *rsa.PublicKey) (output []byte, err error) {
	if input == nil {
		return []byte(""), errors.New(`empty input, encrypt failed`)
	}
	if publicKey == nil {
		return []byte(""), errors.New(`invalid key, encrypt failed`)
	}
	cipher := bytes.NewBuffer(nil)
	err = publicKeyIO(publicKey, bytes.NewBuffer(input), cipher, true)
	if err != nil {
		return []byte(""), err
	}

	return ioutil.ReadAll(cipher)
}

func PrivateDecrypt(input []byte, privateKey *rsa.PrivateKey) (output []byte, err error) {
	if input == nil {
		return []byte(""), errors.New(`empty input, decrypt failed`)
	}
	if privateKey == nil {
		return []byte(""), errors.New(`invalid key, decrypt failed`)
	}
	plain := bytes.NewBuffer(nil)
	err = privateKeyIO(privateKey, bytes.NewBuffer(input), plain, false)
	if err != nil {
		return []byte(""), err
	}
	return ioutil.ReadAll(plain)
}
