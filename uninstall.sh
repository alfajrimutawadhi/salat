#!/bin/bash
version="v1.0.3"
os=$(uname -s)
arch=$(uname -m)

echo "==== Uninstalling salat $version for $os-$arch ====";
rm -r ~/.salat;
rm /usr/local/bin/salat;
echo "==== Salat uninstalled successfully ====";