# Octra Sign Message Standards - OSM-15

[![Network](https://img.shields.io/badge/Network-Octra-blueviolet)](https://octra.network)
[![Protocol](https://img.shields.io/badge/Protocol-OSM--15-success)](https://github.com/dayuwidayadi57/osm15)
[![Team](https://img.shields.io/badge/Developed%20By-QiubitLabs-blue)](https://github.com/QiubitLabs)
[![Coverage](https://img.shields.io/badge/Coverage-90.9%25-brightgreen)](https://github.com/dayuwidayadi57/osm15)
[![Status](https://img.shields.io/badge/Status-Stable-blue)](https://github.com/dayuwidayadi57/osm15)

A high-performance Go library implementing the **OSM-15** (Octra Structured Message) standard. Designed for secure, recursive, and domain-isolated data signing in the Octra ecosystem.

## 🚀 Key Functions for Developers

### 1. Identity & Address
Convert Ed25519 public keys into the official Octra address format.
```go
// Convert public key bytes to "oct..." address
address := osm15.PublicKeyToAddress(publicKeyBytes)
```

### 2. Signing Structured Data
Sign complex messages with recursive type support.
```go
data := osm15.TypedData{
    Domain: osm15.TypedDomain{Name: "OctraPay", Version: "1", ChainID: 1},
    Types: map[string][]osm15.TypedMember{
        "Transaction": {
            {Name: "amount", Type: "uint256"},
            {Name: "to", Type: "address"},
        },
    },
    PrimaryType: "Transaction",
    Message: map[string]interface{}{
        "amount": 5000,
        "to": "oct1abc...",
    },
}

// Returns Base64 signature
signature, err := osm15.SignTypedData(data, privateKeyBase64)
```

### 3. Verification & Signer Recovery
Verify if a signature is valid and recover the signer's address.
```go
// Boolean verification
isValid, err := osm15.VerifyTypedData(data, signature, publicKeyBase64)

// Get signer's Octra address directly from signature
address, err := osm15.GetSignerAddress(data, signature, publicKeyBase64)
```

### 4. JSON Export & Import
Easily export the signed payload to JSON or verify directly from a JSON file.
```go
// Export to formatted JSON
jsonBytes, err := osm15.ExportToJSON(data, signature)

// Verify directly from JSON payload
isValid, err := osm15.VerifyFromJSON(jsonBytes, publicKeyBase64)
```

### 5. Debugging Signing Text
Get the raw string that is being hashed and signed (for debugging or hardware wallet display).
```go
rawText, err := osm15.GetSigningText(data)
fmt.Println(rawText)
```

## ⚙️ Technical Specifications
- **Hashing**: Recursive SHA-256 (EIP-712 Style)
- **Signature**: Ed25519
- **Prefix**: \x19Octra Typed Data:
- **Encoding**: Base64 (Signature & Keys) / Base58 (Address)

## 🛠 Installation
```bash
go get [github.com/dayuwidayadi57/osm15](https://github.com/dayuwidayadi57/osm15)
```

## 💡 Why OSM-15?
If you are familiar with Ethereum's **EIP-712**, you will feel at home. OSM-15 brings:
- **Human-readable signing**: Users see exactly what they sign.
- **EIP-712 Style**: Familiar domain separation and recursive hashing.
- **Enhanced Security**: Optimized for Ed25519 and SHA-256 (Octra Native).
cat <<EOF > README.md
# OSM15

##  Octra Go SDK Standardized Protocol
Professional Open Source Contributor Mode.

## Versioning
- Version: v0.2.2
- Status: Stable
- Date: 2026-01-19

## 📜 License
(c) 2026 QiubitLabs Team. Licensed under the MIT License.
