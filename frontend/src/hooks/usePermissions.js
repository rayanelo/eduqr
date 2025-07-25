import { useContext } from 'react';
import { AuthContext } from '../auth/JwtContext';

// Rôles disponibles
export const ROLES = {
  SUPER_ADMIN: 'super_admin',
  ADMIN: 'admin',
  PROFESSEUR: 'professeur',
  ETUDIANT: 'etudiant',
};

export const usePermissions = () => {
  const { user } = useContext(AuthContext);

  const canManageUsers = user?.role === ROLES.SUPER_ADMIN || user?.role === ROLES.ADMIN;
  const canManageRooms = user?.role === ROLES.SUPER_ADMIN || user?.role === ROLES.ADMIN;
  const canManageSubjects = user?.role === ROLES.SUPER_ADMIN || user?.role === ROLES.ADMIN;
  const canManageCourses = user?.role === ROLES.SUPER_ADMIN || user?.role === ROLES.ADMIN;

  // Vérifier si l'utilisateur peut supprimer des éléments
  const canDelete = () => {
    return user?.role === 'admin' || user?.role === 'super_admin';
  };

  // Vérifier si l'utilisateur peut supprimer un utilisateur spécifique
  const canDeleteUser = (targetUser) => {
    if (!canDelete()) return false;
    
    // Un utilisateur ne peut pas se supprimer lui-même
    if (user?.id === targetUser?.id) return false;
    
    // Seul le Super Admin peut supprimer un Admin
    if (targetUser?.role === 'admin') {
      return user?.role === 'super_admin';
    }
    
    // Les Admins et Super Admins peuvent supprimer les Professeurs et Étudiants
    if (targetUser?.role === 'professeur' || targetUser?.role === 'etudiant') {
      return user?.role === 'admin' || user?.role === 'super_admin';
    }
    
    return false;
  };

  // Vérifier si l'utilisateur peut supprimer une salle
  const canDeleteRoom = () => {
    return canDelete();
  };

  // Vérifier si l'utilisateur peut supprimer une matière
  const canDeleteSubject = () => {
    return canDelete();
  };

  // Vérifier si l'utilisateur peut supprimer un cours
  const canDeleteCourse = () => {
    return canDelete();
  };

  // Vérifier si l'utilisateur peut voir un utilisateur
  const canViewUser = (targetUser) => {
    if (!user) return false;
    
    // Un utilisateur peut toujours se voir lui-même
    if (user.id === targetUser?.id) return true;
    
    // Super Admin peut voir tout le monde
    if (user.role === 'super_admin') return true;
    
    // Admin peut voir les Professeurs et Étudiants, mais pas les autres Admins
    if (user.role === 'admin') {
      return targetUser?.role === 'professeur' || targetUser?.role === 'etudiant';
    }
    
    // Professeur peut voir les autres Professeurs et Étudiants
    if (user.role === 'professeur') {
      return targetUser?.role === 'professeur' || targetUser?.role === 'etudiant';
    }
    
    // Étudiant peut seulement voir les autres Étudiants
    if (user.role === 'etudiant') {
      return targetUser?.role === 'etudiant';
    }
    
    return false;
  };

  // Vérifier si l'utilisateur peut modifier un utilisateur
  const canEditUser = (targetUser) => {
    if (!user) return false;
    
    // Un utilisateur peut toujours se modifier lui-même
    if (user.id === targetUser?.id) return true;
    
    // Super Admin peut modifier tout le monde
    if (user.role === 'super_admin') return true;
    
    // Admin peut modifier les Professeurs et Étudiants, mais pas les autres Admins
    if (user.role === 'admin') {
      return targetUser?.role === 'professeur' || targetUser?.role === 'etudiant';
    }
    
    return false;
  };

  // Obtenir les rôles que l'utilisateur peut créer
  const getCreatableRoles = () => {
    if (!user) return [];
    
    switch (user.role) {
      case 'super_admin':
        return [ROLES.ADMIN, ROLES.PROFESSEUR, ROLES.ETUDIANT];
      case 'admin':
        return [ROLES.PROFESSEUR, ROLES.ETUDIANT];
      default:
        return [];
    }
  };

  // Vérifier si l'utilisateur a un rôle égal ou supérieur
  const hasRoleOrHigher = (requiredRole) => {
    if (!user) return false;
    
    const roleHierarchy = {
      [ROLES.SUPER_ADMIN]: 4,
      [ROLES.ADMIN]: 3,
      [ROLES.PROFESSEUR]: 2,
      [ROLES.ETUDIANT]: 1,
    };
    
    const userLevel = roleHierarchy[user.role] || 0;
    const requiredLevel = roleHierarchy[requiredRole] || 0;
    
    return userLevel >= requiredLevel;
  };

  return {
    canManageUsers,
    canManageRooms,
    canManageSubjects,
    canManageCourses,
    canDelete,
    canDeleteUser,
    canDeleteRoom,
    canDeleteSubject,
    canDeleteCourse,
    canViewUser,
    canEditUser,
    getCreatableRoles,
    hasRoleOrHigher,
    userRole: user?.role,
    userId: user?.id,
  };
}; 