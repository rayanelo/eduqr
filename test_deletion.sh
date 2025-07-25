#!/bin/bash

# Script de test pour les fonctionnalités de suppression sécurisée
echo "=== Test des fonctionnalités de suppression sécurisée ==="

# Configuration
API_BASE="http://localhost:8081"
ADMIN_EMAIL="superadmin@eduqr.com"
ADMIN_PASSWORD="superadmin123"

# Fonction pour obtenir un token
get_token() {
    local response=$(curl -s -X POST "$API_BASE/api/v1/auth/login" \
        -H "Content-Type: application/json" \
        -d "{\"email\":\"$ADMIN_EMAIL\",\"password\":\"$ADMIN_PASSWORD\"}")
    
    echo $response | jq -r '.token'
}

# Fonction pour faire une requête authentifiée
authenticated_request() {
    local method=$1
    local endpoint=$2
    local data=$3
    local token=$4
    
    if [ -n "$data" ]; then
        curl -s -X $method "$API_BASE$endpoint" \
            -H "Authorization: Bearer $token" \
            -H "Content-Type: application/json" \
            -d "$data"
    else
        curl -s -X $method "$API_BASE$endpoint" \
            -H "Authorization: Bearer $token" \
            -H "Content-Type: application/json"
    fi
}

echo "1. Obtention du token d'authentification..."
TOKEN=$(get_token)
if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
    echo "❌ Échec de l'authentification"
    exit 1
fi
echo "✅ Token obtenu avec succès"

echo ""
echo "2. Test de suppression d'un utilisateur avec des cours futurs..."
RESPONSE=$(authenticated_request "DELETE" "/api/v1/users/5" "" "$TOKEN")
echo "Réponse: $RESPONSE"

# Vérifier si la suppression a été bloquée (ce qui est attendu)
if echo "$RESPONSE" | grep -q "future_courses"; then
    echo "✅ Test réussi: Suppression bloquée car l'utilisateur a des cours futurs"
else
    echo "❌ Test échoué: La suppression n'a pas été bloquée comme attendu"
fi

echo ""
echo "3. Test de suppression d'un utilisateur sans cours futurs..."
RESPONSE=$(authenticated_request "DELETE" "/api/v1/users/7" "" "$TOKEN")
echo "Réponse: $RESPONSE"

# Vérifier si la suppression a réussi
if echo "$RESPONSE" | grep -q "success.*true"; then
    echo "✅ Test réussi: Utilisateur supprimé avec succès"
else
    echo "❌ Test échoué: La suppression n'a pas réussi"
fi

echo ""
echo "4. Test de suppression d'une salle avec des cours futurs..."
# D'abord, créons un cours dans une salle
echo "Création d'un cours dans une salle..."
COURSE_DATA='{"name":"Test Course for Room","subject_id":1,"teacher_id":6,"room_id":1,"start_time":"2025-08-15T10:00:00Z","end_time":"2025-08-15T12:00:00Z","duration":120,"description":"Test course for room deletion","is_recurring":false}'
authenticated_request "POST" "/api/v1/admin/courses" "$COURSE_DATA" "$TOKEN" > /dev/null

# Maintenant testons la suppression de la salle
RESPONSE=$(authenticated_request "DELETE" "/api/v1/admin/rooms/1" "" "$TOKEN")
echo "Réponse: $RESPONSE"

if echo "$RESPONSE" | grep -q "future_courses"; then
    echo "✅ Test réussi: Suppression de salle bloquée car elle a des cours futurs"
else
    echo "❌ Test échoué: La suppression de salle n'a pas été bloquée comme attendu"
fi

echo ""
echo "5. Test de suppression d'une matière avec des cours liés..."
RESPONSE=$(authenticated_request "DELETE" "/api/v1/admin/subjects/1" "" "$TOKEN")
echo "Réponse: $RESPONSE"

if echo "$RESPONSE" | grep -q "linked_courses"; then
    echo "✅ Test réussi: Suppression de matière bloquée car elle a des cours liés"
else
    echo "❌ Test échoué: La suppression de matière n'a pas été bloquée comme attendu"
fi

echo ""
echo "=== Résumé des tests ==="
echo "Les tests vérifient que:"
echo "- La suppression d'utilisateurs avec des cours futurs est bloquée"
echo "- La suppression d'utilisateurs sans cours futurs est autorisée"
echo "- La suppression de salles avec des cours futurs est bloquée"
echo "- La suppression de matières avec des cours liés est bloquée" 