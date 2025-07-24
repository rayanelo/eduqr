import { useContext } from 'react';
import { AuthContext } from '../auth/JwtContext';

// Rôles disponibles
export const ROLES = {
  SUPER_ADMIN: 'super_admin',
  ADMIN: 'admin',
  PROFESSEUR: 'professeur',
  ETUDIANT: 'etudiant',
};

// Hiérarchie des rôles (plus l'index est élevé, plus les permissions sont importantes)
const ROLE_HIERARCHY = {
  [ROLES.SUPER_ADMIN]: 4,
  [ROLES.ADMIN]: 3,
  [ROLES.PROFESSEUR]: 2,
  [ROLES.ETUDIANT]: 1,
};

// Fonction pour vérifier si un rôle peut gérer un autre rôle
const canManageRole = (managerRole, targetRole) => {
  // Super Admin peut gérer tout le monde
  if (managerRole === ROLES.SUPER_ADMIN) {
    return true;
  }
  
  // Admin peut gérer Professeur et Étudiant, mais pas d'autres Admins ou Super Admin
  if (managerRole === ROLES.ADMIN) {
    return targetRole === ROLES.PROFESSEUR || targetRole === ROLES.ETUDIANT;
  }
  
  // Professeur et Étudiant ne peuvent gérer personne
  return false;
};

// Fonction pour vérifier si un rôle peut voir un autre rôle
const canViewRole = (viewerRole, targetRole) => {
  // Super Admin peut voir tout le monde
  if (viewerRole === ROLES.SUPER_ADMIN) {
    return true;
  }
  
  // Admin peut voir Professeur et Étudiant, mais pas d'autres Admins (sauf lui-même)
  if (viewerRole === ROLES.ADMIN) {
    return targetRole === ROLES.PROFESSEUR || targetRole === ROLES.ETUDIANT;
  }
  
  // Professeur peut voir d'autres Professeurs et Étudiants
  if (viewerRole === ROLES.PROFESSEUR) {
    return targetRole === ROLES.PROFESSEUR || targetRole === ROLES.ETUDIANT;
  }
  
  // Étudiant ne peut voir que d'autres Étudiants
  if (viewerRole === ROLES.ETUDIANT) {
    return targetRole === ROLES.ETUDIANT;
  }
  
  return false;
};

// Fonction pour obtenir les champs visibles selon le rôle
const getViewableFields = (viewerRole, targetRole) => {
  // Super Admin peut voir tous les champs
  if (viewerRole === ROLES.SUPER_ADMIN) {
    return ['id', 'email', 'first_name', 'last_name', 'phone', 'address', 'avatar', 'role', 'created_at', 'updated_at'];
  }
  
  // Admin peut voir tous les champs pour Professeur et Étudiant
  if (viewerRole === ROLES.ADMIN && (targetRole === ROLES.PROFESSEUR || targetRole === ROLES.ETUDIANT)) {
    return ['id', 'email', 'first_name', 'last_name', 'phone', 'address', 'avatar', 'role', 'created_at', 'updated_at'];
  }
  
  // Professeur peut voir des champs limités pour d'autres Professeurs et Étudiants
  if (viewerRole === ROLES.PROFESSEUR && (targetRole === ROLES.PROFESSEUR || targetRole === ROLES.ETUDIANT)) {
    return ['id', 'first_name', 'last_name', 'role', 'created_at'];
  }
  
  // Étudiant ne peut voir que les noms pour d'autres Étudiants
  if (viewerRole === ROLES.ETUDIANT && targetRole === ROLES.ETUDIANT) {
    return ['id', 'first_name', 'last_name'];
  }
  
  return [];
};

export const usePermissions = () => {
  const { user } = useContext(AuthContext);
  const currentUserRole = user?.role || ROLES.ETUDIANT;

  const canManageUsers = user?.role === ROLES.SUPER_ADMIN || user?.role === ROLES.ADMIN;
  const canManageRooms = user?.role === ROLES.SUPER_ADMIN || user?.role === ROLES.ADMIN;
  const canManageSubjects = user?.role === ROLES.SUPER_ADMIN || user?.role === ROLES.ADMIN;
  const canCreateUser = user?.role === ROLES.SUPER_ADMIN || user?.role === ROLES.ADMIN;

  return {
    // Vérifier si l'utilisateur peut gérer un rôle spécifique
    canManageRole: (targetRole) => canManageRole(currentUserRole, targetRole),
    
    // Vérifier si l'utilisateur peut voir un rôle spécifique
    canViewRole: (targetRole) => canViewRole(currentUserRole, targetRole),
    
    // Vérifier si l'utilisateur peut créer des utilisateurs avec un rôle spécifique
    canCreateUser: (targetRole) => canManageRole(currentUserRole, targetRole),
    
    // Vérifier si l'utilisateur peut modifier un utilisateur avec un rôle spécifique
    canUpdateUser: (targetRole) => canManageRole(currentUserRole, targetRole),
    
    // Vérifier si l'utilisateur peut supprimer un utilisateur avec un rôle spécifique
    canDeleteUser: (targetRole) => canManageRole(currentUserRole, targetRole),
    
    // Vérifier si l'utilisateur peut changer le rôle d'un utilisateur
    canChangeUserRole: (targetRole) => canManageRole(currentUserRole, targetRole),
    
    // Obtenir les champs visibles pour un rôle cible
    getViewableFields: (targetRole) => getViewableFields(currentUserRole, targetRole),
    
    // Vérifier si l'utilisateur peut voir les actions (modifier, supprimer, etc.)
    canSeeActions: (targetRole) => canManageRole(currentUserRole, targetRole),
    
    // Obtenir le rôle actuel de l'utilisateur
    currentUserRole,
    
    // Vérifier si l'utilisateur a un rôle ou plus élevé
    hasRoleOrHigher: (requiredRole) => {
      const currentLevel = ROLE_HIERARCHY[currentUserRole] || 0;
      const requiredLevel = ROLE_HIERARCHY[requiredRole] || 0;
      return currentLevel >= requiredLevel;
    },
    
    // Obtenir les rôles que l'utilisateur peut créer
    getCreatableRoles: () => {
      if (currentUserRole === ROLES.SUPER_ADMIN) {
        return Object.values(ROLES);
      }
      if (currentUserRole === ROLES.ADMIN) {
        return [ROLES.PROFESSEUR, ROLES.ETUDIANT];
      }
      return [];
    },
    
    // Obtenir les rôles que l'utilisateur peut promouvoir vers
    getPromotableRoles: (currentTargetRole) => {
      if (currentUserRole === ROLES.SUPER_ADMIN) {
        return Object.values(ROLES);
      }
      if (currentUserRole === ROLES.ADMIN) {
        if (currentTargetRole === ROLES.ETUDIANT) {
          return [ROLES.PROFESSEUR];
        }
        if (currentTargetRole === ROLES.PROFESSEUR) {
          return [ROLES.ADMIN];
        }
      }
      return [];
    },

    // Permissions simplifiées
    canManageUsers,
    canManageRooms,
    canManageSubjects,
    canCreateUser,
  };
}; 