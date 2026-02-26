package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func runSecureOps(cfg Config) {
	for {
		fmt.Println("=== SecureOps ===")
		fmt.Println("1) Verrouiller un fichier")
		fmt.Println("2) Deverrouiller un fichier")
		fmt.Println("3) Rendre un fichier read-only")
		fmt.Println("0) Retour")
		choice := readLine("Choix: ")
		switch choice {
		case "1":
			path := askPath("Fichier a verrouiller", cfg.DefaultFile)
			if !fileExists(path) {
				fmt.Println("Fichier introuvable ou non valide.")
				break
			}
			if !confirmAction("Confirmer verrouillage") {
				fmt.Println("Annule.")
				break
			}
			lockPath, err := createLock(path, cfg.OutDir)
			if err != nil {
				fmt.Printf("Erreur lock: %v\n", err)
				writeAuditLog(cfg.OutDir, fmt.Sprintf("LOCK FAIL file=%s err=%v", path, err))
				break
			}
			fmt.Printf("Lock cree: %s\n", lockPath)
			writeAuditLog(cfg.OutDir, fmt.Sprintf("LOCK OK file=%s lock=%s", path, lockPath))
		case "2":
			path := askPath("Fichier a deverrouiller", cfg.DefaultFile)
			lockPath := lockPathForFile(path, cfg.OutDir)
			if !fileExists(lockPath) {
				fmt.Println("Aucun lock trouve.")
				break
			}
			if !confirmAction("Confirmer deverrouillage") {
				fmt.Println("Annule.")
				break
			}
			if err := os.Remove(lockPath); err != nil {
				fmt.Printf("Erreur suppression lock: %v\n", err)
				writeAuditLog(cfg.OutDir, fmt.Sprintf("UNLOCK FAIL file=%s err=%v", path, err))
				break
			}
			fmt.Printf("Lock supprime: %s\n", lockPath)
			writeAuditLog(cfg.OutDir, fmt.Sprintf("UNLOCK OK file=%s lock=%s", path, lockPath))
		case "3":
			path := askPath("Fichier a rendre read-only", cfg.DefaultFile)
			if !fileExists(path) {
				fmt.Println("Fichier introuvable ou non valide.")
				break
			}
			if !confirmAction("Confirmer read-only") {
				fmt.Println("Annule.")
				break
			}
			if err := setReadOnly(path); err != nil {
				fmt.Printf("Erreur read-only: %v\n", err)
				writeAuditLog(cfg.OutDir, fmt.Sprintf("READONLY FAIL file=%s err=%v", path, err))
				break
			}
			fmt.Println("Read-only OK.")
			writeAuditLog(cfg.OutDir, fmt.Sprintf("READONLY OK file=%s", path))
		case "0":
			return
		default:
			fmt.Println("Choix invalide.")
		}
		fmt.Println()
	}
}

func createLock(filePath, outDir string) (string, error) {
	if err := ensureDir(outDir); err != nil {
		return "", err
	}
	lockPath := lockPathForFile(filePath, outDir)
	if fileExists(lockPath) {
		return "", fmt.Errorf("deja verrouille (%s)", lockPath)
	}
	if err := os.WriteFile(lockPath, []byte("locked\n"), 0o644); err != nil {
		return "", err
	}
	return lockPath, nil
}

func lockPathForFile(filePath, outDir string) string {
	base := filepath.Base(filePath)
	name := strings.TrimSuffix(base, filepath.Ext(base))
	if name == "" {
		name = "file"
	}
	name = sanitizeFileName(name)
	return filepath.Join(outDir, name+".lock")
}

func setReadOnly(path string) error {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("attrib", "+R", path)
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("%v: %s", err, strings.TrimSpace(string(out)))
		}
		return nil
	default:
		info, err := os.Stat(path)
		if err != nil {
			return err
		}
		mode := info.Mode().Perm()
		newMode := mode &^ 0o222
		return os.Chmod(path, newMode)
	}
}
