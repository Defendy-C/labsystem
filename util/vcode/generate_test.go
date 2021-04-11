package vcode

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"labsystem/configs"
	"os"
	"testing"
)

func TestGenerateRandomVImage(t *testing.T) {
	b, p, e := GenerateRandomVImage(1, 80, configs.CurProjectPath() + "/static/src/2.jpeg")
	require.NoError(t, e)
	f1, e := os.OpenFile("a.png", os.O_CREATE|os.O_RDWR, 0777)
	require.NoError(t, e)
	f2, e := os.OpenFile("b.png", os.O_CREATE|os.O_RDWR, 0777)
	require.NoError(t, e)
	b1, err := ioutil.ReadAll(b)
	require.NoError(t, err)
	b2, err := ioutil.ReadAll(p)
	require.NoError(t, err)
	_, err = f1.Write(b1)
	require.NoError(t, err)
	_, err = f2.Write(b2)
	require.NoError(t, err)
}
