package rules

import (
	"bufio"
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"os"
	"strings"
)

var maliciousIPList =  readFileList("D:\\doan_siem\\correl\\rules\\ip_ignore.txt")

func IsMaliciousIP(ip string) (b bool, alert string,level string)  {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		logrus.Println("IP invalid")
		return false, "", ""
	}

	for _, maliciousIP := range maliciousIPList {
		_, maliciousCIDR, err := net.ParseCIDR(maliciousIP)
		if err != nil {
			continue
		}
		if maliciousCIDR.Contains(parsedIP) {
			return true, fmt.Sprintf("IP %s Malicious in %v",ip, maliciousCIDR), a3
		}
	}

	// Kiểm tra xem địa chỉ IP có thuộc dải IP bất thường hay không
	for _, maliciousIP := range maliciousIPList {
		if ip == maliciousIP {
			logrus.Warnf("IP %v Malicious ", ip)
			return true, fmt.Sprintf("IP %s Malicious ",ip), a3
		}
	}

	return false , "", ""
}

func readFileList(path string) []string {
	// Mở file để đọc
	file, err := os.Open(path)
	if err != nil {
		logrus.Errorf("Error opening file: %v, %v", path, err)
		return nil
	}
	defer file.Close()

	// Tạo một slice để lưu địa chỉ IP từ file
	var readList []string

	// Sử dụng bufio.NewScanner để đọc từng dòng từ file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Sử dụng strings.Fields để tách các địa chỉ IP trong mỗi dòng
		ips := strings.Fields(line)
		// Thêm các địa chỉ IP vào slice
		readList = append(readList, ips...)
	}

	// Kiểm tra lỗi khi đọc file
	if err = scanner.Err(); err != nil {
		logrus.Errorf("Error reading file %v: %v",path, err)
		return readList
	}

	return readList
}
