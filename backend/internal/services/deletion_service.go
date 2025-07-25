package services

import (
	"fmt"
	"time"

	"eduqr-backend/internal/models"
	"eduqr-backend/internal/repositories"
)

type DeletionService struct {
	userRepo    *repositories.UserRepository
	courseRepo  *repositories.CourseRepository
	roomRepo    *repositories.RoomRepository
	subjectRepo *repositories.SubjectRepository
}

func NewDeletionService(
	userRepo *repositories.UserRepository,
	courseRepo *repositories.CourseRepository,
	roomRepo *repositories.RoomRepository,
	subjectRepo *repositories.SubjectRepository,
) *DeletionService {
	return &DeletionService{
		userRepo:    userRepo,
		courseRepo:  courseRepo,
		roomRepo:    roomRepo,
		subjectRepo: subjectRepo,
	}
}

// DeleteUserResponse contient les informations sur la suppression d'un utilisateur
type DeleteUserResponse struct {
	Success       bool                    `json:"success"`
	Message       string                  `json:"message"`
	Warnings      []string                `json:"warnings,omitempty"`
	FutureCourses []models.CourseResponse `json:"future_courses,omitempty"`
	PastCourses   []models.CourseResponse `json:"past_courses,omitempty"`
}

// DeleteUser supprime un utilisateur selon les règles de sécurité
func (s *DeletionService) DeleteUser(userID uint, currentUserRole string, confirmWithCourses bool) (*DeleteUserResponse, error) {
	// Récupérer l'utilisateur à supprimer
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("utilisateur non trouvé")
	}

	// Vérifier les permissions selon le rôle
	if !s.canDeleteUser(currentUserRole, user.Role) {
		return nil, fmt.Errorf("permissions insuffisantes pour supprimer cet utilisateur")
	}

	// Vérifier tous les cours liés à cet utilisateur (futurs et passés)
	allCourses, err := s.courseRepo.GetAllCoursesByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la vérification des cours liés")
	}

	// Debug: afficher le nombre de cours trouvés
	fmt.Printf("DEBUG: Utilisateur %d a %d cours liés\n", userID, len(allCourses))

	// Séparer les cours futurs et passés
	var futureCourses, pastCourses []models.Course
	now := time.Now()

	for _, course := range allCourses {
		if course.StartTime.After(now) {
			futureCourses = append(futureCourses, course)
		} else {
			pastCourses = append(pastCourses, course)
		}
	}

	// Si l'utilisateur a des cours et que la confirmation n'est pas donnée, afficher un avertissement
	if len(allCourses) > 0 && !confirmWithCourses {
		response := &DeleteUserResponse{
			Success:       false,
			Message:       fmt.Sprintf("attention: la suppression de cet utilisateur entraînera la suppression de %d cours liés", len(allCourses)),
			FutureCourses: make([]models.CourseResponse, len(futureCourses)),
			PastCourses:   make([]models.CourseResponse, len(pastCourses)),
		}

		// Convertir les cours futurs
		for i, course := range futureCourses {
			response.FutureCourses[i] = course.ToCourseResponse()
		}

		// Convertir les cours passés
		for i, course := range pastCourses {
			response.PastCourses[i] = course.ToCourseResponse()
		}

		return response, nil
	}

	// Si l'utilisateur a des cours et que la confirmation est donnée, supprimer tous les cours liés
	if len(allCourses) > 0 && confirmWithCourses {
		// Supprimer tous les cours liés
		for _, course := range allCourses {
			if err := s.courseRepo.DeleteCourse(course.ID); err != nil {
				return nil, fmt.Errorf("erreur lors de la suppression du cours %d: %v", course.ID, err)
			}
		}
	}

	// Effectuer la suppression de l'utilisateur (soft delete)
	if err := s.userRepo.DeleteUser(userID); err != nil {
		return nil, fmt.Errorf("erreur lors de la suppression de l'utilisateur")
	}

	response := &DeleteUserResponse{
		Success: true,
		Message: "utilisateur supprimé avec succès",
	}

	// Ajouter des informations sur les cours supprimés si applicable
	if len(allCourses) > 0 {
		response.Warnings = append(response.Warnings,
			fmt.Sprintf("%d cours liés ont également été supprimés", len(allCourses)))
	}

	return response, nil
}

// DeleteRoomResponse contient les informations sur la suppression d'une salle
type DeleteRoomResponse struct {
	Success       bool                    `json:"success"`
	Message       string                  `json:"message"`
	Warnings      []string                `json:"warnings,omitempty"`
	FutureCourses []models.CourseResponse `json:"future_courses,omitempty"`
	ChildRooms    []models.RoomResponse   `json:"child_rooms,omitempty"`
}

