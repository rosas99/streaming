package validation

import (
	"crypto/tls"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/util/validation/field"
	netutils "k8s.io/utils/net"

	"github.com/rosas99/streaming/internal/pkg/known"
)

const (
	DNSName string = `^([a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})*[\._]?$`
)

var rxDNSName = regexp.MustCompile(DNSName)

// IsValiadURL tests that https://host:port is reachble in timeout.
func IsValiadURL(url string, timeout time.Duration) error {
	client := &http.Client{
		// disabel redirect func for import clusternet proxy cluster case
		CheckRedirect: func(rq *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout: timeout,
			}).DialContext,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	rquest, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	_, err = client.Do(rquest)
	if err != nil {
		return err
	}

	return nil
}

func IsValidDNSName(str string) bool {
	if str == "" || len(strings.Replace(str, ".", "", -1)) > 255 {
		return false
	}
	return !IsValidIP(str) && rxDNSName.MatchString(str)
}

func IsValidIP(str string) bool {
	return net.ParseIP(str) != nil
}

func ValidateHostPort(input string, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	hostIP, port, err := net.SplitHostPort(input)
	if err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath, input, "must be IP:port"))
		return allErrs
	}

	if ip := netutils.ParseIPSloppy(hostIP); ip == nil {
		allErrs = append(allErrs, field.Invalid(fldPath, hostIP, "must be a valid IP"))
	}

	if p, err := strconv.Atoi(port); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath, port, "must be a valid port"))
	} else if p < 1 || p > 65535 {
		allErrs = append(allErrs, field.Invalid(fldPath, port, "must be a valid port"))
	}

	return allErrs
}

func IsAdminUser(userID string) bool {
	return userID == known.AdminUserID
}
