package utils

import (
	"fmt"
	"math/rand"
)

func Gen6DigCod() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}