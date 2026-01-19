# Octra OSM-15 Toolchain

A high-performance Go library and CLI for the Octra ecosystem, implementing the **OSM-15** (Octra Structured Message) standard. This tool enables identity creation, secure keystore management, and automated structured data signing.

## Core Features

* **Identity Generation**: Create Ed25519 keypairs and derive Octra addresses (`oct...`).
* **Secure Keystore**: Encrypt private keys using AES-256-GCM and Scrypt (industry standard).
* **OSM-15 Signing**: EIP-712 style structured data hashing with array and recursive type support.
* **Batch & Watcher**: High-throughput processing for manual or automated signing workflows.

---

## Developer Guide: Using the Library

If you are developing a Go application and want to integrate OSM-15, use the following patterns:

### 1. Signing Structured Data
Define your message structure and sign it using a private key.

```go
import "tunnel/osm15"

// 1. Define Data
data := osm15.TypedData{
    Domain: osm15.TypedDomain{
        Name: "OctraPay",
        Version: "1",
        ChainID: 1,
    },
    Types: map[string][]osm15.TypedMember{
        "Transfer": {
            {Name: "to", Type: "string"},
            {Name: "amount", Type: "uint256"},
        },
    },
    PrimaryType: "Transfer",
    Message: map[string]interface{}{
        "to": "octGzrStfZE5ae...",
        "amount": 1000,
    },
}

// 2. Sign (using Base64 private key seed)
signature, err := osm15.SignTypedData(data, "your-private-key-base64")

// 3. Export to JSON for Network Transmission
jsonPayload, _ := osm15.ExportToJSON(data, signature)
```

### 2. Wallet & Keystore Management
Securely encrypt and decrypt keys within your application.

```go
// Encrypt a key to Scrypt-protected JSON
keystoreJSON, err := osm15.EncryptKey(privateKeyB64, "strong-password")

// Decrypt a key from Keystore JSON
privateKey, err := osm15.DecryptKey(keystoreJSON, "strong-password")
```

### 3. Verification
Verify signatures received from other nodes or users.

```go
isValid, err := osm15.VerifyTypedData(data, signature, publicKeyB64)
// Or verify directly from JSON payload
isValid, err := osm15.VerifyFromJSON(payloadJSON, publicKeyB64)
```

---

## Technical Specifications

| Component | Technology | Description |
| :--- | :--- | :--- |
| **Signature** | Ed25519 | Edwards-curve Digital Signature Algorithm |
| **Hash Standard** | SHA3-256 | High-security cryptographic hashing |
| **KDF** | Scrypt | N=32768, R=8, P=1 (Memory-hard against ASIC) |
| **Encryption** | AES-256-GCM | Authenticated encryption with associated data |
| **Address** | Base58 | Octra Address format with `oct` prefix |



---

## Installation & Build

```bash
# Install dependencies
go get github.com/fsnotify/fsnotify
go get github.com/mr-tron/base58
go get golang.org/x/crypto/scrypt

# Build CLI tool
go build -o osm15-tool ./cmd/osm15/main.go
```

## CLI Usage

| Command | Usage |
| :--- | :--- |
| **generate** | `./osm15-tool generate` |
| **encrypt** | `./osm15-tool encrypt -key <KEY> -pass <PW> > wallet.json` |
| **watch-sign**| `./osm15-tool watch-sign -wallet wallet.json -pass <PW>` |
| **batch-sign**| `./osm15-tool batch-sign -in <DIR> -out <DIR> -wallet <KS> -pass <PW>` |

---

## License
Proprietary / Octra Network
