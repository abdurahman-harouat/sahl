# sahl : a package manager for a linux from scratch

it is a simple package manager that can install packages from source, using yaml package definition files on x86_64 & ARM architectures

inspired by [pacman](https://wiki.archlinux.org/title/Pacman) and [Makepkg](https://wiki.archlinux.org/title/Makepkg), it was made for learning purposes.

you can find the list of packages that can be installed [here](https://github.com/abdurahman-harouat/fennec-hub/tree/main/source_files)

## requirements :

### runtime requirements :

- libarchive
- tar

### setting up xorg envirement :

```bash
cat > .xorg_env.sh <<EOF
export XORG_PREFIX="/usr"
export XORG_CONFIG="--prefix=$XORG_PREFIX --sysconfdir=/etc \
    --localstatedir=/var --disable-static"
EOF
```

## Usage :

**note** : make sure that `XORG_PREFIX` and `XORG_CONFIG` are set correctly before installing xorg libraries

```bash
# install a package
sahl -i <package_name>
# verbose output
sahl -v -i <package_name>
# list installed packages
sahl -l
# uninstall a package
sahl -u <package_name>
# display help
sahl -h
# check if a package is installed
sahl -d <package_name>
# force reinstallation of a package
sahl -f -i <package_name>
```

## installation :

**installation requirements :**

- git

run this command to install it:

```bash
# if you have already setup bash startup files
wget -qO- https://raw.githubusercontent.com/abdurahman-harouat/sahl/refs/heads/main/sahl_installer.sh | sh

# if you have not setup bash startup files you can add the --with-startup-files flag
wget -qO- https://raw.githubusercontent.com/abdurahman-harouat/sahl/refs/heads/main/sahl_installer.sh | sh -s -- --with-startup-files

source ~/.bashrc
source ~/.profile

# now install libarchive as some packages require it
sahl -i libarchive
```

## build :

```bash
git clone https://github.com/abdurahman-harouat/sahl.git
cd sahl
go build
```

## Features :

- [x] - install from source
- [x] - md5 checksum verification
- [x] - dependency resolution
- [ ] - install multiple packages at once

## TODO :

- [x] automating the installation process of the package manager
- [ ] checking dependencies in the sahl_installer.sh "making it more automated"
- [ ] add a way to install a group of packages
- [ ] add a way to uninstall packages
- [ ] add a way to update packages
- [ ] add a way to search for packages in the repository
- [ ] add a way to install binary packages
- [ ] add a way to install multiple packages at once
- [x] add a way to chech if a package is installed
- [x] add a way to force reinstallation of a package
- [x] ability to download additional patches or docs
