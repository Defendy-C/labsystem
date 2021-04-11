package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"labsystem/configs"
	"os"
)

/*
    -----generate rsa_key file-----
1. generate rsa_key: rsa.GenerateKey
2. create pem file
3. marshal private key and public key: x509.MarshalPKCS1PrivateKey, x509.MarshalPKCS1PublicKey
4. encode private key and public key to pem file: pem.Encode
*/
func main() {
	priv, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic("generating private key failed")
	}
	privPKCS1 := x509.MarshalPKCS1PrivateKey(priv)
	pubPKIX, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		panic("generating public key failed")
	}
	privPEMFile, err := os.OpenFile(configs.CurProjectPath()+"/configs/private_key.pem", os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		panic("creating/opening \"private_key.pem\" failed")
	}
	pubPEMFile, err := os.OpenFile(configs.CurProjectPath()+"/configs/public_key.pem", os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		panic("creating/opening \"public_key.pem\" failed")
	}
	if err = pem.Encode(privPEMFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privPKCS1,
	}); err != nil {
		panic("writing PRIVATE KEY to private_key.pem failed")
	}

	if err = pem.Encode(pubPEMFile, &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubPKIX,
	}); err != nil {
		panic("writing PUBLIC KEY to public_key.pem failed")
	}
}
