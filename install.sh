#!/bin/bash

binary_path="./lo"
install_path="/usr/local/bin"

echo "installing lo binary..."
echo "You need to have elevated privileges to install the lo binary"

sudo cp "$binary_path" "$install_path/lo"

sudo chmod +x "$install_path/lo"

echo "lo binary installed to $install_path/lo"
echo "installation complete 💯"

echo "reload your shell to use lo"

