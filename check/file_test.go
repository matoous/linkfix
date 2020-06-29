package check

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/matoous/linkfix/models"
	"github.com/matoous/linkfix/models/severity"
)

func mustParseURL(u string) *url.URL {
	uri, err := url.Parse(u)
	if err != nil {
		panic(err)
	}
	return uri
}

func TestFile(t *testing.T) {
	t.Run("existing file", func(t *testing.T) {
		f, err := ioutil.TempFile("", "linkfix-test-file-*")
		assert.NoError(t, err, "shouldn't return error")
		err = f.Close()
		assert.NoError(t, err, "shouldn't return error")
		defer func() {
			err := os.Remove(f.Name())
			assert.NoError(t, err, "shouldn't return error")
		}()

		fix, err := File(models.Link{
			URL: mustParseURL(fmt.Sprintf("file://%s", f.Name())),
		})
		assert.NoError(t, err, "shouldn't return error")
		assert.Equal(t, "", fix.Reason)
	})

	t.Run("non-existent file", func(t *testing.T) {
		fix, err := File(models.Link{
			URL: mustParseURL("file://This/File/For/Sure/Doesnt/Exist"),
		})
		assert.NoError(t, err, "shouldn't return error")
		assert.Equal(t, severity.Error, fix.Severity)
	})
}
