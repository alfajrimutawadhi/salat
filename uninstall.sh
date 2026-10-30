#!/bin/bash
version="v1.0.4"
os=$(uname -s)
arch=$(uname -m)

if [[ "$os" == MINGW* ]] || [[ "$os" == MSYS* ]]; then
    echo "Windows detected - use the PowerShell uninstaller instead:"
    echo '  powershell -c "irm https://raw.githubusercontent.com/alfajrimutawadhi/salat/master/uninstall.ps1 | iex"'
    exit 1
fi

echo "==== Uninstalling salat $version for $os-$arch ====";
rm -r ~/.salat;
rm -f /usr/local/bin/salat;
echo "==== Salat uninstalled successfully ====";