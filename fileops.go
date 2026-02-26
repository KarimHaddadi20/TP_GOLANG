package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"
)

type FileSummary struct {
	Path       string
	Size       int64
	ModTime    time.Time
	Created    time.Time
	HasCreated bool
	Lines      int
	WordCount  int
	AvgWordLen float64
}

func getFileSummary(path string) (FileSummary, []string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return FileSummary{}, nil, err
	}
	if info.IsDir() {
		return FileSummary{}, nil, fmt.Errorf("chemin invalide (dossier)")
	}

	lines, err := readLines(path)
	if err != nil {
		return FileSummary{}, nil, err
	}
	wordCount, avgLen := wordStats(lines)
	created, hasCreated := fileCreationTime(info)

	return FileSummary{
		Path:       path,
		Size:       info.Size(),
		ModTime:    info.ModTime(),
		Created:    created,
		HasCreated: hasCreated,
		Lines:      len(lines),
		WordCount:  wordCount,
		AvgWordLen: avgLen,
	}, lines, nil
}

func printFileSummary(s FileSummary) {
	fmt.Println("Infos fichier:")
	fmt.Printf("Chemin: %s\n", s.Path)
	fmt.Printf("Taille: %d octets\n", s.Size)
	if s.HasCreated {
		fmt.Printf("Creation: %s\n", formatTime(s.Created))
	} else {
		fmt.Println("Creation: n/a")
	}
	fmt.Printf("Modif: %s\n", formatTime(s.ModTime))
	fmt.Printf("Lignes: %d\n", s.Lines)
	fmt.Printf("Mots: %d | Longueur moyenne: %.2f\n", s.WordCount, s.AvgWordLen)
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func wordStats(lines []string) (int, float64) {
	totalLen := 0
	count := 0
	for _, line := range lines {
		for _, token := range strings.Fields(line) {
			cleaned := cleanToken(token)
			if cleaned == "" || isNumeric(cleaned) {
				continue
			}
			count++
			totalLen += len([]rune(cleaned))
		}
	}
	if count == 0 {
		return 0, 0
	}
	return count, float64(totalLen) / float64(count)
}

func cleanToken(token string) string {
	return strings.TrimFunc(token, func(r rune) bool {
		return unicode.IsPunct(r) || unicode.IsSymbol(r)
	})
}

func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func countLinesWithKeyword(lines []string, keyword string) int {
	keyword = strings.ToLower(strings.TrimSpace(keyword))
	if keyword == "" {
		return 0
	}
	count := 0
	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), keyword) {
			count++
		}
	}
	return count
}

func filterLines(lines []string, keyword string, include bool) []string {
	keyword = strings.ToLower(strings.TrimSpace(keyword))
	if keyword == "" {
		return []string{}
	}
	var out []string
	for _, line := range lines {
		has := strings.Contains(strings.ToLower(line), keyword)
		if include && has {
			out = append(out, line)
		}
		if !include && !has {
			out = append(out, line)
		}
	}
	return out
}

func writeLines(path string, lines []string) error {
	content := strings.Join(lines, "\n")
	if content != "" && !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	return os.WriteFile(path, []byte(content), 0o644)
}

func headLines(lines []string, n int) []string {
	if n <= 0 {
		return []string{}
	}
	if n > len(lines) {
		n = len(lines)
	}
	return lines[:n]
}

func tailLines(lines []string, n int) []string {
	if n <= 0 {
		return []string{}
	}
	if n > len(lines) {
		n = len(lines)
	}
	return lines[len(lines)-n:]
}

func listTxtFiles(dir, ext string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.EqualFold(filepath.Ext(d.Name()), ext) {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func batchAnalyze(dir, ext string) ([]FileSummary, error) {
	files, err := listTxtFiles(dir, ext)
	if err != nil {
		return nil, err
	}
	var summaries []FileSummary
	for _, path := range files {
		summary, _, err := getFileSummary(path)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, summary)
	}
	return summaries, nil
}

func writeReport(path string, summaries []FileSummary) error {
	var b strings.Builder
	b.WriteString("Rapport global\n")
	b.WriteString("================\n\n")
	for _, s := range summaries {
		b.WriteString(fmt.Sprintf("Fichier: %s\n", s.Path))
		b.WriteString(fmt.Sprintf("Taille: %d octets\n", s.Size))
		if s.HasCreated {
			b.WriteString(fmt.Sprintf("Creation: %s\n", formatTime(s.Created)))
		}
		b.WriteString(fmt.Sprintf("Modif: %s\n", formatTime(s.ModTime)))
		b.WriteString(fmt.Sprintf("Lignes: %d\n", s.Lines))
		b.WriteString(fmt.Sprintf("Mots: %d | Moyenne: %.2f\n", s.WordCount, s.AvgWordLen))
		b.WriteString("\n")
	}
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

func writeIndex(path string, summaries []FileSummary) error {
	var b strings.Builder
	b.WriteString("Index\n")
	b.WriteString("=====\n")
	for _, s := range summaries {
		b.WriteString(fmt.Sprintf("%s | %d | %s\n", s.Path, s.Size, formatTime(s.ModTime)))
	}
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

func mergeFiles(dir, ext, outPath string) error {
	files, err := listTxtFiles(dir, ext)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("aucun fichier %s dans %s", ext, dir)
	}
	var b strings.Builder
	for i, path := range files {
		lines, err := readLines(path)
		if err != nil {
			return err
		}
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString(strings.Join(lines, "\n"))
	}
	b.WriteString("\n")
	return os.WriteFile(outPath, []byte(b.String()), 0o644)
}
