package common

import "regexp"

// RegexpDate 匹配日期格式：[]string{str1, str2, str3, str4, str5}
func RegexpDate(str string) bool {

	//str1 := "api-2021-01-10.log"
	//str2 := "api.log"
	//str3 := "api-2021-1-1.log"
	//str4 := "api-2021-11-1.log"
	//str5 := "api-2021-1-12.log"

	reg := regexp.MustCompile(`.*\d{4}-\d?\d?-\d?\d?.*`)
	//
	//for _, s := range []string{str1, str2, str3, str4, str5} {
	//	if len(reg.FindAllString(s, -1)) == 0 {
	//		log.Println(s)
	//	}
	//}
	if len(reg.FindAllString(str, -1)) == 0 {
		return false
	}
	return true
}
