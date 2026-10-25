#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
output_dir="${1:-$repo_root/os-payload}"
build_dir="$(mktemp -d)"
raw_persistence="$build_dir/persistence-base.raw"

cleanup() {
    rm -rf "$build_dir"
}
trap cleanup EXIT

mkdir -p "$output_dir"
cp -a "$repo_root/os-recipe/config" "$build_dir/config"
cd "$build_dir"

lb config \
    --distribution trixie \
    --architectures amd64 \
    --binary-images iso-hybrid \
    --debian-installer none \
    --archive-areas "main contrib non-free-firmware" \
    --bootappend-live "boot=live components username=myos hostname=myportableos locales=en_US.UTF-8 keyboard-layouts=us persistence persistence-label=MYOSDATA apparmor=1 security=apparmor" \
    --bootappend-live-failsafe "boot=live components username=myos hostname=myportableos locales=en_US.UTF-8 keyboard-layouts=us persistence persistence-label=MYOSDATA apparmor=1 security=apparmor noapic noacpi noapm nodma nomce nolapic nomodules nosmp"

lb build

iso="$(find . -maxdepth 1 -type f -name '*.iso' -print -quit)"
if [[ -z "$iso" ]]; then
    echo "live-build completed without producing an ISO image" >&2
    exit 1
fi
cp "$iso" "$output_dir/core.iso"

truncate -s 16G "$raw_persistence"
mkfs.ext4 -F -L MYOSDATA \
    -E lazy_itable_init=1,lazy_journal_init=1 \
    -d "$repo_root/os-recipe/config/persistence" "$raw_persistence"
e2label "$raw_persistence" | grep -q '^MYOSDATA$'
qemu-img convert -f raw -O qcow2 -c \
    "$raw_persistence" "$output_dir/persistence-base.qcow2"

qemu-img info "$output_dir/persistence-base.qcow2"
