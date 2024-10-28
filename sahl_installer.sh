#!/bin/bash

# Print a formatted message with timestamp
log_message() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1"
}

# Print a section header
print_section() {
    echo "================================================================"
    log_message "$1"
    echo "================================================================"
}

# setting bash startup files
print_section "Setting up bash startup files"
# Iterate over all arguments
for arg in "$@"; do
    if [ "$arg" == "--with-startup-files" ]; then
        mkdir -p bash_startup_files
        cd bash_startup_files
        wget --no-check-certificate -q https://github.com/abdurahman-harouat/sahl/raw/refs/heads/main/bash_startup_files/.bash_logout
        wget --no-check-certificate -q https://github.com/abdurahman-harouat/sahl/raw/refs/heads/main/bash_startup_files/.bash_profile
        wget --no-check-certificate -q https://github.com/abdurahman-harouat/sahl/raw/refs/heads/main/bash_startup_files/.bashrc
        wget --no-check-certificate -q https://github.com/abdurahman-harouat/sahl/raw/refs/heads/main/bash_startup_files/.profile
        cp .bash_logout .bash_profile .bashrc .profile ~
        mkdir -p etc
        cd etc
        wget --no-check-certificate -q https://github.com/abdurahman-harouat/sahl/raw/refs/heads/main/bash_startup_files/etc/bashrc
        wget --no-check-certificate -q https://github.com/abdurahman-harouat/sahl/raw/refs/heads/main/bash_startup_files/etc/profile
        sudo cp -v bashrc profile /etc/
        mkdir -p profile.d
        cd profile.d
        wget --no-check-certificate -q https://github.com/abdurahman-harouat/sahl/raw/refs/heads/main/bash_startup_files/etc/profile.d/bash_completion.sh
        wget --no-check-certificate -q https://github.com/abdurahman-harouat/sahl/raw/refs/heads/main/bash_startup_files/etc/profile.d/extrapaths.sh
        wget --no-check-certificate -q https://github.com/abdurahman-harouat/sahl/raw/refs/heads/main/bash_startup_files/etc/profile.d/i18n.sh
        wget --no-check-certificate -q https://github.com/abdurahman-harouat/sahl/raw/refs/heads/main/bash_startup_files/etc/profile.d/readline.sh
        wget --no-check-certificate -q https://github.com/abdurahman-harouat/sahl/raw/refs/heads/main/bash_startup_files/etc/profile.d/umask.sh
        wget --no-check-certificate -q https://github.com/abdurahman-harouat/sahl/raw/refs/heads/main/bash_startup_files/etc/profile.d/dircolors.sh
        sudo install --directory --mode=0755 --owner=root --group=root /etc/profile.d
        sudo install --directory --mode=0755 --owner=root --group=root /etc/bash_completion.d
        sudo cp -v bash_completion.sh extrapaths.sh i18n.sh readline.sh umask.sh dircolors.sh /etc/profile.d/
        cd ../../..
        rm -rf bash_startup_files
        log_message "✓ Bash startup files set successfully"
    fi
done


# Create required directories if they don't exist
print_section "Setting up system directories"
if [ ! -d /var ]; then
    log_message "Creating /var directory..."
    sudo mkdir /var
    log_message "✓ /var directory created successfully"
fi

if [ ! -d /var/log ]; then
    log_message "Creating /var/log directory..."
    sudo mkdir /var/log
    log_message "✓ /var/log directory created successfully"
fi

# Create and set up packages log file
print_section "Setting up package logging"
if [ ! -f /var/log/packages.log ]; then
    sudo touch /var/log/packages.log
    log_message "✓ Created /var/log/packages.log"
else
    log_message "→ /var/log/packages.log already exists"
fi

# Set up environment variables
print_section "Configuring environment variables"
if ! grep -q "XORG_PREFIX" ~/.bashrc; then
    log_message "Adding XORG environment variables to .bashrc..."
    {
        echo '# XORG Configuration'
        echo 'export XORG_PREFIX="/usr"'
        echo 'export XORG_CONFIG="--prefix=$XORG_PREFIX --sysconfdir=/etc --localstatedir=/var --disable-static"'
    } >> ~/.bashrc
    log_message "✓ Environment variables added successfully"
else
    log_message "→ XORG environment variables already exist in .bashrc"
fi

# Check for make-ca installation 
print_section "Checking make-ca installation"
if [ -f /usr/sbin/make-ca ]; then
    echo "make-ca is already installed"
