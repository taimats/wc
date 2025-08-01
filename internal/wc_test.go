package internal_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/taimats/wc/internal"
)

func TestCountWords(t *testing.T) {
	f, err := os.Open("testfile1")
	if err != nil {
		t.Fatalf("failed to open file: %s", err)
	}
	defer f.Close()
	got := internal.CountWords(f)
	if got != 17 {
		t.Errorf("Not Equal: (want=%d, got=%d)", 17, got)
	}
}

func TestCmdWC(t *testing.T) {
	t.Run("正常系1:ファイルのみ指定", func(t *testing.T) {
		args := []string{"wc", "./testfile1"}
		err := internal.CmdWC(args)
		if err != nil {
			t.Fatalf("CmdWC failed to Run: (error: %s)", err)
		}
	})
	t.Run("正常系2:ディレクトリを指定", func(t *testing.T) {
		args := []string{"wc", "./testdir2"}
		err := internal.CmdWC(args)
		if err != nil {
			t.Fatalf("CmdWC failed to Run: (error: %s)", err)
		}
	})

	t.Run("異常系:パス名が不明", func(t *testing.T) {
		args := []string{"wc", "./testtest"}
		err := internal.CmdWC(args)
		if err == nil {
			t.Fatalf("CmdWC failed to Run: (error: %s)", err)
		}
	})
}

func TestExecuteWC(t *testing.T) {
	t.Run("正常_01:単一ファイルの指定", func(t *testing.T) {
		path := "./testfile1"
		info, err := os.Stat(path)
		if err != nil {
			t.Fatal(err)
		}
		r := internal.NewRecord(info.ModTime().Format(time.DateOnly), "testfile1", 17)
		want := internal.NewRecordSlice(r)
		a := assert.New(t)

		got, err := internal.ExecuteWC(path)
		a.Equal(want, got)
		a.Nil(err)
	})

	t.Run("正常_02:ディレクトリの指定", func(t *testing.T) {
		path := "./testdir2"
		r := internal.NewRecord("2025-08-01", "test3", 17)
		want := internal.NewRecordSlice(r)
		a := assert.New(t)

		got, err := internal.ExecuteWC(path)
		a.Equal(want, got)
		a.Nil(err)
	})

	t.Run("正常_03:ディレクトリの指定:日付順にソート", func(t *testing.T) {
		path := "./testdir3"
		r1 := internal.NewRecord("2025-07-31", "test3", 17)
		r2 := internal.NewRecord("2025-08-01", "test1", 0)
		r3 := internal.NewRecord("2025-08-01", "test2", 0)
		want := internal.NewRecordSlice(r1, r2, r3)
		a := assert.New(t)

		got, err := internal.ExecuteWC(path)
		a.Equal(want, got)
		a.Nil(err)
	})

	t.Run("異常_01:ディレクトリの指定:不正な拡張子", func(t *testing.T) {
		path := "./test4"
		a := assert.New(t)

		got, err := internal.ExecuteWC(path)
		a.Nil(got)
		a.Nil(err)
	})
}
