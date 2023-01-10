#!/bin/bash
version="v1.0.0"
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
wget https://github.com/alfajrimutawadhi/salat/releases/download/$version/salat-$version-$os-$arch.tar.gz;
tar -xvf salat-$version-$os-$arch.tar.gz -C salat_install;
rm salat-$version-$os-$arch.tar.gz;
chmod +x salat_install/salat;
chmod 666 salat_install/location.toml;
mkdir ~/.salat;
mv salat_install/* ~/.salat;
rm -rf salat_install;
mv ~/.salat/salat /usr/local/bin;
echo "==== Salat installed successfully ====";