package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"os/user"
	"strings"
)

// Greeting returns a greeting message for the given name.
func Greeting(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}

// SumRange returns the sum of integers from start to end (inclusive).
func SumRange(start, end int) int {
	if start > end {
		return 0
	}
	sum := 0
	for i := start; i <= end; i++ {
		sum += i
	}
	return sum
}

// CalculateChecksum computes a simple XOR checksum over the given byte slice.
func CalculateChecksum(data []byte) byte {
	var cs byte
	for _, b := range data {
		cs ^= b
	}
	return cs
}

// HexToBytes converts a hexadecimal string (without spaces) to a byte slice.
// Each pair of hex characters is treated as one byte.
func HexToBytes(hexStr string) ([]byte, error) {
	if len(hexStr)%2 != 0 {
		return nil, fmt.Errorf("hex string length must be even, got %d", len(hexStr))
	}
	result := make([]byte, len(hexStr)/2)
	for i := 0; i < len(hexStr); i += 2 {
		high := hexCharToByte(hexStr[i])
		low := hexCharToByte(hexStr[i+1])
		if high == 0xFF || low == 0xFF {
			return nil, fmt.Errorf("invalid hex character at position %d", i)
		}
		result[i/2] = (high << 4) | low
	}
	return result, nil
}

// hexCharToByte converts a single hex character (0-9, a-f, A-F) to its value.
// Returns 0xFF for invalid characters.
func hexCharToByte(c byte) byte {
	switch {
	case c >= '0' && c <= '9':
		return c - '0'
	case c >= 'a' && c <= 'f':
		return c - 'a' + 10
	case c >= 'A' && c <= 'F':
		return c - 'A' + 10
	default:
		return 0xFF
	}
}

// checkArduinoCLI проверяет, доступен ли arduino-cli в PATH
func checkArduinoCLI() error {
	_, err := exec.LookPath("arduino-cli")
	return err
}

// checkUserInGroup проверяет, состоит ли текущий пользователь в указанной группе
func checkUserInGroup(groupName string) (bool, error) {
	currentUser, err := user.Current()
	if err != nil {
		return false, err
	}

	cmd := exec.Command("groups", currentUser.Username)
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	groups := strings.Fields(string(output))
	for _, g := range groups {
		if g == groupName {
			return true, nil
		}
	}
	return false, nil
}

// getAvailableBoards возвращает список плат через текстовый парсинг arduino-cli board list
func getAvailableBoards() ([]BoardInfo, error) {
	cmd := exec.Command("arduino-cli", "board", "list")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("ошибка запуска arduino-cli: %v", err)
	}

	lines := strings.Split(out.String(), "\n")
	var boards []BoardInfo
	for _, line := range lines {
		if !strings.Contains(line, "/dev/") && !strings.Contains(line, "COM") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		port := fields[0]
		if !strings.HasPrefix(port, "/dev/") && !strings.HasPrefix(port, "COM") {
			continue
		}
		var fqbn string
		for _, field := range fields {
			if strings.Count(field, ":") >= 2 {
				fqbn = field
				break
			}
		}
		boards = append(boards, BoardInfo{Port: port, FQBN: fqbn})
	}
	return boards, nil
}

// uploadHex выполняет загрузку hex-файла через arduino-cli
func uploadHex(hexPath, portPath, fqbn string) error {
	cmd := exec.Command(
		"arduino-cli",
		"upload",
		"-p", portPath,
		"--fqbn", fqbn,
		"--input-file", hexPath,
	)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%v: %s", err, stderr.String())
	}
	return nil
}
