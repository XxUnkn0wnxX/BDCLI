#!/usr/local/bin/zsh

emulate -L zsh
set -euo pipefail

script_dir="${0:A:h}"

if [[ -f "${script_dir}/go.mod" && -f "${script_dir}/main.go" ]]; then
  repo_root="${script_dir}"
elif [[ -f "${script_dir:h}/go.mod" && -f "${script_dir:h}/main.go" ]]; then
  repo_root="${script_dir:h}"
else
  echo "Could not detect repository root from script location: ${script_dir}" >&2
  exit 1
fi

cd "$repo_root"

# Initialize a sane macOS PATH without relying on the user's shell rc files.
if [[ -x /usr/libexec/path_helper ]]; then
  eval "$(/usr/libexec/path_helper -s)"
fi

for path_dir in /usr/local/bin /opt/homebrew/bin /usr/local/opt/go/bin /opt/homebrew/opt/go/bin; do
  if [[ -d "$path_dir" && ":$PATH:" != *":$path_dir:"* ]]; then
    export PATH="$path_dir:$PATH"
  fi
done

if ! command -v go >/dev/null 2>&1; then
  echo "go is not installed or not in PATH" >&2
  exit 1
fi

echo "Using repository root: $repo_root"
echo "Using Go binary: $(command -v go)"
go version

set -x

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
sh -x ./scripts/completions.sh

echo "Building macOS 11 Big Sur Intel binary..."
CGO_ENABLED=0 \
GOOS=darwin \
GOARCH=amd64 \
GOAMD64=v1 \
MACOSX_DEPLOYMENT_TARGET=11.0 \
go build \
  -x \
  -v \
  -trimpath \
  -ldflags "-s -w -X main.version=${short_sha} -X main.commit=${short_sha} -X main.date=${build_date}" \
  -o "$binary_path" \
  ./main.go

cp "$binary_path" "$archive_root/bdcli"
cp README.md LICENSE "$archive_root/"
cp completions/* "$archive_root/completions/"
ls -l "$binary_path" "$archive_root/bdcli" "$archive_root/completions"

echo "Creating artifact archive..."
tar -C "$archive_root" -cvzf "$artifact_path" .

set +x
echo "Build complete."
echo "Binary:   $binary_path"
echo "Artifact: $artifact_path"
