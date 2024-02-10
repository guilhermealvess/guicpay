package entity

import "fmt"

type Money int64

const (
	Cent     Money = 1
	Real     Money = 100 * Cent
	MilReais Money = 1000 * Real
)

func (m Money) String() string {
	return fmt.Sprintf("%.2f BRL", float64(m)/100)
}

func (m Money) Absolute() Money {
	if m < 0 {
		return -1 * m
	}

	return m
}
