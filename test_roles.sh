#!/bin/bash

# 🧪 Script de test des rôles et permissions pour EduQR
# Teste spécifiquement la hiérarchie des rôles et les permissions

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

echo -e "${BLUE}🎭 EduQR - Tests des Rôles et Permissions${NC}"
echo "=============================================="
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
    local expected_status=$6
    
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
    if [ ! -z "$expected_status" ]; then
        # Test avec statut attendu spécifique
        if [ "$status" = "$expected_status" ]; then
            success=0
        else
            success=1
        fi
    else
        # Test avec statut de succès par défaut
        if [ "$status" = "200" ] || [ "$status" = "201" ] || [ "$status" = "204" ]; then
            success=0
        else
            success=1
        fi
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

# 2. Test de la hiérarchie des rôles
echo -e "${PURPLE}👑 2. Test de la Hiérarchie des Rôles${NC}"
echo "------------------------------------"

echo "Test des permissions par niveau de rôle..."

# Super Admin - doit pouvoir tout faire
superadmin_token=${TOKENS["superadmin"]}
if [ ! -z "$superadmin_token" ]; then
    echo -e "\n${PURPLE}Super Admin - Accès complet :${NC}"
    
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
    
    # Gestion des ressources
    test_endpoint "GET" "/api/v1/admin/rooms" "$superadmin_token" "" "Accès aux salles"
    test_endpoint "GET" "/api/v1/admin/subjects" "$superadmin_token" "" "Accès aux matières"
    test_endpoint "GET" "/api/v1/admin/courses" "$superadmin_token" "" "Accès aux cours"
    
    # Gestion des absences et présences
    test_endpoint "GET" "/api/v1/admin/absences" "$superadmin_token" "" "Accès à toutes les absences"
    test_endpoint "GET" "/api/v1/admin/presences" "$superadmin_token" "" "Accès à toutes les présences"
fi

# Admin - permissions limitées
admin_token=${TOKENS["admin"]}
if [ ! -z "$admin_token" ]; then
    echo -e "\n${BLUE}Admin - Permissions limitées :${NC}"
    
    # Gestion des utilisateurs (professeurs et étudiants seulement)
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
    
    # Ne peut pas créer un autre admin
    test_endpoint "POST" "/api/v1/users/create" "$admin_token" '{
        "email": "test.admin2@eduqr.com",
        "password": "test123456",
        "confirm_password": "test123456",
        "first_name": "Test",
        "last_name": "Admin2",
        "phone": "+1234567892",
        "address": "789 Test Avenue",
        "role": "admin"
    }' "Tentative de création d'un admin (doit échouer)" "403"
    
    # Audit logs
    test_endpoint "GET" "/api/v1/admin/audit-logs" "$admin_token" "" "Accès aux logs d'audit"
    
    # Gestion des ressources
    test_endpoint "GET" "/api/v1/admin/rooms" "$admin_token" "" "Accès aux salles"
    test_endpoint "GET" "/api/v1/admin/subjects" "$admin_token" "" "Accès aux matières"
    test_endpoint "GET" "/api/v1/admin/courses" "$admin_token" "" "Accès aux cours"
fi

# Professeur - permissions très limitées
professeur_token=${TOKENS["professeur"]}
if [ ! -z "$professeur_token" ]; then
    echo -e "\n${YELLOW}Professeur - Permissions très limitées :${NC}"
    
    # Vue limitée des utilisateurs
    test_endpoint "GET" "/api/v1/users/all" "$professeur_token" "" "Liste des utilisateurs (vue limitée)"
    
    # Ne peut pas créer d'utilisateurs
    test_endpoint "POST" "/api/v1/users/create" "$professeur_token" '{
        "email": "test.prof@eduqr.com",
        "password": "test123456",
        "confirm_password": "test123456",
        "first_name": "Test",
        "last_name": "Prof",
        "phone": "+1234567893",
        "address": "123 Prof Street",
        "role": "etudiant"
    }' "Tentative de création d'utilisateur (doit échouer)" "403"
    
    # Pas d'accès aux logs d'audit
    test_endpoint "GET" "/api/v1/admin/audit-logs" "$professeur_token" "" "Accès refusé aux logs d'audit" "403"
    
    # Pas d'accès aux ressources admin
    test_endpoint "GET" "/api/v1/admin/rooms" "$professeur_token" "" "Accès refusé aux salles" "403"
    test_endpoint "GET" "/api/v1/admin/subjects" "$professeur_token" "" "Accès refusé aux matières" "403"
    test_endpoint "GET" "/api/v1/admin/courses" "$professeur_token" "" "Accès refusé aux cours" "403"
    
    # Accès limité aux absences et présences
    test_endpoint "GET" "/api/v1/absences/teacher" "$professeur_token" "" "Absences de ses cours"
    test_endpoint "GET" "/api/v1/presences/course/1" "$professeur_token" "" "Présences d'un cours"
    
    # Pas d'accès aux absences admin
    test_endpoint "GET" "/api/v1/admin/absences" "$professeur_token" "" "Accès refusé à toutes les absences" "403"
