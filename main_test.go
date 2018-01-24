package main

import (
	"testing"
	"fmt"
	"time"
)

func Test(t *testing.T) {
	var (
		numeroSecreto = 13509
		p             = 1
		q             = 2
		m             = 0
	)

	for q < numeroSecreto {
		q *= 2
	}

	for m != numeroSecreto && p != q {
		m = (p + q) / 2
		fmt.Println(p, q, m)
		if m > numeroSecreto {
			q = m
		} else {
			p = m
		}
		time.Sleep(time.Millisecond * 25)
	}
}