package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var stdinReader = bufio.NewReader(os.Stdin)

func readLine(prompt string) string {
	fmt.Print(prompt)
	line, _ := stdinReader.ReadString('\n')
	return strings.TrimSpace(line)
}

func confirmAction(prompt string) bool {
	answer := strings.ToLower(readLine(prompt + " (yes/no): "))
	return answer == "yes"
}

func readNonEmpty(prompt string) string {
	for {
		line := readLine(prompt)
		if line != "" {
			return line
		}
		fmt.Println("Valeur vide, recommencez.")
	}
}

func readIntWithDefault(prompt string, def int) int {
	line := readLine(fmt.Sprintf("%s [%d]: ", prompt, def))
	if line == "" {
		return def
	}
	n, err := strconv.Atoi(strings.TrimSpace(line))
	if err != nil || n < 0 {
		return def
	}
	return n
}

func askPath(prompt, def string) string {
	if def != "" {
		prompt = fmt.Sprintf("%s [%s]: ", prompt, def)
	} else {
		prompt = prompt + ": "
	}
	line := readLine(prompt)
	if line == "" {
		return def
	}
	return filepath.Clean(line)
}

func ensureDir(path string) error {
	return os.MkdirAll(path, 0o755)
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func sanitizeFileName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "article"
	}
	var b strings.Builder
	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' || r == '(' || r == ')' {
			b.WriteRune(r)
		} else {
			b.WriteByte('_')
		}
	}
	return b.String()
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return "n/a"
	}
	return t.Format("2006-01-02 15:04:05")
}

func writeAuditLog(outDir, message string) {
	if outDir == "" {
		return
	}
	if err := ensureDir(outDir); err != nil {
		return
	}
	path := filepath.Join(outDir, "audit.log")
	line := fmt.Sprintf("%s | %s\n", time.Now().Format("2006-01-02 15:04:05"), message)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return
	}
	defer file.Close()
	_, _ = file.WriteString(line)
}
