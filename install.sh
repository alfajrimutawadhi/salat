#!/bin/bash
version="v1.0.3"
os=$(uname -s)
arch=$(uname -m)

if [ "$os" == "Linux" ]; then
    os="linux"
elif [ "$os" == "Darwin" ]; then
    os="darwin"
else
    echo "Unsupported OS: $os"
    exit 1
fi

if [ "$arch" == "x86_64" ]; then
    arch="amd64"
elif [ "$arch" == "arm64" ]; then
    arch="arm64"
else
    echo "Unsupported architecture: $arch"
    exit 1
fi

mkdir salat_install;
echo "==== Downloading salat $version for $os-$arch ====";
url="https://github.com/alfajrimutawadhi/salat/releases/download/$version/salat-$version-$os-$arch.tar.gz"
if command -v curl >/dev/null 2>&1; then
    curl -fSL -o salat-$version-$os-$arch.tar.gz "$url" || exit 1
elif command -v wget >/dev/null 2>&1; then
    wget -O salat-$version-$os-$arch.tar.gz "$url" || exit 1
else
    echo "Error: neither curl nor wget is installed"
    exit 1
fi
if [ ! -d "$HOME/.salat" ]; then
    mkdir "$HOME/.salat";
fi
tar -xf salat-$version-$os-$arch.tar.gz -C salat_install;
rm salat-$version-$os-$arch.tar.gz;
mv salat_install/* "$HOME/.salat";
rm -rf salat_install;
chmod +x "$HOME/.salat/salat";
chmod 666 "$HOME/.salat/config.json"

# Move binary to a directory in PATH; use sudo only if needed
if [ -w /usr/local/bin ] && [ -d /usr/local/bin ]; then
    mv "$HOME/.salat/salat" /usr/local/bin;
elif command -v sudo >/dev/null 2>&1; then
    echo "==== SALAT_PASSWORD_NOTE: requesting sudo access to install to /usr/local/bin ===="
    sudo mv "$HOME/.salat/salat" /usr/local/bin;
else
    echo "Warning: cannot write to /usr/local/bin and sudo not available."
    echo "Installing to $HOME/.local/bin instead."
    mkdir -p "$HOME/.local/bin";
    mv "$HOME/.salat/salat" "$HOME/.local/bin";
    case ":$PATH:" in
        *":$HOME/.local/bin:"*) ;;
        *) echo "Note: add $HOME/.local/bin to your PATH (e.g. add to ~/.bashrc or ~/.zshrc):" ; echo '  export PATH="$HOME/.local/bin:$PATH"' ;;
    esac
fi
echo "==== Salat installed successfully ====";