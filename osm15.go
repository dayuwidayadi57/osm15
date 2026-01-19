// Package osm15 implements the Octra Structured Message Standard (OSM-15).
//
// OSM-15 is a high-performance signing standard for the Octra Network,
// featuring recursive struct hashing, domain separation, and Ed25519 
// signature schemes. It is designed to be human-readable and 
// cryptographically secure, similar to Ethereum's EIP-712.
//
// Specification: https://github.com/dayuwidayadi57/osm15
// Reference Implementation: Go 1.19+
//
// (c) 2026 QiubitLabs Team. All rights reserved.
// Licensed under the MIT License.

package osm15

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/mr-tron/base58"
)

type TypedMember struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type TypedDomain struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	ChainID int    `json:"chainId"`
}

type TypedData struct {
	Domain      TypedDomain              `json:"domain"`
	Types       map[string][]TypedMember `json:"types"`
	PrimaryType string                   `json:"primaryType"`
	Message     map[string]interface{}   `json:"message"`
}

func (d TypedDomain) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":    d.Name,
		"version": d.Version,
		"chainId": d.ChainID,
	}
}

type SignedPayload struct {
	Data      TypedData `json:"data"`
	Signature string    `json:"signature"`
}

func HashTypedData(data TypedData) ([]byte, error) {
	if data.Types == nil {
		data.Types = make(map[string][]TypedMember)
	}
	
<<<<<<< HEAD
=======
	// WAJIB: Daftarkan tipe TypedDomain agar hashStruct bisa memproses field-nya
>>>>>>> 2ac639c (Update via GoSmartPush v17.1 - 2026-01-19 14:47:02)
	data.Types["TypedDomain"] = []TypedMember{
		{Name: "name", Type: "string"},
		{Name: "version", Type: "string"},
		{Name: "chainId", Type: "uint256"},
	}

	domainHash, _ := hashStruct("TypedDomain", data.Domain.ToMap(), data.Types)
	messageHash, _ := hashStruct(data.PrimaryType, data.Message, data.Types)

	// Data binary yang akan di-sign (gabungan domain & message hash)
	payloadBinary := append(domainHash, messageHash...)
	
	const TypedPrefix = "\x19Octra Typed Data:"
	
	//  (\nLength\n)
	signingBody := fmt.Sprintf("%s\n%d\n%s", 
		TypedPrefix, 
		len(payloadBinary), 
		string(payloadBinary),
	)

<<<<<<< HEAD
=======
	payloadBinary := append(domainHash, messageHash...)
	
	const TypedPrefix = "\x19Octra Typed Data:\n"
	
	// Mengikuti standar OSM-1: Prefix + Length + \n + BinaryData
	signingBody := fmt.Sprintf("%s%d\n%s", 
		TypedPrefix, 
		len(payloadBinary), 
		string(payloadBinary),
	)

>>>>>>> 2ac639c (Update via GoSmartPush v17.1 - 2026-01-19 14:47:02)
	hash := sha256.Sum256([]byte(signingBody))
	return hash[:], nil
}

func hashStruct(typeName string, data map[string]interface{}, types map[string][]TypedMember) ([]byte, error) {
	typeString := encodeType(typeName, types)
	typeHash := sha256.Sum256([]byte(typeString))
	
	var buf bytes.Buffer
	buf.Write(typeHash[:])

	members := types[typeName]
	for _, member := range members {
		val := data[member.Name]
		encodedVal, err := encodeValue(member.Type, val, types)
		if err != nil { return nil, err }
		buf.Write(encodedVal)
	}

	finalHash := sha256.Sum256(buf.Bytes())
	return finalHash[:], nil
}

func encodeValue(typeName string, value interface{}, types map[string][]TypedMember) ([]byte, error) {
	if strings.HasSuffix(typeName, "[]") {
		baseType := typeName[:len(typeName)-2]
		rv := reflect.ValueOf(value)
		if rv.Kind() != reflect.Slice {
			return nil, fmt.Errorf("expected slice for type %s", typeName)
		}
		var buf bytes.Buffer
		for i := 0; i < rv.Len(); i++ {
			encoded, err := encodeValue(baseType, rv.Index(i).Interface(), types)
			if err != nil { return nil, err }
			buf.Write(encoded)
		}
		hash := sha256.Sum256(buf.Bytes())
		return hash[:], nil
	}

	if _, ok := types[typeName]; ok {
		mapVal, _ := value.(map[string]interface{})
		return hashStruct(typeName, mapVal, types)
	}

	var data []byte
	switch typeName {
	case "string":
		s, _ := value.(string)
		h := sha256.Sum256([]byte(s))
		return h[:], nil
	case "address":
		s, _ := value.(string)
		data = []byte(s)
	case "uint256", "int", "chainId":
		data = []byte(fmt.Sprintf("%v", value))
	default:
		data, _ = json.Marshal(value)
	}
	h := sha256.Sum256(data)
	return h[:], nil
}

