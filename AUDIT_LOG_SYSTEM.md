# Système de Journal d'Activité (Audit Log)

## Vue d'ensemble

Le système de journal d'activité d'EduQR permet de tracer et d'enregistrer toutes les actions importantes effectuées dans l'application. Il fournit une transparence complète et facilite les enquêtes internes en cas de besoin.

## Objectifs

- **Sécurité et transparence** : Tracer qui a fait quoi, quand et sur quelle ressource
- **Enquêtes internes** : Faciliter les investigations en cas d'incident
- **Conformité** : Respecter les exigences de traçabilité des actions sensibles
- **Monitoring** : Surveiller l'activité des utilisateurs et détecter les comportements anormaux

## Architecture

### Backend (Go/Gin)

#### Modèles
- `AuditLog` : Structure principale pour stocker les entrées de journal
- `AuditLogRequest` : Structure pour les requêtes de création
- `AuditLogResponse` : Structure pour les réponses API
- `AuditLogFilter` : Structure pour les filtres de recherche

#### Composants
- **Repository** (`audit_log_repository.go`) : Couche d'accès aux données
- **Service** (`audit_log_service.go`) : Logique métier
- **Controller** (`audit_log_controller.go`) : Endpoints API
- **Middleware** (`audit_middleware.go`) : Interception automatique des actions

#### Endpoints API
```
GET    /api/v1/admin/audit-logs           # Liste des logs avec filtres
GET    /api/v1/admin/audit-logs/:id       # Détails d'un log
GET    /api/v1/admin/audit-logs/stats     # Statistiques
GET    /api/v1/admin/audit-logs/recent    # Logs récents
GET    /api/v1/admin/audit-logs/user/:id/activity  # Activité utilisateur
GET    /api/v1/admin/audit-logs/resource/:type/:id # Historique ressource
DELETE /api/v1/admin/audit-logs/clean     # Nettoyage anciens logs
```

### Frontend (React/MUI)

#### Composants
- `AuditLogPage` : Page principale
- `AuditLogTable` : Tableau des logs avec pagination
- `AuditLogFilters` : Filtres de recherche
- `AuditLogStats` : Statistiques et graphiques
- `AuditLogDetailDialog` : Modal de détails

#### Hook
- `useAuditLogs` : Gestion des appels API et état

## Types d'Actions Traçées

### Actions CRUD
- **CREATE** : Création d'utilisateurs, salles, matières, cours
- **UPDATE** : Modification de données existantes
- **DELETE** : Suppression (soft delete) d'éléments

### Actions d'Authentification
- **LOGIN** : Connexion utilisateur
- **LOGOUT** : Déconnexion utilisateur

### Ressources Surveillées
- **user** : Utilisateurs
- **room** : Salles
- **subject** : Matières
- **course** : Cours
- **event** : Événements

## Informations Enregistrées

### Métadonnées de Base
- ID unique du log
- Timestamp de création
- Description humaine de l'action

### Informations Utilisateur
- ID de l'utilisateur
- Email de l'utilisateur
- Rôle de l'utilisateur

### Informations d'Action
- Type d'action (create, update, delete, login, logout)
- Type de ressource affectée
- ID de la ressource (si applicable)

### Données de Changement
- **Anciennes valeurs** : État avant modification (JSON)
- **Nouvelles valeurs** : État après modification (JSON)

### Informations Techniques
- Adresse IP de l'utilisateur
- User Agent du navigateur

## Permissions d'Accès

### Rôles Autorisés
- **Super Admin** : Accès complet à tous les logs
- **Admin** : Accès complet à tous les logs
- **Professeur** : Aucun accès
- **Étudiant** : Aucun accès

### Sécurité
- Authentification JWT requise
- Middleware de vérification des rôles
- Logs de connexion automatiques

## Fonctionnalités Frontend

### Interface Utilisateur
- **Tableau paginé** : Affichage des logs avec tri et pagination
- **Filtres avancés** : Recherche par action, ressource, utilisateur, dates
- **Statistiques** : Graphiques de répartition par action, ressource, rôle
- **Détails complets** : Modal avec toutes les informations d'un log

### Filtres Disponibles
- **Recherche textuelle** : Dans descriptions, emails, rôles
- **Type d'action** : create, update, delete, login, logout
- **Type de ressource** : user, room, subject, course, event
- **ID de ressource** : Filtrage par ID spécifique
- **ID d'utilisateur** : Activité d'un utilisateur spécifique
- **Période** : Filtrage par dates de début et fin

### Fonctionnalités Avancées
- **Export** : Possibilité d'exporter en PDF/CSV (futur)
- **Nettoyage** : Suppression automatique des anciens logs
- **Recherche** : Moteur de recherche textuelle (futur)

## Implémentation Technique

### Middleware d'Audit
Le `AuditMiddleware` intercepte automatiquement :
- Les actions de connexion/déconnexion
- Les requêtes POST/PUT/DELETE sur les ressources sensibles
- Les modifications de données utilisateur

### Logging Automatique
```go
// Exemple d'utilisation dans un service
auditService.LogUserAction(
    userID,
    userEmail,
    userRole,
    "create",
    "user",
    &newUserID,
    "Création d'un nouvel utilisateur",
    ipAddress,
    userAgent,
    nil, // oldValues
    newUser, // newValues
)
```

### Base de Données
- Table `audit_logs` avec index sur les champs fréquemment utilisés
- Soft delete pour conserver l'historique
- Stockage JSON pour les anciennes/nouvelles valeurs

## Maintenance et Performance

### Nettoyage Automatique
- Suppression des logs de plus d'un an
- Endpoint de nettoyage manuel
- Configuration de rétention personnalisable

### Optimisations
- Index sur les champs de recherche fréquents
- Pagination pour éviter les requêtes lourdes
- Cache des statistiques (futur)

### Monitoring
- Alertes sur les actions sensibles
- Détection d'activité anormale
- Rapports d'audit périodiques

## Utilisation

### Accès à l'Interface
1. Se connecter avec un compte Admin ou Super Admin
2. Naviguer vers "Journal d'Activité" dans le menu
3. Utiliser les filtres pour rechercher des actions spécifiques
4. Cliquer sur l'icône "œil" pour voir les détails d'un log

### Exemples de Recherche
- **Toutes les connexions** : Filtrer par action "login"
- **Modifications d'utilisateurs** : Filtrer par action "update" et ressource "user"
- **Activité d'un utilisateur** : Utiliser le filtre ID utilisateur
- **Actions récentes** : Utiliser les dates de début/fin

## Sécurité et Conformité

### Protection des Données
- Chiffrement des données sensibles
- Anonymisation optionnelle des logs
- Respect du RGPD pour la rétention

### Audit de Sécurité
- Logs d'accès aux logs d'audit
- Traçabilité complète des consultations
- Alertes sur les accès non autorisés

## Évolutions Futures

### Fonctionnalités Planifiées
- **Export PDF/CSV** : Génération de rapports
- **Alertes en temps réel** : Notifications sur actions sensibles
- **Analyse comportementale** : Détection d'anomalies
- **API webhook** : Intégration avec systèmes externes
- **Rétention configurable** : Politiques de conservation flexibles

### Améliorations Techniques
- **Cache Redis** : Performance des requêtes fréquentes
- **Elasticsearch** : Recherche full-text avancée
- **Streaming** : Logs en temps réel
- **Compression** : Optimisation du stockage 