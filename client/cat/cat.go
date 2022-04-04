package cat

import (
	"bufio"
	"os/exec"
)

// based on https://medium.com/@arunprabhu.1/tailing-a-file-in-golang-72944204f22b
// follows a given file and passes lines from it to the given func
func Cat(file string, token string, callback func(s string, token string)) {
	cmd := exec.Command("cat", file)

	reader, err := cmd.StdoutPipe()
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(reader)
	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			// passing the token here is obviously bad, but I got lazy
			callback(line, token)
		}
	}()

	err = cmd.Start()
	if err != nil {
		return
	}

	err = cmd.Wait()
	if err != nil {
		return
	}
}
