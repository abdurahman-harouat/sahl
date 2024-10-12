# todo : check dependencies , and if they do not exist , install them
echo "Downloading repository ...."
git clone https://github.com/abdurahman-harouat/sahl
cd sahl

echo "building ...."
go build -o sahl


echo "installing ...."

# delete older sahl version if exists
sudo rm /usr/local/bin/sahl

# copy sahl to bin folder
sudo mv sahl /usr/local/bin/

# allow sahl file to run as a program
sudo chmod +x /usr/local/bin/sahl
