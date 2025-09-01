# Décisions Techniques - Fire Salamander

## Backend: Go
- Raison: Performance pour crawler concurrent
- Alternative écartée: Rust (courbe apprentissage)

## ML: CamemBERT/DistilCamemBERT  
- Raison: Optimisé pour le français
- Mode fast: DistilCamemBERT
- Mode deep: CamemBERT complet

## LLM: Mistral 7B local
- Raison: Coût maîtrisé vs API
- Quantization 4bit pour RAM
- Cache suggestions 3 jours

## Streaming: JSON-RPC
- Raison: Meilleur pour streaming temps réel
- Alternative écartée: REST (pas adapté au streaming)

## Templates: Go natifs + Alpine.js
- Raison: Préservation travail existant
- Performance optimale
- Branding SEPTEO intégré

## Base données: SQLite → PostgreSQL
- MVP: SQLite pour simplicité
- Production: Migration PostgreSQL
- Schema défini dans scripts/init-project.sh