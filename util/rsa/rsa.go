package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"go.uber.org/zap"
	"labsystem/configs"
	"labsystem/model/srverr"
	"labsystem/util"
	"labsystem/util/logger"
)

const priviteKey = `
-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQCiSQsbadO2FKbyQf6L7UpJmL86dtdBh+6oW+FEyUjwXStc/9Ij
4DtWr6r9/ckGTKC3YxtUy01uwwrZkU8shi9wcbU1NlkT/rMd9G0azQ7S8ftJLh9p
+vvInDPOoAp+46j+5hNx2zcbPzeVSWQrwgw91dlivCWqohy5pFqKNMC37wIDAQAB
AoGAQy2LbIClQ/VOmnAZJ+FFZCvX1SLGvni/oW/PU4MjL2TSA4Beg71ZsuurE9MP
eQ0VZnSfJr21wwj4H9IX/S1ayz9lSVZvMw2jc/QOzqYMcl0jX/zCsqg+RyZfpnUg
PqrfYA67Gtxq+aWY2MdulFwXTd2zAFPMDxz0Mff6q2F87jECQQDGHoUNqLOSesWi
XqWY5vSefTFAAeZ3/0lE8soTlvOuWg1Kd8o8b/vHf8vLtf8QA8WEbli9BCx0zGKT
Erqa68MJAkEA0bJ7F1481vDjE/LInQfsIEY1v3SNFKqAOMHm03teVknF9qiFOqPs
r1wUo8f5xTOaqp8jG9eWfpYrkGbQx2qJNwJBALH62vbRVjjCOGfMdSGuYZU7zJG+
40drvGvHXQhKNLGtzoDWqoSu6VZolTF1VyUQ7lExPGDZrgJjVi7fxWnFYNkCQHPQ
w7MwVSHRgVkvKiWwjnV00ZKgaJsVYCP40+fRp1S5AiEtCghNKRfhH3UlqYmE6tN5
eQaTl4bqobzrTDwXVxkCQH3Jgwsn2IfoglmP8YT4+PMr9Wc5ICAF9RBFAaYTe1Tk
GTcNTgGS4ULDriZt7I9V+d7MGwF++VLdsSeqpTUlkoM=
-----END RSA PRIVATE KEY-----

`

func GetRSAPrivateKey(privPEM []byte) (priv *rsa.PrivateKey, err error) {
	// decode pem
	privBlocks, _ := pem.Decode(privPEM)
	if privBlocks == nil || privBlocks.Type != "RSA PRIVATE KEY" {
		logger.Log.Warn("pem parse failed", zap.String("pem", string(privPEM)))
		return nil, srverr.ErrInvalidPEM
	}
	// get private key
	if priv, err = x509.ParsePKCS1PrivateKey(privBlocks.Bytes); err != nil {
		// TODO record log
		return
	}
	return
}

func GetRSAPublicKey(pubPEM []byte) (pub *rsa.PublicKey, err error) {
	// decode pem
	pubBlocks, _ := pem.Decode(pubPEM)
	if pubBlocks == nil || pubBlocks.Type != "RSA PUBLIC KEY" {
		logger.Log.Warn("pem parse failed")
		return nil, srverr.ErrInvalidPEM
	}
	// get public key
	var val interface{}
	val, err = x509.ParsePKIXPublicKey(pubBlocks.Bytes)
	if err != nil {
		logger.Log.Warn("parse public key failed", zap.Error(err))
		return
	}
	pub = val.(*rsa.PublicKey)

	return
}

func Encrypt(raw string) (enText string, err error) {
	publicKeyPem, err := util.ReadAll(configs.CurProjectPath() + configs.RSAPublicKeyPEM)
	if err != nil {
		return "", err
	}
	publicKey, err := GetRSAPublicKey(publicKeyPem)
	if err != nil {
		return "", err
	}
	enBytes, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(raw))
	if err != nil {
		return "", err
	}
	enText = base64.StdEncoding.EncodeToString(enBytes)

	return
}

func Decrypt(enText string) (deText string, err error) {
	// base64 decode
	var deBase64Text, deByteText []byte
	if deBase64Text, err = base64.StdEncoding.DecodeString(enText); err != nil {
		return
	}
	var privateKeyPem []byte
	if privateKeyPem, err = util.ReadAll(configs.CurProjectPath() + configs.RSAPrivateKeyPEM); err != nil {
		return
	}
	// get private key
	var privateKey *rsa.PrivateKey
	if privateKey, err = GetRSAPrivateKey(privateKeyPem); err != nil {
		return
	}
	if deByteText, err = rsa.DecryptPKCS1v15(rand.Reader, privateKey, deBase64Text); err != nil {
		return
	}
	deText = string(deByteText)

	return
}

func Compare(enText1, enText2 string) bool {
	txt1, err := Decrypt(enText1)
	if err != nil {
		return false
	}
	txt2, err := Decrypt(enText2)
	if err != nil || txt1 != txt2 {
		return false
	}

	return true
}