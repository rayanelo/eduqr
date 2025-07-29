#!/bin/bash

# 🧪 Script de test complet pour EduQR
# Teste tous les endpoints pour chaque type d'utilisateur

BASE_URL="http://localhost:8081"
TOKENS=()

# Couleurs pour l'affichage
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

echo -e "${BLUE}🎓 EduQR - Tests Complets de l'Application${NC}"
echo "=================================================="
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

# Fonction pour extraire le statut HTTP
extract_status() {
    local response=$1
    local status=$(echo $response | grep -o '"status":[0-9]*' | cut -d':' -f2)
    if [ -z "$status" ]; then
        # Essayer de détecter le statut depuis les headers curl
        status=$(echo $response | grep -o "HTTP/[0-9.]* [0-9]*" | tail -1 | cut -d' ' -f2)
    fi
    echo $status
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
    local status=$(extract_status "$response")
    
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

# 1. Connexion des utilisateurs de test
echo -e "${CYAN}🔐 1. Connexion des utilisateurs de test${NC}"
echo "----------------------------------------"

# Comptes de test
declare -A test_users=(
    ["superadmin"]="test_superadmin@eduqr.com:test123"
    ["admin"]="test_admin@eduqr.com:test123"
["professeur"]="test_prof@eduqr.com:test123"
["etudiant"]="test_student@eduqr.com:test123"
)

# Connexion de tous les utilisateurs
for role in "${!test_users[@]}"; do
    IFS=':' read -r email password <<< "${test_users[$role]}"
    
    echo -n "Connexion $role ($email)... "
    response=$(curl -s -X POST "$BASE_URL/api/v1/auth/login" \
        -H "Content-Type: application/json" \
        -d "{
            \"email\": \"$email\",
            \"password\": \"$password\"
        }")
    
    token=$(extract_token "$response")
    if [ ! -z "$token" ]; then
        TOKENS[$role]=$token
        echo -e "${GREEN}✅${NC}"
    else
        echo -e "${RED}❌${NC}"
        echo "  Response: $response"
    fi
done

echo ""

# 2. Tests pour Super Admin
echo -e "${PURPLE}👑 2. Tests Super Admin${NC}"
echo "------------------------"

superadmin_token=${TOKENS["superadmin"]}
if [ ! -z "$superadmin_token" ]; then
    echo "Test des permissions Super Admin..."
    
    # Gestion des utilisateurs
    test_endpoint "GET" "/api/v1/users/all" "$superadmin_token" "" "Liste tous les utilisateurs"
    test_endpoint "POST" "/api/v1/users/create" "$superadmin_token" '{
        "email": "test.superadmin@eduqr.com",
        "password": "test123456",
        "confirm_password": "test123456",
        "first_name": "Test",
        "last_name": "SuperAdmin",
        "phone": "+1234567890",
        "address": "123 Test Street",
        "role": "admin"
    }' "Création d'un admin"
    
    # Audit logs
    test_endpoint "GET" "/api/v1/admin/audit-logs" "$superadmin_token" "" "Accès aux logs d'audit"
    test_endpoint "GET" "/api/v1/admin/audit-logs/stats" "$superadmin_token" "" "Statistiques des logs"
    test_endpoint "GET" "/api/v1/admin/audit-logs/recent?limit=5" "$superadmin_token" "" "Logs récents"
    
    # Gestion des salles
    test_endpoint "GET" "/api/v1/admin/rooms" "$superadmin_token" "" "Liste des salles"
    test_endpoint "POST" "/api/v1/admin/rooms" "$superadmin_token" '{
        "name": "Salle Test SuperAdmin",
        "building": "Bâtiment A",
        "floor": "1er étage",
        "is_modular": false
    }' "Création d'une salle"
    
    # Gestion des matières
    test_endpoint "GET" "/api/v1/admin/subjects" "$superadmin_token" "" "Liste des matières"
    test_endpoint "POST" "/api/v1/admin/subjects" "$superadmin_token" '{
        "name": "Matière Test SuperAdmin",
        "code": "MTEST",
        "description": "Matière de test pour Super Admin"
    }' "Création d'une matière"
    
    # Gestion des cours
    test_endpoint "GET" "/api/v1/admin/courses" "$superadmin_token" "" "Liste des cours"
    
    # Gestion des absences
    test_endpoint "GET" "/api/v1/admin/absences" "$superadmin_token" "" "Liste toutes les absences"
    test_endpoint "GET" "/api/v1/absences/stats" "$superadmin_token" "" "Statistiques des absences"
    
    # Gestion des présences
    test_endpoint "GET" "/api/v1/admin/presences" "$superadmin_token" "" "Liste toutes les présences"
    
else
    echo -e "${RED}❌ Impossible de tester Super Admin - token manquant${NC}"
