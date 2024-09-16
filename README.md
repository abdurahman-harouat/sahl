# pie : a package manager for a linux from scratch

it is a simple package manager that can install packages from source, using yaml package definition files.

inspired by [pacman](https://wiki.archlinux.org/title/Pacman) and [Makepkg](https://wiki.archlinux.org/title/Makepkg), it was made for learning purposes.

## requirements :

### runtime requirements :

- libarchive
- tar

### build requirements :

- go

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
pie -i <package_name>
# verbose output
pie -v -i <package_name>
# list installed packages
pie -l
# uninstall a package
pie -u <package_name>
# display help
pie -h
# check if a package is installed
pie -d <package_name>
# force reinstallation of a package
pie -f -i <package_name>
```

## build :

```bash
git clone https://github.com/abdurahman-harouat/pie.git
cd pie
go build
```

## Features :

- [x] - install from source
- [x] - md5 checksum verification
- [x] - dependency resolution
- [ ] - install multiple packages at once

## TODO :

- [ ] automating the installation process of the pacakge manager
- [ ] add a way to uninstall packages
- [ ] add a way to update packages
- [ ] add a way to search for packages in the repository
- [ ] add a way to install binary packages
- [ ] add a way to install multiple packages at once
- [x] add a way to chech if a package is installed
- [x] add a way to force reinstallation of a package
- [x] ability to download additional patches or docs