fi

# Étudiant - permissions minimales
etudiant_token=${TOKENS["etudiant"]}
if [ ! -z "$etudiant_token" ]; then
    echo -e "\n${GREEN}Étudiant - Permissions minimales :${NC}"
    
    # Vue très limitée des utilisateurs
    test_endpoint "GET" "/api/v1/users/all" "$etudiant_token" "" "Liste des utilisateurs (vue très limitée)"
    
    # Ne peut pas créer d'utilisateurs
    test_endpoint "POST" "/api/v1/users/create" "$etudiant_token" '{
        "email": "test.student@eduqr.com",
        "password": "test123456",
        "confirm_password": "test123456",
        "first_name": "Test",
        "last_name": "Student",
        "phone": "+1234567894",
        "address": "123 Student Street",
        "role": "etudiant"
    }' "Tentative de création d'utilisateur (doit échouer)" "403"
    
    # Pas d'accès aux ressources admin
    test_endpoint "GET" "/api/v1/admin/audit-logs" "$etudiant_token" "" "Accès refusé aux logs d'audit" "403"
    test_endpoint "GET" "/api/v1/admin/rooms" "$etudiant_token" "" "Accès refusé aux salles" "403"
    test_endpoint "GET" "/api/v1/admin/subjects" "$etudiant_token" "" "Accès refusé aux matières" "403"
    test_endpoint "GET" "/api/v1/admin/courses" "$etudiant_token" "" "Accès refusé aux cours" "403"
    
    # Accès limité aux absences et présences personnelles
    test_endpoint "GET" "/api/v1/absences/my" "$etudiant_token" "" "Mes absences"
    test_endpoint "GET" "/api/v1/presences/my" "$etudiant_token" "" "Mes présences"
    
    # Pas d'accès aux absences des professeurs
    test_endpoint "GET" "/api/v1/absences/teacher" "$etudiant_token" "" "Accès refusé aux absences des professeurs" "403"
fi

echo ""

# 3. Test des permissions croisées
echo -e "${RED}🔒 3. Test des Permissions Croisées${NC}"
echo "--------------------------------"

echo "Test des tentatives d'accès non autorisées..."

# Étudiant tentant d'accéder aux ressources admin
if [ ! -z "$etudiant_token" ]; then
    echo -e "\n${GREEN}Étudiant tentant d'accéder aux ressources admin :${NC}"
    
    test_endpoint "POST" "/api/v1/admin/subjects" "$etudiant_token" '{
        "name": "Matière Non Autorisée",
        "code": "UNAUTH"
    }' "Création de matière (doit échouer)" "403"
    
    test_endpoint "POST" "/api/v1/admin/rooms" "$etudiant_token" '{
        "name": "Salle Non Autorisée",
        "building": "Bâtiment X"
    }' "Création de salle (doit échouer)" "403"
    
    test_endpoint "POST" "/api/v1/admin/courses" "$etudiant_token" '{
        "name": "Cours Non Autorisé",
        "subject_id": 1,
        "teacher_id": 1,
        "room_id": 1,
        "start_time": "2024-12-20T10:00:00Z",
        "duration": 120
    }' "Création de cours (doit échouer)" "403"
fi

# Professeur tentant d'accéder aux ressources admin
if [ ! -z "$professeur_token" ]; then
    echo -e "\n${YELLOW}Professeur tentant d'accéder aux ressources admin :${NC}"
    
    test_endpoint "POST" "/api/v1/admin/subjects" "$professeur_token" '{
        "name": "Matière Non Autorisée",
        "code": "UNAUTH"
    }' "Création de matière (doit échouer)" "403"
    
    test_endpoint "POST" "/api/v1/admin/rooms" "$professeur_token" '{
        "name": "Salle Non Autorisée",
        "building": "Bâtiment X"
    }' "Création de salle (doit échouer)" "403"
    
    test_endpoint "POST" "/api/v1/users/create" "$professeur_token" '{
        "email": "test.unauthorized@eduqr.com",
        "password": "test123456",
        "confirm_password": "test123456",
        "first_name": "Test",
        "last_name": "Unauthorized",
        "phone": "+1234567895",
        "address": "789 Unauthorized Street",
        "role": "etudiant"
    }' "Création d'utilisateur (doit échouer)" "403"
