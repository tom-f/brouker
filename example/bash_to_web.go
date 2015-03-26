package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os/exec"
)

func main() {

	cmd := exec.Command("./test2.sh")

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		panic(err)
	}

	cmd.Start()

	buf := make([]byte, 1024)
	n, err := stdout.Read(buf)

	for err == nil {
		fmt.Print(string(buf[:n]))
		go send(buf)
		// read some more
		_, err = stdout.Read(buf)
	}

	cmd.Wait()

}

func send(buf []byte) {
	body := bytes.NewBuffer(buf)
	r, err := http.Post("http://localhost:8050/msg", "text/plain", body)
	if err != nil {
		panic(err)
	}
	fmt.Println(r.Status)
	r.Body.Close()
}
