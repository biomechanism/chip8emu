package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func mainEntry() {

	file, err := os.Open("/tmp/chip8asm.asm")

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		ops := strings.Split(line, " ")

		fmt.Println(ops[0])
		fmt.Println("------")
		fmt.Println(ops[1])

	}

}

func parseLosdOps(operands string) bool {

	switch string(operands[0]) {
	case "$":
		fmt.Println("Hex number")
	case "%":
		fmt.Println("Bin number")
	default:
		fmt.Println("Decimal number")
	}

	return true
}
