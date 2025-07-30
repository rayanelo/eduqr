#!/bin/bash

echo "🚀 Démarrage du script de seed complet pour EduQR..."

# Vérifier que nous sommes dans le bon répertoire
if [ ! -f "go.mod" ]; then
    echo "❌ Erreur: Ce script doit être exécuté depuis le répertoire backend/"
    exit 1
fi

# Vérifier que Go est installé
if ! command -v go &> /dev/null; then
    echo "❌ Erreur: Go n'est pas installé ou n'est pas dans le PATH"
    exit 1
fi

echo "📦 Compilation du script de seed..."
go build -o seed_all cmd/seed_all/main.go

if [ $? -ne 0 ]; then
    echo "❌ Erreur lors de la compilation du script de seed"
    exit 1
fi

echo "🎯 Exécution du script de seed..."
./seed_all

if [ $? -ne 0 ]; then
    echo "❌ Erreur lors de l'exécution du script de seed"
    exit 1
fi

echo "🧹 Nettoyage..."
rm -f seed_all

echo "✅ Script de seed terminé avec succès!"
echo ""
echo "📋 Données créées :"
echo "   👥 15 utilisateurs (1 super admin, 1 admin, 3 professeurs, 10 étudiants)"
echo "   📚 6 matières (Mathématiques, Physique, Informatique, Histoire, Anglais, Chimie)"
echo "   🏫 6 salles (Salles A101/A102, B201/B202, Labos C301/C302)"
echo "   📅 15 cours (5 passés + 10 futurs)"
echo "   ✅ Présences simulées pour les cours passés"
echo "   ❌ Absences justifiées simulées"
echo ""
echo "🔑 Comptes de test :"
echo "   Super Admin: superadmin@eduqr.com / password123"
echo "   Admin: admin@eduqr.com / password123"
echo "   Professeur: jean.dupont@eduqr.com / password123"
echo "   Étudiant: alice.bernard@eduqr.com / password123" 