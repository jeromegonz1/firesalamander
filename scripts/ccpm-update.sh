#!/bin/bash
# Script pour mettre à jour CCPM après chaque session

DATE=$(date +%Y%m%d_%H%M)
SESSION_FILE=".claude/memory/session_${DATE}.md"

echo "# Session ${DATE}" > $SESSION_FILE
echo "" >> $SESSION_FILE
echo "## Changements" >> $SESSION_FILE
git diff --name-only >> $SESSION_FILE
echo "" >> $SESSION_FILE
echo "## État actuel" >> $SESSION_FILE
echo "Voir .claude/context/current_state.md" >> $SESSION_FILE

echo "✅ Session logged to $SESSION_FILE"