#!/bin/bash

# Couleurs pour l'affichage
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}🔍 Test des améliorations QR codes par salle${NC}"

# URL de base
BASE_URL="http://localhost:8080"

# Fonction pour afficher les résultats
print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✅ $2${NC}"
    else
        echo -e "${RED}❌ $2${NC}"
    fi
}

# Test de connexion à l'API
echo -e "\n${YELLOW}1. Test de connexion à l'API...${NC}"
curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/api/v1/health" > /tmp/health_status
HEALTH_STATUS=$(cat /tmp/health_status)
if [ "$HEALTH_STATUS" = "200" ]; then
    print_result 0 "API accessible"
else
    print_result 1 "API non accessible (code: $HEALTH_STATUS)"
    exit 1
fi

# Test avec authentification (utilisateur admin)
echo -e "\n${YELLOW}2. Test avec authentification admin...${NC}"

# Login admin
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d '{
        "email": "admin@eduqr.com",
        "password": "admin123"
    }')

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    print_result 1 "Échec de l'authentification admin"
    echo -e "   Réponse: $LOGIN_RESPONSE"
    exit 1
else
    print_result 0 "Authentification admin réussie"
fi

# Test des salles
echo -e "\n${YELLOW}3. Test des salles...${NC}"
ROOMS_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/admin/rooms" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json")

echo -e "   Réponse salles: $ROOMS_RESPONSE"

# Extraire l'ID de la première salle
ROOM_ID=$(echo $ROOMS_RESPONSE | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)

if [ -z "$ROOM_ID" ]; then
    print_result 1 "Aucune salle trouvée"
else
    print_result 0 "Salle trouvée (ID: $ROOM_ID)"
    
    # Test des cours par salle pour aujourd'hui
    echo -e "\n${YELLOW}4. Test des cours par salle pour aujourd'hui...${NC}"
    TODAY=$(date +%Y-%m-%d)
    COURSES_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/admin/courses/by-room/$ROOM_ID?date=$TODAY" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json")
    
    echo -e "   Réponse cours: $COURSES_RESPONSE"
    
    # Compter le nombre de cours
    COURSE_COUNT=$(echo $COURSES_RESPONSE | grep -o '"id":[0-9]*' | wc -l)
    
    if [ "$COURSE_COUNT" -eq 0 ]; then
        print_result 1 "Aucun cours trouvé pour aujourd'hui"
    else
        print_result 0 "Cours trouvés: $COURSE_COUNT"
        
        # Extraire l'ID du premier cours
        COURSE_ID=$(echo $COURSES_RESPONSE | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)
        
        if [ -n "$COURSE_ID" ]; then
            # Test du QR code pour ce cours
            echo -e "\n${YELLOW}5. Test du QR code pour le cours $COURSE_ID...${NC}"
            QR_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/qr-codes/course/$COURSE_ID" \
                -H "Authorization: Bearer $TOKEN" \
                -H "Content-Type: application/json")
            
            echo -e "   Réponse QR code: $QR_RESPONSE"
            
            # Vérifier si le QR code contient des données
            if echo "$QR_RESPONSE" | grep -q '"qr_code_data"'; then
                print_result 0 "QR code généré avec succès"
                
                # Vérifier la structure des données
                if echo "$QR_RESPONSE" | grep -q '"course_name"'; then
                    print_result 0 "Structure des données QR code correcte"
                else
                    print_result 1 "Structure des données QR code incorrecte"
                fi
            else
                print_result 1 "QR code non généré ou invalide"
            fi
        fi
    fi
fi

# Test des cours futurs (sans date spécifique)
echo -e "\n${YELLOW}6. Test des cours futurs (sans date)...${NC}"
FUTURE_COURSES_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/admin/courses/by-room/$ROOM_ID" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json")

echo -e "   Réponse cours futurs: $FUTURE_COURSES_RESPONSE"

# Nettoyage
rm -f /tmp/health_status

echo -e "\n${GREEN}🎉 Test des améliorations terminé !${NC}"
echo -e "${YELLOW}📋 Résumé des améliorations :${NC}"
echo -e "   • Affichage des cours dans l'ordre chronologique"
echo -e "   • Filtrage des cours futurs uniquement"
echo -e "   • Interface utilisateur améliorée"
echo -e "   • Informations détaillées dans le modal QR code"
echo -e "   • Indicateurs visuels pour les cours en cours" 