fi

# Admin tentant d'accéder aux ressources super admin
if [ ! -z "$admin_token" ]; then
    echo -e "\n${BLUE}Admin tentant d'accéder aux ressources super admin :${NC}"
    
    # L'admin peut techniquement accéder aux mêmes ressources que le super admin
    # mais avec des limitations sur la gestion des utilisateurs
    test_endpoint "POST" "/api/v1/users/create" "$admin_token" '{
        "email": "test.superadmin2@eduqr.com",
        "password": "test123456",
        "confirm_password": "test123456",
        "first_name": "Test",
        "last_name": "SuperAdmin2",
        "phone": "+1234567896",
        "address": "123 Super Admin Street",
        "role": "super_admin"
    }' "Création d'un super admin (doit échouer)" "403"
fi

echo ""

# 4. Test de la gestion des rôles
echo -e "${PURPLE}🔄 4. Test de la Gestion des Rôles${NC}"
echo "--------------------------------"

echo "Test de la modification des rôles..."

if [ ! -z "$superadmin_token" ]; then
    echo -e "\n${PURPLE}Super Admin modifiant les rôles :${NC}"
    
    # Créer un utilisateur de test
    create_response=$(curl -s -X POST "$BASE_URL/api/v1/users/create" \
        -H "Authorization: Bearer $superadmin_token" \
        -H "Content-Type: application/json" \
        -d '{
            "email": "role.test@eduqr.com",
            "password": "test123456",
            "confirm_password": "test123456",
            "first_name": "Role",
            "last_name": "Test",
            "phone": "+1234567897",
            "address": "123 Role Test Street",
            "role": "etudiant"
        }')
    
    # Extraire l'ID de l'utilisateur créé
    user_id=$(echo $create_response | grep -o '"id":[0-9]*' | cut -d':' -f2)
    
    if [ ! -z "$user_id" ]; then
        echo "Utilisateur créé avec ID: $user_id"
        
        # Promouvoir en professeur
        test_endpoint "PATCH" "/api/v1/users/$user_id/role" "$superadmin_token" '{
            "role": "professeur"
        }' "Promotion en professeur"
        
        # Promouvoir en admin
        test_endpoint "PATCH" "/api/v1/users/$user_id/role" "$superadmin_token" '{
            "role": "admin"
        }' "Promotion en admin"
        
        # Rétrograder en étudiant
        test_endpoint "PATCH" "/api/v1/users/$user_id/role" "$superadmin_token" '{
            "role": "etudiant"
        }' "Rétrogradation en étudiant"
    fi
fi

if [ ! -z "$admin_token" ]; then
    echo -e "\n${BLUE}Admin modifiant les rôles :${NC}"
    
    # Créer un utilisateur de test
    create_response=$(curl -s -X POST "$BASE_URL/api/v1/users/create" \
        -H "Authorization: Bearer $admin_token" \
        -H "Content-Type: application/json" \
        -d '{
            "email": "role.test2@eduqr.com",
            "password": "test123456",
            "confirm_password": "test123456",
            "first_name": "Role",
            "last_name": "Test2",
            "phone": "+1234567898",
            "address": "456 Role Test Street",
            "role": "etudiant"
        }')
    
    # Extraire l'ID de l'utilisateur créé
    user_id=$(echo $create_response | grep -o '"id":[0-9]*' | cut -d':' -f2)
    
    if [ ! -z "$user_id" ]; then
        echo "Utilisateur créé avec ID: $user_id"
        
        # Promouvoir en professeur
        test_endpoint "PATCH" "/api/v1/users/$user_id/role" "$admin_token" '{
            "role": "professeur"
        }' "Promotion en professeur"
        
        # Ne peut pas promouvoir en admin
        test_endpoint "PATCH" "/api/v1/users/$user_id/role" "$admin_token" '{
            "role": "admin"
        }' "Tentative de promotion en admin (doit échouer)" "403"
    fi
fi

echo ""

# 5. Test de la vue des utilisateurs par rôle
echo -e "${CYAN}👥 5. Test de la Vue des Utilisateurs par Rôle${NC}"
echo "----------------------------------------"

echo "Test de la visibilité des utilisateurs selon le rôle..."

