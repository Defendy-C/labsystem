package rsa

import (
	"github.com/stretchr/testify/require"
	"labsystem/configs"
	"labsystem/util"
	"testing"
)

func TestGetRSAKey(t *testing.T) {
	var privPEM []byte
	var err error
	// get private key pem
	if privPEM, err = util.ReadAll(configs.CurProjectPath() + configs.RSAPrivateKeyPEM); err != nil {
		panic("util.ReadAll is failed, err: " + err.Error())
	}
	priv, err := GetRSAPrivateKey(privPEM)
	require.NoError(t, err)
	var pubPEM []byte
	// get public key pem
	if pubPEM, err = util.ReadAll(configs.CurProjectPath() + configs.RSAPublicKeyPEM); err != nil {
		panic("util.ReadAll is failed, err: " + err.Error())
	}
	pub, err := GetRSAPublicKey(pubPEM)
	require.NoError(t, err)
	require.Equal(t, &priv.PublicKey, pub)
}

func TestRSADecrypt(t *testing.T) {
	pwd, err := Decrypt("HpM7mF4U3e2mqMfb2ApX7kHVru6Yz2VRJ6o5RhLDtHjt8ZLTn7qxtKFXGuqBFFyMtHD1nmI7o5hqozEd4XeFYaXC67dY04brl5N/I3vv4zAu87Ega7g8yzjW4i1LReVEWaPIBDCFZTLwDmLHQDuX7BczY/R4TSrJ/cIOcNFfeV8=")
	require.NoError(t, err)
	require.Equal(t, "123", pwd)
}

func TestCompare(t *testing.T) {
	pwd1 := `MpJucmZlKGm5CCqKKxrQNUhr8zT+lMvNKowpo7LKfKJVcPcwpVu52DhAxdFkxcqrlkcH3B5fwOqfPLrYwFl0W60JIYYJ/kHw6n7wKEby9Bw0nXuBhEbtHb3o3eQMLBS1LEJ4HVyzZrB9byno+2DE3NTf+HcN7DqfU8ev3GyZpVA=`
	pwd2 := `aj6kNOi+Sja6N2UHYRqtHiK21WVqCid3VcJujNpOLaTfr+L4tP9MVE+QR8mvjRsnpXa5kY57LpN7IwMs/yyI8CA4gvKa+/L5f/N5os+ZHdZCXj2j5Ku8kAh93k3BhYESg/tiQ6++YtiyDVIe5802Hs7KKjummwXtStgm5DB4rvA=`
	require.True(t, Compare(pwd1, pwd2))
}

// 123456: LYTQwRSUb6o+dMpa5uZJixl6tXirmTXSdWl690r051CCJYFnLvhDx5YfmhFq5d+zMOai8mDF/vGTfqW0VjFi73/hvjUmz3DZDQ9U0Q1dKkcT4ucFNkADG95qEnOGl1OxW7DJymiASMdU76fctKFV+g5MWJSsh+Q0Cel4KYjaFGU=