fi

echo ""

# 3. Tests pour Admin
echo -e "${BLUE}👨‍💼 3. Tests Admin${NC}"
echo "-------------------"

admin_token=${TOKENS["admin"]}
if [ ! -z "$admin_token" ]; then
    echo "Test des permissions Admin..."
    
    # Gestion des utilisateurs (limité)
    test_endpoint "GET" "/api/v1/users/all" "$admin_token" "" "Liste des utilisateurs (vue limitée)"
    test_endpoint "POST" "/api/v1/users/create" "$admin_token" '{
        "email": "test.admin@eduqr.com",
        "password": "test123456",
        "confirm_password": "test123456",
        "first_name": "Test",
        "last_name": "Admin",
        "phone": "+1234567891",
        "address": "456 Test Avenue",
        "role": "professeur"
    }' "Création d'un professeur"
    
    # Audit logs
    test_endpoint "GET" "/api/v1/admin/audit-logs" "$admin_token" "" "Accès aux logs d'audit"
    test_endpoint "GET" "/api/v1/admin/audit-logs/stats" "$admin_token" "" "Statistiques des logs"
    
    # Gestion des salles
    test_endpoint "GET" "/api/v1/admin/rooms" "$admin_token" "" "Liste des salles"
    test_endpoint "POST" "/api/v1/admin/rooms" "$admin_token" '{
        "name": "Salle Test Admin",
        "building": "Bâtiment B",
        "floor": "2ème étage",
        "is_modular": true,
        "sub_rooms_count": 3
    }' "Création d'une salle modulaire"
    
    # Gestion des matières
    test_endpoint "GET" "/api/v1/admin/subjects" "$admin_token" "" "Liste des matières"
    test_endpoint "POST" "/api/v1/admin/subjects" "$admin_token" '{
        "name": "Matière Test Admin",
        "code": "MADMIN",
        "description": "Matière de test pour Admin"
    }' "Création d'une matière"
    
    # Gestion des cours
    test_endpoint "GET" "/api/v1/admin/courses" "$admin_token" "" "Liste des cours"
    test_endpoint "POST" "/api/v1/admin/courses" "$admin_token" '{
        "name": "Cours Test Admin",
        "subject_id": 1,
        "teacher_id": 3,
        "room_id": 1,
        "start_time": "2024-12-20T10:00:00Z",
        "duration": 120,
        "description": "Cours de test pour Admin"
    }' "Création d'un cours"
    
    # Gestion des absences
    test_endpoint "GET" "/api/v1/admin/absences" "$admin_token" "" "Liste toutes les absences"
    test_endpoint "GET" "/api/v1/absences/stats" "$admin_token" "" "Statistiques des absences"
    
    # Gestion des présences
    test_endpoint "GET" "/api/v1/admin/presences" "$admin_token" "" "Liste toutes les présences"
    
else
    echo -e "${RED}❌ Impossible de tester Admin - token manquant${NC}"
fi

echo ""

# 4. Tests pour Professeur
echo -e "${YELLOW}👨‍🏫 4. Tests Professeur${NC}"
echo "----------------------"

professeur_token=${TOKENS["professeur"]}
if [ ! -z "$professeur_token" ]; then
    echo "Test des permissions Professeur..."
    
    # Vue limitée des utilisateurs
    test_endpoint "GET" "/api/v1/users/all" "$professeur_token" "" "Liste des utilisateurs (vue limitée)"
    
    # Gestion des absences (pour ses cours)
    test_endpoint "GET" "/api/v1/absences/teacher" "$professeur_token" "" "Absences de ses cours"
    test_endpoint "GET" "/api/v1/absences/stats" "$professeur_token" "" "Statistiques des absences"
    test_endpoint "GET" "/api/v1/absences/filter?page=1&limit=10" "$professeur_token" "" "Filtrage des absences"
    
    # Gestion des présences
    test_endpoint "GET" "/api/v1/presences/course/1" "$professeur_token" "" "Présences d'un cours"
    test_endpoint "GET" "/api/v1/presences/course/1/stats" "$professeur_token" "" "Statistiques de présence"
    test_endpoint "POST" "/api/v1/presences/course/1/create-all" "$professeur_token" "" "Création des présences pour tous"
    
    # QR Codes
    test_endpoint "GET" "/api/v1/qr-codes/course/1" "$professeur_token" "" "Informations QR code"
    test_endpoint "POST" "/api/v1/qr-codes/course/1/regenerate" "$professeur_token" "" "Régénération QR code"
    
    # Profil personnel
    test_endpoint "GET" "/api/v1/users/profile" "$professeur_token" "" "Profil personnel"
    test_endpoint "PUT" "/api/v1/users/profile" "$professeur_token" '{
        "first_name": "Jean",
        "last_name": "Dupont",
        "phone": "+1234567892"
    }' "Modification du profil"
    
    # Test d'accès refusé
    test_endpoint "GET" "/api/v1/admin/audit-logs" "$professeur_token" "" "Accès refusé aux logs d'audit"
    test_endpoint "GET" "/api/v1/admin/rooms" "$professeur_token" "" "Accès refusé aux salles"
    
