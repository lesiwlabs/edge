#!/usr/bin/env sh

PROGRAM="$BIN_URL/$APP-$(uname -s)-$(uname -m)"
INSTALL_DIR="$HOME/.local/bin"

if [ "$(id -u)" = "0" ]
then
    INSTALL_DIR="/usr/local/bin"
fi

if ! TEMP="$(mktemp)"
then
    >&2 echo "Failed to create temporary file."
    exit 1
fi

if ! curl -fsSL -o "$TEMP" "$PROGRAM"
then
    >&2 echo "Failed to download $APP for platform ($(uname -s), $(uname -m))."
    >&2 echo "Please file a ticket: $TICKET_URL"
    exit 1
fi

if ! mkdir -p "$INSTALL_DIR"
then
    >&2 echo "Failed to mkdir \"$INSTALL_DIR\"; do you have permissions?"
    exit 1
fi

if ! install -m 755 "$TEMP" "$INSTALL_DIR/$APP"
then
    >&2 echo "Failed to install $APP; do you have write permissions for $INSTALL_DIR?"
    exit 1
fi

rm -f "$TEMP"
hash -r

if ! command -v "$APP" >/dev/null 2>&1
then
    >&2 echo "WARNING: $INSTALL_DIR is not on your \$PATH."
    case "$SHELL" in
        */zsh) >&2 echo "Add the following line to your ~/.zshrc:"
               ;;
        */bash) >&2 echo "Add the following line to your ~/.bashrc:"
                ;;
        *)      >&2 echo "Add the following line to your ~/.profile:"
    esac
    if [ "$(id -u)" = "0" ]
    then
        echo "  PATH=\"$INSTALL_DIR:\$PATH\""
    else
        echo "  PATH=\"\$HOME/.local/bin:\$PATH\""
    fi
fi

echo "$INSTALL_DIR/$APP"
