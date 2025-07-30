# EduQR Frontend

Interface utilisateur React pour l'application EduQR - Système de gestion des présences par QR Code.

## 🚀 Technologies

- **React 18** - Framework JavaScript
- **Material-UI (MUI)** - Composants UI
- **React Router** - Navigation
- **React Hook Form** - Gestion des formulaires
- **FullCalendar** - Calendrier interactif
- **Axios** - Client HTTP
- **JWT** - Authentification

## 📦 Installation

```bash
# Installer les dépendances
npm install
# ou
yarn install
```

## 🏃‍♂️ Démarrage

```bash
# Démarrer en mode développement
npm start
# ou
yarn start
```

L'application sera accessible sur `http://localhost:3000`

## 🔧 Configuration

### Variables d'environnement

Créer un fichier `.env` à la racine du projet :

```env
REACT_APP_API_URL=http://localhost:8081/api/v1
```

### Configuration de l'API

L'URL de l'API backend est configurée dans `src/utils/api.js`

## 🏗️ Structure du projet

```
src/
├── auth/           # Authentification JWT
├── components/     # Composants réutilisables
├── hooks/          # Hooks personnalisés
├── layouts/        # Layouts de l'application
├── pages/          # Pages de l'application
├── routes/         # Configuration des routes
├── sections/       # Sections de pages
├── theme/          # Configuration du thème MUI
└── utils/          # Utilitaires
```

## 👥 Rôles utilisateurs

- **Super Admin** : Accès complet à toutes les fonctionnalités
- **Admin** : Gestion des utilisateurs, cours, salles, matières
- **Professeur** : Gestion de ses cours, présences, absences
- **Étudiant** : Consultation de ses cours et présences

## 🎨 Fonctionnalités principales

- **Authentification** : Login/Logout avec JWT
- **Tableau de bord** : Vue d'ensemble des statistiques
- **Calendrier** : Visualisation et gestion des cours
- **Gestion des utilisateurs** : CRUD des utilisateurs
- **Gestion des cours** : Création, modification, suppression
- **Gestion des salles** : CRUD des salles
- **Gestion des matières** : CRUD des matières
- **Gestion des absences** : Validation des demandes d'absence
- **QR Codes** : Génération et affichage des QR codes par salle
- **Journal d'activité** : Audit des actions utilisateurs

## 🔗 API Backend

Ce frontend nécessite le backend EduQR qui doit être en cours d'exécution sur `http://localhost:8081`

## 📝 Scripts disponibles

```bash
npm start          # Démarrer en mode développement
npm run build      # Construire pour la production
npm test           # Lancer les tests
npm run eject      # Éjecter la configuration (irréversible)
```

## 🤝 Contribution

1. Fork le projet
2. Créer une branche feature (`git checkout -b feature/AmazingFeature`)
3. Commit les changements (`git commit -m 'Add some AmazingFeature'`)
4. Push vers la branche (`git push origin feature/AmazingFeature`)
5. Ouvrir une Pull Request

## 📄 Licence

Ce projet est sous licence MIT. 