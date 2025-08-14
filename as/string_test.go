package as_test

import (
	"testing"

	. "github.com/blugnu/test"

	"github.com/blugnu/env/as"
)

func TestString(t *testing.T) {
	With(t)

	// arrange
	var sut = "value"

	// act
	result, err := as.String(sut)

	// assert
	Expect(err).IsNil()
	Expect(result).To(Equal("value"))
}
