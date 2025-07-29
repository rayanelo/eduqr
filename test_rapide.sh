#!/bin/bash

# 🧪 Script de test rapide pour EduQR
# Teste les fonctionnalités principales rapidement

BASE_URL="http://localhost:8081"

# Couleurs pour l'affichage
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

echo -e "${BLUE}⚡ EduQR - Test Rapide${NC}"
echo "=========================="
echo ""

# Fonction pour afficher les résultats
print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "  ${GREEN}✅ $2${NC}"
    else
        echo -e "  ${RED}❌ $2${NC}"
        if [ ! -z "$3" ]; then
            echo -e "    ${YELLOW}Response: $3${NC}"
        fi
    fi
}

# Fonction pour extraire le token de la réponse
extract_token() {
    echo $1 | grep -o '"token":"[^"]*"' | cut -d'"' -f4
}

# Fonction pour tester un endpoint
test_endpoint() {
    local method=$1
    local endpoint=$2
    local token=$3
    local data=$4
    local description=$5
    
    local headers="Content-Type: application/json"
    if [ ! -z "$token" ]; then
        headers="$headers -H \"Authorization: Bearer $token\""
    fi
    
    local curl_cmd="curl -s -X $method \"$BASE_URL$endpoint\" -H \"$headers\""
    if [ ! -z "$data" ]; then
        curl_cmd="$curl_cmd -d '$data'"
    fi
    
    local response=$(eval $curl_cmd)
    local status=$(echo $response | grep -o '"status":[0-9]*' | cut -d':' -f2)
    
    if [ -z "$status" ]; then
        status=200
    fi
    
    local success=0
    if [ "$status" = "200" ] || [ "$status" = "201" ] || [ "$status" = "204" ]; then
        success=0
    else
        success=1
    fi
    
    print_result $success "$description" "$response"
    return $success
}

# 1. Test de santé de l'API
echo -e "${CYAN}🏥 1. Test de santé de l'API${NC}"
echo "---------------------------"

test_endpoint "GET" "/health" "" "" "Health check de l'API"

echo ""

# 2. Test d'authentification
echo -e "${CYAN}🔐 2. Test d'authentification${NC}"
echo "---------------------------"

# Connexion Super Admin
echo "Connexion Super Admin..."
response=$(curl -s -X POST "$BASE_URL/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d '{
        "email": "test_superadmin@eduqr.com",
        "password": "test123"
    }')

token=$(extract_token "$response")
if [ ! -z "$token" ]; then
    echo -e "  ${GREEN}✅ Connexion Super Admin réussie${NC}"
    SUPERADMIN_TOKEN=$token
else
    echo -e "  ${RED}❌ Échec de connexion Super Admin${NC}"
    echo "    Response: $response"
    exit 1
fi

echo ""

# 3. Test des endpoints principaux
echo -e "${CYAN}🎯 3. Test des endpoints principaux${NC}"
echo "--------------------------------"

# Gestion des utilisateurs
test_endpoint "GET" "/api/v1/users/all" "$SUPERADMIN_TOKEN" "" "Liste des utilisateurs"
test_endpoint "GET" "/api/v1/users/profile" "$SUPERADMIN_TOKEN" "" "Profil utilisateur"

# Audit logs
test_endpoint "GET" "/api/v1/admin/audit-logs" "$SUPERADMIN_TOKEN" "" "Logs d'audit"
test_endpoint "GET" "/api/v1/admin/audit-logs/stats" "$SUPERADMIN_TOKEN" "" "Statistiques des logs"

# Gestion des ressources
test_endpoint "GET" "/api/v1/admin/rooms" "$SUPERADMIN_TOKEN" "" "Liste des salles"
test_endpoint "GET" "/api/v1/admin/subjects" "$SUPERADMIN_TOKEN" "" "Liste des matières"
test_endpoint "GET" "/api/v1/admin/courses" "$SUPERADMIN_TOKEN" "" "Liste des cours"

# Gestion des absences et présences
test_endpoint "GET" "/api/v1/admin/absences" "$SUPERADMIN_TOKEN" "" "Liste des absences"
test_endpoint "GET" "/api/v1/absences/stats" "$SUPERADMIN_TOKEN" "" "Statistiques des absences"
test_endpoint "GET" "/api/v1/admin/presences" "$SUPERADMIN_TOKEN" "" "Liste des présences"

