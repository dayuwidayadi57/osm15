package osm15

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestOSM15_ArrayAndRecursiveDebug(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(nil)
	privB64 := base64.StdEncoding.EncodeToString(priv.Seed())
	pubB64 := base64.StdEncoding.EncodeToString(pub)

	types := map[string][]TypedMember{
		"Wallet": {
			{Name: "owner", Type: "address"},
			{Name: "assets", Type: "Asset[]"},
		},
		"Asset": {
			{Name: "name", Type: "string"},
			{Name: "amount", Type: "uint256"},
		},
	}

	domain := TypedDomain{Name: "OctraVault", Version: "1", ChainID: 1}

	// Message dengan Array of Structs
	message := map[string]interface{}{
		"owner": "oct123",
		"assets": []map[string]interface{}{
			{"name": "OCT", "amount": 1000},
			{"name": "GOLD", "amount": 50},
		},
	}

	data := TypedData{
		Domain:      domain,
		Types:       types,
		PrimaryType: "Wallet",
		Message:     message,
	}

	// 1. Test Hashing
	digest, err := HashTypedData(data)
	if err != nil {
		t.Fatalf("Hash error: %v", err)
	}
	fmt.Printf("\n[DEBUG] Wallet Digest: %s\n", hex.EncodeToString(digest))

	// 2. Test Signing
	sig, err := SignTypedData(data, privB64)
	if err != nil {
		t.Fatalf("Sign error: %v", err)
	}
	fmt.Printf("[DEBUG] Array Signature: %s\n", sig)

	// 3. Test Verification
	valid, err := VerifyTypedData(data, sig, pubB64)
	if err != nil || !valid {
		t.Fatalf("Verification failed for array data")
	}

	// 4. Test Tamper Resistance (Ubah isi array sedikit saja)
	tamperedMessage := map[string]interface{}{
		"owner": "oct123",
		"assets": []map[string]interface{}{
			{"name": "OCT", "amount": 1001}, // Beda 1 unit
			{"name": "GOLD", "amount": 50},
		},
	}
	tamperedData := data
	tamperedData.Message = tamperedMessage

	validTampered, _ := VerifyTypedData(tamperedData, sig, pubB64)
	if validTampered {
		t.Error("SECURITY FAIL: Tampered array data should not be valid")
	}

	fmt.Println("[DEBUG] Array & Recursive check: SUCCESS")
}

func TestOSM15_DomainIsolation(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(nil)
	privB64 := base64.StdEncoding.EncodeToString(priv.Seed())
	pubB64 := base64.StdEncoding.EncodeToString(pub)

	types := map[string][]TypedMember{
		"Mail": {{Name: "content", Type: "string"}},
	}

	msg := map[string]interface{}{"content": "Octra is awesome"}

	dataA := TypedData{
		Domain: TypedDomain{Name: "AppA", Version: "1", ChainID: 1},
		Types: types, PrimaryType: "Mail", Message: msg,
	}
	dataB := TypedData{
		Domain: TypedDomain{Name: "AppB", Version: "1", ChainID: 1},
		Types: types, PrimaryType: "Mail", Message: msg,
	}

	digestA, _ := HashTypedData(dataA)
	digestB, _ := HashTypedData(dataB)

	fmt.Printf("[DEBUG] Digest AppA: %s\n", hex.EncodeToString(digestA))
	fmt.Printf("[DEBUG] Digest AppB: %s\n", hex.EncodeToString(digestB))

	if hex.EncodeToString(digestA) == hex.EncodeToString(digestB) {
		t.Fatal("Domain isolation failed: Digests are identical")
	}

	sig, _ := SignTypedData(dataA, privB64)
	valid, _ := VerifyTypedData(dataB, sig, pubB64)
	if valid {
		t.Error("Security Breach: Cross-domain signature accepted")
	}
}

func TestOSM15_JSONFlowAndSigner(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(nil)
	privB64 := base64.StdEncoding.EncodeToString(priv.Seed())
	pubB64 := base64.StdEncoding.EncodeToString(pub)
	expectedAddr := PublicKeyToAddress(pub)

	data := TypedData{
		Domain: TypedDomain{Name: "OctraNetwork", Version: "1", ChainID: 1},
		Types: map[string][]TypedMember{
			"Msg": {{Name: "text", Type: "string"}},
		},
		PrimaryType: "Msg",
		Message:     map[string]interface{}{"text": "Hello Octra"},
	}

	// 1. Sign
	sig, _ := SignTypedData(data, privB64)

	// 2. Export to JSON
	jsonBytes, err := ExportToJSON(data, sig)
	if err != nil { t.Fatalf("Export error: %v", err) }
	fmt.Printf("\n[DEBUG] Exported JSON:\n%s\n", string(jsonBytes))

	// 3. Verify from JSON
	valid, _ := VerifyFromJSON(jsonBytes, pubB64)
	if !valid { t.Error("Failed to verify from JSON") }

	// 4. Get Signer Address
	signerAddr, err := GetSignerAddress(data, sig, pubB64)
	if err != nil || signerAddr != expectedAddr {
		t.Errorf("Signer address mismatch! Got: %s, Expected: %s", signerAddr, expectedAddr)
	}
	fmt.Printf("[DEBUG] Signer Address Verified: %s\n", signerAddr)
}

func TestOSM15_KeystoreFlow(t *testing.T) {
	privKey := "OQJS215Wyy0PbMpw1M3hOi7LucCeJLI6AX8Jx484mxc="
	password := "octra-secure-123"

	// 1. Test Encrypt
	keystoreJSON, err := EncryptKey(privKey, password)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}
	fmt.Printf("\n[DEBUG] Keystore Generated:\n%s\n", string(keystoreJSON))

	// 2. Test Decrypt Success
	decryptedKey, err := DecryptKey(keystoreJSON, password)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	if decryptedKey != privKey {
		t.Errorf("Mismatch! Expected %s, got %s", privKey, decryptedKey)
	}

	// 3. Test Decrypt Wrong Password
	_, err = DecryptKey(keystoreJSON, "password-salah")
	if err == nil {
		t.Error("Security Breach: Decryption should fail with wrong password")
	}

	fmt.Println("[DEBUG] Keystore Security: PASSED")
}
