package main

import (
	"fmt"
	"strings"
	"text/scanner"
)

func main() {
	var s scanner.Scanner
	input := "5 + 10"
	s.Init(strings.NewReader(input))
	for {
		tok := s.Scan()
		if tok == scanner.EOF {
			break
		}
		fmt.Printf("Token: %s,\n", s.TokenText())
	}
}
