package source_test

import (
    "io/ioutil"
    "os"
    "testing"

    "github.com/stretchr/testify/assert"

    "github.com/matoous/linkfix/source"
)

func TestFilesystemSource_List(t *testing.T) {
    tmpDir, err := ioutil.TempDir("", "linkfixtest-*")
    assert.NoError(t, err, "couldn't create tmp dir")
    defer func() {
        err = os.RemoveAll(tmpDir)
        assert.NoError(t, err, "couldn't cleanup tmp directory")
    }()

    s := source.Filesystem(tmpDir, []string{})
    _, _ = s.List()
}