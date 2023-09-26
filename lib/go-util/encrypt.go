package util

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

// AESEncrypt takes in a 16, 24, or 32 byte AES key (128, 192, 256 bit encryption respectively) and plain text. It
// returns the encrypted text. In case of error, the text returned is an empty string. AES requires input text to
// be greater than 12 bytes in length.
func AESEncrypt(bytesToEncrypt []byte, aesKey []byte) ([]byte, error) {
	cipherBlock, err := aes.NewCipher(aesKey)
	if err != nil {
		return []byte{}, err
	}

	gcm, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return []byte{}, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return []byte{}, err
	}

	return gcm.Seal(nonce, nonce, bytesToEncrypt, nil), nil
}

// AESDecrypt takes in a 16, 24, or 32 byte AES key (128, 192, 256 bit encryption respectively) and encrypted text. It
// returns the resulting decrypted text. In case of error, the text returned is an empty string. AES requires input
// text to be greater than 12 bytes in length.
func AESDecrypt(bytesToDecrypt []byte, aesKey []byte) ([]byte, error) {
	cipherBlock, err := aes.NewCipher(aesKey)
	if err != nil {
		return []byte{}, err
	}

	gcm, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return []byte{}, err
	}

	nonceSize := gcm.NonceSize()
	if len(bytesToDecrypt) < nonceSize {
		return []byte{}, err
	}

	nonce := bytesToDecrypt[:nonceSize]
	ciphertext := bytesToDecrypt[nonceSize:]
	decryptedString, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return []byte{}, err
	}
	return decryptedString, nil
}

// ValidateAESKey takes in a byte slice and tests if it's a valid AES key (16,
// 24, or 32 bytes), returning an error if it isn't.
func ValidateAESKey(keyBytes []byte) error {
	_, err := aes.NewCipher(keyBytes)
	return err
}
