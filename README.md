# MyPortableOS

MyPortableOS is a Debian-based live desktop that runs in QEMU. It includes an
XFCE desktop, Firefox ESR, and a small collection of everyday apps. It does not
install itself to the host computer.

## Requirements

- Linux AppImage: install QEMU (`qemu-system-x86_64`) and `qemu-utils`
  (`qemu-img`).
- Windows executable: install QEMU for Windows, including `qemu-img.exe`.
- At least 4 GB of memory available to the virtual machine.

The virtual machine uses QEMU user-mode networking and a dedicated 16 GB
virtual disk. It does not pass through or repartition the host's physical
disks. The virtual disk stores guest changes and files persistently on the
host:

- Linux: `${XDG_DATA_HOME:-~/.local/share}/myportableos/`
- Windows: `%LOCALAPPDATA%\MyPortableOS\`

Deleting `persistence.qcow2` deletes the virtual machine's saved changes and
files. The base disk and each Windows OS image are also stored in this
application data directory. They do not replace or delete files elsewhere on
the host.

## Privacy and security

The OS image contains no advertising or bundled third-party browser
extensions. Firefox ESR is configured to disable telemetry and sponsored
content and to enable strict tracking protection. The live system enables
AppArmor and a default-deny incoming firewall.

These are security and privacy measures, not a guarantee that any system is
"hyper-secure." Keep the host, QEMU, and the guest image updated, and do not
open untrusted files or websites just because they are inside a virtual
machine.

## Releases

Push a version tag such as `v1.0.0` to build and publish a GitHub Release. Each release includes:

- `MyOS-x86_64.AppImage` for Linux
- `MyOS-windows-x86_64.exe` for Windows
- GitHub-generated source code ZIP and TAR archives for the tagged revision