echo ""

# 4. Test de création de données
echo -e "${CYAN}➕ 4. Test de création de données${NC}"
echo "--------------------------------"

# Créer une salle
test_endpoint "POST" "/api/v1/admin/rooms" "$SUPERADMIN_TOKEN" '{
    "name": "Salle Test Rapide",
    "building": "Bâtiment Test",
    "floor": "Rez-de-chaussée",
    "is_modular": false
}' "Création d'une salle"

# Créer une matière
test_endpoint "POST" "/api/v1/admin/subjects" "$SUPERADMIN_TOKEN" '{
    "name": "Matière Test Rapide",
    "code": "MTR",
    "description": "Matière de test pour le test rapide"
}' "Création d'une matière"

# Créer un utilisateur
test_endpoint "POST" "/api/v1/users/create" "$SUPERADMIN_TOKEN" '{
    "email": "test.rapide@eduqr.com",
    "password": "test123456",
    "confirm_password": "test123456",
    "first_name": "Test",
    "last_name": "Rapide",
    "phone": "+1234567890",
    "address": "123 Test Rapide Street",
    "role": "etudiant"
}' "Création d'un utilisateur"

echo ""

# 5. Test de sécurité
echo -e "${CYAN}🔒 5. Test de sécurité${NC}"
echo "------------------------"

# Test sans authentification
test_endpoint "GET" "/api/v1/users/all" "" "" "Accès sans authentification (doit échouer)"

# Test avec token invalide
test_endpoint "GET" "/api/v1/users/all" "invalid_token" "" "Accès avec token invalide (doit échouer)"

echo ""

# 6. Test des endpoints publics
echo -e "${CYAN}🌐 6. Test des endpoints publics${NC}"
echo "--------------------------------"

# Inscription
test_endpoint "POST" "/api/v1/auth/register" "" '{
    "email": "newuser.rapide@eduqr.com",
    "password": "newuser123",
    "confirm_password": "newuser123",
    "first_name": "Nouveau",
    "last_name": "Utilisateur",
    "phone": "+1234567891",
    "address": "123 New User Street"
}' "Inscription d'un nouvel utilisateur"

echo ""

# 7. Test de validation
echo -e "${CYAN}📝 7. Test de validation${NC}"
echo "---------------------------"

# Email invalide
test_endpoint "POST" "/api/v1/users/create" "$SUPERADMIN_TOKEN" '{
    "email": "invalid-email",
    "password": "test123456",
    "confirm_password": "test123456",
    "first_name": "Test",
    "last_name": "Invalid",
    "phone": "+1234567892",
    "address": "123 Invalid Street",
    "role": "etudiant"
}' "Création avec email invalide (doit échouer)"

# Mot de passe trop court
test_endpoint "POST" "/api/v1/users/create" "$SUPERADMIN_TOKEN" '{
    "email": "shortpass@eduqr.com",
    "password": "123",
    "confirm_password": "123",
    "first_name": "Test",
    "last_name": "ShortPass",
    "phone": "+1234567893",
    "address": "123 Short Pass Street",
    "role": "etudiant"
}' "Création avec mot de passe trop court (doit échouer)"

echo ""

# 8. Résumé du test rapide
echo -e "${GREEN}📊 8. Résumé du test rapide${NC}"
echo "---------------------------"

echo "✅ Fonctionnalités testées :"
echo "  • Santé de l'API"
echo "  • Authentification"
echo "  • Gestion des utilisateurs"
echo "  • Logs d'audit"
echo "  • Gestion des ressources (salles, matières, cours)"
echo "  • Gestion des absences et présences"
echo "  • Création de données"
echo "  • Sécurité et authentification"
echo "  • Endpoints publics"
echo "  • Validation des données"

echo ""
echo -e "${GREEN}🎉 Test rapide terminé !${NC}"
echo ""
echo "📋 Pour des tests plus complets :"
echo "  • Test complet : ./test_complet_eduqr.sh"
echo "  • Test des rôles : ./test_roles.sh"
echo "  • Test des logs d'audit : ./test_audit_logs.sh"
echo ""
echo "🔧 Pour redémarrer l'application :"
echo "  • Backend : cd backend && go run cmd/server/main.go"
echo "  • Frontend : cd frontend && npm start" 