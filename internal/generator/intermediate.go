package generator

import (
	"crypto/x509"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gigvault/shared/pkg/crypto"
)

// SignIntermediateCA signs an intermediate CA certificate with the root CA
func SignIntermediateCA(rootCertPath, rootKeyPath, csrPath, outputDir string, validityYears int) error {
	// Read root CA certificate
	rootCertPEM, err := os.ReadFile(rootCertPath)
	if err != nil {
		return fmt.Errorf("failed to read root certificate: %w", err)
	}

	rootCert, err := crypto.ParseCertificate(rootCertPEM)
	if err != nil {
		return fmt.Errorf("failed to parse root certificate: %w", err)
	}

	// Read root CA private key
	rootKeyPEM, err := os.ReadFile(rootKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read root private key: %w", err)
	}

	rootKey, err := crypto.ParsePrivateKeyFromPEM(rootKeyPEM)
	if err != nil {
		return fmt.Errorf("failed to decode root private key: %w", err)
	}

	// Read CSR
	csrPEM, err := os.ReadFile(csrPath)
	if err != nil {
		return fmt.Errorf("failed to read CSR: %w", err)
	}

	csr, err := crypto.ParseCSR(csrPEM)
	if err != nil {
		return fmt.Errorf("failed to parse CSR: %w", err)
	}

	// Verify CSR signature
	if err := csr.CheckSignature(); err != nil {
		return fmt.Errorf("CSR signature verification failed: %w", err)
	}

	// Generate serial number
	serial, err := crypto.GenerateSerialNumber()
	if err != nil {
		return fmt.Errorf("failed to generate serial number: %w", err)
	}

	// Create intermediate CA certificate template
	notBefore := time.Now()
	notAfter := notBefore.AddDate(validityYears, 0, 0)

	template := &x509.Certificate{
		SerialNumber:          serial,
		Subject:               csr.Subject,
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            0, // No further subordinate CAs
		MaxPathLenZero:        true,
	}

	// Sign the certificate
	certPEM, err := crypto.SignCertificate(template, rootCert, csr.PublicKey, rootKey)
	if err != nil {
		return fmt.Errorf("failed to sign certificate: %w", err)
	}

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write intermediate CA certificate
	certPath := filepath.Join(outputDir, "intermediate-ca.crt")
	if err := os.WriteFile(certPath, certPEM, 0644); err != nil {
		return fmt.Errorf("failed to write certificate: %w", err)
	}

	fmt.Printf("Intermediate CA certificate signed: %s\n", certPath)
	fmt.Printf("Valid from %s to %s\n", notBefore.Format(time.RFC3339), notAfter.Format(time.RFC3339))
	fmt.Printf("Issuer: %s\n", rootCert.Subject.CommonName)
	fmt.Printf("Subject: %s\n", csr.Subject.CommonName)

	return nil
}
