package postgres

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
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
)

func aesEncrypt(bytesToEncrypt []byte, aesKey []byte) (string, error) {
	cipherBlock, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	return string(gcm.Seal(nonce, nonce, bytesToEncrypt, nil)), nil
}

func aesDecrypt(bytesToDecrypt []byte, aesKey []byte) ([]byte, error) {
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

func readKeyFromFile(fileLocation string) ([]byte, error) {
	keyBase64, err := ioutil.ReadFile(fileLocation)
	if err != nil {
		return []byte{}, errors.New("reading file '" + fileLocation + "':" + err.Error())
	}

	keyBase64String := string(keyBase64)

	key, err := base64.StdEncoding.DecodeString(keyBase64String)
	if err != nil {
		return []byte{}, errors.New("AES key cannot be decoded")
	}

	// verify the key works
	_, err = aes.NewCipher(key)
	if err != nil {
		return []byte{}, err
	}

	return key, nil
}
