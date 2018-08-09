package utils

import (
	"testing"
	"fmt"
)

func TestCrossHistory (t *testing.T) {

}

func TestCrossQuery (t *testing.T) {
	re,_ := CrossQuery("51a82a2159d8714f1046c39e1aa73c556529a2d452800044ecf26d930b9e4831")
	fmt.Println(re.Value)
}
