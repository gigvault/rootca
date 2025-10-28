package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gigvault/rootca/internal/generator"
)

func main() {
	var (
		action     = flag.String("action", "", "Action to perform: generate, sign, revoke")
		cn         = flag.String("cn", "", "Common Name for the certificate")
		outputDir  = flag.String("output", "./certs", "Output directory for certificates")
		validYears = flag.Int("validity", 20, "Validity period in years")
	)

	flag.Parse()

	if *action == "" {
		fmt.Println("GigVault Root CA Tool")
		fmt.Println("Usage: rootca -action <generate|sign|revoke> [options]")
		fmt.Println("")
		fmt.Println("Actions:")
		fmt.Println("  generate  Generate a new root CA certificate (air-gapped)")
		fmt.Println("  sign      Sign an intermediate CA certificate")
		fmt.Println("  revoke    Revoke a certificate")
		fmt.Println("")
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	switch *action {
	case "generate":
		if *cn == "" {
			log.Fatal("Common Name (-cn) is required for generate action")
		}
		if err := generator.GenerateRootCA(*cn, *outputDir, *validYears); err != nil {
			log.Fatalf("Failed to generate root CA: %v", err)
		}
		fmt.Printf("Root CA generated successfully in %s\n", *outputDir)

	case "sign":
		fmt.Println("Sign intermediate CA (TODO: implement)")
		// TODO: Implement intermediate CA signing

	case "revoke":
		fmt.Println("Revoke certificate (TODO: implement)")
		// TODO: Implement certificate revocation

	default:
		log.Fatalf("Unknown action: %s", *action)
	}
}
