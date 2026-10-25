@echo off
title Launching Portable OS Window...
set "BIN_PATH=%~dp0os-payload\core.iso"

:: Look for standard QEMU installations on Windows
if exist "%ProgramFiles%\qemu\qemu-system-x86_64.exe" (
    "%ProgramFiles%\qemu\qemu-system-x86_64.exe" -cdrom "%BIN_PATH%" -m 1024 -vga std -boot d
) else if exist "%ProgramFiles(x86)%\qemu\qemu-system-x86_64.exe" (
    "%ProgramFiles(x86)%\qemu\qemu-system-x86_64.exe" -cdrom "%BIN_PATH%" -m 1024 -vga std -boot d
) else (
    echo [ERROR] QEMU emulator was not detected on this system.
    echo Please install QEMU for Windows to use this portable operating system.
    pause
)
