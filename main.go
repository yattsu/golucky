package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	mathrand "math/rand/v2"
	"os"
	"os/exec"
	"reflect"
	"runtime"
)

var clear map[string]func()

func init() {
	log.SetFlags(0)

	clear = make(map[string]func())
	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func main() {
	callClear()

	lives := 253
	filePath, err := pathToCat()
	if err != nil {
		log.Println("A cat is not present...")
		os.Exit(0)
	}
	if !checkSum() {
		log.Println("The cat...")
		log.Println("...it's...dust.")
		os.Exit(0)
	}

	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		log.Panicln("Could not generate volatile.")
	}

	if err := encryptFile(key, filePath); err != nil {
		log.Panicln("Could not encrypt cat.")
	}

	log.Println("The cat's soul is jumbled (by AES256 encryption...)")
	log.Println("Restore its mangled soul by winning the slots")
	log.Println("Closing this program will cause the cat to cease forever")
	log.Println(fmt.Sprintf("You have %d lives left", lives))
	var pause string
	log.Println("\nPress Enter to spin.")
	fmt.Scanf("\n", &pause)

	win := slots(lives)

	if win {
		decryptFile(key, filePath)
		log.Println("\nGambled a cat's soul and won")
		log.Println("It is safe...")
		log.Println("               ...for now.")
		os.Exit(0)
	} else {
		callClear()
		log.Println("The cat is forever lost...")
		log.Println("         ...it is only a remnant of what once was...")
		log.Println("                  ...forever in your memory.")
		os.Exit(0)
	}
}

func callClear() {
	value, ok := clear[runtime.GOOS]
	if ok {
		value()
	}
}

func checkSum() bool {
	const sha = "c732789e3469f68d210e6fba42839a3f61420e71ffcfc606f385d6de7d8d926c"
	var match bool
	match = false

	f, err := os.Open("ender.jpg")
	if err != nil {
		log.Panicln("Could not check the cat's soul.")
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Panicln("Could not check the cat's soul.")
	}

	if fmt.Sprintf("%x", h.Sum(nil)) == sha {
		match = true
	}

	return match
}

func slots(lives int) bool {
	results := make([]int, 3)
	win := false

	var pause string
	for {
		callClear()
		for idx := range results {
			results[idx] = mathrand.IntN(10)
		}
		lives -= 1
		log.Println("_ - _ - _ - _")
		log.Println(fmt.Sprintf("| %d | %d | %d |", results[0], results[1], results[2]))
		log.Println("= = = = = = =")
		log.Println(fmt.Sprintf("   ❤️ %d ❤️", lives))
		if reflect.DeepEqual(results, []int{4, 4, 4}) {
			win = true
			break
		}
		if lives == 0 {
			log.Println("\nOh no...")
		} else {
			log.Println("\nPress Enter to spin.")
		}
		fmt.Scanf("\n", &pause)
		if lives <= 0 {
			break
		}
	}

	return win
}

func pathToCat() (string, error) {
	fileName := "ender.jpg"
	var path string
	var err error

	if _, err = os.Stat(fileName); err == nil {
		path = fileName
	}

	return path, err
}

func encryptFile(key []byte, path string) (err error) {
	fileContents, err := os.ReadFile(path)
	if err != nil || len(path) == 0 {
		log.Println("File not found.")
	}

	encrypted, _ := encrypt(key, fileContents)

	if err := os.WriteFile(path, encrypted, 0644); err != nil {
		log.Panicln("Could not overwrite cat.")
	}

	return err
}

func decryptFile(key []byte, path string) (err error) {
	fileContents, err := os.ReadFile(path)
	if err != nil || len(path) == 0 {
		log.Println("Cat not found.")
	}

	decrypted, _ := decrypt(key, fileContents)

	if err := os.WriteFile(path, decrypted, 0644); err != nil {
		log.Panicln("Could not overwrite cat.")
	}

	return err
}

func makeCipher(key []byte) cipher.Block {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Panicln(err)
	}

	return block
}

func encrypt(key []byte, fileContents []byte) (encrypted []byte, err error) {
	cipherBlock := makeCipher(key)
	cipherText := make([]byte, aes.BlockSize+len(fileContents))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Panicln("Could not fill iv with dust.")
	}

	stream := cipher.NewCFBEncrypter(cipherBlock, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], fileContents)

	return cipherText, err
}

func decrypt(key []byte, cipherText []byte) (decrypted []byte, err error) {
	cipherBlock := makeCipher(key)
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(cipherBlock, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return cipherText, err
}
