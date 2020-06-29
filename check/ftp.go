package check

import (
	"fmt"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"

	"github.com/matoous/linkfix/models"
	"github.com/matoous/linkfix/models/severity"
)

func FTP(link models.Link) (models.Fix, error) {
	fix := models.Fix{Link: link}
	c, err := ftp.Dial(link.URL.Host, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		fix.Reason = fmt.Sprintf("establish connection with the server: %s", err.Error())
		fix.Severity = severity.Error
		return fix, nil
	}
	pw, _ := link.URL.User.Password()
	err = c.Login(link.URL.User.Username(), pw)
	if err != nil {
		fix.Reason = fmt.Sprintf("authorize with the server: %s", err.Error())
		fix.Severity = severity.Error
		return fix, nil
	}

	if link.URL.Path == "" || strings.HasSuffix(link.URL.Path, "/") {
		// directory
		err := c.ChangeDir(link.URL.Path)
		if err != nil {
			fix.Reason = fmt.Sprintf("couldn't find the directory `%s`: %s", link.URL.Path, err.Error())
			fix.Severity = severity.Error
			return fix, nil
		}
	} else {
		// file
		_, err := c.FileSize(link.URL.Path)
		if err != nil {
			fix.Reason = fmt.Sprintf("couldn't find the file `%s`: %s", link.URL.Path, err.Error())
			fix.Severity = severity.Error
			return fix, nil
		}
	}

	if err := c.Quit(); err != nil {
		return fix, err
	}
	return fix, nil
}
