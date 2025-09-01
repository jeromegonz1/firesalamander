#!/bin/bash
# Générer le contexte complet pour Claude

echo "🔥 FIRE SALAMANDER - CONTEXTE COMPLET"
echo "======================================"
echo ""
cat .claude/context/project.md
echo ""
echo "--- ÉTAT ACTUEL ---"
cat .claude/context/current_state.md
echo ""
echo "--- DÉCISIONS ---"
cat .claude/context/decisions.md
echo ""
echo "======================================"
echo "Copier ce contexte dans Claude Code ☝️"