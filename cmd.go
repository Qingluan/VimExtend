package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	// var output []rune
	buffer := bytes.NewBuffer([]byte{})
	buf := make([]byte, 4096)
	io.CopyBuffer(buffer, reader, buf)

	// if out, err := utils.SendHTTPFromReader(buffer); err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Println(string(out))
	// }
	for _, v := range strings.Split(buffer.String(), "\n") {
		fmt.Println("-> ", v)
	}
}
