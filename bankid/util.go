package bankid

import (
	"crypto/md5"
	"fmt"
	"io"
	"strings"
)

func md5sum(in string) string {
	hasher := md5.New()
	io.WriteString(hasher, in)
	hash := hasher.Sum(nil)
	return fmt.Sprintf("%x", hash)
}

func IsMobileUserAgent(userAgent string) bool {
	mobileBrowsers := []string{"Mobile", "Android", "iPhone", "iPad"}
	for _, mobileBrowser := range mobileBrowsers {
		if strings.Contains(userAgent, mobileBrowser) {
			return true
		}
	}
	return false
}
