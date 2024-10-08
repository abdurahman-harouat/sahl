# todo : check dependencies , and if they do not exist , install them
echo "Downloading repository ...."
git clone https://github.com/abdurahman-harouat/pie
cd pie

echo "building ...."
go build -o pie


echo "installing ...."

# delete older pie version if exists
sudo rm /usr/local/bin/pie

# copy pie to bin folder
sudo mv pie /usr/local/bin/

# allow pie file to run as a program
sudo chmod +x /usr/local/bin/pie