/*
 * Octra Go SDK - OSM-15 (Octra Structured Message)
 *
 * Copyright (c) 2026 Qiubit Team
 * Licensed under the MIT License
 * Version: 0.2.1
 *
 * This module handles secure keystore management using 
 * AES-256-GCM and Scrypt key derivation for the Octra network.
 */

package osm15

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519" 
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"

	"golang.org/x/crypto/scrypt"
)

type Keystore struct {
	Address string `json:"address"`
	Crypto  Crypto `json:"crypto"`
}

type Crypto struct {
	Cipher     string `json:"cipher"`
	CipherText string `json:"ciphertext"`
	CipherParams struct {
		IV string `json:"iv"`
	} `json:"cipherparams"`
	Kdf       string `json:"kdf"`
	KdfParams struct {
		N     int    `json:"n"`
		R     int    `json:"r"`
		P     int    `json:"p"`
		Salt  string `json:"salt"`
	} `json:"kdfparams"`
}

func EncryptKey(privateKeyB64 string, password string) ([]byte, error) {
	seed, _ := base64.StdEncoding.DecodeString(privateKeyB64)
	
	// Generate Address from private key
	priv := ed25519.NewKeyFromSeed(seed)
	address := PublicKeyToAddress(priv.Public().(ed25519.PublicKey))

	salt := make([]byte, 32)
	io.ReadFull(rand.Reader, salt)

	derivedKey, _ := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)

	block, _ := aes.NewCipher(derivedKey)
	gcm, _ := cipher.NewGCM(block)
	iv := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, iv)

	cipherText := gcm.Seal(nil, iv, seed, nil)

	ks := Keystore{
		Address: address, 
		Crypto: Crypto{
			Cipher:     "aes-256-gcm",
			CipherText: base64.StdEncoding.EncodeToString(cipherText),
			Kdf:        "scrypt",
		},
	}
	ks.Crypto.CipherParams.IV = base64.StdEncoding.EncodeToString(iv)
	ks.Crypto.KdfParams.N = 32768
	ks.Crypto.KdfParams.R = 8
	ks.Crypto.KdfParams.P = 1
	ks.Crypto.KdfParams.Salt = base64.StdEncoding.EncodeToString(salt)

	return json.MarshalIndent(ks, "", "  ")
}

func DecryptKey(keystoreJSON []byte, password string) (string, error) {
	var ks Keystore
	if err := json.Unmarshal(keystoreJSON, &ks); err != nil {
		return "", err
	}

	salt, _ := base64.StdEncoding.DecodeString(ks.Crypto.KdfParams.Salt)
	iv, _ := base64.StdEncoding.DecodeString(ks.Crypto.CipherParams.IV)
	rawCipher, _ := base64.StdEncoding.DecodeString(ks.Crypto.CipherText)

	derivedKey, _ := scrypt.Key([]byte(password), salt, ks.Crypto.KdfParams.N, ks.Crypto.KdfParams.R, ks.Crypto.KdfParams.P, 32)

	block, _ := aes.NewCipher(derivedKey)
	gcm, _ := cipher.NewGCM(block)

	seed, err := gcm.Open(nil, iv, rawCipher, nil)
	if err != nil {
		return "", fmt.Errorf("invalid password")
	}

	return base64.StdEncoding.EncodeToString(seed), nil
}
