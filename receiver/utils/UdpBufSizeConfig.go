package utils

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
)

func UdpBufSizeConfig(bufferSize int) error {
	switch os := runtime.GOOS; os {
	case "linux":
		if err := configureLinux(bufferSize); err != nil {
			return err
		}
	case "darwin":
		if err := configureMacOS(bufferSize); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported OS: %s", os)
	}
	return nil
}

func configureLinux(bufferSize int) error {
	fmt.Println("Running on Linux")
	cmdStr := fmt.Sprintf("sudo sysctl -w net.core.rmem_max=%d && sudo sysctl -w net.core.rmem_default=%d && sudo sysctl -w net.core.wmem_max=%d && sudo sysctl -w net.core.wmem_default=%d", bufferSize, bufferSize, bufferSize, bufferSize)
	if err := executeCommand("bash", "-c", cmdStr); err != nil {
		return err
	}
	return nil
}

func configureMacOS(bufferSize int) error {
	fmt.Println("Running on macOS")
	cmdStr := fmt.Sprintf("sudo sysctl -w kern.ipc.maxsockbuf=%d && sudo sysctl -w net.inet.udp.recvspace=%d && sudo sysctl -w net.inet.udp.maxdgram=%d", bufferSize, bufferSize, bufferSize)
	if err := executeCommand("bash", "-c", cmdStr); err != nil {
		return err
	}
	return nil
}

func executeCommand(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error executing command: %v\n", err)
		return err
	}
	fmt.Printf("Output: %s\n", output)
	return nil
}
