package utils

import (
	"fmt"
	"math/big"
)

type Signature struct {
	R *big.Int // public key x coordinate
	S *big.Int // the signature
}

func (s *Signature) String() string {
	return fmt.Sprintf("%x%x", s.R, s.S)
}
