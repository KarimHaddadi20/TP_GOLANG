package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"strings"
)

func main() {
	configPath := flag.String("config", "config.json", "chemin vers le fichier de config")
	flag.Parse()

	cfg, err := loadConfig(*configPath)
	if err != nil {
		fmt.Printf("Config introuvable (%s), utilisation des valeurs par defaut.\n", *configPath)
	}

	if err := ensureDir(cfg.OutDir); err != nil {
		fmt.Printf("Erreur creation dossier out: %v\n", err)
		return
	}

	currentFile := cfg.DefaultFile
	if !fileExists(currentFile) {
		fmt.Printf("Fichier par defaut introuvable: %s\n", currentFile)
	}

	for {
		printMenu(currentFile)
		choice := strings.ToUpper(readLine("Votre choix: "))
		switch choice {
		case "1":
			path := askPath("Chemin du fichier courant", currentFile)
			if !fileExists(path) {
				fmt.Println("Fichier introuvable.")
				break
			}
			currentFile = path
			fmt.Printf("Fichier courant: %s\n", currentFile)
		case "A":
			runFileAnalysis(cfg, &currentFile)
		case "B":
			runMultiAnalysis(cfg)
		case "C":
			runWikiAnalysis(cfg)
		case "D":
			runProcOps(cfg)
		case "E":
			runSecureOps(cfg)
		case "Q":
			fmt.Println("Fin du programme.")
			return
		default:
			fmt.Println("Choix invalide.")
		}
		fmt.Println()
	}
}

func printMenu(currentFile string) {
	fmt.Println("=== Menu principal ===")
	fmt.Printf("Fichier courant: %s\n", currentFile)
	fmt.Println("1) Choisir fichier courant")
	fmt.Println("A) Analyse sur fichier courant")
	fmt.Println("B) Analyse multi-fichiers")
	fmt.Println("C) Analyser une page Wikipedia")
	fmt.Println("D) ProcessOps")
	fmt.Println("E) SecureOps")
	fmt.Println("Q) Quitter")
}

func runFileAnalysis(cfg Config, currentFile *string) {
	path := askPath("Chemin du fichier", *currentFile)
	if !fileExists(path) {
		fmt.Println("Fichier introuvable ou non valide.")
		return
	}
	*currentFile = path

	summary, lines, err := getFileSummary(path)
	if err != nil {
		fmt.Printf("Erreur lecture fichier: %v\n", err)
		return
	}

	printFileSummary(summary)

	keyword := readNonEmpty("Mot-cle pour filtres: ")
	count := countLinesWithKeyword(lines, keyword)
	fmt.Printf("Lignes contenant \"%s\": %d\n", keyword, count)

	filteredPath := filepath.Join(cfg.OutDir, "filtered"+cfg.DefaultExt)
	filteredNotPath := filepath.Join(cfg.OutDir, "filtered_not"+cfg.DefaultExt)

	if err := writeLines(filteredPath, filterLines(lines, keyword, true)); err != nil {
		fmt.Printf("Erreur ecriture %s: %v\n", filteredPath, err)
	} else {
		fmt.Printf("OK: %s\n", filteredPath)
	}

	if err := writeLines(filteredNotPath, filterLines(lines, keyword, false)); err != nil {
		fmt.Printf("Erreur ecriture %s: %v\n", filteredNotPath, err)
	} else {
		fmt.Printf("OK: %s\n", filteredNotPath)
	}

	nHead := readIntWithDefault("N pour head", 5)
	nTail := readIntWithDefault("N pour tail", 5)

	headPath := filepath.Join(cfg.OutDir, "head"+cfg.DefaultExt)
	tailPath := filepath.Join(cfg.OutDir, "tail"+cfg.DefaultExt)

	if err := writeLines(headPath, headLines(lines, nHead)); err != nil {
		fmt.Printf("Erreur ecriture %s: %v\n", headPath, err)
	} else {
		fmt.Printf("OK: %s\n", headPath)
	}

	if err := writeLines(tailPath, tailLines(lines, nTail)); err != nil {
		fmt.Printf("Erreur ecriture %s: %v\n", tailPath, err)
	} else {
		fmt.Printf("OK: %s\n", tailPath)
	}
}

func runMultiAnalysis(cfg Config) {
	dir := askPath("Repertoire a analyser", cfg.BaseDir)
	if !dirExists(dir) {
		fmt.Println("Repertoire introuvable ou non valide.")
		return
	}

	summaries, err := batchAnalyze(dir, cfg.DefaultExt)
	if err != nil {
		fmt.Printf("Erreur analyse batch: %v\n", err)
		return
	}
	if len(summaries) == 0 {
		fmt.Println("Aucun fichier .txt trouve.")
		return
	}

	fmt.Printf("Fichiers analyses: %d\n", len(summaries))
	for _, s := range summaries {
		fmt.Printf("- %s | lignes: %d | mots: %d\n", s.Path, s.Lines, s.WordCount)
	}

	reportPath := filepath.Join(cfg.OutDir, "report"+cfg.DefaultExt)
	indexPath := filepath.Join(cfg.OutDir, "index"+cfg.DefaultExt)
	mergedPath := filepath.Join(cfg.OutDir, "merged"+cfg.DefaultExt)

	if err := writeReport(reportPath, summaries); err != nil {
		fmt.Printf("Erreur ecriture report: %v\n", err)
	} else {
		fmt.Printf("OK: %s\n", reportPath)
	}

	if err := writeIndex(indexPath, summaries); err != nil {
		fmt.Printf("Erreur ecriture index: %v\n", err)
	} else {
		fmt.Printf("OK: %s\n", indexPath)
	}

	if err := mergeFiles(cfg.BaseDir, cfg.DefaultExt, mergedPath); err != nil {
		fmt.Printf("Erreur fusion: %v\n", err)
	} else {
		fmt.Printf("OK: %s\n", mergedPath)
	}
}