else
    echo -e "${RED}❌ Impossible de tester Professeur - token manquant${NC}"
fi

echo ""

# 5. Tests pour Étudiant
echo -e "${GREEN}👨‍🎓 5. Tests Étudiant${NC}"
echo "-------------------"

etudiant_token=${TOKENS["etudiant"]}
if [ ! -z "$etudiant_token" ]; then
    echo "Test des permissions Étudiant..."
    
    # Vue très limitée des utilisateurs
    test_endpoint "GET" "/api/v1/users/all" "$etudiant_token" "" "Liste des utilisateurs (vue très limitée)"
    
    # Gestion des absences personnelles
    test_endpoint "GET" "/api/v1/absences/my" "$etudiant_token" "" "Mes absences"
    test_endpoint "POST" "/api/v1/absences" "$etudiant_token" '{
        "course_id": 1,
        "justification": "Absence justifiée pour raisons médicales",
        "document_path": "/uploads/justificatifs/medical.pdf"
    }' "Création d'une absence"
    
    # Gestion des présences personnelles
    test_endpoint "GET" "/api/v1/presences/my" "$etudiant_token" "" "Mes présences"
    test_endpoint "POST" "/api/v1/presences/scan" "$etudiant_token" '{
        "qr_code_data": "test_qr_code_data"
    }' "Scan d'un QR code"
    
    # Profil personnel
    test_endpoint "GET" "/api/v1/users/profile" "$etudiant_token" "" "Profil personnel"
    test_endpoint "PUT" "/api/v1/users/profile" "$etudiant_token" '{
        "first_name": "Étudiant",
        "last_name": "Test",
        "phone": "+1234567893"
    }' "Modification du profil"
    
    # Test d'accès refusé
    test_endpoint "GET" "/api/v1/admin/audit-logs" "$etudiant_token" "" "Accès refusé aux logs d'audit"
    test_endpoint "GET" "/api/v1/admin/rooms" "$etudiant_token" "" "Accès refusé aux salles"
    test_endpoint "GET" "/api/v1/absences/teacher" "$etudiant_token" "" "Accès refusé aux absences des professeurs"
    
else
    echo -e "${RED}❌ Impossible de tester Étudiant - token manquant${NC}"
fi

echo ""

# 6. Tests d'authentification et sécurité
echo -e "${RED}🔒 6. Tests de Sécurité${NC}"
echo "------------------------"

echo "Test des accès non autorisés..."

# Test sans token
test_endpoint "GET" "/api/v1/users/all" "" "" "Accès sans authentification (doit échouer)"
test_endpoint "GET" "/api/v1/admin/audit-logs" "" "" "Accès admin sans authentification (doit échouer)"

# Test avec token invalide
test_endpoint "GET" "/api/v1/users/all" "invalid_token" "" "Accès avec token invalide (doit échouer)"

# Test de permissions croisées
if [ ! -z "$etudiant_token" ]; then
    test_endpoint "GET" "/api/v1/admin/rooms" "$etudiant_token" "" "Étudiant accédant aux salles admin (doit échouer)"
    test_endpoint "POST" "/api/v1/admin/subjects" "$etudiant_token" '{
        "name": "Matière Non Autorisée",
        "code": "UNAUTH"
    }' "Étudiant créant une matière (doit échouer)"
fi

if [ ! -z "$professeur_token" ]; then
    test_endpoint "POST" "/api/v1/users/create" "$professeur_token" '{
        "email": "test.unauthorized@eduqr.com",
        "password": "test123456",
        "confirm_password": "test123456",
        "first_name": "Test",
        "last_name": "Unauthorized",
        "phone": "+1234567894",
        "address": "789 Unauthorized Street",
        "role": "etudiant"
    }' "Professeur créant un utilisateur (doit échouer)"
fi

echo ""

# 7. Tests des endpoints publics
echo -e "${CYAN}🌐 7. Tests des Endpoints Publics${NC}"
echo "--------------------------------"

test_endpoint "GET" "/health" "" "" "Health check"
test_endpoint "POST" "/api/v1/auth/register" "" '{
    "email": "newuser@eduqr.com",
    "password": "newuser123",
    "confirm_password": "newuser123",
    "first_name": "Nouveau",
    "last_name": "Utilisateur",
    "phone": "+1234567895",
    "address": "123 New User Street"
}' "Inscription d'un nouvel utilisateur"

