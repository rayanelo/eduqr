package services

import (
	"eduqr-backend/internal/models"
	"eduqr-backend/internal/repositories"
	"errors"
	"fmt"
)

type RoomService struct {
	roomRepo *repositories.RoomRepository
}

func NewRoomService(roomRepo *repositories.RoomRepository) *RoomService {
	return &RoomService{roomRepo: roomRepo}
}

// GetAllRooms récupère toutes les salles avec filtres
func (s *RoomService) GetAllRooms(filter *models.RoomFilter) ([]models.RoomResponse, error) {
	rooms, err := s.roomRepo.GetAllRooms(filter)
	if err != nil {
		return nil, err
	}

	responses := make([]models.RoomResponse, len(rooms))
	for i, room := range rooms {
		responses[i] = room.ToRoomResponse()
	}

	return responses, nil
}

// GetRoomByID récupère une salle par son ID
func (s *RoomService) GetRoomByID(id uint) (*models.RoomResponse, error) {
	room, err := s.roomRepo.GetRoomByID(id)
	if err != nil {
		return nil, err
	}

	response := room.ToRoomResponse()
	return &response, nil
}

// CreateRoom crée une nouvelle salle avec gestion des sous-salles
func (s *RoomService) CreateRoom(req *models.CreateRoomRequest) (*models.RoomResponse, error) {
	// Vérifier si le nom existe déjà
	exists, err := s.roomRepo.CheckRoomExists(req.Name, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("une salle avec ce nom existe déjà")
	}

	// Créer la salle principale
	room := &models.Room{
		Name:      req.Name,
		Building:  req.Building,
		Floor:     req.Floor,
		IsModular: req.IsModular,
	}

	err = s.roomRepo.CreateRoom(room)
	if err != nil {
		return nil, err
	}

	// Si la salle est modulable et qu'on a spécifié un nombre de sous-salles
	if req.IsModular && req.SubRoomsCount >= 2 {
		err = s.createSubRooms(room.ID, req.Name, req.SubRoomsCount, req.Building, req.Floor)
		if err != nil {
			return nil, err
		}
	}

	// Récupérer la salle avec ses sous-salles
	createdRoom, err := s.roomRepo.GetRoomByID(room.ID)
	if err != nil {
		return nil, err
	}

	response := createdRoom.ToRoomResponse()
	return &response, nil
}

// UpdateRoom met à jour une salle
func (s *RoomService) UpdateRoom(id uint, req *models.UpdateRoomRequest) (*models.RoomResponse, error) {
	// Vérifier si la salle existe
	room, err := s.roomRepo.GetRoomByID(id)
	if err != nil {
		return nil, err
	}

	// Vérifier si le nouveau nom existe déjà (sauf pour cette salle)
	exists, err := s.roomRepo.CheckRoomExists(req.Name, &id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("une salle avec ce nom existe déjà")
	}

	// Mettre à jour les champs
	room.Name = req.Name
	room.Building = req.Building
	room.Floor = req.Floor
	room.IsModular = req.IsModular

	err = s.roomRepo.UpdateRoom(room)
	if err != nil {
		return nil, err
	}

	// Récupérer la salle mise à jour
	updatedRoom, err := s.roomRepo.GetRoomByID(id)
	if err != nil {
		return nil, err
	}

	response := updatedRoom.ToRoomResponse()
	return &response, nil
}

// DeleteRoom supprime une salle
func (s *RoomService) DeleteRoom(id uint) error {
	// Vérifier si la salle existe
	_, err := s.roomRepo.GetRoomByID(id)
	if err != nil {
		return err
	}

	// TODO: Vérifier si la salle est utilisée dans des cours
	// Si oui, empêcher la suppression ou faire un soft delete

	return s.roomRepo.DeleteRoom(id)
}

// createSubRooms crée les sous-salles pour une salle modulable
func (s *RoomService) createSubRooms(parentID uint, parentName string, count int, building string, floor string) error {
	suffixes := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}

	for i := 0; i < count && i < len(suffixes); i++ {
		subRoomName := fmt.Sprintf("%s %s", parentName, suffixes[i])

		// Vérifier si le nom existe déjà
		exists, err := s.roomRepo.CheckRoomExists(subRoomName, nil)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("une sous-salle avec le nom '%s' existe déjà", subRoomName)
		}

		subRoom := &models.Room{
			Name:     subRoomName,
			Building: building, // Hériter du bâtiment de la salle parent
			Floor:    floor,    // Hériter de l'étage de la salle parent
			ParentID: &parentID,
		}

		err = s.roomRepo.CreateRoom(subRoom)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetModularRooms récupère toutes les salles modulables
func (s *RoomService) GetModularRooms() ([]models.RoomResponse, error) {
	rooms, err := s.roomRepo.GetModularRooms()
	if err != nil {
		return nil, err
	}

	responses := make([]models.RoomResponse, len(rooms))
	for i, room := range rooms {
		responses[i] = room.ToRoomResponse()
	}

	return responses, nil
}