else
    log_message "make-ca not found. Installing make-ca..."

    wget --no-check-certificate -q "https://github.com/lfs-book/make-ca/archive/v1.14/make-ca-1.14.tar.gz" &&
    tar -xzf make-ca-1.14.tar.gz &&
    cd make-ca-1.14 &&
    sudo make install &&
    sudo install -vdm755 /etc/ssl/local &&
    sudo /usr/sbin/make-ca -g &&
    sudo systemctl enable update-pki.timer &&

    wget "http://www.cacert.org/certs/root.crt" &&
    wget "http://www.cacert.org/certs/class3.crt" &&
    sudo openssl x509 -in root.crt -text -fingerprint -setalias "CAcert Class 1 root" \
        -addtrust serverAuth -addtrust emailProtection -addtrust codeSigning \
        | sudo tee /etc/ssl/local/CAcert_Class_1_root.pem > /dev/null &&
    sudo openssl x509 -in class3.crt -text -fingerprint -setalias "CAcert Class 3 root" \
        -addtrust serverAuth -addtrust emailProtection -addtrust codeSigning \
        | sudo tee /etc/ssl/local/CAcert_Class_3_root.pem > /dev/null &&
    sudo /usr/sbin/make-ca -r &&

    sudo mkdir -pv /etc/profile.d &&
    sudo tee /etc/profile.d/pythoncerts.sh > /dev/null << "EOF6"
# Begin /etc/profile.d/pythoncerts.sh

export _PIP_STANDALONE_CERT=/etc/pki/tls/certs/ca-bundle.crt

# End /etc/profile.d/pythoncerts.sh
EOF6

    log_message "make-ca installation and configuration completed successfully."
    
    # Clean up downloaded and extracted files
    cd .. &&
    rm -rf make-ca-1.14 make-ca-1.14.tar.gz root.crt class3.crt
fi


# Check for libunistring installation 
print_section "Checking libunistring installation"
if [ -d /usr/include/unistring ]; then
    echo "libunistring is already installed"
else
    log_message "libunistring not found. Installing libunistring..."

    wget -q "https://ftp.gnu.org/gnu/libunistring/libunistring-1.2.tar.xz" &&
    tar -xzf libunistring-1.2.tar.xz &&
    cd libunistring-1.2 &&
    ./configure --prefix=/usr    \
            --disable-static \
            --docdir=/usr/share/doc/libunistring-1.2 &&
    make &&
    sudo make install

    log_message "libunistring installation completed successfully."
    
    # Clean up downloaded and extracted files
    cd .. &&
    rm -rf libunistring-1.2 libunistring-1.2.tar.xz
fi

# Check for libidn2 installation 
print_section "Checking libidn2 installation"
if [ -d /usr/share/gtk-doc/html/libidn2 ]; then
    echo "libidn2 is already installed"
else
    log_message "libidn2 not found. Installing libidn2..."

    wget -q "https://ftp.gnu.org/gnu/libidn/libidn2-2.3.7.tar.gz" &&
    tar -xzf libidn2-2.3.7.tar.gz &&
    cd libidn2-2.3.7 &&
    ./configure --prefix=/usr --disable-static &&
    make &&
    sudo make install

    log_message "libidn2 installation completed successfully."
    
    # Clean up downloaded and extracted files
    cd .. &&
    rm -rf libidn2-2.3.7 libidn2-2.3.7.tar.gz
fi

# Check for libpsl installation 
print_section "Checking libpsl installation"
if command -v psl &> /dev/null; then
    echo "libpsl is already installed"
else
    log_message "libpsl not found. Installing libpsl..."

    wget -q "https://github.com/rockdaboot/libpsl/releases/download/0.21.5/libpsl-0.21.5.tar.gz" &&
    tar -xzf libpsl-0.21.5.tar.gz &&
    cd libpsl-0.21.5 &&
    mkdir -p build &&
    cd build &&
    meson setup --prefix=/usr --buildtype=release &&
    ninja &&
    sudo ninja install

    log_message "libpsl installation completed successfully."
    
    # Clean up downloaded and extracted files
    cd ../.. &&
    rm -rf libpsl-0.21.5 libpsl-0.21.5.tar.gz
fi

# Check for curl installation 
print_section "Checking curl installation"
if command -v curl &> /dev/null; then
    echo "curl is already installed"
else
    log_message "curl not found. Installing curl..."

    wget -q "https://curl.se/download/curl-8.9.1.tar.xz" &&
    tar -xzf curl-8.9.1.tar.xz &&
    cd curl-8.9.1 &&
    ./configure --prefix=/usr                           \
            --disable-static                        \
            --with-openssl                          \
            --enable-threaded-resolver              \
            --with-ca-path=/etc/ssl/certs &&
    make &&
    sudo make install &&
    sudo rm -rf docs/examples/.deps &&
    sudo find docs \( -name Makefile\* -o  \
                -name \*.1       -o  \
                -name \*.3       -o  \
                -name CMakeLists.txt \) -delete &&
    sudo cp -v -R docs -T /usr/share/doc/curl-8.9.1


    log_message "curl installation completed successfully."
    
    # Clean up downloaded and extracted files
    cd .. &&
    rm -rf curl-8.9.1 curl-8.9.1.tar.xz
fi


# Determine system architecture
print_section "Detecting system architecture"
ARCH=$(uname -m)
log_message "Detected architecture: $ARCH"

