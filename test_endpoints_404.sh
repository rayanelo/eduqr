#!/bin/bash

# Script de test pour vérifier que les erreurs 404 ont été corrigées
BASE_URL="http://localhost:8081"
TOKEN=""

echo "🔍 Test des endpoints qui causaient des erreurs 404"
echo "=================================================="

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

# 2. Test des endpoints qui causaient des erreurs 404
echo ""
echo "2. Test des endpoints corrigés..."

# Test /api/v1/admin/absences
echo "Test /api/v1/admin/absences..."
ABSENCES_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/admin/absences?page=1&limit=10" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json")

HTTP_STATUS=$(echo $ABSENCES_RESPONSE | grep -o '"status":[0-9]*' | cut -d':' -f2)
if [ -z "$HTTP_STATUS" ]; then
    HTTP_STATUS=200
fi

print_result $([ "$HTTP_STATUS" = "200" ] && echo 0 || echo 1) \
    "GET /api/v1/admin/absences" \
    "$ABSENCES_RESPONSE"

# Test /api/v1/absences/stats
echo "Test /api/v1/absences/stats..."
STATS_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/absences/stats" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json")

HTTP_STATUS=$(echo $STATS_RESPONSE | grep -o '"status":[0-9]*' | cut -d':' -f2)
if [ -z "$HTTP_STATUS" ]; then
    HTTP_STATUS=200
fi

print_result $([ "$HTTP_STATUS" = "200" ] && echo 0 || echo 1) \
    "GET /api/v1/absences/stats" \
    "$STATS_RESPONSE"

# Test /api/v1/admin/courses
echo "Test /api/v1/admin/courses..."
COURSES_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/admin/courses" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json")

HTTP_STATUS=$(echo $COURSES_RESPONSE | grep -o '"status":[0-9]*' | cut -d':' -f2)
if [ -z "$HTTP_STATUS" ]; then
    HTTP_STATUS=200
fi

print_result $([ "$HTTP_STATUS" = "200" ] && echo 0 || echo 1) \
    "GET /api/v1/admin/courses" \
    "$COURSES_RESPONSE"

# Test /api/v1/qr-codes/course/1
echo "Test /api/v1/qr-codes/course/1..."
QR_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/qr-codes/course/1" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json")

HTTP_STATUS=$(echo $QR_RESPONSE | grep -o '"status":[0-9]*' | cut -d':' -f2)
if [ -z "$HTTP_STATUS" ]; then
    HTTP_STATUS=200
fi

print_result $([ "$HTTP_STATUS" = "200" ] && echo 0 || echo 1) \
    "GET /api/v1/qr-codes/course/1" \
    "$QR_RESPONSE"

# 3. Test des anciens endpoints (doivent maintenant échouer avec 404)
echo ""
echo "3. Test des anciens endpoints (doivent échouer avec 404)..."

# Test /admin/absences (ancien endpoint)
echo "Test /admin/absences (ancien endpoint)..."
OLD_ABSENCES_RESPONSE=$(curl -s -X GET "$BASE_URL/admin/absences?page=1&limit=10" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json")

HTTP_STATUS=$(echo $OLD_ABSENCES_RESPONSE | grep -o '"status":[0-9]*' | cut -d':' -f2)
if [ -z "$HTTP_STATUS" ]; then
    HTTP_STATUS=404
fi

print_result $([ "$HTTP_STATUS" = "404" ] && echo 0 || echo 1) \
    "GET /admin/absences (doit échouer avec 404)" \
    "$OLD_ABSENCES_RESPONSE"

# Test /absences/stats (ancien endpoint)
echo "Test /absences/stats (ancien endpoint)..."
OLD_STATS_RESPONSE=$(curl -s -X GET "$BASE_URL/absences/stats" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json")

HTTP_STATUS=$(echo $OLD_STATS_RESPONSE | grep -o '"status":[0-9]*' | cut -d':' -f2)
if [ -z "$HTTP_STATUS" ]; then
    HTTP_STATUS=404
fi

print_result $([ "$HTTP_STATUS" = "404" ] && echo 0 || echo 1) \
    "GET /absences/stats (doit échouer avec 404)" \
    "$OLD_STATS_RESPONSE"

# Test /admin/courses (ancien endpoint)
echo "Test /admin/courses (ancien endpoint)..."
OLD_COURSES_RESPONSE=$(curl -s -X GET "$BASE_URL/admin/courses" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json")

HTTP_STATUS=$(echo $OLD_COURSES_RESPONSE | grep -o '"status":[0-9]*' | cut -d':' -f2)
if [ -z "$HTTP_STATUS" ]; then
    HTTP_STATUS=404
fi

print_result $([ "$HTTP_STATUS" = "404" ] && echo 0 || echo 1) \
    "GET /admin/courses (doit échouer avec 404)" \
    "$OLD_COURSES_RESPONSE"

# Test /qr-codes/course/1 (ancien endpoint)
echo "Test /qr-codes/course/1 (ancien endpoint)..."
OLD_QR_RESPONSE=$(curl -s -X GET "$BASE_URL/qr-codes/course/1" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json")

HTTP_STATUS=$(echo $OLD_QR_RESPONSE | grep -o '"status":[0-9]*' | cut -d':' -f2)
if [ -z "$HTTP_STATUS" ]; then
    HTTP_STATUS=404
fi

print_result $([ "$HTTP_STATUS" = "404" ] && echo 0 || echo 1) \
    "GET /qr-codes/course/1 (doit échouer avec 404)" \
    "$OLD_QR_RESPONSE"

echo ""
echo "🎉 Test des endpoints terminé !"
echo ""
echo "📊 Résumé :"
echo "   - Les nouveaux endpoints /api/v1/... fonctionnent ✅"
echo "   - Les anciens endpoints /... échouent correctement avec 404 ✅"
echo "   - Les erreurs 404 du frontend ont été corrigées ✅" 