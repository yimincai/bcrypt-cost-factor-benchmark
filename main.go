package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"time"
	"unsafe"

	_ "go.uber.org/automaxprocs"
	"golang.org/x/crypto/bcrypt"
)

var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@_#"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func main() {
	valStart := bcrypt.MinCost
	valEnd := bcrypt.MaxCost
	if os.Getenv("PASSWORD_LENGTH") == "" {
		os.Setenv("PASSWORD_LENGTH", "18")
	}
	lengthOfPassword, err := strconv.Atoi(os.Getenv("PASSWORD_LENGTH"))
	if err != nil {
		fmt.Println("Error during conversion")
		return
	}
	if lengthOfPassword <= 0 {
		fmt.Printf("Invalid number : %v\n", lengthOfPassword)
		return
	}

	log.Printf("PASSWORD_LENGTH = %v", lengthOfPassword)
	fmt.Println("========================================================================================")
	for cost := valStart; cost <= valEnd; cost++ {
		password := RandString(lengthOfPassword)
		log.Printf("COST FACTOR : %v\n", cost)
		Total(password, cost)
		fmt.Println("========================================================================================")
	}
	forever := make(chan bool)

	log.Printf("[*] Benchmark finished, to exit please press CTRL+C")
	<-forever
}

func Total(password string, costFactor int) {
	encryptedPassword := HashPassword(password, costFactor)
	err := ComparePassword([]byte(encryptedPassword), []byte(password))
	if err != nil {
		log.Fatalf("Failed to compare password:%v\n", err)
	}
}

func HashPassword(password string, costFactor int) string {
	defer TimeTracker(time.Now())
	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		costFactor,
	)
	if err != nil {
		log.Fatal(err)
	}
	return string(encryptedPassword)
}

func ComparePassword(hashedPassword, password []byte) error {
	defer TimeTracker(time.Now())

	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}

// RandString
func RandString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

func TimeTracker(start time.Time) {
	elapsed := time.Since(start)

	// Skip this function, and fetch the PC and file for its parent.
	pc, _, _, _ := runtime.Caller(1)

	// Retrieve a function object this functions parent.
	funcObj := runtime.FuncForPC(pc)

	// Regex to extract just the function name (and not the module path).
	runtimeFunc := regexp.MustCompile(`^.*\.(.*)$`)
	name := runtimeFunc.ReplaceAllString(funcObj.Name(), "$1")

	log.Printf("%s\ttook\t%s", name, elapsed)
}
