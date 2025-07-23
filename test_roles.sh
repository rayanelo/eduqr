#!/bin/bash

# Configuration
API_BASE="http://localhost:8081/api/v1"
SUPER_ADMIN_EMAIL="superadmin@eduqr.com"
SUPER_ADMIN_PASSWORD="superadmin123"
ADMIN_EMAIL="admin@eduqr.com"
ADMIN_PASSWORD="admin123"
PROF_EMAIL="prof1@eduqr.com"
PROF_PASSWORD="prof123"
ETUDIANT_EMAIL="etudiant1@eduqr.com"
ETUDIANT_PASSWORD="student123"

echo "🧪 Test du système de rôles EduQR"
echo "=================================="
echo ""

# Fonction pour obtenir un token
get_token() {
    local email=$1
    local password=$2
    local role_name=$3
    
    echo "🔐 Connexion en tant que $role_name..."
    response=$(curl -s -X POST "$API_BASE/auth/login" \
        -H "Content-Type: application/json" \
        -d "{
            \"email\": \"$email\",
            \"password\": \"$password\"
        }")
    
    token=$(echo $response | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    
    if [ -z "$token" ]; then
        echo "❌ Échec de connexion pour $role_name"
        echo "Response: $response"
        return 1
    else
        echo "✅ Connexion réussie pour $role_name"
        echo "Token: ${token:0:20}..."
        echo ""
    fi
}

# Fonction pour tester les permissions
test_permissions() {
    local token=$1
    local role_name=$2
    
    echo "🔍 Test des permissions pour $role_name"
    echo "----------------------------------------"
    
    # Test 1: Lister tous les utilisateurs
    echo "📋 Test: Lister tous les utilisateurs"
    response=$(curl -s -X GET "$API_BASE/users/all" \
        -H "Authorization: Bearer $token")
    
    if echo "$response" | grep -q "users"; then
        echo "✅ Peut lister les utilisateurs"
        user_count=$(echo "$response" | grep -o '"users":\[[^]]*\]' | grep -o '\[.*\]' | jq length 2>/dev/null || echo "N/A")
        echo "   Nombre d'utilisateurs visibles: $user_count"
    else
        echo "❌ Ne peut pas lister les utilisateurs"
        echo "   Erreur: $(echo "$response" | grep -o '"error":"[^"]*"' | cut -d'"' -f4)"
    fi
    
    # Test 2: Créer un nouvel utilisateur
    echo ""
    echo "➕ Test: Créer un nouvel utilisateur"
    response=$(curl -s -X POST "$API_BASE/users/create" \
        -H "Authorization: Bearer $token" \
        -H "Content-Type: application/json" \
        -d "{
            \"email\": \"test_${role_name,,}@eduqr.com\",
            \"password\": \"test123\",
            \"confirm_password\": \"test123\",
            \"first_name\": \"Test\",
            \"last_name\": \"$role_name\",
            \"phone\": \"+1234567899\",
            \"address\": \"Test Address\",
            \"role\": \"etudiant\"
        }")
    
    if echo "$response" | grep -q "id"; then
        echo "✅ Peut créer un utilisateur"
        user_id=$(echo "$response" | grep -o '"id":[0-9]*' | cut -d':' -f2)
        echo "   ID utilisateur créé: $user_id"
    else
        echo "❌ Ne peut pas créer d'utilisateur"
        echo "   Erreur: $(echo "$response" | grep -o '"error":"[^"]*"' | cut -d'"' -f4)"
    fi
    
    # Test 3: Voir un utilisateur spécifique
    echo ""
    echo "👤 Test: Voir un utilisateur spécifique (ID: 2)"
    response=$(curl -s -X GET "$API_BASE/users/2" \
        -H "Authorization: Bearer $token")
    
    if echo "$response" | grep -q "id"; then
        echo "✅ Peut voir l'utilisateur"
        user_name=$(echo "$response" | grep -o '"first_name":"[^"]*"' | cut -d'"' -f4)
        user_role=$(echo "$response" | grep -o '"role":"[^"]*"' | cut -d'"' -f4)
        echo "   Nom: $user_name, Rôle: $user_role"
    else
        echo "❌ Ne peut pas voir l'utilisateur"
        echo "   Erreur: $(echo "$response" | grep -o '"error":"[^"]*"' | cut -d'"' -f4)"
    fi
    
    # Test 4: Modifier un utilisateur
    echo ""
    echo "✏️ Test: Modifier un utilisateur (ID: 2)"
    response=$(curl -s -X PUT "$API_BASE/users/2" \
        -H "Authorization: Bearer $token" \
        -H "Content-Type: application/json" \
        -d "{
            \"first_name\": \"Modified\",
            \"last_name\": \"User\"
        }")
    
    if echo "$response" | grep -q "id"; then
        echo "✅ Peut modifier l'utilisateur"
    else
        echo "❌ Ne peut pas modifier l'utilisateur"
        echo "   Erreur: $(echo "$response" | grep -o '"error":"[^"]*"' | cut -d'"' -f4)"
    fi
    
    # Test 5: Changer le rôle d'un utilisateur
    echo ""
    echo "🔄 Test: Changer le rôle d'un utilisateur (ID: 2)"
    response=$(curl -s -X PATCH "$API_BASE/users/2/role" \
        -H "Authorization: Bearer $token" \
        -H "Content-Type: application/json" \
        -d "{
            \"role\": \"professeur\"
        }")
    
    if echo "$response" | grep -q "id"; then
        echo "✅ Peut changer le rôle"
        new_role=$(echo "$response" | grep -o '"role":"[^"]*"' | cut -d'"' -f4)
        echo "   Nouveau rôle: $new_role"
    else
        echo "❌ Ne peut pas changer le rôle"
        echo "   Erreur: $(echo "$response" | grep -o '"error":"[^"]*"' | cut -d'"' -f4)"
    fi
    
    # Test 6: Supprimer un utilisateur
    echo ""
    echo "🗑️ Test: Supprimer un utilisateur (ID: 2)"
    response=$(curl -s -X DELETE "$API_BASE/users/2" \
        -H "Authorization: Bearer $token")
    
    if echo "$response" | grep -q "deleted successfully"; then
        echo "✅ Peut supprimer l'utilisateur"
    else
        echo "❌ Ne peut pas supprimer l'utilisateur"
        echo "   Erreur: $(echo "$response" | grep -o '"error":"[^"]*"' | cut -d'"' -f4)"
    fi
    
    echo ""
    echo "=================================="
    echo ""
}

# Test Super Admin
echo "🚀 Test du Super Admin"
echo "======================"
get_token "$SUPER_ADMIN_EMAIL" "$SUPER_ADMIN_PASSWORD" "Super Admin"
if [ $? -eq 0 ]; then
    test_permissions "$token" "Super Admin"
fi

# Test Admin
echo "👨‍💼 Test de l'Admin"
echo "==================="
get_token "$ADMIN_EMAIL" "$ADMIN_PASSWORD" "Admin"
if [ $? -eq 0 ]; then
    test_permissions "$token" "Admin"
fi

# Test Professeur
echo "👨‍🏫 Test du Professeur"
echo "======================"
get_token "$PROF_EMAIL" "$PROF_PASSWORD" "Professeur"
if [ $? -eq 0 ]; then
    test_permissions "$token" "Professeur"
fi

# Test Étudiant
echo "👨‍🎓 Test de l'Étudiant"
echo "======================"
get_token "$ETUDIANT_EMAIL" "$ETUDIANT_PASSWORD" "Étudiant"
if [ $? -eq 0 ]; then
    test_permissions "$token" "Étudiant"
fi

echo "🎉 Tests terminés !"
echo ""
echo "📊 Résumé des permissions:"
echo "=========================="
echo "Super Admin: Toutes les permissions"
echo "Admin: Gestion des Professeurs et Étudiants"
echo "Professeur: Lecture seule des Professeurs et Étudiants"
echo "Étudiant: Lecture seule des autres Étudiants" 