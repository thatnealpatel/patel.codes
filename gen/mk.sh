#!/bin/sh
# patel.codes/mk.sh: vanity shim for setting up new hosts
set -e
src='https://raw.githubusercontent.com/thatnealpatel/dotfiles2/main/.bootstrap/bootstrap.sh'
exec bash -c "$(curl -fsSL "$src")" bootstrap.sh "$@"
