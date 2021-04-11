package jwt

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"labsystem/configs"
	"labsystem/model"
	"labsystem/model/srverr"
	"labsystem/util"
	rsa2 "labsystem/util/rsa"
	"time"
)

// note: payload key "exp" is used
func Token(payload map[string]interface{}) (raw string, err error) {
	token := jwt.New(jwt.SigningMethodRS256)
	claims := make(map[string]interface{})
	// pack payload to jwt token
	for k, v := range payload {
		claims[k] = v
	}
	claims["exp"] = time.Now().Add(model.TokenExp)
	token.Claims = jwt.MapClaims(claims)
	// get rsa_private_key pem
	var privateKeyPem []byte
	if privateKeyPem, err = util.ReadAll(configs.CurProjectPath() + configs.RSAPrivateKeyPEM); err != nil {
		return
	}
	// get rsa_keys
	var privateKey *rsa.PrivateKey
	if privateKey, err = rsa2.GetRSAPrivateKey(privateKeyPem); err != nil {
		return
	}
	// get raw token
	raw, err = token.SignedString(privateKey)
	return
}

func ParseToken(rawToken string) (payload map[string]interface{}, err error) {
	// get rsa_public_key pem
	var publicKeyPem []byte
	if publicKeyPem, err = util.ReadAll(configs.CurProjectPath() + configs.RSAPublicKeyPEM); err != nil {
		return
	}
	token, err := jwt.Parse(rawToken, func(token *jwt.Token) (key interface{}, err error) {
		return rsa2.GetRSAPublicKey(publicKeyPem)
	})
	if err != nil {
		return nil, srverr.ErrInvalidToken
	}
	payload = token.Claims.(jwt.MapClaims)
	// check is token outdated
	if exp, ok := payload["exp"]; ok {
		if t, err := time.Parse(time.RFC3339Nano, exp.(string)); err != nil || time.Now().After(t) {
			return nil, srverr.ErrInvalidToken
		}
	} else {
		return nil, srverr.ErrInvalidToken
	}

	delete(payload, "exp")
	return
}
