# 🧪 Scripts de Test pour EduQR

Ce dossier contient plusieurs scripts de test pour l'application EduQR, permettant de tester toutes les fonctionnalités de l'application de manière automatisée.

## 📋 Scripts Disponibles

### 1. `test_rapide.sh` - Test Rapide ⚡
**Durée estimée : 2-3 minutes**

Teste les fonctionnalités principales de l'application :
- ✅ Santé de l'API
- ✅ Authentification
- ✅ Gestion des utilisateurs
- ✅ Logs d'audit
- ✅ Gestion des ressources (salles, matières, cours)
- ✅ Gestion des absences et présences
- ✅ Création de données
- ✅ Sécurité et authentification
- ✅ Endpoints publics
- ✅ Validation des données

```bash
./test_rapide.sh
```

### 2. `test_roles.sh` - Test des Rôles et Permissions 🎭
**Durée estimée : 5-7 minutes**

Teste spécifiquement la hiérarchie des rôles et les permissions :
- ✅ Hiérarchie des rôles (Super Admin → Admin → Professeur → Étudiant)
- ✅ Permissions par niveau de rôle
- ✅ Permissions croisées
- ✅ Gestion des rôles (promotion/rétrogradation)
- ✅ Vue des utilisateurs selon le rôle
- ✅ Suppression sécurisée
- ✅ Tentatives d'accès non autorisées

```bash
./test_roles.sh
```

### 3. `test_complet_eduqr.sh` - Test Complet 🎓
**Durée estimée : 10-15 minutes**

Teste TOUTES les fonctionnalités de l'application pour chaque type d'utilisateur :
- ✅ Tests pour Super Admin (accès complet)
- ✅ Tests pour Admin (permissions limitées)
- ✅ Tests pour Professeur (permissions très limitées)
- ✅ Tests pour Étudiant (permissions minimales)
- ✅ Tests de sécurité et authentification
- ✅ Tests des endpoints publics
- ✅ Tests de validation des données
- ✅ Tests de performance et limites

```bash
./test_complet_eduqr.sh
```

### 4. `test_audit_logs.sh` - Test des Logs d'Audit 📊
**Durée estimée : 3-5 minutes**

Teste spécifiquement le système de logs d'audit :
- ✅ Récupération des logs d'audit
- ✅ Statistiques des logs
- ✅ Logs récents
- ✅ Filtrage des logs
- ✅ Accès non autorisé

```bash
./test_audit_logs.sh
```

## 🚀 Prérequis

### 1. Application en cours d'exécution
Assurez-vous que l'application EduQR est démarrée :

```bash
# Terminal 1 - Backend
cd backend
go run cmd/server/main.go

# Terminal 2 - Frontend (optionnel pour les tests API)
cd frontend
npm start
```

### 2. Base de données initialisée
Exécutez les scripts de seed pour créer les données de test :

```bash
# Créer les utilisateurs de base
cd backend
go run cmd/seed/main.go

# Ou créer toutes les données de test
go run cmd/seed_all/main.go
```

### 3. Comptes de test disponibles
Les scripts utilisent les comptes de test suivants :

| Rôle | Email | Mot de passe |
|------|-------|--------------|
| **Super Admin** | `superadmin@eduqr.com` | `superadmin123` |
| **Admin** | `admin@eduqr.com` | `admin123` |
| **Professeur** | `prof1@eduqr.com` | `prof123` |
| **Étudiant** | `etudiant1@eduqr.com` | `student123` |

## 🎯 Utilisation

### Test rapide pour vérifier que tout fonctionne
```bash
./test_rapide.sh
```

### Test complet pour valider toutes les fonctionnalités
```bash
./test_complet_eduqr.sh
```

### Test spécifique des rôles et permissions
```bash
./test_roles.sh
```

### Test du système d'audit
```bash
./test_audit_logs.sh
```

## 📊 Interprétation des Résultats

### Codes de couleur
- 🟢 **Vert** : Test réussi
- 🔴 **Rouge** : Test échoué
- 🟡 **Jaune** : Informations supplémentaires

