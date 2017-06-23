package helpers

import (
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopyFileContents(t *testing.T) {
	f, err := ioutil.TempFile("", "ioutil-test")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	err = CopyFileContents(f.Name(), "test/test.txt", 0x644)
	require.NoError(t, err)

	b, err := ioutil.ReadFile(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, "content\n", string(b))

	err = CopyFileContents(f.Name(), "test/missing.txt", 0x644)
	require.Error(t, err)
}

func TestReadFileMagic(t *testing.T) {
	b, err := ReadFileMagic("test/test.txt")
	require.NoError(t, err)
	require.Equal(t, "cont", string(b))

	b, err = ReadFileMagic("test/empty.txt")
	require.NoError(t, err)
	require.Equal(t, []byte{0, 0, 0, 0}, b)

	b, err = ReadFileMagic("test/missing.txt")
	require.Error(t, err)
}

func TestVisitCreatedFile(t *testing.T) {
	f, err := ioutil.TempFile("", "ioutil-test")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	VisitCreatedFile(f.Name(), func(w io.Writer) error {
		w.Write([]byte("test"))
		return nil
	})

	b, err := ioutil.ReadFile(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, "test", string(b))
}
