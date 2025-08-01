#!/bin/bash

# Setup Fire Salamander sur Infomaniak
echo "🔥 Configuration Fire Salamander sur Infomaniak..."

APP_NAME="firesalamander"
APP_DIR="$HOME/apps/$APP_NAME"
SERVICE_NAME="fire-salamander"

# Création des répertoires
echo "📁 Création des répertoires..."
mkdir -p $APP_DIR/{logs,data,backups}

# Configuration de la base de données MySQL
echo "🗄️  Configuration MySQL..."
mysql -u root -p << EOF
CREATE DATABASE IF NOT EXISTS firesalamander CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER IF NOT EXISTS 'firesalamander'@'localhost' IDENTIFIED BY 'strong_password_here';
GRANT ALL PRIVILEGES ON firesalamander.* TO 'firesalamander'@'localhost';
FLUSH PRIVILEGES;
EOF

# Variables d'environnement
echo "🔧 Configuration des variables d'environnement..."
cat > $APP_DIR/.env << EOF
ENV=production
DB_NAME=firesalamander
DB_USER=firesalamander
DB_PASS=strong_password_here
OPENAI_API_KEY=your_openai_api_key_here
EOF

# Script de démarrage
echo "🚀 Création du script de démarrage..."
cat > $APP_DIR/start-firesalamander.sh << EOF
#!/bin/bash
cd $APP_DIR
export \$(cat .env | xargs)
nohup ./$APP_NAME > logs/app.log 2>&1 &
echo \$! > $APP_NAME.pid
echo "Fire Salamander démarré (PID: \$(cat $APP_NAME.pid))"
EOF

# Script d'arrêt
cat > $APP_DIR/stop-firesalamander.sh << EOF
#!/bin/bash
cd $APP_DIR
if [ -f $APP_NAME.pid ]; then
    PID=\$(cat $APP_NAME.pid)
    kill \$PID
    rm $APP_NAME.pid
    echo "Fire Salamander arrêté (PID: \$PID)"
else
    echo "Aucun processus Fire Salamander trouvé"
fi
EOF

# Script de redémarrage
cat > $APP_DIR/restart-firesalamander.sh << EOF
#!/bin/bash
cd $APP_DIR
./stop-firesalamander.sh
sleep 2
./start-firesalamander.sh
EOF

# Permissions
chmod +x $APP_DIR/*.sh

# Configuration du service systemd (optionnel)
echo "⚙️  Configuration du service systemd..."
sudo tee /etc/systemd/system/$SERVICE_NAME.service > /dev/null << EOF
[Unit]
Description=Fire Salamander SEO Analyzer
After=network.target mysql.service

[Service]
Type=simple
User=$USER
WorkingDirectory=$APP_DIR
EnvironmentFile=$APP_DIR/.env
ExecStart=$APP_DIR/$APP_NAME
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

# Activation du service
sudo systemctl daemon-reload
sudo systemctl enable $SERVICE_NAME

# Configuration de la rotation des logs
echo "📋 Configuration de la rotation des logs..."
sudo tee /etc/logrotate.d/$SERVICE_NAME > /dev/null << EOF
$APP_DIR/logs/*.log {
    daily
    missingok
    rotate 7
    compress
    delaycompress
    notifempty
    create 644 $USER $USER
    postrotate
        systemctl reload $SERVICE_NAME > /dev/null 2>&1 || true
    endscript
}
EOF

# Backup automatique quotidien
echo "💾 Configuration du backup automatique..."
(crontab -l 2>/dev/null; echo "0 2 * * * $APP_DIR/backup-db.sh") | crontab -

cat > $APP_DIR/backup-db.sh << EOF
#!/bin/bash
cd $APP_DIR
DATE=\$(date +%Y%m%d_%H%M%S)
mysqldump -u firesalamander -pstrong_password_here firesalamander > backups/firesalamander_\$DATE.sql
# Garde seulement les 7 derniers backups
find backups/ -name "firesalamander_*.sql" -mtime +7 -delete
EOF

chmod +x $APP_DIR/backup-db.sh

echo "✅ Configuration terminée!"
echo ""
echo "📋 Étapes suivantes:"
echo "1. Modifier le mot de passe MySQL dans .env"
echo "2. Ajouter votre clé OpenAI dans .env"
echo "3. Déployer l'application avec deploy.sh"
echo "4. Démarrer le service: sudo systemctl start $SERVICE_NAME"
echo ""
echo "🔧 Commandes utiles:"
echo "   Statut: sudo systemctl status $SERVICE_NAME"
echo "   Logs: journalctl -u $SERVICE_NAME -f"
echo "   Redémarrage: sudo systemctl restart $SERVICE_NAME"