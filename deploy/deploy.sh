#!/bin/bash

# Fire Salamander Deployment Script
echo "🔥 Déploiement Fire Salamander..."

# Configuration
APP_NAME="firesalamander"
BUILD_DIR="./build"
REMOTE_HOST="${DEPLOY_HOST:-your-server.com}"
REMOTE_USER="${DEPLOY_USER:-user}"
REMOTE_PATH="${DEPLOY_PATH:-~/apps/firesalamander}"

# Vérifications préalables
if [ -z "$DEPLOY_HOST" ] || [ -z "$DEPLOY_USER" ]; then
    echo "❌ Variables d'environnement manquantes:"
    echo "   DEPLOY_HOST: serveur de destination"
    echo "   DEPLOY_USER: utilisateur SSH"
    echo "   DEPLOY_PATH: chemin de déploiement (optionnel)"
    exit 1
fi

# Création du répertoire de build
mkdir -p $BUILD_DIR

# Build pour Linux
echo "🔨 Build de l'application..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $BUILD_DIR/$APP_NAME

if [ $? -ne 0 ]; then
    echo "❌ Échec du build"
    exit 1
fi

# Copie des fichiers de configuration
cp -r config $BUILD_DIR/
cp -r deploy/setup-infomaniak.sh $BUILD_DIR/

# Upload vers Infomaniak
echo "📤 Upload vers le serveur..."
scp -r $BUILD_DIR/* $REMOTE_USER@$REMOTE_HOST:$REMOTE_PATH/

if [ $? -ne 0 ]; then
    echo "❌ Échec de l'upload"
    exit 1
fi

# Restart du service
echo "🔄 Redémarrage du service..."
ssh $REMOTE_USER@$REMOTE_HOST "cd $REMOTE_PATH && ./restart-firesalamander.sh"

if [ $? -eq 0 ]; then
    echo "🦎 Fire Salamander déployé avec succès!"
    echo "🌐 Application disponible sur: http://$REMOTE_HOST"
else
    echo "⚠️  Déploiement terminé mais le redémarrage a échoué"
    echo "   Vérifiez manuellement le service sur le serveur"
fi

# Nettoyage
rm -rf $BUILD_DIR
echo "✅ Nettoyage terminé"