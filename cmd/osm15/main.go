package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"tunnel/osm15"

	"github.com/fsnotify/fsnotify"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: osm15 <command> [<args>]")
		fmt.Println("Commands: generate, sign, batch-sign, watch-sign, encrypt, decrypt")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "generate":
		pub, priv, _ := ed25519.GenerateKey(nil)
		fmt.Printf("Private Key (Base64): %s\n", base64.StdEncoding.EncodeToString(priv.Seed()))
		fmt.Printf("Public Key (Base64):  %s\n", base64.StdEncoding.EncodeToString(pub))
		fmt.Printf("Address:              %s\n", osm15.PublicKeyToAddress(pub))

	case "sign":
		signCmd := flag.NewFlagSet("sign", flag.ExitOnError)
		signFile := signCmd.String("file", "", "TypedData JSON file")
		walletFile := signCmd.String("wallet", "", "Keystore file path")
		password := signCmd.String("pass", "", "Keystore password")
		signCmd.Parse(os.Args[2:])

		if *signFile == "" || *walletFile == "" || *password == "" {
			fmt.Println("Usage: sign -file <data.json> -wallet <wallet.json> -pass <password>")
			os.Exit(1)
		}

		ksData, _ := ioutil.ReadFile(*walletFile)
		privKeyB64, err := osm15.DecryptKey(ksData, *password)
		if err != nil {
			fmt.Println("Error: Invalid password")
			os.Exit(1)
		}

		fileData, _ := ioutil.ReadFile(*signFile)
		var typedData osm15.TypedData
		json.Unmarshal(fileData, &typedData)

		sig, _ := osm15.SignTypedData(typedData, privKeyB64)
		output, _ := osm15.ExportToJSON(typedData, sig)
		fmt.Println(string(output))

	case "batch-sign":
		batchCmd := flag.NewFlagSet("batch-sign", flag.ExitOnError)
		inDir := batchCmd.String("in", "", "Input directory")
		outDir := batchCmd.String("out", "", "Output directory")
		walletFile := batchCmd.String("wallet", "", "Keystore file")
		password := batchCmd.String("pass", "", "Password")
		batchCmd.Parse(os.Args[2:])

		ksData, _ := ioutil.ReadFile(*walletFile)
		privKey, err := osm15.DecryptKey(ksData, *password)
		if err != nil {
			fmt.Println("Error: Invalid password")
			os.Exit(1)
		}

		files, _ := ioutil.ReadDir(*inDir)
		os.MkdirAll(*outDir, 0755)

		for _, f := range files {
			if !f.IsDir() && strings.HasSuffix(f.Name(), ".json") {
				processFile(*inDir+"/"+f.Name(), *outDir, privKey)
			}
		}
		fmt.Println("Batch signing completed!")

	case "watch-sign":
		watchCmd := flag.NewFlagSet("watch-sign", flag.ExitOnError)
		inDir := watchCmd.String("in", "pending_tx", "Input directory")
		outDir := watchCmd.String("out", "signed_tx", "Output directory")
		walletFile := watchCmd.String("wallet", "", "Keystore file")
		password := watchCmd.String("pass", "", "Password")
		watchCmd.Parse(os.Args[2:])

		if *walletFile == "" || *password == "" {
			fmt.Println("Usage: watch-sign -in <dir> -out <dir> -wallet <ks.json> -pass <pw>")
			os.Exit(1)
		}

		ksData, _ := ioutil.ReadFile(*walletFile)
		privKey, err := osm15.DecryptKey(ksData, *password)
		if err != nil {
			fmt.Println("Error: Invalid password")
			os.Exit(1)
		}

		os.MkdirAll(*inDir, 0755)
		os.MkdirAll(*outDir, 0755)

		watcher, _ := fsnotify.NewWatcher()
		defer watcher.Close()

		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok { return }
					if event.Op&fsnotify.Create == fsnotify.Create && strings.HasSuffix(event.Name, ".json") {
						fmt.Printf("Detected: %s\n", filepath.Base(event.Name))
						processFile(event.Name, *outDir, privKey)
					}
				case err, ok := <-watcher.Errors:
					if !ok { return }
					fmt.Println("Watcher error:", err)
				}
			}
		}()

		watcher.Add(*inDir)
		fmt.Printf("Watcher active on ./%s. Press Ctrl+C to stop.\n", *inDir)
		select {}

	case "encrypt":
		encCmd := flag.NewFlagSet("encrypt", flag.ExitOnError)
		key := encCmd.String("key", "", "Private key (Base64)")
		pass := encCmd.String("pass", "", "Password")
		encCmd.Parse(os.Args[2:])

		ks, _ := osm15.EncryptKey(*key, *pass)
		fmt.Println(string(ks))

	case "decrypt":
		decCmd := flag.NewFlagSet("decrypt", flag.ExitOnError)
		file := decCmd.String("file", "", "Keystore file")
		pass := decCmd.String("pass", "", "Password")
		decCmd.Parse(os.Args[2:])

		data, _ := ioutil.ReadFile(*file)
		key, err := osm15.DecryptKey(data, *pass)
		if err != nil {
			fmt.Println("Error: Invalid password")
			os.Exit(1)
		}
		fmt.Printf("Decrypted Private Key: %s\n", key)

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func processFile(filePath, outDir, privKey string) {
	fileData, _ := ioutil.ReadFile(filePath)
	var typedData osm15.TypedData
	if err := json.Unmarshal(fileData, &typedData); err != nil {
		fmt.Printf("Skip %s: Invalid format\n", filePath)
		return
	}

	sig, _ := osm15.SignTypedData(typedData, privKey)
	output, _ := osm15.ExportToJSON(typedData, sig)
	
	outPath := filepath.Join(outDir, "signed_"+filepath.Base(filePath))
	ioutil.WriteFile(outPath, output, 0644)
	fmt.Printf("Signed and saved to %s\n", outPath)
}
