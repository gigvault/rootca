package generator

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gigvault/shared/pkg/crypto"
)

// GenerateRootCA generates a new root CA certificate and private key
func GenerateRootCA(cn, outputDir string, validityYears int) error {
	// Create output directory
	if err := os.MkdirAll(outputDir, 0700); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate ECDSA P-384 key for root CA (higher security)
	privateKey, err := crypto.GenerateP384Key()
	if err != nil {
		return fmt.Errorf("failed to generate private key: %w", err)
	}

	// Generate serial number
	serial, err := crypto.GenerateSerialNumber()
	if err != nil {
		return fmt.Errorf("failed to generate serial number: %w", err)
	}

	// Create certificate template
	notBefore := time.Now()
	notAfter := notBefore.AddDate(validityYears, 0, 0)

	template := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:   cn,
			Organization: []string{"GigVault"},
			Country:      []string{"US"},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            1,
		MaxPathLenZero:        false,
	}

	// Self-sign the certificate
	certPEM, err := crypto.SignCertificate(template, template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return fmt.Errorf("failed to sign certificate: %w", err)
	}

	// Encode private key to PEM
	keyPEM, err := crypto.EncodePrivateKeyToPEM(privateKey)
	if err != nil {
		return fmt.Errorf("failed to encode private key: %w", err)
	}

	// Write certificate to file
	certPath := filepath.Join(outputDir, "root-ca.crt")
	if err := os.WriteFile(certPath, certPEM, 0644); err != nil {
		return fmt.Errorf("failed to write certificate: %w", err)
	}

	// Write private key to file (restricted permissions)
	keyPath := filepath.Join(outputDir, "root-ca.key")
	if err := os.WriteFile(keyPath, keyPEM, 0600); err != nil {
		return fmt.Errorf("failed to write private key: %w", err)
	}

	fmt.Printf("Root CA certificate: %s\n", certPath)
	fmt.Printf("Root CA private key: %s\n", keyPath)
	fmt.Printf("Valid from %s to %s\n", notBefore.Format(time.RFC3339), notAfter.Format(time.RFC3339))
	fmt.Printf("\n⚠️  IMPORTANT: Store the private key in a secure, air-gapped location!\n")

	return nil
}
