package main

import "fmt"

func main() {
	phone := "13888888888"
	phone = PhoneMasking(phone)
	fmt.Printf("phone: %v\n", phone)

	var v1, v2 = 0.1, 0.2
	fmt.Println(v1 + v2)
}

// PhoneMasking 手机号码脱敏
func PhoneMasking(phone string) string {
	if len(phone) > 6 {
		phone = phone[0:2] + "****" + phone[len(phone)-4:]
	}
	return phone
}