### Statuts HTTP attendus
- `200` : Succès
- `201` : Création réussie
- `204` : Suppression réussie
- `400` : Erreur de validation (attendu pour certains tests)
- `401` : Non authentifié (attendu pour certains tests)
- `403` : Non autorisé (attendu pour certains tests)
- `404` : Ressource non trouvée
- `500` : Erreur serveur

## 🔧 Dépannage

### Problème : "Connection refused"
```bash
# Vérifiez que le backend est démarré
curl http://localhost:8081/health
```

### Problème : "User not found"
```bash
# Recréez les données de test
cd backend
go run cmd/seed_all/main.go
```

### Problème : "Permission denied"
```bash
# Rendez les scripts exécutables
chmod +x *.sh
```

### Problème : "Token invalid"
```bash
# Vérifiez que les comptes de test existent
cd backend
go run cmd/check_users/main.go
```

## 📈 Fonctionnalités Testées

### 🔐 Authentification et Autorisation
- Connexion avec JWT
- Gestion des sessions
- Middleware d'authentification
- Middleware d'autorisation par rôle
- Validation des tokens

### 👥 Gestion des Utilisateurs
- Création d'utilisateurs
- Modification des profils
- Suppression sécurisée
- Gestion des rôles
- Vue filtrée selon les permissions

### 🏫 Gestion des Ressources
- **Salles** : CRUD, salles modulaires
- **Matières** : CRUD, codes et descriptions
- **Cours** : CRUD, planification, conflits
- **Événements** : CRUD, calendrier

### ❌ Gestion des Absences
- Déclaration d'absences par les étudiants
- Validation par les professeurs/admins
- Justificatifs et documents
- Statistiques et rapports

### ✅ Gestion des Présences
- Scan de QR codes
- Enregistrement des présences
- Statistiques de présence
- Génération de QR codes

### 📊 Logs d'Audit
- Traçabilité des actions
- Statistiques d'activité
- Filtrage et recherche
- Nettoyage automatique

## 🎭 Hiérarchie des Rôles

### Super Admin (Niveau 4)
- ✅ Accès complet à toutes les fonctionnalités
- ✅ Gestion de tous les utilisateurs
- ✅ Suppression de tous les éléments
- ✅ Accès aux logs d'audit

### Admin (Niveau 3)
- ✅ Gestion des professeurs et étudiants
- ✅ Gestion des ressources (salles, matières, cours)
- ✅ Accès aux logs d'audit
- ❌ Ne peut pas gérer d'autres admins

### Professeur (Niveau 2)
- ✅ Gestion des absences de ses cours
- ✅ Gestion des présences de ses cours
- ✅ Génération de QR codes
- ❌ Pas d'accès aux ressources admin

### Étudiant (Niveau 1)
- ✅ Déclaration de ses absences
- ✅ Scan de QR codes
- ✅ Consultation de ses données
- ❌ Accès très limité

## 📝 Personnalisation

### Modifier l'URL de base
Éditez la variable `BASE_URL` dans chaque script :
```bash
BASE_URL="http://localhost:8081"  # URL par défaut
```

### Ajouter de nouveaux tests
Créez un nouveau script en vous basant sur les modèles existants :
```bash
cp test_rapide.sh mon_test_personnalise.sh
chmod +x mon_test_personnalise.sh
```

### Modifier les comptes de test
Éditez les variables dans les scripts :
```bash
declare -A test_users=(
    ["superadmin"]="votre_email:votre_mot_de_passe"
    # ...
)
```

## 🔗 Liens Utiles

- [Documentation de l'API](backend/curl_examples.md)
- [Système de logs d'audit](AUDIT_LOG_SYSTEM.md)
- [État de l'application](EDUQR_APPLICATION_STATE.md)
- [Exemples curl](backend/curl_examples.md)

## 📞 Support

En cas de problème avec les tests :
1. Vérifiez que l'application est démarrée
2. Vérifiez que les données de test sont créées
3. Consultez les logs du backend
4. Vérifiez la connectivité réseau

---

**🎉 Bon test !** 