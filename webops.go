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
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "fr-FR,fr;q=0.9,en;q=0.8")

	resp, err := client.Do(req)
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
