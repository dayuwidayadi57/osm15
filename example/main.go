package main

import (
"encoding/base64"
"fmt"
"log"
"github.com/dayuwidayadi57/osm15""
)

func main() {
fmt.Println("=== Octra OSM-15 Toolchain Demo ===")

// 1. Generate Keypair
privB64, pubB64, err := osm15.GenerateKeypair()
if err != nil {
log.Fatal(err)
}

pubBytes, _ := base64.StdEncoding.DecodeString(pubB64)
address := osm15.PublicKeyToAddress(pubBytes)

fmt.Printf("[+] New Account: %s\n", address)

// 2. Data
data := osm15.TypedData{
Domain: osm15.TypedDomain{Name: "OctraPay", Version: "1", ChainID: 1},
Types: map[string][]osm15.TypedMember{
"Transaction": {
{Name: "from", Type: "string"},
{Name: "to", Type: "string"},
{Name: "amount", Type: "uint256"},
},
},
PrimaryType: "Transaction",
Message: map[string]interface{}{
"from":   address,
"to":     "oct1exampledest...",
"amount": 1000,
},
}

// 3. Sign
signature, err := osm15.SignTypedData(data, privB64)
if err != nil {
log.Fatal(err)
}

// 4. Verify
recoveredAddr, err := osm15.GetSignerAddress(data, signature, pubB64)
if err != nil {
log.Fatal(err)
}

if recoveredAddr == address {
fmt.Println("[+] Verification: SUCCESS")
} else {
fmt.Println("[-] Verification: FAILED")
}

// 5. Final Export (Print biar gak error unused)
jsonPayload, _ := osm15.ExportToJSON(data, signature)
fmt.Printf("\n[+] Payload JSON:\n%s\n", string(jsonPayload))
}
