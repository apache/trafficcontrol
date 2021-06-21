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
	"math/rand"
	"reflect"
	"testing"
	"time"
)

var (
	validKeySizes = [...]int{16, 24, 32}
)

func randomByteArray(size int) []byte {
	data := make([]byte, size)
	for i, _ := range data {
		data[i] = byte('a' + rune(rand.Intn(26)))
	}
	return data
}

func TestInvalidKey(t *testing.T) {
	rand.Seed(time.Now().Unix())

	text := []byte("this is my favorite test on the citadel")
	emptyKey := []byte{}

	err := ValidateAESKey(emptyKey)
	if err == nil {
		t.Fatal("expected ValidateAESKey to return error with empty key")
	}

	_, err = AESEncrypt(emptyKey, text)
	if err == nil {
		t.Fatal("expected AESEncrypt to return error with empty key")
	}
	_, err = AESDecrypt(emptyKey, text)
	if err == nil {
		t.Fatal("expected AESDecrypt to return error with empty key")
	}

	for _, keySize := range validKeySizes {
		shortKey := randomByteArray(keySize - 1)
		longKey := randomByteArray(keySize + 1)
		if len(shortKey) > keySize-1 {
			t.Fatal("expected shortKey must be less than 32 characters")
		}
		if len(longKey) < keySize {
			t.Fatal("expected longKey to be more than 32 characters")
		}

		err = ValidateAESKey(shortKey)
		if err == nil {
			t.Fatal("expected ValidateAESKey to return error with short key")
		}

		err = ValidateAESKey(longKey)
		if err == nil {
			t.Fatal("expected ValidateAESKey to return error with long key")
		}

		_, err = AESEncrypt(shortKey, text)
		if err == nil {
			t.Fatal("expected AESEncrypt to return error with short key")
		}
		_, err = AESDecrypt(shortKey, text)
		if err == nil {
			t.Fatal("expected AESDecrypt to return error with short key")
		}

		_, err = AESEncrypt(longKey, text)
		if err == nil {
			t.Fatal("expected AESEncrypt to return err with long key")
		}
		_, err = AESDecrypt(longKey, text)
		if err == nil {
			t.Fatal("expected AESDecrypt to return err with long key")
		}
	}
}

func TestInvalidText(t *testing.T) {
	for _, keySize := range validKeySizes {
		validKey := randomByteArray(keySize)
		shortText := []byte("hello")
		_, err := AESEncrypt(validKey, shortText)
		if err == nil {
			t.Fatal("expected AESEncrypt to return error with short text")
		}
		_, err = AESDecrypt(validKey, shortText)
		if err == nil {
			t.Fatal("expected AESDecrypt to return error with short text")
		}
	}
}

func TestValidKey(t *testing.T) {
	text := []byte("this is my favorite test on the citadel")
	for _, keySize := range validKeySizes {
		validKey := randomByteArray(keySize)

		err := ValidateAESKey(validKey)
		if err != nil {
			t.Fatal("expected ValidatedAESKey to succeed, got : " + err.Error())
		}

		encText, err := AESEncrypt(text, validKey)
		if err != nil {
			t.Fatal("expected AESEncrypt to succeed, got: " + err.Error())
		}

		if reflect.DeepEqual(encText, text) {
			t.Fatal("expected AESEncrypt to encrypt text")
		}

		decText, err := AESDecrypt(encText, validKey)
		if err != nil {
			t.Fatal("expected AESDecrypt to succeed, got: " + err.Error())
		}

		if reflect.DeepEqual(encText, decText) {
			t.Fatal("expected AESDecrypt to change encrypted text")
		}

		if !reflect.DeepEqual(text, decText) {
			t.Fatal("expected AESDecrypt to return the original plain text, got: " + string(decText))
		}
	}
}
