#!/usr/bin/env bash
# Builds template.apk and template.aab from AndroidManifest.xml and res/ with aapt2.
# Usage: ANDROID_HOME=/opt/android-sdk ./make.bash
set -euo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SDK="${ANDROID_HOME:-${ANDROID_SDK_ROOT:-/opt/android-sdk}}"
AAPT2="$(ls -d "$SDK"/build-tools/*/aapt2 | sort -V | tail -1)"
JAR="$(ls -d "$SDK"/platforms/android-*/android.jar | sort -V | tail -1)"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

"$AAPT2" compile --dir "$HERE/res" -o "$TMP/res.zip"
LINK=(link -I "$JAR" --manifest "$HERE/AndroidManifest.xml" --min-sdk-version 24 --target-sdk-version 36 "$TMP/res.zip")
"$AAPT2" "${LINK[@]}" -o "$HERE/template.apk"
"$AAPT2" "${LINK[@]}" --proto-format -o "$TMP/proto.apk"

mkdir -p "$TMP/aab/base/manifest"
(cd "$TMP/aab/base" && unzip -q ../../proto.apk && mv AndroidManifest.xml manifest/)
printf '\x0a\x08\x12\x061.18.3\x12\x0c\x0a\x00\x12\x04\x08\x01\x10\x02\x32\x02\x08\x01' > "$TMP/aab/BundleConfig.pb"
rm -f "$HERE/template.aab"
(cd "$TMP/aab" && TZ=UTC find . -exec touch -t 198001010000 {} + && zip -qrXD "$HERE/template.aab" BundleConfig.pb base)
