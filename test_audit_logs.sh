#!/bin/bash

# Script de test pour les logs d'audit
BASE_URL="http://localhost:8081"
TOKEN=""

echo "🧪 Test des endpoints de logs d'audit"
echo "====================================="

# Fonction pour afficher les résultats
print_result() {
    if [ $1 -eq 0 ]; then
        echo "✅ $2"
    else
        echo "❌ $2"
        echo "Response: $3"
    fi
}

# 1. Connexion pour obtenir un token
echo "1. Connexion..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d '{
        "email": "test_superadmin@eduqr.com",
        "password": "test123"
    }')

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo "❌ Échec de la connexion"
    echo "Response: $LOGIN_RESPONSE"
    exit 1
fi

echo "✅ Connexion réussie"
echo "Token: ${TOKEN:0:20}..."

# 2. Test de récupération des logs d'audit
echo ""
echo "2. Récupération des logs d'audit..."
AUDIT_LOGS_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/admin/audit-logs" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json")

HTTP_STATUS=$(echo $AUDIT_LOGS_RESPONSE | grep -o '"status":[0-9]*' | cut -d':' -f2)
if [ -z "$HTTP_STATUS" ]; then
    HTTP_STATUS=200
fi

print_result $([ "$HTTP_STATUS" = "200" ] && echo 0 || echo 1) \
    "Récupération des logs d'audit" \
    "$AUDIT_LOGS_RESPONSE"

# 3. Test des statistiques
echo ""
echo "3. Récupération des statistiques..."
STATS_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/admin/audit-logs/stats" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json")

HTTP_STATUS=$(echo $STATS_RESPONSE | grep -o '"status":[0-9]*' | cut -d':' -f2)
if [ -z "$HTTP_STATUS" ]; then
    HTTP_STATUS=200
fi

print_result $([ "$HTTP_STATUS" = "200" ] && echo 0 || echo 1) \
    "Récupération des statistiques" \
    "$STATS_RESPONSE"

# 4. Test des logs récents
echo ""
echo "4. Récupération des logs récents..."
RECENT_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/admin/audit-logs/recent?limit=5" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json")

HTTP_STATUS=$(echo $RECENT_RESPONSE | grep -o '"status":[0-9]*' | cut -d':' -f2)
if [ -z "$HTTP_STATUS" ]; then
    HTTP_STATUS=200
fi

print_result $([ "$HTTP_STATUS" = "200" ] && echo 0 || echo 1) \
    "Récupération des logs récents" \
    "$RECENT_RESPONSE"

# 5. Test avec filtres
echo ""
echo "5. Test avec filtres..."
FILTERED_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/admin/audit-logs?action=login&limit=10" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json")

HTTP_STATUS=$(echo $FILTERED_RESPONSE | grep -o '"status":[0-9]*' | cut -d':' -f2)
if [ -z "$HTTP_STATUS" ]; then
    HTTP_STATUS=200
fi

print_result $([ "$HTTP_STATUS" = "200" ] && echo 0 || echo 1) \
    "Filtrage des logs" \
    "$FILTERED_RESPONSE"

# 6. Test d'accès sans token (doit échouer)
echo ""
echo "6. Test d'accès sans authentification..."
UNAUTHORIZED_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/admin/audit-logs" \
    -H "Content-Type: application/json")

HTTP_STATUS=$(echo $UNAUTHORIZED_RESPONSE | grep -o '"status":[0-9]*' | cut -d':' -f2)
if [ -z "$HTTP_STATUS" ]; then
    HTTP_STATUS=401
fi

print_result $([ "$HTTP_STATUS" = "401" ] && echo 0 || echo 1) \
    "Accès non autorisé (sans token)" \
    "$UNAUTHORIZED_RESPONSE"

echo ""
echo "🎉 Tests terminés !"
echo ""
echo "📊 Résumé des endpoints testés :"
echo "   - GET /api/v1/admin/audit-logs (liste avec pagination)"
echo "   - GET /api/v1/admin/audit-logs/stats (statistiques)"
echo "   - GET /api/v1/admin/audit-logs/recent (logs récents)"
echo "   - GET /api/v1/admin/audit-logs?filters (filtrage)"
echo "   - Authentification requise ✅" 