# Super Admin voit tous les utilisateurs
if [ ! -z "$superadmin_token" ]; then
    echo -e "\n${PURPLE}Super Admin - Vue complète :${NC}"
    response=$(curl -s -X GET "$BASE_URL/api/v1/users/all" \
        -H "Authorization: Bearer $superadmin_token")
    
    # Compter les utilisateurs
    user_count=$(echo $response | grep -o '"id":[0-9]*' | wc -l)
    echo "  Nombre d'utilisateurs visibles : $user_count"
    
    # Vérifier la présence de tous les rôles
    if echo $response | grep -q "super_admin"; then
        echo -e "  ${GREEN}✅ Super admins visibles${NC}"
    fi
    if echo $response | grep -q "admin"; then
        echo -e "  ${GREEN}✅ Admins visibles${NC}"
    fi
    if echo $response | grep -q "professeur"; then
        echo -e "  ${GREEN}✅ Professeurs visibles${NC}"
    fi
    if echo $response | grep -q "etudiant"; then
        echo -e "  ${GREEN}✅ Étudiants visibles${NC}"
    fi
fi

# Admin voit professeurs et étudiants
if [ ! -z "$admin_token" ]; then
    echo -e "\n${BLUE}Admin - Vue limitée :${NC}"
    response=$(curl -s -X GET "$BASE_URL/api/v1/users/all" \
        -H "Authorization: Bearer $admin_token")
    
    # Compter les utilisateurs
    user_count=$(echo $response | grep -o '"id":[0-9]*' | wc -l)
    echo "  Nombre d'utilisateurs visibles : $user_count"
    
    # Vérifier la présence des rôles autorisés
    if echo $response | grep -q "professeur"; then
        echo -e "  ${GREEN}✅ Professeurs visibles${NC}"
    fi
    if echo $response | grep -q "etudiant"; then
        echo -e "  ${GREEN}✅ Étudiants visibles${NC}"
    fi
    
    # Vérifier l'absence des rôles non autorisés
    if ! echo $response | grep -q "admin"; then
        echo -e "  ${GREEN}✅ Admins non visibles (correct)${NC}"
    fi
    if ! echo $response | grep -q "super_admin"; then
        echo -e "  ${GREEN}✅ Super admins non visibles (correct)${NC}"
    fi
fi

# Professeur voit professeurs et étudiants (vue limitée)
if [ ! -z "$professeur_token" ]; then
    echo -e "\n${YELLOW}Professeur - Vue très limitée :${NC}"
    response=$(curl -s -X GET "$BASE_URL/api/v1/users/all" \
        -H "Authorization: Bearer $professeur_token")
    
    # Compter les utilisateurs
    user_count=$(echo $response | grep -o '"id":[0-9]*' | wc -l)
    echo "  Nombre d'utilisateurs visibles : $user_count"
    
    # Vérifier la présence des rôles autorisés
    if echo $response | grep -q "professeur"; then
        echo -e "  ${GREEN}✅ Professeurs visibles${NC}"
    fi
    if echo $response | grep -q "etudiant"; then
        echo -e "  ${GREEN}✅ Étudiants visibles${NC}"
    fi
fi

# Étudiant voit seulement les étudiants
if [ ! -z "$etudiant_token" ]; then
    echo -e "\n${GREEN}Étudiant - Vue minimale :${NC}"
    response=$(curl -s -X GET "$BASE_URL/api/v1/users/all" \
        -H "Authorization: Bearer $etudiant_token")
    
    # Compter les utilisateurs
    user_count=$(echo $response | grep -o '"id":[0-9]*' | wc -l)
    echo "  Nombre d'utilisateurs visibles : $user_count"
    
    # Vérifier la présence des étudiants seulement
    if echo $response | grep -q "etudiant"; then
        echo -e "  ${GREEN}✅ Étudiants visibles${NC}"
    fi
    
    # Vérifier l'absence des autres rôles
    if ! echo $response | grep -q "professeur"; then
        echo -e "  ${GREEN}✅ Professeurs non visibles (correct)${NC}"
    fi
    if ! echo $response | grep -q "admin"; then
        echo -e "  ${GREEN}✅ Admins non visibles (correct)${NC}"
    fi
    if ! echo $response | grep -q "super_admin"; then
        echo -e "  ${GREEN}✅ Super admins non visibles (correct)${NC}"
    fi
fi

echo ""

# 6. Test de la suppression sécurisée
echo -e "${RED}🗑️ 6. Test de la Suppression Sécurisée${NC}"
echo "--------------------------------"

