# 🧪 Exemples de requêtes curl pour tester le système de rôles EduQR

## 📋 Comptes de test créés

```
Super Admin: superadmin@eduqr.com / superadmin123
Admin: admin@eduqr.com / admin123
Professeur: prof1@eduqr.com / prof123
Étudiant: etudiant1@eduqr.com / student123
```

## 🔐 1. Connexion des utilisateurs

### Super Admin
```bash
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "superadmin@eduqr.com",
    "password": "superadmin123"
  }'
```

### Admin
```bash
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@eduqr.com",
    "password": "admin123"
  }'
```

### Professeur
```bash
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "prof1@eduqr.com",
    "password": "prof123"
  }'
```

### Étudiant
```bash
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "etudiant1@eduqr.com",
    "password": "student123"
  }'
```

## 📋 2. Lister tous les utilisateurs

Remplacez `YOUR_TOKEN` par le token obtenu lors de la connexion.

```bash
curl -X GET http://localhost:8081/api/v1/users/all \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## ➕ 3. Créer un nouvel utilisateur

```bash
curl -X POST http://localhost:8081/api/v1/users/create \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "nouveau@eduqr.com",
    "password": "password123",
    "confirm_password": "password123",
    "first_name": "Nouveau",
    "last_name": "Utilisateur",
    "phone": "+1234567890",
    "address": "123 Rue Test",
    "role": "etudiant"
  }'
```

## 👤 4. Voir un utilisateur spécifique

```bash
curl -X GET http://localhost:8081/api/v1/users/2 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## ✏️ 5. Modifier un utilisateur

```bash
curl -X PUT http://localhost:8081/api/v1/users/2 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Modifié",
    "last_name": "Utilisateur"
  }'
```

## 🔄 6. Changer le rôle d'un utilisateur

```bash
curl -X PATCH http://localhost:8081/api/v1/users/2/role \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "role": "professeur"
  }'
```

## 🗑️ 7. Supprimer un utilisateur

```bash
curl -X DELETE http://localhost:8081/api/v1/users/2 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 🧪 Test complet avec script

Pour exécuter un test complet automatique :

```bash
# Test rapide
./quick_test.sh

# Test complet avec toutes les permissions
./test_roles.sh
```

## 📊 Résumé des permissions

| Rôle | Peut voir | Peut créer | Peut modifier | Peut supprimer | Peut changer rôle |
|------|-----------|------------|---------------|----------------|-------------------|
| **Super Admin** | Tous les utilisateurs | Tous les rôles | Tous les utilisateurs | Tous les utilisateurs | Tous les rôles |
| **Admin** | Professeurs + Étudiants | Professeurs + Étudiants | Professeurs + Étudiants | Professeurs + Étudiants | Professeurs + Étudiants |
| **Professeur** | Professeurs + Étudiants (lecture seule) | ❌ | ❌ | ❌ | ❌ |
| **Étudiant** | Étudiants seulement (lecture seule) | ❌ | ❌ | ❌ | ❌ |

## 🔍 Champs visibles par rôle

### Super Admin
- Tous les champs pour tous les utilisateurs

### Admin
- Tous les champs pour Professeurs et Étudiants
- Ne voit pas les autres Admins

### Professeur
- ID, nom, prénom, rôle, date de création pour Professeurs et Étudiants

### Étudiant
- ID, nom, prénom pour les autres Étudiants seulement 