package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
)

var MatsukenSamba = "http://a588.phobos.apple.com/us/r20/Music/v4/82/a0/27/82a02731-c7f8-fec1-c686-dad12eca272f/mzaf_8898382402336846733.m4a"

func main() {

	response, err := http.Get(MatsukenSamba)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	fmt.Println("status:", response.Status)

	file, err := ioutil.TempFile(os.TempDir(), "tmp_bgm_")
	if err != nil {
		panic(err)
	}

	defer file.Close()
	defer os.Remove(file.Name())

	fmt.Print("▶️　Playing Matsuken Samba💃")

	io.Copy(file, response.Body)

	out, _ := exec.Command("afplay", file.Name(), "-r", "1").CombinedOutput()
	fmt.Print(string(out))

}