// DeleteRoom supprime une salle selon les règles de sécurité
func (s *DeletionService) DeleteRoom(roomID uint) (*DeleteRoomResponse, error) {
	// Récupérer la salle
	room, err := s.roomRepo.GetRoomByID(roomID)
	if err != nil {
		return nil, fmt.Errorf("salle non trouvée")
	}

	// Vérifier les cours futurs pour cette salle et ses sous-salles
	futureCourses, err := s.courseRepo.GetFutureCoursesByRoom(roomID)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la vérification des cours futurs")
	}

	// Si la salle a des cours futurs, empêcher la suppression
	if len(futureCourses) > 0 {
		response := &DeleteRoomResponse{
			Success: false,
			Message: fmt.Sprintf("impossible de supprimer la salle car elle a %d cours à venir", len(futureCourses)),
		}

		return response, nil
	}

	// Si c'est une salle modulable, vérifier les sous-salles
	var childRooms []models.RoomResponse
	if room.IsModular {
		children, err := s.roomRepo.GetChildRooms(roomID)
		if err != nil {
			return nil, fmt.Errorf("erreur lors de la récupération des sous-salles")
		}

		// Vérifier les cours futurs dans les sous-salles
		for _, child := range children {
			childFutureCourses, err := s.courseRepo.GetFutureCoursesByRoom(child.ID)
			if err != nil {
				return nil, fmt.Errorf("erreur lors de la vérification des cours futurs dans la sous-salle")
			}

			if len(childFutureCourses) > 0 {
				return &DeleteRoomResponse{
					Success:       false,
					Message:       fmt.Sprintf("impossible de supprimer la salle car la sous-salle '%s' a %d cours à venir", child.Name, len(childFutureCourses)),
					FutureCourses: make([]models.CourseResponse, len(childFutureCourses)),
				}, nil
			}

			childRooms = append(childRooms, child.ToRoomResponse())
		}
	}

	// Effectuer la suppression (soft delete)
	if err := s.roomRepo.DeleteRoom(roomID); err != nil {
		return nil, fmt.Errorf("erreur lors de la suppression de la salle")
	}

	response := &DeleteRoomResponse{
		Success:    true,
		Message:    "salle supprimée avec succès",
		ChildRooms: childRooms,
	}

	// Ajouter des avertissements si nécessaire
	if room.IsModular && len(childRooms) > 0 {
		response.Warnings = append(response.Warnings,
			fmt.Sprintf("%d sous-salles ont également été supprimées", len(childRooms)))
	}

	return response, nil
}

// DeleteSubjectResponse contient les informations sur la suppression d'une matière
type DeleteSubjectResponse struct {
	Success       bool                    `json:"success"`
	Message       string                  `json:"message"`
	Warnings      []string                `json:"warnings,omitempty"`
	LinkedCourses []models.CourseResponse `json:"linked_courses,omitempty"`
}

// DeleteSubject supprime une matière selon les règles de sécurité
func (s *DeletionService) DeleteSubject(subjectID uint) (*DeleteSubjectResponse, error) {
	// Vérifier les cours liés à cette matière
	linkedCourses, err := s.courseRepo.GetCoursesBySubject(subjectID)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la vérification des cours liés")
	}

	// Si la matière a des cours liés, empêcher la suppression
	if len(linkedCourses) > 0 {
		response := &DeleteSubjectResponse{
			Success: false,
			Message: fmt.Sprintf("impossible de supprimer la matière car elle est liée à %d cours", len(linkedCourses)),
		}

		return response, nil
	}

	// Effectuer la suppression (soft delete)
	if err := s.subjectRepo.DeleteSubject(subjectID); err != nil {
		return nil, fmt.Errorf("erreur lors de la suppression de la matière")
	}

	return &DeleteSubjectResponse{
		Success: true,
		Message: "matière supprimée avec succès",
	}, nil
}

// DeleteCourseResponse contient les informations sur la suppression d'un cours
type DeleteCourseResponse struct {
	Success            bool     `json:"success"`
	Message            string   `json:"message"`
	Warnings           []string `json:"warnings,omitempty"`
	DeletedOccurrences int      `json:"deleted_occurrences,omitempty"`
}

// DeleteCourse supprime un cours selon les règles de sécurité
func (s *DeletionService) DeleteCourse(courseID uint, deleteRecurring bool) (*DeleteCourseResponse, error) {
	// Récupérer le cours
	course, err := s.courseRepo.GetCourseByID(courseID)
	if err != nil {
		return nil, fmt.Errorf("cours non trouvé")
	}

	// Vérifier s'il y a des présences enregistrées
	hasAttendance, err := s.courseRepo.HasAttendance(courseID)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la vérification des présences")
	}

	response := &DeleteCourseResponse{}

	// Ajouter un avertissement si des présences existent
	if hasAttendance {
		response.Warnings = append(response.Warnings,
			"attention: des données de présence seront perdues")
	}

	// Gérer la suppression selon le type de cours
	if course.IsRecurring {
		if deleteRecurring {
			// Supprimer toute la série récurrente
			if err := s.courseRepo.DeleteRecurringCourses(courseID); err != nil {
				return nil, fmt.Errorf("erreur lors de la suppression de la série récurrente")
			}
			response.Message = "série de cours récurrents supprimée avec succès"
			response.DeletedOccurrences = -1 // Indique toute la série
		} else {
			// Supprimer seulement cette occurrence
			if err := s.courseRepo.DeleteCourse(courseID); err != nil {
				return nil, fmt.Errorf("erreur lors de la suppression du cours")
			}
			response.Message = "occurrence du cours supprimée avec succès"
			response.DeletedOccurrences = 1
		}
	} else {
		// Cours ponctuel
		if err := s.courseRepo.DeleteCourse(courseID); err != nil {
			return nil, fmt.Errorf("erreur lors de la suppression du cours")
		}
		response.Message = "cours supprimé avec succès"
		response.DeletedOccurrences = 1
	}

	response.Success = true
	return response, nil
}

// canDeleteUser vérifie si l'utilisateur actuel peut supprimer l'utilisateur cible
func (s *DeletionService) canDeleteUser(currentUserRole, targetUserRole string) bool {
	// Seul le Super Admin peut supprimer un Admin
	if targetUserRole == models.RoleAdmin {
		return currentUserRole == models.RoleSuperAdmin
	}

	// Les Admins et Super Admins peuvent supprimer les Professeurs et Étudiants
	if targetUserRole == models.RoleProfesseur || targetUserRole == models.RoleEtudiant {
		return currentUserRole == models.RoleAdmin || currentUserRole == models.RoleSuperAdmin
	}

	return false
}
