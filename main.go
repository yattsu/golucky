package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
)

func init() {
	log.SetFlags(0)
}

func main() {
	filePath, err := pathToDummy()
	if err != nil {
		log.Panicln("Could not find file.")
	}

	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		log.Panicln("Could not generate volatile.")
	}

	encryptFile(key, filePath)
	var pause string
	fmt.Scan(&pause)
	decryptFile(key, filePath)
}

func pathToDummy() (string, error) {
	fileName := "cat.jpg"
	var path string
	var err error

	switch runtime.GOOS {
	case "linux":
		if _, err = os.Stat(fmt.Sprintf("/tmp/%s", fileName)); err == nil {
			path = fmt.Sprintf("/tmp/%s", fileName)
		}
	default:
		path = ""
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
		log.Panicln("Could not overwrite file.")
	}

	return err
}

func decryptFile(key []byte, path string) (err error) {
	fileContents, err := os.ReadFile(path)
	if err != nil || len(path) == 0 {
		log.Println("File not found.")
	}

	decrypted, _ := decrypt(key, fileContents)

	if err := os.WriteFile(path, decrypted, 0644); err != nil {
		log.Panicln("Could not overwrite file.")
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
		log.Panicln("Could not fill iv with random bytes.")
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
