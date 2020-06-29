package check

import (
	"fmt"
	"net/mail"
	"net/url"
	"strings"

	"github.com/mcnijman/go-emailaddress"

	"github.com/matoous/linkfix/models"
	"github.com/matoous/linkfix/models/severity"
)

var cgiAddresses = []string{"to", "cc", "bcc"}

func MailTo(link models.Link) (models.Fix, error) {
	fix := models.Fix{Link: link}
	addresses, err := parseAllAddresses(link.URL)
	if err != nil {
		return models.Fix{}, err
	}
	var errs []error
	for i := range addresses {
		err = validateMail(addresses[i].Address)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		fix.Reason = fmt.Sprintf("%v", errs)
		fix.Severity = severity.Error
	}
	return fix, nil
}

func validateMail(m string) error {
	email, err := emailaddress.Parse(m)
	if err != nil {
		return fmt.Errorf("invalid email `%s`: %s", m, err.Error())
	}
	err = email.ValidateIcanSuffix()
	if err != nil {
		return fmt.Errorf("not an icann suffix: %s", err.Error())
	}
	err = email.ValidateHost()
	if err != nil {
		// https://serversmtp.com/smtp-error/
		switch {
		case strings.Contains(err.Error(), "550"):
			return fmt.Errorf("invalid host: email account `%s` doesn't exist", m)
		case strings.Contains(err.Error(), "513"):
			return fmt.Errorf("invalid host: address type inccorect for `%s`", m)
		default:
			return fmt.Errorf("invalid host: %s", err.Error())
		}
	}
	return nil
}

func parseAllAddresses(uri *url.URL) ([]*mail.Address, error) {
	addresses, err := mail.ParseAddressList(uri.Opaque)
	if err != nil {
		return nil, err
	}
	params, err := url.ParseQuery(uri.RawQuery)
	if err != nil {
		// if the mailto is invalid we want to report it, so return nil, err while we also could return addresses from
		// first step
		return nil, err
	}
	for i := range cgiAddresses {
		l := params.Get(cgiAddresses[i])
		if l == "" {
			continue
		}
		addrl, err := mail.ParseAddressList(l)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, addrl...)
	}
	return addresses, nil
}