func encodeType(primaryType string, types map[string][]TypedMember) string {
	unsortedDeps := findDependencies(primaryType, types, make(map[string]bool))
	deps := make([]string, 0, len(unsortedDeps))
	for dep := range unsortedDeps {
		if dep != primaryType { deps = append(deps, dep) }
	}
	sort.Strings(deps)
	var sb strings.Builder
	sb.WriteString(formatMainType(primaryType, types))
	for _, dep := range deps {
		sb.WriteString(formatMainType(dep, types))
	}
	return sb.String()
}

func findDependencies(primaryType string, types map[string][]TypedMember, found map[string]bool) map[string]bool {
	if found[primaryType] { return found }
	if _, ok := types[primaryType]; !ok { return found }
	found[primaryType] = true
	for _, member := range types[primaryType] {
		cleanType := strings.TrimSuffix(member.Type, "[]")
		findDependencies(cleanType, types, found)
	}
	return found
}

func formatMainType(name string, types map[string][]TypedMember) string {
	var sb strings.Builder
	sb.WriteString(name)
	sb.WriteString("(")
	for i, member := range types[name] {
		if i > 0 { sb.WriteString(",") }
		sb.WriteString(fmt.Sprintf("%s %s", member.Type, member.Name))
	}
	sb.WriteString(")")
	return sb.String()
}

func SignTypedData(data TypedData, privateKeyB64 string) (string, error) {
	digest, err := HashTypedData(data)
	if err != nil { return "", err }
	seed, _ := base64.StdEncoding.DecodeString(privateKeyB64)
	priv := ed25519.NewKeyFromSeed(seed)
	sig := ed25519.Sign(priv, digest)
	return base64.StdEncoding.EncodeToString(sig), nil
}

func VerifyTypedData(data TypedData, signatureB64 string, publicKeyB64 string) (bool, error) {
	digest, err := HashTypedData(data)
	if err != nil { return false, err }
	sig, _ := base64.StdEncoding.DecodeString(signatureB64)
	pk, _ := base64.StdEncoding.DecodeString(publicKeyB64)
	return ed25519.Verify(pk, digest, sig), nil
}

func PublicKeyToAddress(publicKey []byte) string {
	hash := sha256.Sum256(publicKey)
	return "oct" + base58.Encode(hash[:])
}

func ExportToJSON(data TypedData, signature string) ([]byte, error) {
	payload := SignedPayload{
		Data:      data,
		Signature: signature,
	}
	return json.MarshalIndent(payload, "", "  ")
}

func VerifyFromJSON(payloadJSON []byte, publicKeyB64 string) (bool, error) {
	var payload SignedPayload
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		return false, err
	}
	return VerifyTypedData(payload.Data, payload.Signature, publicKeyB64)
}

func GetSignerAddress(data TypedData, signatureB64 string, publicKeyB64 string) (string, error) {
	valid, err := VerifyTypedData(data, signatureB64, publicKeyB64)
	if err != nil { return "", err }
	if !valid { return "", fmt.Errorf("invalid signature") }
	
	pk, _ := base64.StdEncoding.DecodeString(publicKeyB64)
	return PublicKeyToAddress(pk), nil
}

func GetSigningText(data TypedData) (string, error) {
	domainHash, _ := hashStruct("TypedDomain", data.Domain.ToMap(), data.Types)
	messageHash, _ := hashStruct(data.PrimaryType, data.Message, data.Types)
	payloadBinary := append(domainHash, messageHash...)
	
	const TypedPrefix = "\x19Octra Typed Data:\n"
	return fmt.Sprintf("%s%d\n%s", TypedPrefix, len(payloadBinary), string(payloadBinary)), nil
}

<<<<<<< HEAD
=======
func GenerateKeypair() (string, string, error) {
    pub, priv, err := ed25519.GenerateKey(nil)
    if err != nil {
        return "", "", err
    }
    
    privBase64 := base64.StdEncoding.EncodeToString(priv.Seed())
    pubBase64 := base64.StdEncoding.EncodeToString(pub)
    
    return privBase64, pubBase64, nil
}
>>>>>>> 2ac639c (Update via GoSmartPush v17.1 - 2026-01-19 14:47:02)
