package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func runWikiAnalysis(cfg Config) {
	article := readNonEmpty("Article Wikipedia (ex: Go_(langage)): ")
	paragraphs, err := fetchWikiParagraphs(cfg.WikiLang, article)
	if err != nil {
		fmt.Printf("Erreur Wikipedia: %v\n", err)
		return
	}

	keyword := readLine("Mot-cle (vide => Go): ")
	if strings.TrimSpace(keyword) == "" {
		keyword = "Go"
	}

	wordCount, avgLen := wordStats(paragraphs)
	matched := countLinesWithKeyword(paragraphs, keyword)

	outName := "wiki_" + sanitizeFileName(article) + cfg.DefaultExt
	outPath := filepath.Join(cfg.OutDir, outName)
	if err := writeLines(outPath, paragraphs); err != nil {
		fmt.Printf("Erreur ecriture %s: %v\n", outPath, err)
		return
	}

	fmt.Printf("OK: %s\n", outPath)
	fmt.Printf("Mots: %d | Longueur moyenne: %.2f\n", wordCount, avgLen)
	fmt.Printf("Lignes contenant \"%s\": %d\n", keyword, matched)
}

func fetchWikiParagraphs(lang, article string) ([]string, error) {
	url := fmt.Sprintf("https://%s.wikipedia.org/wiki/%s", lang, article)
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status HTTP %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var paragraphs []string
	doc.Find("#mw-content-text p").Each(func(_ int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			paragraphs = append(paragraphs, text)
		}
	})
	if len(paragraphs) == 0 {
		return nil, fmt.Errorf("aucun paragraphe trouve")
	}
	return paragraphs, nil
}
