package main

import (
	"fmt"
	"io/ioutil"
	"labsystem/configs"
	"labsystem/util/vcode"
	"os"
	"strconv"
)

const puzzleNo = 1
func main() {
	path := configs.CurProjectPath() + "/static/"
	srcPath := path + "src/"
	for i := 1;i <= 4;i++ {
		x := vcode.PuzzleWidth
		for ;vcode.CheckVCode(x);x += 10 {
			puzBuf, imgBuf, err := vcode.GenerateRandomVImage(puzzleNo, x, srcPath+ strconv.Itoa(i) + ".jpeg")
			if err != nil {
				fmt.Println("generate verify code failed:", err)
				continue
			}
			f1, err := os.OpenFile(path + strconv.Itoa(i) + strconv.Itoa(x) + ".jpeg", os.O_CREATE|os.O_RDWR, 777)
			if err != nil {
				panic("file open error:" + err.Error())
			}
			f2, err := os.OpenFile(path + strconv.Itoa(i) + strconv.Itoa(x) + vcode.PuzzleSuffix + ".jpeg", os.O_CREATE|os.O_RDWR, 777)
			if err != nil {
				panic("file open error:" + err.Error())
			}
			b1, err := ioutil.ReadAll(imgBuf)
			if err != nil {
				panic("buffer read error:" + err.Error())
			}
			if _, err = f1.Write(b1); err != nil {
				panic("file write error:" + err.Error())
			}
			b2, err := ioutil.ReadAll(puzBuf)
			if err != nil {
				panic("buffer read error:" + err.Error())
			}
			if _, err := f2.Write(b2); err != nil {
				panic("file write error:" + err.Error())
			}
			f1.Close()
			f2.Close()
		}
	}
}
