# GigVault Root CA

Offline root Certificate Authority tools and scripts for GigVault PKI.

## ⚠️ Security Notice

The Root CA is the trust anchor of the entire PKI infrastructure. It should:
- Be generated on an **air-gapped system**
- Have its private key stored **offline** in a secure location (HSM, secure vault)
- Only be brought online for signing intermediate CA certificates
- Use ECDSA P-384 for maximum security

## Usage

### Generate Root CA

```bash
# Using the script
./scripts/generate_root.sh

# Or directly with the tool
go run ./cmd/rootca -action generate \
  -cn "GigVault Root CA" \
  -output ./certs \
  -validity 20
```

### Sign Intermediate CA Certificate

```bash
# TODO: Implement
go run ./cmd/rootca -action sign \
  -csr intermediate-ca.csr \
  -root-cert ./certs/root-ca.crt \
  -root-key ./certs/root-ca.key
```

## Directory Structure

```
rootca/
├── cmd/rootca/          # CLI tool
├── internal/
│   ├── generator/       # Root CA generation logic
│   └── storage/         # Secure key storage adapters
├── scripts/
│   └── generate_root.sh # Air-gapped generation script
└── certs/               # Generated certificates (gitignored)
```

## Development

```bash
# Build
make build

# Test
make test

# Build Docker image
make docker
```

## Operational Security

1. **Generation**: Always generate on an air-gapped system
2. **Storage**: Store private key in HSM or encrypted offline storage
3. **Access**: Limit access to root key to authorized personnel only
4. **Rotation**: Plan for root CA rotation every 10-20 years
5. **Backup**: Maintain multiple encrypted backups in separate locations

## License

Copyright © 2025 GigVault

