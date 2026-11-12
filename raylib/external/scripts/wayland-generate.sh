#!/bin/sh
# Generate the Wayland protocol headers used by GLFW and the native Wayland platform into ../wayland

set -e

OUT="$(cd "$(dirname "$0")/../wayland" && pwd)"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

git clone -q --depth 1 https://gitlab.freedesktop.org/wayland/wayland.git "$TMP/wayland"
git clone -q --depth 1 https://gitlab.freedesktop.org/wayland/wayland-protocols.git "$TMP/wayland-protocols"

generate() {
    wayland-scanner client-header "$1" "$OUT/$2-client-protocol.h"
    wayland-scanner private-code "$1" "$OUT/$2-client-protocol-code.h"
}

rm -f "$OUT"/*-client-protocol.h "$OUT"/*-client-protocol-code.h

generate "$TMP/wayland/protocol/wayland.xml" wayland

P="$TMP/wayland-protocols"
generate "$P/stable/xdg-shell/xdg-shell.xml" xdg-shell
generate "$P/stable/viewporter/viewporter.xml" viewporter
generate "$P/stable/tablet/tablet-v2.xml" tablet-v2
generate "$P/staging/fractional-scale/fractional-scale-v1.xml" fractional-scale-v1
generate "$P/staging/xdg-activation/xdg-activation-v1.xml" xdg-activation-v1
generate "$P/staging/cursor-shape/cursor-shape-v1.xml" cursor-shape-v1
generate "$P/staging/xdg-toplevel-icon/xdg-toplevel-icon-v1.xml" xdg-toplevel-icon-v1
generate "$P/staging/alpha-modifier/alpha-modifier-v1.xml" alpha-modifier-v1
generate "$P/unstable/xdg-decoration/xdg-decoration-unstable-v1.xml" xdg-decoration-unstable-v1
generate "$P/unstable/relative-pointer/relative-pointer-unstable-v1.xml" relative-pointer-unstable-v1
generate "$P/unstable/pointer-constraints/pointer-constraints-unstable-v1.xml" pointer-constraints-unstable-v1
generate "$P/unstable/idle-inhibit/idle-inhibit-unstable-v1.xml" idle-inhibit-unstable-v1
generate "$P/unstable/xdg-output/xdg-output-unstable-v1.xml" xdg-output-unstable-v1
