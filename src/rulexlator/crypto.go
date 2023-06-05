/*
	SPDX-License-Identifier: Apache-2.0

	Copyright (c) 2023 Cisco Systems, Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
)

var key = []byte("c0blDa1VLjpQeJIy")
var masterPassword = []byte("6SJrn6p7RO2D17AY")

func getToken(forString string) string {
	fullToken, err := EncryptMessage(key, forString)
	if err != nil {
		log.Printf("Cannot encrypt string: %v\n", err)
		return ""
	}
	return fullToken[0:16]
}

func verifyToken(forString string, token string) bool {
	expectedToken := getToken(forString)
	return expectedToken == token
}

func EncryptMessage(key []byte, message string) (string, error) {
	str := key
	msg := []byte(message)
	str = append(str, msg...)
	hasher := sha256.New()
	hasher.Write(str)
	sha256 := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha256, nil
}

func DecryptMessage(key []byte, message string) (string, error) {
	cipherText, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return "", fmt.Errorf("could not base64 decode: %v", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("could not create new cipher: %v", err)
	}

	if len(cipherText) < aes.BlockSize {
		return "", fmt.Errorf("invalid ciphertext block size")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}

func test() {
	key := []byte("0123456789abcdef") // must be of 16 bytes for this example to work
	message := "1234567890123456"

	encrypted, _ := EncryptMessage(key, message)
	fmt.Println(encrypted)
	fmt.Println(len(encrypted))

	decrypted, _ := DecryptMessage(key, encrypted)
	fmt.Println(decrypted)
}
