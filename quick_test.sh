#!/bin/bash

echo "🚀 Test rapide du système de rôles EduQR"
echo "========================================"
echo ""

# Configuration
API_BASE="http://localhost:8081/api/v1"

# 1. Test Super Admin
echo "1️⃣ Test Super Admin (superadmin@eduqr.com)"
echo "-------------------------------------------"

# Connexion Super Admin
echo "🔐 Connexion..."
SUPER_ADMIN_RESPONSE=$(curl -s -X POST "$API_BASE/auth/login" \
    -H "Content-Type: application/json" \
    -d '{
        "email": "superadmin@eduqr.com",
        "password": "superadmin123"
    }')

SUPER_ADMIN_TOKEN=$(echo $SUPER_ADMIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -n "$SUPER_ADMIN_TOKEN" ]; then
    echo "✅ Connexion réussie"
    
    # Lister les utilisateurs
    echo "📋 Liste des utilisateurs..."
    USERS_RESPONSE=$(curl -s -X GET "$API_BASE/users/all" \
        -H "Authorization: Bearer $SUPER_ADMIN_TOKEN")
    
    if echo "$USERS_RESPONSE" | grep -q "users"; then
        echo "✅ Peut voir tous les utilisateurs"
        USER_COUNT=$(echo "$USERS_RESPONSE" | grep -o '"users":\[[^]]*\]' | grep -o '\[.*\]' | jq length 2>/dev/null || echo "N/A")
        echo "   Nombre d'utilisateurs: $USER_COUNT"
    else
        echo "❌ Erreur lors de la récupération des utilisateurs"
    fi
else
    echo "❌ Échec de connexion Super Admin"
fi

echo ""
echo "2️⃣ Test Admin (admin@eduqr.com)"
echo "--------------------------------"

# Connexion Admin
echo "🔐 Connexion..."
ADMIN_RESPONSE=$(curl -s -X POST "$API_BASE/auth/login" \
    -H "Content-Type: application/json" \
    -d '{
        "email": "admin@eduqr.com",
        "password": "admin123"
    }')

ADMIN_TOKEN=$(echo $ADMIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -n "$ADMIN_TOKEN" ]; then
    echo "✅ Connexion réussie"
    
    # Lister les utilisateurs
    echo "📋 Liste des utilisateurs..."
    USERS_RESPONSE=$(curl -s -X GET "$API_BASE/users/all" \
        -H "Authorization: Bearer $ADMIN_TOKEN")
    
    if echo "$USERS_RESPONSE" | grep -q "users"; then
        echo "✅ Peut voir les utilisateurs autorisés"
        USER_COUNT=$(echo "$USERS_RESPONSE" | grep -o '"users":\[[^]]*\]' | grep -o '\[.*\]' | jq length 2>/dev/null || echo "N/A")
        echo "   Nombre d'utilisateurs visibles: $USER_COUNT"
    else
        echo "❌ Erreur lors de la récupération des utilisateurs"
    fi
else
    echo "❌ Échec de connexion Admin"
fi

echo ""
echo "3️⃣ Test Professeur (prof1@eduqr.com)"
echo "------------------------------------"

# Connexion Professeur
echo "🔐 Connexion..."
PROF_RESPONSE=$(curl -s -X POST "$API_BASE/auth/login" \
    -H "Content-Type: application/json" \
    -d '{
        "email": "prof1@eduqr.com",
        "password": "prof123"
    }')

PROF_TOKEN=$(echo $PROF_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -n "$PROF_TOKEN" ]; then
    echo "✅ Connexion réussie"
    
    # Lister les utilisateurs
    echo "📋 Liste des utilisateurs..."
    USERS_RESPONSE=$(curl -s -X GET "$API_BASE/users/all" \
        -H "Authorization: Bearer $PROF_TOKEN")
    
    if echo "$USERS_RESPONSE" | grep -q "users"; then
        echo "✅ Peut voir les utilisateurs autorisés"
        USER_COUNT=$(echo "$USERS_RESPONSE" | grep -o '"users":\[[^]]*\]' | grep -o '\[.*\]' | jq length 2>/dev/null || echo "N/A")
        echo "   Nombre d'utilisateurs visibles: $USER_COUNT"
    else
        echo "❌ Erreur lors de la récupération des utilisateurs"
    fi
else
    echo "❌ Échec de connexion Professeur"
fi

echo ""
echo "4️⃣ Test Étudiant (etudiant1@eduqr.com)"
echo "--------------------------------------"

# Connexion Étudiant
echo "🔐 Connexion..."
ETUDIANT_RESPONSE=$(curl -s -X POST "$API_BASE/auth/login" \
    -H "Content-Type: application/json" \
    -d '{
        "email": "etudiant1@eduqr.com",
        "password": "student123"
    }')

ETUDIANT_TOKEN=$(echo $ETUDIANT_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -n "$ETUDIANT_TOKEN" ]; then
    echo "✅ Connexion réussie"
    
    # Lister les utilisateurs
    echo "📋 Liste des utilisateurs..."
    USERS_RESPONSE=$(curl -s -X GET "$API_BASE/users/all" \
        -H "Authorization: Bearer $ETUDIANT_TOKEN")
    
    if echo "$USERS_RESPONSE" | grep -q "users"; then
        echo "✅ Peut voir les utilisateurs autorisés"
        USER_COUNT=$(echo "$USERS_RESPONSE" | grep -o '"users":\[[^]]*\]' | grep -o '\[.*\]' | jq length 2>/dev/null || echo "N/A")
        echo "   Nombre d'utilisateurs visibles: $USER_COUNT"
    else
        echo "❌ Erreur lors de la récupération des utilisateurs"
    fi
else
    echo "❌ Échec de connexion Étudiant"
fi

echo ""
echo "🎉 Test rapide terminé !"
echo ""
echo "📊 Comptes de test créés:"
echo "=========================="
echo "Super Admin: superadmin@eduqr.com / superadmin123"
echo "Admin: admin@eduqr.com / admin123"
echo "Professeur: prof1@eduqr.com / prof123"
echo "Étudiant: etudiant1@eduqr.com / student123"
echo ""
echo "🔧 Pour tester manuellement:"
echo "curl -X POST http://localhost:8081/api/v1/auth/login \\"
echo "  -H 'Content-Type: application/json' \\"
echo "  -d '{\"email\": \"superadmin@eduqr.com\", \"password\": \"superadmin123\"}'" 