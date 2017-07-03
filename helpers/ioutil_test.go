package helpers

import (
	"io"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

var testDataDir = filepath.Join("testdata")

func testFile(name string) string {
	return filepath.Join(testDataDir, name)
}

func TestCopyFileContents(t *testing.T) {
	f, err := ioutil.TempFile("", "ioutil-test")
	if err != nil {
		t.Fatal(err)
	}
	err = f.Close()
	require.NoError(t, err)
	err = CopyFileContents(f.Name(), testFile("test.txt"), 0x644)
	require.NoError(t, err)

	b, err := ioutil.ReadFile(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, "content\n", string(b))

	err = CopyFileContents(f.Name(), testFile("missing.txt"), 0x644)
	require.Error(t, err)
}

func TestReadFileMagic(t *testing.T) {
	b, err := ReadFileMagic(testFile("test.txt"))
	require.NoError(t, err)
	require.Equal(t, "cont", string(b))

	b, err = ReadFileMagic(testFile("empty.txt"))
	require.NoError(t, err)
	require.Equal(t, []byte{0, 0, 0, 0}, b)

	_, err = ReadFileMagic(testFile("missing.txt"))
	require.Error(t, err)
}

func TestVisitCreatedFile(t *testing.T) {
	f, err := ioutil.TempFile("", "ioutil-test")
	if err != nil {
		t.Fatal(err)
	}
	err = f.Close()
	require.NoError(t, err)
	err = VisitCreatedFile(f.Name(), func(w io.Writer) error {
		_, e := w.Write([]byte("test"))
		return e
	})
	require.NoError(t, err)

	b, err := ioutil.ReadFile(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, "test", string(b))
}
