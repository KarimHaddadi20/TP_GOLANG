# Projet Go - FileOps / WebOps / ProcOps / SecureOps

## Objectif
Outil console en Go pour analyser des fichiers texte, des pages Wikipedia,
des processus systeme, et appliquer des actions de securite (lock/read-only).

## Procedure d'execution
1. go run .
   (ou go run . --config config.json)
2. Suivre le menu interactif.

## Fichiers attendus
- config.json (obligatoire)
- data/input.txt (exemple fourni)
- out/ (sorties generees)

## Structure du projet
- main.go : menu principal
- config.go : lecture config JSON
- fileops.go : FileOps (niveau 10)
- webops.go : WebOps (niveau 12)
- procops.go : ProcOps (niveau 14)
- secureops.go : SecureOps (niveau 16)
- utils.go : fonctions utilitaires
- data/ : fichiers d'entree
- out/ : sorties + audit.log

## Configuration (exemple)
config.json :
{
  "default_file": "data/input.txt",
  "base_dir": "data",
  "out_dir": "out",
  "default_ext": ".txt",
  "wiki_lang": "fr",
  "process_top_n": 10
}

## Fonctionnalites par niveau
- Niveau 10: menu, config, choix fichier, stats, filtres,
  head/tail, report/index/merge.
- Niveau 12: Wikipedia via goquery + stats.
- Niveau 14: processus (liste/filtre/kill) Windows/macOS.
- Niveau 16: SecureOps (lock/unlock, read-only, audit.log).

## Utilisation rapide
- A: analyse fichier -> out/filtered*.txt, out/head.txt, out/tail.txt
- B: analyse dossier -> out/report.txt, out/index.txt, out/merged.txt
- C: Wikipedia -> out/wiki_<article>.txt
- D: ProcessOps -> liste/filtre/kill
- E: SecureOps -> out/<nom>.lock + out/audit.log

## Scenario de test rapide
1. Lancer: go run .
2. A -> mot-cle "lorem" -> head/tail 3 -> verifier out/filtered*.txt, head.txt, tail.txt
3. B -> dossier "data" -> verifier out/report.txt, out/index.txt, out/merged.txt
4. C -> article "Go_(langage)" -> verifier out/wiki_Go_(langage).txt
5. E -> lock/unlock + read-only -> verifier out/audit.log

## Niveau vise
- 16/20

## Notes
- Si un champ manque dans la config, la valeur par defaut est utilisee.
- Toutes les sorties sont dans out/.
- Actions sensibles confirmees et loggees dans out/audit.log.
