package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if scanner.Scan() {
			num, err := strconv.Atoi(scanner.Text())
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(num)
			}
		}
	}
}
