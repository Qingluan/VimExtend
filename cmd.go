package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/Qingluan/VimExtend/utils"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var output []rune
	r, w := io.Pipe()
	for {
		input, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		output = append(output, input)

	}
	w.Write([]byte(string(output)))
	if out, err := utils.SendHTTPFromReader(r); err != nil {
		fmt.Println(err.errors())
	} else {
		fmt.Println(string(out))
	}
	// for _, v := range strings.Split(string(output), "\n") {
	// 	fmt.Println("-> ", v)
	// }
}
