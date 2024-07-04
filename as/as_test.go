package as

import (
	"errors"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/blugnu/env"
	"github.com/blugnu/test"
)

func TestDuration(t *testing.T) {
	// ARRANGE
	var sut = "1s"

	// ACT
	result, err := Duration(sut)

	// ASSERT
	test.That(t, err).IsNil()
	test.That(t, result).Equals(time.Second)
}

func TestDuration_WithSpecifiedUnit(t *testing.T) {
	// ARRANGE
	var sut = "10"

	// ACT
	result, err := Duration(sut, time.Second)

	// ASSERT
	test.That(t, err).IsNil()
	test.That(t, result).Equals(10 * time.Second)
}

func TestDuration_DurationWithSpecifiedUnit(t *testing.T) {
	// ARRANGE
	var sut = "1h"

	// ACT
	result, err := Duration(sut, time.Second)

	// ASSERT
	test.Error(t, err).Is(strconv.ErrSyntax)
	test.That(t, result).Equals(0)
}

func TestInt(t *testing.T) {
	// ARRANGE
	var sut = "123"

	// ACT
	result, err := Int(sut)

	// ASSERT
	test.That(t, err).IsNil()
	test.That(t, result).Equals(123)
}

func TestInt_WhenConversionFails(t *testing.T) {
	// ARRANGE
	var sut = "not-a-number"

	// ACT
	result, err := Int(sut)

	// ASSERT
	test.Error(t, err).Is(strconv.ErrSyntax)
	test.That(t, result).Equals(0)
}

func TestPortNo(t *testing.T) {
	// ARRANGE
	var sut = "123"

	// ACT
	result, err := PortNo(sut)

	// ASSERT
	test.That(t, err).IsNil()
	test.That(t, result).Equals(123)
}

func TestPortNo_WhenConversionFails(t *testing.T) {
	// ARRANGE
	var sut = "not-a-number"

	// ACT
	result, err := PortNo(sut)

	// ASSERT
	test.Error(t, err).Is(strconv.ErrSyntax)
	test.That(t, result).Equals(0)
}

func TestPortNo_WhenOutOfRange(t *testing.T) {
	// ARRANGE
	var sut = "65536"

	// ACT
	result, err := PortNo(sut)

	// ASSERT
	test.Error(t, err).Is(env.RangeError[int]{Min: 0, Max: 65535})
	test.That(t, result).Equals(0)
}

func TestString(t *testing.T) {
	// ARRANGE
	var sut = "value"

	// ACT
	result, err := String(sut)

	// ASSERT
	test.That(t, err).IsNil()
	test.That(t, result).Equals("value")
}

func TestAbsoluteURL(t *testing.T) {
	// ARRANGE
	var sut = "http://example.com"

	// ACT
	result, err := AbsoluteURL(sut)

	// ASSERT
	test.That(t, err).IsNil()
	test.That(t, result).Equals(&url.URL{Scheme: "http", Host: "example.com"})
}

func TestAbsoluteURL_WhenNotAValidURI(t *testing.T) {
	// ARRANGE
	converr := errors.New("url error")
	defer test.Using(&urlParse, func(string) (*url.URL, error) { return nil, converr })()

	// ACT
	result, err := AbsoluteURL("any")

	// ASSERT
	test.Error(t, err).Is(converr)
	test.That(t, result).IsNil()
}

func TestAbsoluteURL_WhenNotAnAbsoluteURL(t *testing.T) {
	// ARRANGE
	var sut = "/path"

	// ACT
	result, err := AbsoluteURL(sut)

	// ASSERT
	test.Error(t, err).Is(ErrNotAnAbsoluteURL)
	test.That(t, result).IsNil()
}
