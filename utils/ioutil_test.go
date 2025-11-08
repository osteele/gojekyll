package utils

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

var testDataDir = filepath.Join("testdata")

func testFile(name string) string {
	return filepath.Join(testDataDir, name)
}

func TestCopyFileContents(t *testing.T) {
	f, err := os.CreateTemp("", "ioutil-test")
	if err != nil {
		t.Fatal(err)
	}
	err = f.Close()
	require.NoError(t, err)
	err = CopyFileContents(f.Name(), testFile("test.txt"), 0x644)
	require.NoError(t, err)

	b, err := os.ReadFile(f.Name())
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
	f, err := os.CreateTemp("", "ioutil-test")
	if err != nil {
		t.Fatal(err)
	}
	err = f.Close()
	require.NoError(t, err)
	err = VisitCreatedFile(f.Name(), func(w io.Writer) error {
		_, e := io.WriteString(w, "test")
		return e
	})
	require.NoError(t, err)

	b, err := os.ReadFile(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, "test", string(b))
}

func TestRemoveEmptyDirectories(t *testing.T) {
	// Create a temporary directory structure
	tmpDir, err := os.MkdirTemp("", "remove-empty-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create subdirectories
	emptyDir := filepath.Join(tmpDir, "empty")
	require.NoError(t, os.Mkdir(emptyDir, 0755))

	nestedDir := filepath.Join(tmpDir, "nested", "deep")
	require.NoError(t, os.MkdirAll(nestedDir, 0755))

	dirWithFile := filepath.Join(tmpDir, "withfile")
	require.NoError(t, os.Mkdir(dirWithFile, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(dirWithFile, "file.txt"), []byte("content"), 0644))

	// Run RemoveEmptyDirectories
	err = RemoveEmptyDirectories(tmpDir)
	require.NoError(t, err)

	// Verify the root directory still exists
	_, err = os.Stat(tmpDir)
	require.NoError(t, err, "root directory should not be removed")

	// Verify empty subdirectories were removed
	_, err = os.Stat(emptyDir)
	require.True(t, os.IsNotExist(err), "empty subdirectory should be removed")

	_, err = os.Stat(nestedDir)
	require.True(t, os.IsNotExist(err), "nested empty subdirectory should be removed")

	// Verify directory with files still exists
	_, err = os.Stat(dirWithFile)
	require.NoError(t, err, "directory with files should not be removed")
}

func TestRemoveEmptyDirectoriesWithSymlink(t *testing.T) {
	// Create a temporary directory to be the symlink target
	targetDir, err := os.MkdirTemp("", "symlink-target")
	require.NoError(t, err)
	defer os.RemoveAll(targetDir)

	// Create a symlink to the target
	tmpDir, err := os.MkdirTemp("", "symlink-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	symlinkPath := filepath.Join(tmpDir, "dest")
	err = os.Symlink(targetDir, symlinkPath)
	require.NoError(t, err)

	// Create some empty subdirectories in the target through the symlink
	emptyDir := filepath.Join(symlinkPath, "empty")
	require.NoError(t, os.Mkdir(emptyDir, 0755))

	// Run RemoveEmptyDirectories on the symlink
	err = RemoveEmptyDirectories(symlinkPath)
	require.NoError(t, err)

	// Verify the symlink still exists
	info, err := os.Lstat(symlinkPath)
	require.NoError(t, err, "symlink should not be removed")
	require.True(t, info.Mode()&os.ModeSymlink != 0, "should still be a symlink")

	// Verify the symlink target still exists
	_, err = os.Stat(targetDir)
	require.NoError(t, err, "symlink target should still exist")

	// Verify empty subdirectories in the target were removed
	_, err = os.Stat(emptyDir)
	require.True(t, os.IsNotExist(err), "empty subdirectory should be removed")
}
