#!/usr/local/bin/zsh

emulate -L zsh
set -euo pipefail

script_dir="${0:A:h}"
repo_root="${script_dir:h}"

cd "$repo_root"

if ! command -v go >/dev/null 2>&1; then
  echo "go is not installed or not in PATH" >&2
  exit 1
fi

short_sha="$(git rev-parse --short HEAD 2>/dev/null || echo manual)"
if ! git diff-index --quiet HEAD -- 2>/dev/null; then
  short_sha="${short_sha}-dirty"
fi

build_date="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
binary_name="bdcli-darwin-amd64-macos11"
binary_path="dist/${binary_name}"
archive_root="dist/archive"
artifact_name="bdcli_${short_sha}_darwin_amd64_macos11.tar.gz"
artifact_path="dist/${artifact_name}"

rm -rf "$archive_root"
mkdir -p dist "$archive_root/completions"

echo "Generating shell completions..."
sh ./scripts/completions.sh

echo "Building macOS 11 Big Sur Intel binary..."
CGO_ENABLED=0 \
GOOS=darwin \
GOARCH=amd64 \
GOAMD64=v1 \
MACOSX_DEPLOYMENT_TARGET=11.0 \
go build \
  -trimpath \
  -ldflags "-s -w -X main.version=${short_sha} -X main.commit=${short_sha} -X main.date=${build_date}" \
  -o "$binary_path" \
  ./main.go

cp "$binary_path" "$archive_root/bdcli"
cp README.md LICENSE "$archive_root/"
cp completions/* "$archive_root/completions/"

echo "Creating artifact archive..."
tar -C "$archive_root" -czf "$artifact_path" .

echo "Build complete."
echo "Binary:   $binary_path"
echo "Artifact: $artifact_path"
