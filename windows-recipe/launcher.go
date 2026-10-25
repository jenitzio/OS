package main

import (
	"embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

//go:embed core.iso persistence-base.qcow2
var bundledImages embed.FS

var imageVersion = "dev"

func main() {
	if err := launch(); err != nil {
		fmt.Fprintln(os.Stderr, "MyOS could not start:", err)
		os.Exit(1)
	}
}

func launch() error {
	qemuPath, err := findQEMU()
	if err != nil {
		return err
	}

	dataDir := os.Getenv("LOCALAPPDATA")
	if dataDir == "" {
		var err error
		dataDir, err = os.UserConfigDir()
		if err != nil {
			return fmt.Errorf("find the user data directory: %w", err)
		}
	}
	dataDir = filepath.Join(dataDir, "MyPortableOS")
	if err := os.MkdirAll(dataDir, 0700); err != nil {
		return fmt.Errorf("create persistent data directory: %w", err)
	}

	isoPath := filepath.Join(dataDir, "core-"+imageVersion+".iso")
	if err := extractBundledFile("core.iso", isoPath); err != nil {
		return fmt.Errorf("prepare OS image: %w", err)
	}

	basePath := filepath.Join(dataDir, "persistence-base.qcow2")
	if err := extractBundledFile("persistence-base.qcow2", basePath); err != nil {
		return fmt.Errorf("prepare persistence base image: %w", err)
	}

	persistencePath := filepath.Join(dataDir, "persistence.qcow2")
	if _, err := os.Stat(persistencePath); os.IsNotExist(err) {
		qemuImg, err := findQEMUImage(qemuPath)
		if err != nil {
			return err
		}
		create := exec.Command(qemuImg, "create", "-f", "qcow2", "-F", "qcow2",
			"-b", basePath, persistencePath)
		if output, err := create.CombinedOutput(); err != nil {
			return fmt.Errorf("create persistent virtual disk: %w: %s", err, output)
		}
	} else if err != nil {
		return fmt.Errorf("check persistent virtual disk: %w", err)
	}

	command := exec.Command(qemuPath,
		"-cdrom", isoPath,
		"-m", "4096",
		"-smp", "2",
		"-vga", "std",
		"-nic", "user,model=e1000",
		"-drive", "file="+persistencePath+",format=qcow2,if=virtio",
		"-boot", "d",
	)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin
	if err := command.Run(); err != nil {
		return fmt.Errorf("run QEMU: %w", err)
	}

	return nil
}

func findQEMUImage(qemuPath string) (string, error) {
	candidate := filepath.Join(filepath.Dir(qemuPath), "qemu-img.exe")
	if _, err := os.Stat(candidate); err == nil {
		return candidate, nil
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("check for qemu-img.exe at %q: %w", candidate, err)
	}

	if qemuImg, err := exec.LookPath("qemu-img.exe"); err == nil {
		return qemuImg, nil
	}
	return "", fmt.Errorf("qemu-img.exe was not detected next to QEMU or on PATH")
}

func extractBundledFile(name, destination string) error {
	if _, err := os.Stat(destination); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("check %q: %w", destination, err)
	}

	source, err := bundledImages.Open(name)
	if err != nil {
		return fmt.Errorf("read bundled file %q: %w", name, err)
	}
	defer source.Close()

	file, err := os.CreateTemp(filepath.Dir(destination), ".myos-image-*")
	if err != nil {
		return fmt.Errorf("create temporary image in %q: %w", filepath.Dir(destination), err)
	}
	temporaryPath := file.Name()
	defer os.Remove(temporaryPath)

	if _, err := io.Copy(file, source); err != nil {
		file.Close()
		return fmt.Errorf("extract %q: %w", name, err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close %q: %w", destination, err)
	}
	if err := os.Rename(temporaryPath, destination); err != nil {
		return fmt.Errorf("install %q: %w", destination, err)
	}
	return nil
}

func findQEMU() (string, error) {
	for _, base := range []string{os.Getenv("ProgramFiles"), os.Getenv("ProgramFiles(x86)")} {
		if base == "" {
			continue
		}
		candidate := filepath.Join(base, "qemu", "qemu-system-x86_64.exe")
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		} else if !os.IsNotExist(err) {
			return "", fmt.Errorf("check for QEMU at %q: %w", candidate, err)
		}
	}

	if qemuPath, err := exec.LookPath("qemu-system-x86_64.exe"); err == nil {
		return qemuPath, nil
	}

	return "", fmt.Errorf("QEMU was not detected; install QEMU for Windows and try again")
}
