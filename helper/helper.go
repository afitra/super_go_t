package helper

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func DelayProses(detik int) {
	fmt.Printf("Menunda proses selama %d detik...\n", detik)
	time.Sleep(time.Duration(detik) * time.Second)
}

func Convert_string_to_bool(s string) bool {
	if result, err := strconv.ParseBool(s); err == nil {
		return result
	}
	return false
}

func Convert_string_to_duration(s string) (time.Duration, error) {
	var i int64
	var err error
	if i, err = strconv.ParseInt(s, 10, 64); err != nil {
		return 0, err
	}
	return time.Duration(i) * time.Hour, err
}

func Split_string_to_array(s string) []string {
	return strings.Split(s, "\n")
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func VerifyPassword(hashedPassword, inputPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
	return err
}

func Generate_code(prefix string, length int) string {
	var sb strings.Builder
	sb.WriteString(prefix)
	sb.WriteString("-")

	for i := 0; i < length; i++ {
		sb.WriteByte(letterBytes[rand.Intn(len(letterBytes))])
	}

	return sb.String()
}
