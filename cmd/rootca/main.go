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
		action      = flag.String("action", "", "Action to perform: generate, sign, inspect")
		cn          = flag.String("cn", "", "Common Name for the certificate")
		outputDir   = flag.String("output", "./certs", "Output directory for certificates")
		validYears  = flag.Int("validity", 20, "Validity period in years")
		rootCert    = flag.String("root-cert", "", "Path to root CA certificate (for sign action)")
		rootKey     = flag.String("root-key", "", "Path to root CA private key (for sign action)")
		csr         = flag.String("csr", "", "Path to CSR file (for sign action)")
		certPath    = flag.String("cert", "", "Path to certificate file (for inspect action)")
	)

	flag.Parse()

	if *action == "" {
		fmt.Println("GigVault Root CA Tool - Offline Certificate Authority")
		fmt.Println("======================================================")
		fmt.Println("")
		fmt.Println("Usage: rootca -action <generate|sign|inspect> [options]")
		fmt.Println("")
		fmt.Println("Actions:")
		fmt.Println("  generate  Generate a new root CA certificate (air-gapped)")
		fmt.Println("  sign      Sign an intermediate CA CSR with the root CA")
		fmt.Println("  inspect   Display certificate information")
		fmt.Println("")
		fmt.Println("Options:")
		flag.PrintDefaults()
		fmt.Println("")
		fmt.Println("Examples:")
		fmt.Println("  # Generate a new root CA")
		fmt.Println("  rootca -action generate -cn \"GigVault Root CA\" -output ./root-ca")
		fmt.Println("")
		fmt.Println("  # Sign an intermediate CA CSR")
		fmt.Println("  rootca -action sign -root-cert ./root-ca/root-ca.crt \\")
		fmt.Println("    -root-key ./root-ca/root-ca.key -csr ./intermediate.csr \\")
		fmt.Println("    -output ./intermediate-ca -validity 10")
		fmt.Println("")
		fmt.Println("  # Inspect a certificate")
		fmt.Println("  rootca -action inspect -cert ./root-ca/root-ca.crt")
		fmt.Println("")
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
		fmt.Printf("\n✅ Root CA generated successfully in %s\n", *outputDir)

	case "sign":
		if *rootCert == "" || *rootKey == "" || *csr == "" {
			log.Fatal("Root certificate (-root-cert), root key (-root-key), and CSR (-csr) are required for sign action")
		}
		if err := generator.SignIntermediateCA(*rootCert, *rootKey, *csr, *outputDir, *validYears); err != nil {
			log.Fatalf("Failed to sign intermediate CA: %v", err)
		}
		fmt.Printf("\n✅ Intermediate CA signed successfully in %s\n", *outputDir)

	case "inspect":
		if *certPath == "" {
			log.Fatal("Certificate path (-cert) is required for inspect action")
		}
		if err := generator.InspectCertificate(*certPath); err != nil {
			log.Fatalf("Failed to inspect certificate: %v", err)
		}

	default:
		log.Fatalf("Unknown action: %s", *action)
	}
}
