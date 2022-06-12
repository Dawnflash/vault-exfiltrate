package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Dawnflash/vault-exfiltrate/lib"
)

func usage() {
	os.Stderr.WriteString("usage: vault-exfiltrate extract vault_pid keyring_file\n")
	os.Stderr.WriteString("or:    vault-exfiltrate extract-core core_file keyring_file\n")
	os.Stderr.WriteString("or:    vault-exfiltrate decrypt keyring.json path/In/Vault data_file\n")
	os.Stderr.WriteString("or:    vault-exfiltrate split base64_encoded_value num_shares\n")
	os.Stderr.WriteString("or:    vault-exfiltrate combine shares...\n")
}

func main_() int {
	if len(os.Args) == 1 {
		usage()
		return 1
	}

	subcommand := strings.ToLower(os.Args[1])
	args := os.Args[2:]

	switch subcommand {
	case "extract-core":
		if len(args) < 2 {
			usage()
			return 1
		}
		keyring, err := lib.FindMasterKeyInCore(args[0], args[1])
		if err != nil {
			panic(err)
		}
		os.Stdout.Write(keyring)
		return 0
	case "extract":
		if len(args) < 2 {
			usage()
			return 1
		}
		keyring, err := lib.FindMasterKeyLive(args[0], args[1])
		if err != nil {
			panic(err)
		}
		os.Stdout.Write(keyring)
		return 0
	case "decrypt":
		if len(args) < 3 {
			usage()
			return 1
		}
		plaintext, err := lib.DecryptFile(args[0], args[1], args[2])
		if err != nil {
			panic(err)
		}
		os.Stdout.Write(plaintext)
		return 0
	case "split":
		if len(args) < 2 {
			usage()
			return 1
		}
		shares, err := lib.SecretShares(args[0], args[1])
		if err != nil {
			panic(err)
		}
		for _, share := range shares {
			fmt.Println(share)
		}
		return 0
	case "combine":
		combined, err := lib.CombineShares(args)
		if err != nil {
			panic(err)
		}
		fmt.Println(combined)
		return 0
	}

	usage()
	return 1
}

func main() {
	os.Exit(main_())
}
