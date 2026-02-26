package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

type ProcessInfo struct {
	PID  int
	Name string
}

func runProcOps(cfg Config) {
	for {
		fmt.Println("=== ProcessOps ===")
		fmt.Println("1) Lister les processus")
		fmt.Println("2) Rechercher/filtrer")
		fmt.Println("3) Kill securise")
		fmt.Println("0) Retour")
		choice := readLine("Choix: ")
		switch choice {
		case "1":
			topN := readIntWithDefault("Top N", cfg.ProcessTopN)
			procs, err := listProcesses(topN)
			if err != nil {
				fmt.Printf("Erreur liste: %v\n", err)
				break
			}
			printProcesses(procs)
		case "2":
			term := readNonEmpty("Mot de recherche: ")
			procs, err := listProcesses(0)
			if err != nil {
				fmt.Printf("Erreur liste: %v\n", err)
				break
			}
			filtered := filterProcesses(procs, term)
			printProcesses(filtered)
		case "3":
			pidStr := readNonEmpty("PID a tuer: ")
			pid, err := strconv.Atoi(pidStr)
			if err != nil {
				fmt.Println("PID invalide.")
				break
			}
			procs, _ := listProcesses(0)
			proc, found := findProcessByPID(procs, pid)
			procName := "unknown"
			if found {
				fmt.Printf("Processus: %d | %s\n", proc.PID, proc.Name)
				procName = proc.Name
			} else {
				fmt.Printf("Processus: %d | (nom inconnu)\n", pid)
			}
			if !confirmAction("Confirmer kill") {
				fmt.Println("Annule.")
				break
			}
			force := strings.ToLower(readLine("Forcer? (y/n): ")) == "y"
			if err := killProcess(pid, force); err != nil {
				fmt.Printf("Erreur kill: %v\n", err)
				writeAuditLog(cfg.OutDir, fmt.Sprintf("KILL FAIL pid=%d name=%s err=%v", pid, procName, err))
			} else {
				fmt.Println("Kill OK.")
				writeAuditLog(cfg.OutDir, fmt.Sprintf("KILL OK pid=%d name=%s force=%t", pid, procName, force))
			}
		case "0":
			return
		default:
			fmt.Println("Choix invalide.")
		}
		fmt.Println()
	}
}

func listProcesses(topN int) ([]ProcessInfo, error) {
	switch runtime.GOOS {
	case "windows":
		return listProcessesWindows(topN)
	case "darwin":
		return listProcessesDarwin(topN)
	default:
		return nil, fmt.Errorf("OS non supporte: %s", runtime.GOOS)
	}
}

func listProcessesWindows(topN int) ([]ProcessInfo, error) {
	cmd := exec.Command("tasklist", "/FO", "CSV")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(bytes.NewReader(out))
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	var procs []ProcessInfo
	for i, rec := range records {
		if i == 0 {
			continue
		}
		if len(rec) < 2 {
			continue
		}
		pid, err := strconv.Atoi(strings.TrimSpace(rec[1]))
		if err != nil {
			continue
		}
		procs = append(procs, ProcessInfo{
			PID:  pid,
			Name: strings.TrimSpace(rec[0]),
		})
		if topN > 0 && len(procs) >= topN {
			break
		}
	}
	return procs, nil
}

func listProcessesDarwin(topN int) ([]ProcessInfo, error) {
	cmd := exec.Command("ps", "-Ao", "pid,comm")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(out), "\n")
	var procs []ProcessInfo
	for i, line := range lines {
		if i == 0 || strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		pid, err := strconv.Atoi(fields[0])
		if err != nil {
			continue
		}
		procs = append(procs, ProcessInfo{
			PID:  pid,
			Name: fields[1],
		})
		if topN > 0 && len(procs) >= topN {
			break
		}
	}
	return procs, nil
}

func filterProcesses(procs []ProcessInfo, term string) []ProcessInfo {
	term = strings.ToLower(strings.TrimSpace(term))
	if term == "" {
		return procs
	}
	var out []ProcessInfo
	for _, p := range procs {
		if strings.Contains(strings.ToLower(p.Name), term) {
			out = append(out, p)
		}
	}
	return out
}

func findProcessByPID(procs []ProcessInfo, pid int) (ProcessInfo, bool) {
	for _, p := range procs {
		if p.PID == pid {
			return p, true
		}
	}
	return ProcessInfo{}, false
}

func printProcesses(procs []ProcessInfo) {
	if len(procs) == 0 {
		fmt.Println("Aucun processus.")
		return
	}
	for _, p := range procs {
		fmt.Printf("%d | %s\n", p.PID, p.Name)
	}
}

func killProcess(pid int, force bool) error {
	switch runtime.GOOS {
	case "windows":
		args := []string{"/PID", strconv.Itoa(pid), "/T"}
		if force {
			args = append(args, "/F")
		}
		cmd := exec.Command("taskkill", args...)
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("%v: %s", err, strings.TrimSpace(string(out)))
		}
		return nil
	case "darwin":
		args := []string{strconv.Itoa(pid)}
		if force {
			args = []string{"-9", strconv.Itoa(pid)}
		}
		cmd := exec.Command("kill", args...)
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("%v: %s", err, strings.TrimSpace(string(out)))
		}
		return nil
	default:
		return fmt.Errorf("OS non supporte: %s", runtime.GOOS)
	}
}
