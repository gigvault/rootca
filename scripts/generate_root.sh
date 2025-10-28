#!/bin/bash
# Air-gapped Root CA Generation Script
# This script should be run on a secure, air-gapped system

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"
OUTPUT_DIR="${OUTPUT_DIR:-$ROOT_DIR/certs}"
CN="${CN:-GigVault Root CA}"
VALIDITY_YEARS="${VALIDITY_YEARS:-20}"

echo "==================================="
echo "GigVault Root CA Generation"
echo "==================================="
echo ""
echo "Common Name: $CN"
echo "Validity: $VALIDITY_YEARS years"
echo "Output Directory: $OUTPUT_DIR"
echo ""

# Ensure we're running on an air-gapped system
echo "⚠️  WARNING: This should be run on an air-gapped, secure system!"
echo ""
read -p "Continue? (yes/no): " confirm
if [ "$confirm" != "yes" ]; then
    echo "Aborted."
    exit 1
fi

# Build the rootca tool
echo "Building rootca tool..."
cd "$ROOT_DIR"
go build -o ./bin/rootca ./cmd/rootca

# Generate the root CA
echo "Generating root CA..."
./bin/rootca -action generate -cn "$CN" -output "$OUTPUT_DIR" -validity "$VALIDITY_YEARS"

echo ""
echo "==================================="
echo "Root CA generation complete!"
echo "==================================="
echo ""
echo "Next steps:"
echo "1. Verify the generated certificates"
echo "2. Store the private key securely offline"
echo "3. Transfer the root certificate (root-ca.crt) to the intermediate CA system"
echo "4. DESTROY this private key after proper backup!"

