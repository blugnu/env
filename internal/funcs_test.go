package internal_test

import (
	"os"
	"testing"
	"time"

	. "github.com/blugnu/test"

	"github.com/blugnu/env/internal"
)

type mockFileInfo struct{ dir bool }

func (*mockFileInfo) Name() string       { return "mock" }
func (*mockFileInfo) Size() int64        { return 0 }
func (*mockFileInfo) Mode() os.FileMode  { return 0 }
func (*mockFileInfo) ModTime() time.Time { return time.Time{} }
func (m *mockFileInfo) IsDir() bool      { return m.dir }
func (*mockFileInfo) Sys() any           { return nil }

func TestFileExists(t *testing.T) {
	With(t)

	// act + assert
	Expect(internal.FileExists(".env")).To(BeFalse())
	Expect(internal.FileExists("consts.go")).To(BeTrue())

	Run(Test("when Stat returns an error other than ErrNotExist", func() {
		// arrange
		errStat := os.ErrPermission
		defer Restore(Original(&internal.Stat).ReplacedBy(func(string) (os.FileInfo, error) {
			return nil, errStat
		}))

		// act + assert
		Expect(internal.FileExists("anyfile")).To(BeFalse())
	}))

	Run(Test("when Stat returns a non-nil file info and an error", func() {
		// arrange
		errStat := os.ErrPermission
		defer Restore(Original(&internal.Stat).ReplacedBy(func(string) (os.FileInfo, error) {
			return &mockFileInfo{}, errStat
		}))

		// act + assert
		Expect(internal.FileExists("anyfile")).To(BeTrue())
	}))

	Run(Test("when Stat returns a non-nil file info and ErrNotExist", func() {
		// arrange
		defer Restore(Original(&internal.Stat).ReplacedBy(func(string) (os.FileInfo, error) {
			return &mockFileInfo{}, os.ErrNotExist
		}))

		// act + assert
		Expect(internal.FileExists("anyfile")).To(BeFalse())
	}))

	Run(Test("when filename is a directory", func() {
		// arrange
		defer Restore(Original(&internal.Stat).ReplacedBy(func(string) (os.FileInfo, error) {
			return &mockFileInfo{dir: true}, nil
		}))

		// act + assert
		Expect(internal.FileExists("anydir")).To(BeFalse())
	}))
}

func TestNewFileReader(t *testing.T) {
	With(t)

	Run(Test("when Open is successful", func() {
		// arrange
		defer Restore(Original(&internal.Open).ReplacedBy(func(string) (*os.File, error) {
			f, err := os.CreateTemp(t.TempDir(), "test:open-ok-*")
			Require(err).IsNil()
			return f, nil
		}))

		// act
		result, err := internal.NewFileReader("/path/to/file")
		if result != nil {
			defer func() { _ = result.Close() }()
		}

		// assert
		Expect(err).IsNil()
		Expect(result).IsNotNil()
	}))

	Run(Test("when Open fails", func() {
		// arrange
		errOpen := os.ErrNotExist
		defer Restore(Original(&internal.Open).ReplacedBy(func(string) (*os.File, error) {
			return nil, errOpen
		}))

		// act
		result, err := internal.NewFileReader("/path/to/file")

		// assert
		Expect(err).Is(errOpen)
		Expect(result).IsNil()
	}))
}
