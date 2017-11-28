#!/bin/sh

TMPDIR="/tmp"
GLGLFW_PATH="`pwd`/../glfw/src"

cd $TMPDIR
git clone https://github.com/wayland-project/wayland-protocols
cd wayland-protocols

rm -f "$GLGLFW_PATH"/wayland-pointer-constraints-unstable-v1-client-protocol.{h,c}
rm -f "$GLGLFW_PATH"/wayland-relative-pointer-unstable-v1-client-protocol.{h,c}
rm -f "$GLGLFW_PATH"/wayland-idle-inhibit-unstable-v1-client-protocol.{h,c}

wayland-scanner code ./unstable/pointer-constraints/pointer-constraints-unstable-v1.xml "$GLGLFW_PATH"/wayland-pointer-constraints-unstable-v1-client-protocol.c
wayland-scanner client-header ./unstable/pointer-constraints/pointer-constraints-unstable-v1.xml "$GLGLFW_PATH"/wayland-pointer-constraints-unstable-v1-client-protocol.h

wayland-scanner code ./unstable/relative-pointer/relative-pointer-unstable-v1.xml "$GLGLFW_PATH"/wayland-relative-pointer-unstable-v1-client-protocol.c
wayland-scanner client-header ./unstable/relative-pointer/relative-pointer-unstable-v1.xml "$GLGLFW_PATH"/wayland-relative-pointer-unstable-v1-client-protocol.h

wayland-scanner code ./unstable/idle-inhibit/idle-inhibit-unstable-v1.xml "$GLGLFW_PATH"/wayland-idle-inhibit-unstable-v1-client-protocol.c
wayland-scanner client-header ./unstable/idle-inhibit/idle-inhibit-unstable-v1.xml "$GLGLFW_PATH"/wayland-idle-inhibit-unstable-v1-client-protocol.h

# Patch for cgo
sed -i "s|types|wl_pc_types|g" "$GLGLFW_PATH"/wayland-pointer-constraints-unstable-v1-client-protocol.c
sed -i "s|types|wl_rp_types|g" "$GLGLFW_PATH"/wayland-relative-pointer-unstable-v1-client-protocol.c
sed -i "s|types|wl_ii_types|g" "$GLGLFW_PATH"/wayland-idle-inhibit-unstable-v1-client-protocol.c
