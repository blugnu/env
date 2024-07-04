package env

import (
	"errors"
	"os"
	"testing"

	"github.com/blugnu/test"
)

func TestFileReader(t *testing.T) {
	// ARRANGE
	openerr := errors.New("open error")
	osOpenCalled := false
	defer test.Using(&osOpen, func(string) (*os.File, error) {
		osOpenCalled = true
		return nil, openerr
	})()

	// ACT
	_, err := newFileReader("")

	// ASSERT
	test.That(t, err).Equals(openerr)
	test.IsTrue(t, osOpenCalled)
}
