/**
*FileName: utils
*Create on 2018/12/7 下午9:29
*Create by mok
*/

package utils

import (
	"strings"
	"strconv"
	"errors"
	"fmt"
)

//将ip转变为int类型
func  IPToInt(IP string) (int, error) {
	Intip, err := ConvertToIntIP(IP)
	if err != nil {
		return 0, err
	}
	return Intip, nil
}

func ConvertToIntIP(IP string) (int, error) {
	ips := strings.Split(IP, ".")
	E := errors.New("Not A IP.")
	if len(ips) != 4 {
		return 0, E
	}
	var intIP int
	for k, v := range ips {
		i, err := strconv.Atoi(v)
		if err != nil || i > 255 {
			return 0, E
		}
		intIP = intIP | i<< uint(8*(3-k))
	}
	return intIP, nil
}

//将int类型转化为ip
func IntToIP(intIP int) (string, error) {
	i4 := intIP & 255
	i3 := intIP >> 8 & 255
	i2 := intIP >> 16 & 255
	i1 := intIP >> 24 & 255
	if i1 > 255 || i2 > 255 || i3 > 255 || i4 > 255 {
		return "", errors.New("Isn't a IntIP Type.")
	}
	ipstring := fmt.Sprintf("%d.%d.%d.%d", i4, i3, i2, i1)
	return ipstring, nil
}