# Set Go version and architecture-specific URL
case $ARCH in
    x86_64)
        GO_VERSION="1.23.2"
        GO_ARCH="amd64"
        ;;
    aarch64)
        GO_VERSION="1.23.0"
        GO_ARCH="arm64"
        ;;
    *)
        log_message "ERROR: Unsupported architecture: $ARCH"
        exit 1
        ;;
esac
log_message "Selected Go version ${GO_VERSION} for ${GO_ARCH} architecture"

# Check for existing Go installation
print_section "Checking for existing Go installation"
if [ -d "/usr/local/go" ]; then
    EXISTING_GO_VERSION=$(/usr/local/go/bin/go version 2>/dev/null | awk '{print $3}' | sed 's/go//')
    log_message "Found existing Go installation: version ${EXISTING_GO_VERSION}"
    log_message "→ Skipping Go installation"
else
    # Install Go
    print_section "Installing Go"
    GO_URL="https://go.dev/dl/go${GO_VERSION}.linux-${GO_ARCH}.tar.gz"
    log_message "Downloading Go from: ${GO_URL}"
    
    if wget --no-check-certificate -q "$GO_URL"; then
        log_message "✓ Download completed successfully"
        
        log_message "Extracting Go archive..."
        if sudo tar -C /usr/local -xzf "go${GO_VERSION}.linux-${GO_ARCH}.tar.gz"; then
            log_message "✓ Go extracted successfully to /usr/local/go"
            
            # Set up Go environment variables
            if ! grep -q "/usr/local/go/bin" ~/.bashrc; then
                echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
                echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
                log_message "✓ Added Go to PATH in .bashrc and .profile"
            fi
            
            # Log installation
            echo "go ${GO_VERSION} installed on $(date)" | sudo tee -a /var/log/packages.log
            log_message "✓ Installation logged to packages.log"
            
            # Clean up downloaded archive
            rm "go${GO_VERSION}.linux-${GO_ARCH}.tar.gz"
            log_message "✓ Cleaned up installation files"
        else
            log_message "ERROR: Failed to extract Go archive"
            exit 1
        fi
    else
        log_message "ERROR: Failed to download Go"
        exit 1
    fi
fi


# Check for Git installation
print_section "Checking Git installation"
if command -v git &> /dev/null; then
    GIT_VERSION=$(git --version | awk '{print $3}')
    log_message "Git is already installed (version ${GIT_VERSION})"
else
    log_message "Git not found. Installing Git..."

    wget -q https://www.kernel.org/pub/software/scm/git/git-2.46.0.tar.xz &&
    tar -xf git-2.46.0.tar.xz &&
    cd git-2.46.0 &&
    ./configure --prefix=/usr \
            --with-gitconfig=/etc/gitconfig \
            --with-python=python3 &&
    make &&
    sudo make perllibdir=/usr/lib/perl5/5.38.2/site_perl install

    if command -v git &> /dev/null; then
        GIT_VERSION=$(git --version | awk '{print $3}')
        log_message "✓ Git installed successfully (version ${GIT_VERSION})"
        echo "git ${GIT_VERSION} installed on $(date)" | sudo tee -a /var/log/packages.log

        # Cleanup: remove downloaded files and extracted folder
        cd .. &&
        rm -rf git-2.46.0 git-2.46.0.tar.xz
    else
        log_message "ERROR: Git installation failed"
        exit 1
    fi
fi



# Install SAHL
print_section "Installing SAHL"
log_message "Cloning SAHL repository..."
if git clone https://github.com/abdurahman-harouat/sahl; then
    cd sahl || exit
    log_message "✓ Repository cloned successfully"
    
    log_message "Building SAHL..."
    if go build -o sahl; then
        log_message "✓ Build completed successfully"
        
        # Remove existing installation if present
        if [ -e /usr/local/bin/sahl ]; then
            log_message "Removing previous SAHL installation..."
            sudo rm /usr/local/bin/sahl
            log_message "✓ Previous installation removed"
        fi
        
        # Install new version
        log_message "Installing SAHL to /usr/local/bin..."
        if sudo mv sahl /usr/local/bin/ && sudo chmod +x /usr/local/bin/sahl; then
            log_message "✓ SAHL installed successfully"
            
            # Move out of the cloned directory and remove it
            cd ..
            log_message "Cleaning up SAHL repository..."
            rm -rf sahl
            log_message "✓ SAHL repository removed successfully"
        else
            log_message "ERROR: Failed to install SAHL"
            exit 1
        fi
    else
        log_message "ERROR: Failed to build SAHL"
        exit 1
    fi
else
    log_message "ERROR: Failed to clone SAHL repository"
    exit 1
fi

print_section "Installation Complete"
log_message "Please run 'source ~/.bashrc' or restart your terminal to apply the changes"
log_message "You can now use SAHL by running the 'sahl' command"
