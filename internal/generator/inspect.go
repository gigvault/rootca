package generator

import (
	"crypto/x509"
	"fmt"
	"os"

	"github.com/gigvault/shared/pkg/crypto"
)

// InspectCertificate displays certificate information
func InspectCertificate(certPath string) error {
	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		return fmt.Errorf("failed to read certificate: %w", err)
	}

	cert, err := crypto.ParseCertificate(certPEM)
	if err != nil {
		return fmt.Errorf("failed to parse certificate: %w", err)
	}

	fmt.Println("Certificate Information:")
	fmt.Println("========================")
	fmt.Printf("Subject: %s\n", cert.Subject.String())
	fmt.Printf("Issuer: %s\n", cert.Issuer.String())
	fmt.Printf("Serial Number: %s\n", cert.SerialNumber.Text(16))
	fmt.Printf("Not Before: %s\n", cert.NotBefore.Format("2006-01-02 15:04:05 MST"))
	fmt.Printf("Not After: %s\n", cert.NotAfter.Format("2006-01-02 15:04:05 MST"))
	fmt.Printf("Is CA: %v\n", cert.IsCA)
	fmt.Printf("Max Path Length: %d\n", cert.MaxPathLen)
	fmt.Printf("Key Usage: %s\n", keyUsageToString(cert.KeyUsage))

	if len(cert.ExtKeyUsage) > 0 {
		fmt.Printf("Extended Key Usage:\n")
		for _, eku := range cert.ExtKeyUsage {
			fmt.Printf("  - %s\n", extKeyUsageToString(eku))
		}
	}

	if len(cert.DNSNames) > 0 {
		fmt.Printf("DNS Names:\n")
		for _, dns := range cert.DNSNames {
			fmt.Printf("  - %s\n", dns)
		}
	}

	if len(cert.EmailAddresses) > 0 {
		fmt.Printf("Email Addresses:\n")
		for _, email := range cert.EmailAddresses {
			fmt.Printf("  - %s\n", email)
		}
	}

	return nil
}

func keyUsageToString(usage x509.KeyUsage) string {
	var usages []string
	if usage&x509.KeyUsageDigitalSignature != 0 {
		usages = append(usages, "DigitalSignature")
	}
	if usage&x509.KeyUsageContentCommitment != 0 {
		usages = append(usages, "ContentCommitment")
	}
	if usage&x509.KeyUsageKeyEncipherment != 0 {
		usages = append(usages, "KeyEncipherment")
	}
	if usage&x509.KeyUsageDataEncipherment != 0 {
		usages = append(usages, "DataEncipherment")
	}
	if usage&x509.KeyUsageKeyAgreement != 0 {
		usages = append(usages, "KeyAgreement")
	}
	if usage&x509.KeyUsageCertSign != 0 {
		usages = append(usages, "CertSign")
	}
	if usage&x509.KeyUsageCRLSign != 0 {
		usages = append(usages, "CRLSign")
	}
	if usage&x509.KeyUsageEncipherOnly != 0 {
		usages = append(usages, "EncipherOnly")
	}
	if usage&x509.KeyUsageDecipherOnly != 0 {
		usages = append(usages, "DecipherOnly")
	}

	if len(usages) == 0 {
		return "None"
	}

	result := usages[0]
	for i := 1; i < len(usages); i++ {
		result += ", " + usages[i]
	}
	return result
}

func extKeyUsageToString(usage x509.ExtKeyUsage) string {
	switch usage {
	case x509.ExtKeyUsageAny:
		return "Any"
	case x509.ExtKeyUsageServerAuth:
		return "ServerAuth"
	case x509.ExtKeyUsageClientAuth:
		return "ClientAuth"
	case x509.ExtKeyUsageCodeSigning:
		return "CodeSigning"
	case x509.ExtKeyUsageEmailProtection:
		return "EmailProtection"
	case x509.ExtKeyUsageIPSECEndSystem:
		return "IPSECEndSystem"
	case x509.ExtKeyUsageIPSECTunnel:
		return "IPSECTunnel"
	case x509.ExtKeyUsageIPSECUser:
		return "IPSECUser"
	case x509.ExtKeyUsageTimeStamping:
		return "TimeStamping"
	case x509.ExtKeyUsageOCSPSigning:
		return "OCSPSigning"
	case x509.ExtKeyUsageMicrosoftServerGatedCrypto:
		return "MicrosoftServerGatedCrypto"
	case x509.ExtKeyUsageNetscapeServerGatedCrypto:
		return "NetscapeServerGatedCrypto"
	case x509.ExtKeyUsageMicrosoftCommercialCodeSigning:
		return "MicrosoftCommercialCodeSigning"
	case x509.ExtKeyUsageMicrosoftKernelCodeSigning:
		return "MicrosoftKernelCodeSigning"
	default:
		return fmt.Sprintf("Unknown(%d)", usage)
	}
}
