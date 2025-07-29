#!/bin/bash

# Couleurs pour l'affichage
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}🔍 Test des statistiques des absences${NC}"

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

# Test des statistiques des absences (sans authentification)
echo -e "\n${YELLOW}2. Test des statistiques sans authentification...${NC}"
STATS_RESPONSE=$(curl -s -w "%{http_code}" "$BASE_URL/api/v1/absences/stats" -o /tmp/stats_response)
STATS_CODE=$(echo $STATS_RESPONSE | tail -c 4)
STATS_BODY=$(cat /tmp/stats_response)

if [ "$STATS_CODE" = "401" ]; then
    print_result 0 "Authentification requise (normal)"
    echo -e "   Réponse: $STATS_BODY"
else
    print_result 1 "Code de réponse inattendu: $STATS_CODE"
    echo -e "   Réponse: $STATS_BODY"
fi

# Test avec authentification (utilisateur admin)
echo -e "\n${YELLOW}3. Test des statistiques avec authentification admin...${NC}"

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
else
    print_result 0 "Authentification admin réussie"
    
    # Test des statistiques avec token
    STATS_AUTH_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/absences/stats" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json")
    
    echo -e "   Statistiques: $STATS_AUTH_RESPONSE"
    
    # Vérifier la structure de la réponse
    if echo "$STATS_AUTH_RESPONSE" | grep -q '"total_absences"'; then
        print_result 0 "Structure des statistiques correcte"
    else
        print_result 1 "Structure des statistiques incorrecte"
    fi
fi

# Nettoyage
rm -f /tmp/health_status /tmp/stats_response

echo -e "\n${GREEN}🎉 Test terminé !${NC}" 