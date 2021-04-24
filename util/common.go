package util

import (
	cRand "crypto/rand"
	"errors"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"time"
)

var (
	ErrInvalidRegexp = errors.New("invalid regexp")
	ErrVerifyFailed = errors.New("this string isn't accord with regexp")
)

func ReadAll(filename string) (data []byte, err error) {
	var f *os.File
	if f, err = os.Open(filename); err != nil {
		return
	}
	data, err = ioutil.ReadAll(f)
	return
}

func RandIntN(n int) int {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	switch n {
	case 0:
		return r.Int()
	default:
		return r.Intn(n)
	}
}

func StringFormatVerify(adminName string, exp string) error {
	//limit format
	re, err := regexp.Compile(exp)
	if err != nil {
		return ErrInvalidRegexp
	}
	if !re.MatchString(adminName) {
		return ErrVerifyFailed
	}

	return nil
}

func Uuid() string {
	now := time.Now().UnixNano()
	randBytes := make([]byte, 16)
	cRand.Read(randBytes)

	return strconv.Itoa(int(now)) + string(randBytes)
}