package rules

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/url"
	"strings"
)

const (
	a0 = "nomad"
	a1 = "low"
	a2 = "moderate"
	a3 = "high"
	a4 = "extreme"

)

var maliciousDomainList = readFileList("D:\\doan_siem\\correl\\rules\\domain_black")


func IsMaliciousURL(inputURL string) (b bool, alert string,level string) {
	parsedURL, err := url.Parse(inputURL)
	logrus.Println(parsedURL)
	if err != nil {
		logrus.Error("fail parse URL:", err)
		return false, "", a0
	}

	host := parsedURL.Hostname()

	for _, maliciousDomain := range maliciousDomainList {
		if strings.Contains(host, maliciousDomain) {
			return true, fmt.Sprintf(" Malicious URL: %s ", inputURL) , a2
		}
	}
	return false, "",  a0
}