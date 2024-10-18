#!/bin/bash

# TODO: Check dependencies, and if they do not exist, install them
echo "Downloading repository ...."
git clone https://github.com/abdurahman-harouat/sahl
cd sahl

echo "Building ...."
go build -o sahl

echo "Installing ...."

# Delete older sahl version if it exists
if [ -e /usr/local/bin/sahl ]; then
    sudo rm /usr/local/bin/sahl
    echo "Old version of sahl deleted."
else
    echo "No old version of sahl found."
fi

# Copy sahl to bin folder
sudo mv sahl /usr/local/bin/

# Allow sahl file to run as a program
sudo chmod +x /usr/local/bin/sahl
echo "sahl installed successfully."

# Set environment variables permanently
echo "Setting environment variables permanently..."

# Check if the exports already exist in .bashrc
if ! grep -q "XORG_PREFIX" ~/.bashrc; then
    {
        echo 'export XORG_PREFIX="/usr"'
        echo 'export XORG_CONFIG="--prefix=$XORG_PREFIX --sysconfdir=/etc --localstatedir=/var --disable-static"'
    } >> ~/.bashrc
    echo "Environment variables added to .bashrc."
else
    echo "Environment variables already exist in .bashrc."
fi

# Create /var/log/packages.log if it does not exist
if [ ! -d /var ]; then
    echo "Creating /var directory..."
    sudo mkdir /var
fi

if [ ! -d /var/log ]; then
    echo "Creating /var/log directory..."
    sudo mkdir /var/log
fi

# Create /var/log/packages.log if it does not exist
if [ ! -f /var/log/packages.log ]; then
    sudo touch /var/log/packages.log
    echo "/var/log/packages.log created."
else
    echo "/var/log/packages.log already exists."
fi

# Notify user to apply changes
echo "Please run 'source ~/.bashrc' or restart your terminal to apply the changes."