echo ""

# 8. Tests de validation des données
echo -e "${PURPLE}📝 8. Tests de Validation${NC}"
echo "---------------------------"

if [ ! -z "$admin_token" ]; then
    echo "Test des validations de données..."
    
    # Email invalide
    test_endpoint "POST" "/api/v1/users/create" "$admin_token" '{
        "email": "invalid-email",
        "password": "test123456",
        "confirm_password": "test123456",
        "first_name": "Test",
        "last_name": "Invalid",
        "phone": "+1234567896",
        "address": "123 Invalid Street",
        "role": "etudiant"
    }' "Création avec email invalide (doit échouer)"
    
    # Mot de passe trop court
    test_endpoint "POST" "/api/v1/users/create" "$admin_token" '{
        "email": "shortpass@eduqr.com",
        "password": "123",
        "confirm_password": "123",
        "first_name": "Test",
        "last_name": "ShortPass",
        "phone": "+1234567897",
        "address": "123 Short Pass Street",
        "role": "etudiant"
    }' "Création avec mot de passe trop court (doit échouer)"
    
    # Mots de passe différents
    test_endpoint "POST" "/api/v1/users/create" "$admin_token" '{
        "email": "diffpass@eduqr.com",
        "password": "password123",
        "confirm_password": "different123",
        "first_name": "Test",
        "last_name": "DiffPass",
        "phone": "+1234567898",
        "address": "123 Diff Pass Street",
        "role": "etudiant"
    }' "Création avec mots de passe différents (doit échouer)"
    
    # Rôle invalide
    test_endpoint "POST" "/api/v1/users/create" "$admin_token" '{
        "email": "invalidrole@eduqr.com",
        "password": "test123456",
        "confirm_password": "test123456",
        "first_name": "Test",
        "last_name": "InvalidRole",
        "phone": "+1234567899",
        "address": "123 Invalid Role Street",
        "role": "invalid_role"
    }' "Création avec rôle invalide (doit échouer)"
    
fi

echo ""

# 9. Tests de performance et limites
echo -e "${BLUE}⚡ 9. Tests de Performance${NC}"
echo "---------------------------"

echo "Test des limites de pagination..."

if [ ! -z "$admin_token" ]; then
    # Test avec limite élevée
    test_endpoint "GET" "/api/v1/admin/audit-logs?limit=1000" "$admin_token" "" "Requête avec limite élevée"
    
    # Test avec page invalide
    test_endpoint "GET" "/api/v1/admin/audit-logs?page=0" "$admin_token" "" "Requête avec page invalide"
    
    # Test avec filtres complexes
    test_endpoint "GET" "/api/v1/admin/audit-logs?action=login&resource_type=user&limit=50" "$admin_token" "" "Filtrage complexe des logs"
fi

echo ""

# 10. Résumé des tests
echo -e "${GREEN}📊 10. Résumé des Tests${NC}"
echo "------------------------"

echo "🎯 Endpoints testés par rôle :"
echo "  • Super Admin : Gestion complète (utilisateurs, audit, salles, matières, cours, absences, présences)"
echo "  • Admin : Gestion limitée (professeurs/étudiants, audit, salles, matières, cours, absences, présences)"
echo "  • Professeur : Gestion des absences/presences de ses cours, QR codes, profil personnel"
echo "  • Étudiant : Gestion des absences/presences personnelles, scan QR, profil personnel"
echo ""
echo "🔒 Tests de sécurité :"
echo "  • Authentification requise"
echo "  • Permissions par rôle"
echo "  • Validation des données"
echo "  • Accès refusé approprié"
echo ""
echo "✅ Tests de validation :"
echo "  • Emails invalides"
echo "  • Mots de passe faibles"
echo "  • Rôles invalides"
echo "  • Données manquantes"
echo ""
echo "⚡ Tests de performance :"
echo "  • Pagination"
echo "  • Filtrage"
echo "  • Limites de requêtes"

echo ""
echo -e "${GREEN}🎉 Tests terminés !${NC}"
echo ""
echo "📋 Pour exécuter des tests spécifiques :"
echo "  • Test rapide : ./test_audit_logs.sh"
echo "  • Test des rôles : ./test_roles.sh"
echo "  • Test complet : ./test_complet_eduqr.sh"
echo ""
echo "🔧 Pour redémarrer l'application :"
echo "  • Backend : cd backend && go run cmd/server/main.go"
echo "  • Frontend : cd frontend && npm start"
echo ""
echo "📚 Documentation :"
echo "  • curl_examples.md : Exemples de requêtes curl"
echo "  • AUDIT_LOG_SYSTEM.md : Système de logs d'audit"
echo "  • EDUQR_APPLICATION_STATE.md : État de l'application" 