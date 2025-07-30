package models

import (
	"time"
)

// Room représente une salle dans le système
type Room struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	Name      string     `json:"name" gorm:"uniqueIndex;not null"`
	Building  string     `json:"building"`
	Floor     string     `json:"floor"`
	IsModular bool       `json:"is_modular" gorm:"default:false"`
	ParentID  *uint      `json:"parent_id" gorm:"index"`
	Parent    *Room      `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children  []Room     `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// RoomResponse représente la réponse pour une salle
type RoomResponse struct {
	ID        uint           `json:"id"`
	Name      string         `json:"name"`
	Building  string         `json:"building"`
	Floor     string         `json:"floor"`
	IsModular bool           `json:"is_modular"`
	ParentID  *uint          `json:"parent_id"`
	Parent    *RoomResponse  `json:"parent,omitempty"`
	Children  []RoomResponse `json:"children,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// CreateRoomRequest représente la requête de création d'une salle
type CreateRoomRequest struct {
	Name      string `json:"name" binding:"required"`
	Building  string `json:"building"`
	Floor     string `json:"floor"`
	IsModular bool   `json:"is_modular"`
	// Si modulable, nombre de sous-salles à créer
	SubRoomsCount int `json:"sub_rooms_count,omitempty"`
}

// UpdateRoomRequest représente la requête de modification d'une salle
type UpdateRoomRequest struct {
	Name      string `json:"name" binding:"required"`
	Building  string `json:"building"`
	Floor     string `json:"floor"`
	IsModular bool   `json:"is_modular"`
}

// RoomFilter représente les filtres pour la recherche de salles
type RoomFilter struct {
	Name     string `form:"name"`
	Building string `form:"building"`
	Floor    string `form:"floor"`
	Modular  *bool  `form:"modular"`
	ParentID *uint  `form:"parent_id"`
}

// ToRoomResponse convertit un Room en RoomResponse
func (r *Room) ToRoomResponse() RoomResponse {
	response := RoomResponse{
		ID:        r.ID,
		Name:      r.Name,
		Building:  r.Building,
		Floor:     r.Floor,
		IsModular: r.IsModular,
		ParentID:  r.ParentID,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}

	if r.Parent != nil {
		parentResponse := r.Parent.ToRoomResponse()
		response.Parent = &parentResponse
	}

	if len(r.Children) > 0 {
		children := make([]RoomResponse, len(r.Children))
		for i, child := range r.Children {
			children[i] = child.ToRoomResponse()
		}
		response.Children = children
	}

	return response
}
