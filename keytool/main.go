package main

import (
	"flag"
	"fmt"
	rsa "github.com/wunderkind2k1/auth-server-task/pkg/keys/internal"
	"log/slog"
	"os"
	"path/filepath"
)

func main() {
	// Set up logging
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Define commands
	generateCmd := flag.NewFlagSet("generate", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)

	// Define flags
	keysDir := flag.String("dir", "keys", "Directory to store keys")
	bits := generateCmd.Int("bits", 2048, "Number of bits for RSA key pair")
	keyID := deleteCmd.String("id", "", "Key ID to delete")

	// Parse command
	if len(os.Args) < 2 {
		fmt.Println("expected 'generate', 'list', or 'delete' command")
		os.Exit(1)
	}

	// Create key manager
	manager, err := rsa.NewManager(*keysDir)
	if err != nil {
		slog.Error("Failed to create key manager", "error", err)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "generate":
		err := generateCmd.Parse(os.Args[2:])
		if err != nil {
			slog.Error("Failed to parse args for keypair to generate", "error", err)
			os.Exit(1)
		}
		if err := handleGenerate(manager, *bits); err != nil {
			slog.Error("Failed to generate key pair", "error", err)
			os.Exit(1)
		}

	case "list":
		err := listCmd.Parse(os.Args[2:])
		if err != nil {
			slog.Error("Failed to parse key pair to list", "error", err)
			os.Exit(1)
		}
		if err := handleList(manager); err != nil {
			slog.Error("Failed to list key pairs", "error", err)
			os.Exit(1)
		}

	case "delete":
		err := deleteCmd.Parse(os.Args[2:])
		if err != nil {
			slog.Error("Failed to parse key pair to delete", "error", err)
			os.Exit(1)
		}
		if *keyID == "" {
			fmt.Println("error: -id flag is required")
			deleteCmd.PrintDefaults()
			os.Exit(1)
		}
		if err := handleDelete(manager, *keyID); err != nil {
			slog.Error("Failed to delete key pair", "error", err)
			os.Exit(1)
		}

	default:
		fmt.Printf("unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func handleGenerate(manager *rsa.Manager, bits int) error {
	// Generate new key pair
	keyPair, err := manager.GenerateKeyPair(bits)
	if err != nil {
		return fmt.Errorf("failed to generate key pair: %w", err)
	}

	// Save the key pair
	if err := manager.SaveKeyPair(keyPair); err != nil {
		return fmt.Errorf("failed to save key pair: %w", err)
	}

	slog.Info("Generated new key pair",
		"keyID", keyPair.KeyID,
		"bits", bits,
		"privateKey", filepath.Join(manager.KeysDir(), fmt.Sprintf("%s.private.pem", keyPair.KeyID)),
		"publicKey", filepath.Join(manager.KeysDir(), fmt.Sprintf("%s.public.pem", keyPair.KeyID)),
	)

	return nil
}

func handleList(manager *rsa.Manager) error {
	keyIDs, err := manager.ListKeyPairs()
	if err != nil {
		return fmt.Errorf("failed to list key pairs: %w", err)
	}

	if len(keyIDs) == 0 {
		slog.Info("No key pairs found")
		return nil
	}

	slog.Info("Available key pairs", "count", len(keyIDs))
	for _, keyID := range keyIDs {
		slog.Info("Key pair",
			"keyID", keyID,
			"privateKey", filepath.Join(manager.KeysDir(), fmt.Sprintf("%s.private.pem", keyID)),
			"publicKey", filepath.Join(manager.KeysDir(), fmt.Sprintf("%s.public.pem", keyID)),
		)
	}

	return nil
}

func handleDelete(manager *rsa.Manager, keyID string) error {
	if err := manager.DeleteKeyPair(keyID); err != nil {
		return fmt.Errorf("failed to delete key pair: %w", err)
	}
	slog.Info("Deleted key pair", "keyID", keyID)
	return nil
}
