#!/bin/bash
# Script pour prendre des captures d'écran de Fire Salamander

echo "📸 Prise de captures d'écran Fire Salamander..."

# Créer le dossier screenshots s'il n'existe pas
mkdir -p screenshots

# Attendre que le serveur soit prêt
sleep 2

# Page d'accueil
echo "Capture page d'accueil..."
screencapture -x screenshots/01_home.png

echo "Pour capturer les autres pages :"
echo "1. Naviguez vers http://localhost:8080/analyze?url=https://example.com"
echo "2. Appuyez sur Cmd+Shift+4 et sélectionnez la zone"
echo "3. Sauvegardez dans le dossier screenshots/"

echo "✅ Screenshots disponibles dans ./screenshots/"