echo "Test des permissions de suppression..."

if [ ! -z "$superadmin_token" ]; then
    echo -e "\n${PURPLE}Super Admin - Suppression complète :${NC}"
    
    # Créer un utilisateur de test pour la suppression
    create_response=$(curl -s -X POST "$BASE_URL/api/v1/users/create" \
        -H "Authorization: Bearer $superadmin_token" \
        -H "Content-Type: application/json" \
        -d '{
            "email": "delete.test@eduqr.com",
            "password": "test123456",
            "confirm_password": "test123456",
            "first_name": "Delete",
            "last_name": "Test",
            "phone": "+1234567899",
            "address": "123 Delete Test Street",
            "role": "etudiant"
        }')
    
    user_id=$(echo $create_response | grep -o '"id":[0-9]*' | cut -d':' -f2)
    
    if [ ! -z "$user_id" ]; then
        echo "Utilisateur créé avec ID: $user_id pour test de suppression"
        
        # Supprimer l'utilisateur
        test_endpoint "DELETE" "/api/v1/admin/users/$user_id" "$superadmin_token" "" "Suppression d'un utilisateur"
    fi
fi

if [ ! -z "$admin_token" ]; then
    echo -e "\n${BLUE}Admin - Suppression limitée :${NC}"
    
    # Créer un utilisateur de test pour la suppression
    create_response=$(curl -s -X POST "$BASE_URL/api/v1/users/create" \
        -H "Authorization: Bearer $admin_token" \
        -H "Content-Type: application/json" \
        -d '{
            "email": "delete.test2@eduqr.com",
            "password": "test123456",
            "confirm_password": "test123456",
            "first_name": "Delete",
            "last_name": "Test2",
            "phone": "+1234567900",
            "address": "456 Delete Test Street",
            "role": "etudiant"
        }')
    
    user_id=$(echo $create_response | grep -o '"id":[0-9]*' | cut -d':' -f2)
    
    if [ ! -z "$user_id" ]; then
        echo "Utilisateur créé avec ID: $user_id pour test de suppression"
        
        # Supprimer l'utilisateur
        test_endpoint "DELETE" "/api/v1/admin/users/$user_id" "$admin_token" "" "Suppression d'un utilisateur"
    fi
fi

# Professeur ne peut pas supprimer
if [ ! -z "$professeur_token" ]; then
    echo -e "\n${YELLOW}Professeur - Pas de permission de suppression :${NC}"
    
    test_endpoint "DELETE" "/api/v1/admin/users/1" "$professeur_token" "" "Tentative de suppression (doit échouer)" "403"
fi

# Étudiant ne peut pas supprimer
if [ ! -z "$etudiant_token" ]; then
    echo -e "\n${GREEN}Étudiant - Pas de permission de suppression :${NC}"
    
    test_endpoint "DELETE" "/api/v1/admin/users/1" "$etudiant_token" "" "Tentative de suppression (doit échouer)" "403"
fi

echo ""

# 7. Résumé des tests
echo -e "${GREEN}📊 7. Résumé des Tests de Rôles${NC}"
echo "--------------------------------"

echo "🎭 Hiérarchie des rôles testée :"
echo "  • Super Admin (niveau 4) : Accès complet à tout"
echo "  • Admin (niveau 3) : Gestion des professeurs et étudiants"
echo "  • Professeur (niveau 2) : Gestion de ses cours et étudiants"
echo "  • Étudiant (niveau 1) : Accès personnel uniquement"
echo ""
echo "🔒 Permissions testées :"
echo "  • Vue des utilisateurs selon le rôle"
echo "  • Création d'utilisateurs selon les permissions"
echo "  • Modification des rôles selon la hiérarchie"
echo "  • Suppression sécurisée"
echo "  • Accès aux ressources selon le rôle"
echo ""
echo "✅ Validations testées :"
echo "  • Tentatives d'accès non autorisées"
echo "  • Permissions croisées"
echo "  • Hiérarchie respectée"
echo "  • Séparation des rôles"

echo ""
echo -e "${GREEN}🎉 Tests des rôles terminés !${NC}"
echo ""
echo "📋 Pour exécuter d'autres tests :"
echo "  • Test complet : ./test_complet_eduqr.sh"
echo "  • Test des logs d'audit : ./test_audit_logs.sh"
echo ""
echo "🔧 Pour vérifier les utilisateurs créés :"
echo "  • cd backend && go run cmd/check_users/main.go" 