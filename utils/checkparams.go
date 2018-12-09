/**
*FileName: check
*Create on 2018/12/6 上午4:43
*Create by mok
*/

package utils

import (
	"strings"
	"regexp"
	"errors"
)

//检查参数
func CheckParams(params []*string,ops... operation)error{
	for _, op  := range ops{
		if err := op(params);err != nil{
			return err
		}
	}
	return nil
}
type operation func(params []*string)error

//否为空
func WithNotEmpty()operation{
	return func(params []*string) error {
		for _,param := range params{
			if *param == ""{
				return errors.New("传入参数错误，参数不能为空")
			}
		}
		return nil
	}
}

//去掉空格
func WithTram()operation{
	return func(params []*string) error {
		for _,param := range params{
			*param = strings.Trim(*param," ")
		}
		return nil
	}
}

//邮箱格式检验
func WithEmailRule(loc int)operation{
	return func(params []*string) error {
			emailregstr := `^[a-zA-Z0-9_.-]+@[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*\.[a-zA-Z0-9]{2,6}$`
			emailreg := regexp.MustCompile(emailregstr)
			if !emailreg.MatchString(*params[loc]){
				return errors.New("邮箱格式不正确")
			}
		return nil
	}
}