# Projet Go - FileOps / WebOps / ProcOps

## Procedure d'execution
1. go run .
2. Suivre le menu interactif.

## Fichiers attendus
- config.txt (obligatoire)
- data/input.txt (exemple fourni)
- out/ (sorties generees)

## Fonctionnalites implementees
- Niveau 10: menu, config txt, analyse fichier, analyse multi-fichiers,
  head/tail, filtres, report/index/merge.
- Niveau 12: analyse Wikipedia (goquery).
- Niveau 14: gestion processus (liste, filtre, kill) Windows/macOS.

## Niveau vise
- 14/20

## Description du travail effectue
- Lecture de la config, gestion des erreurs et sorties dans out/.
- FileOps complet (analyse, filtres, head/tail, batch, report/index/merge).
- WebOps (telechargement Wikipedia + stats).
- ProcOps (liste/filtre/kill avec confirmation).

## Notes
- Si un chemin est vide, la valeur par defaut est utilisee.
- Tous les fichiers generes vont dans out/.
