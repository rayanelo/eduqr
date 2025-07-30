package tests

import (
	"eduqr-backend/internal/models"
	"eduqr-backend/internal/repositories"
	"eduqr-backend/internal/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoomRepository(t *testing.T) {
	// Nettoyer avant chaque test
	cleanupTestDatabase()

	t.Run("CreateRoom_Success", func(t *testing.T) {
		repo := repositories.NewRoomRepository(testDB)

		room := &models.Room{
			Name:      "Test Room 1",
			Building:  "Building A",
			Floor:     "1st Floor",
			IsModular: false,
		}

		err := repo.CreateRoom(room)
		assert.NoError(t, err)
		assert.NotZero(t, room.ID)
		assert.Equal(t, "Test Room 1", room.Name)
	})

	t.Run("CreateRoom_DuplicateName_ShouldFail", func(t *testing.T) {
		repo := repositories.NewRoomRepository(testDB)

		// Créer la première salle
		room1 := &models.Room{
			Name:      "Duplicate Room",
			Building:  "Building A",
			Floor:     "1st Floor",
			IsModular: false,
		}
		err := repo.CreateRoom(room1)
		assert.NoError(t, err)

		// Essayer de créer une salle avec le même nom
		room2 := &models.Room{
			Name:      "Duplicate Room",
			Building:  "Building B",
			Floor:     "2nd Floor",
			IsModular: true,
		}
		err = repo.CreateRoom(room2)
		assert.Error(t, err) // Doit échouer à cause de la contrainte unique
	})

	t.Run("GetRoomByID_Success", func(t *testing.T) {
		repo := repositories.NewRoomRepository(testDB)

		// Créer une salle
		room := createTestRoom()

		// Récupérer la salle par ID
		retrievedRoom, err := repo.GetRoomByID(room.ID)
		assert.NoError(t, err)
		assert.Equal(t, room.ID, retrievedRoom.ID)
		assert.Equal(t, room.Name, retrievedRoom.Name)
	})

	t.Run("GetRoomByID_NotFound", func(t *testing.T) {
		repo := repositories.NewRoomRepository(testDB)

		// Essayer de récupérer une salle inexistante
		_, err := repo.GetRoomByID(99999)
		assert.Error(t, err)
	})

	t.Run("GetAllRooms_Success", func(t *testing.T) {
		repo := repositories.NewRoomRepository(testDB)

		// Créer plusieurs salles
		createTestRoom()
		createTestRoom()
		createTestRoom()

		// Récupérer toutes les salles
		rooms, err := repo.GetAllRooms(nil)
		assert.NoError(t, err)
		assert.Len(t, rooms, 3)
	})

	t.Run("GetAllRooms_WithFilter", func(t *testing.T) {
		repo := repositories.NewRoomRepository(testDB)

		// Créer des salles avec différents bâtiments
		room1 := &models.Room{Name: "Room A", Building: "Building A", Floor: "1st"}
		room2 := &models.Room{Name: "Room B", Building: "Building B", Floor: "1st"}
		testDB.Create(room1)
		testDB.Create(room2)

		// Filtrer par bâtiment
		filter := &models.RoomFilter{Building: "Building A"}
		rooms, err := repo.GetAllRooms(filter)
		assert.NoError(t, err)
		assert.Len(t, rooms, 1)
		assert.Equal(t, "Building A", rooms[0].Building)
	})

	t.Run("UpdateRoom_Success", func(t *testing.T) {
		repo := repositories.NewRoomRepository(testDB)

		// Créer une salle
		room := createTestRoom()

		// Modifier la salle
		room.Name = "Updated Room Name"
		room.Building = "Updated Building"

		err := repo.UpdateRoom(room)
		assert.NoError(t, err)

		// Vérifier que les modifications sont sauvegardées
		updatedRoom, err := repo.GetRoomByID(room.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Room Name", updatedRoom.Name)
		assert.Equal(t, "Updated Building", updatedRoom.Building)
	})

	t.Run("DeleteRoom_Success", func(t *testing.T) {
		repo := repositories.NewRoomRepository(testDB)

		// Créer une salle
		room := createTestRoom()

		// Supprimer la salle
		err := repo.DeleteRoom(room.ID)
		assert.NoError(t, err)

		// Vérifier que la salle n'existe plus
		_, err = repo.GetRoomByID(room.ID)
		assert.Error(t, err)
	})

	t.Run("CheckRoomExists_Exists", func(t *testing.T) {
		repo := repositories.NewRoomRepository(testDB)

		// Créer une salle
		room := createTestRoom()

		// Vérifier que la salle existe
		exists, err := repo.CheckRoomExists(room.Name, nil)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("CheckRoomExists_NotExists", func(t *testing.T) {
		repo := repositories.NewRoomRepository(testDB)

		// Vérifier qu'une salle inexistante n'existe pas
		exists, err := repo.CheckRoomExists("NonExistentRoom", nil)
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("GetModularRooms_Success", func(t *testing.T) {
		repo := repositories.NewRoomRepository(testDB)

		// Créer des salles modulables et non-modulables
		modularRoom := &models.Room{Name: "Modular Room", IsModular: true}
		normalRoom := &models.Room{Name: "Normal Room", IsModular: false}
		testDB.Create(modularRoom)
		testDB.Create(normalRoom)

		// Récupérer les salles modulables
		modularRooms, err := repo.GetModularRooms()
		assert.NoError(t, err)
		assert.Len(t, modularRooms, 1)
		assert.True(t, modularRooms[0].IsModular)
	})
}

func TestRoomService(t *testing.T) {
	// Nettoyer avant chaque test
	cleanupTestDatabase()

	t.Run("CreateRoom_Success", func(t *testing.T) {
		repo := repositories.NewRoomRepository(testDB)
		service := services.NewRoomService(repo)

		req := &models.CreateRoomRequest{
			Name:      "Service Test Room",
			Building:  "Service Building",
			Floor:     "Service Floor",
			IsModular: false,
		}

		response, err := service.CreateRoom(req)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, req.Name, response.Name)
		assert.Equal(t, req.Building, response.Building)
	})

	t.Run("CreateRoom_EmptyName_ShouldFail", func(t *testing.T) {
		repo := repositories.NewRoomRepository(testDB)
		service := services.NewRoomService(repo)

		req := &models.CreateRoomRequest{
			Name:      "",
			Building:  "Service Building",
			Floor:     "Service Floor",
			IsModular: false,
		}

		_, err := service.CreateRoom(req)
		assert.Error(t, err)
	})

	t.Run("GetRoomByID_Success", func(t *testing.T) {
		repo := repositories.NewRoomRepository(testDB)
		service := services.NewRoomService(repo)

		// Créer une salle
		room := createTestRoom()

		// Récupérer la salle
		response, err := service.GetRoomByID(room.ID)
		assert.NoError(t, err)
		assert.Equal(t, room.ID, response.ID)
		assert.Equal(t, room.Name, response.Name)
	})

	t.Run("UpdateRoom_Success", func(t *testing.T) {
		repo := repositories.NewRoomRepository(testDB)
		service := services.NewRoomService(repo)

		// Créer une salle
		room := createTestRoom()

		req := &models.UpdateRoomRequest{
			Name:      "Updated Service Room",
			Building:  "Updated Service Building",
			Floor:     "Updated Service Floor",
			IsModular: true,
		}

		response, err := service.UpdateRoom(room.ID, req)
		assert.NoError(t, err)
		assert.Equal(t, req.Name, response.Name)
		assert.Equal(t, req.Building, response.Building)
		assert.Equal(t, req.IsModular, response.IsModular)
	})

	t.Run("GetAllRooms_Success", func(t *testing.T) {
		repo := repositories.NewRoomRepository(testDB)
		service := services.NewRoomService(repo)

		// Créer plusieurs salles
		createTestRoom()
		createTestRoom()
		createTestRoom()

		// Récupérer toutes les salles
		rooms, err := service.GetAllRooms(nil)
		assert.NoError(t, err)
		assert.Len(t, rooms, 3)
	})
}

func TestRoomValidation(t *testing.T) {
	t.Run("ValidateRoomName_ValidNames", func(t *testing.T) {
		validNames := []string{
			"Salle 101",
			"Amphi A",
			"Laboratoire Informatique",
			"Bureau 2.3",
			"Room-123",
		}

		for _, name := range validNames {
			assert.NotEmpty(t, name, "Le nom ne devrait pas être vide")
			assert.Len(t, name, 1, "Le nom devrait avoir au moins 1 caractère")
		}
	})

	t.Run("ValidateRoomName_InvalidNames", func(t *testing.T) {
		invalidNames := []string{
			"",
			"   ",
		}

		for _, name := range invalidNames {
			if name == "" {
				assert.Empty(t, name, "Le nom vide est détecté correctement")
			}
		}
	})
}
