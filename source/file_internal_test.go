package source

import (
    "fmt"
    "io/ioutil"
    "os"
    "testing"

    "github.com/stretchr/testify/assert"
)

func Test_listFiles(t *testing.T) {
    tmpDir, err := ioutil.TempDir("", "linkfixtest-*")
    assert.NoError(t, err, "couldn't create tmp dir")
    defer func() {
        err = os.RemoveAll(tmpDir)
        assert.NoError(t, err, "couldn't cleanup tmp directory")
    }()

    testFile1 := fmt.Sprintf("%s/test_file1.txt", tmpDir)
    err = ioutil.WriteFile(testFile1, nil, 0700)
    assert.NoError(t, err, "couldn't create tmp file")

    err= os.MkdirAll(fmt.Sprintf("%s/deeply/nested/directory", tmpDir), 0700)
    assert.NoError(t, err, "couldn't create nested tmp directory")

    testFile2 := fmt.Sprintf("%s/deeply/nested/directory/test_file2.txt", tmpDir)
    err = ioutil.WriteFile(testFile2, nil, 0700)
    assert.NoError(t, err, "couldn't create tmp file in nested directory")

    want := []string{testFile1, testFile2}
    got, err := listFiles(tmpDir)
    assert.NoError(t, err, "shouldn't return error")
    assert.ElementsMatch(t, want, got, "should list all files")

}