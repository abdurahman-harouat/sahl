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


# Create /etc/profile.d directory with proper permissions and ownership
# log_message "Setting up /etc/profile.d directory..."
# sudo install --directory --mode=0755 --owner=root --group=root /etc/profile.d
# log_message "✓ Created /etc/profile.d directory with proper permissions and ownership"


# Making sure that bash_competion.d directory is present
# log_message "Adding bash_completion.d..."
# sudo install --directory --mode=0755 --owner=root --group=root /etc/bash_completion.d
# log_message "✓ Added bash_completion.d directory succesfully"


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
    
    if wget "$GO_URL"; then
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