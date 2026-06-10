package service

import (
	"errors"
	"strings"
)

func cleanPhone(phone string) (string, error) {
	index := strings.IndexRune(phone, '8')
	if index < 0 {
		return "", errors.New("invalid phone number")
	}

	return phone[index:], nil
}
