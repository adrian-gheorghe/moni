#!/bin/bash

# Get final URL after curl is redirected
RELEASE_URL=$(curl -Ls -o /dev/null -w %{url_effective} https://github.com/adrian-gheorghe/moni/releases/latest)
# Extract tag after the last forward slash 
TAG="${RELEASE_URL##*/}"

# Check if moni is currently installed
LOCAL_PATH=$(which moni)
 if [ -x "$LOCAL_PATH" ]; then
    echo "Moni is already Installed"
    echo "Try running $(moni --help)"
    exit 0
fi

echo "Attempting to download moni v${TAG}"

if [[ "$OSTYPE" == "linux-gnu" ]]; then
    BINARY_PATH="moni-linux"
elif [[ "$OSTYPE" == "darwin"* ]]; then
    BINARY_PATH="moni-darwin"
else
    BINARY_PATH="moni"
fi

curl -L "https://github.com/adrian-gheorghe/moni/releases/download/$TAG/$BINARY_PATH" --output $BINARY_PATH
chmod +x $BINARY_PATH
mv $BINARY_PATH /usr/local/bin/moni

echo "Moni installed successfully!"
moni --help