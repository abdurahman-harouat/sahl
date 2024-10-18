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
