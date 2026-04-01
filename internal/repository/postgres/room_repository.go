package postgres

import (
	"context"
	"time"

	"github.com/AdityaAWP/IdiomaMate/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type roomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) domain.RoomRepository {
	return &roomRepository{db: db}
}

func (r *roomRepository) Create(ctx context.Context, room *domain.Room) error {
	return r.db.WithContext(ctx).Create(room).Error
}

func (r *roomRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Room, error) {
	var room domain.Room
	err := r.db.WithContext(ctx).Preload("Participants.User").Where("id = ?", id).First(&room).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *roomRepository) Update(ctx context.Context, room *domain.Room) error {
	return r.db.WithContext(ctx).Save(room).Error
}

func (r *roomRepository) Close(ctx context.Context, roomID uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&domain.Room{}).
		Where("id = ?", roomID).
		Updates(map[string]interface{}{
			"status":    domain.RoomStatusClosed,
			"closed_at": &now,
		}).Error
}

func (r *roomRepository) ListActiveLobbies(ctx context.Context, targetLanguage, proficiencyLevel string, offset, limit int) ([]domain.Room, int64, error) {
	var rooms []domain.Room
	var count int64

	query := r.db.WithContext(ctx).
		Model(&domain.Room{}).
		Where("type = ? AND status IN (?, ?)", domain.RoomTypeLobby, domain.RoomStatusWaiting, domain.RoomStatusActive)

	if targetLanguage != "" {
		query = query.Where("target_language = ?", targetLanguage)
	}
	if proficiencyLevel != "" {
		query = query.Where("proficiency_level = ?", proficiencyLevel)
	}

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Preload("Participants.User").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&rooms).Error; err != nil {
		return nil, 0, err
	}

	return rooms, count, nil
}

func (r *roomRepository) AddParticipant(ctx context.Context, participant *domain.RoomParticipant) error {
	return r.db.WithContext(ctx).Create(participant).Error
}

func (r *roomRepository) RemoveParticipant(ctx context.Context, roomID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("room_id = ? AND user_id = ?", roomID, userID).
		Delete(&domain.RoomParticipant{}).Error
}

func (r *roomRepository) GetParticipant(ctx context.Context, roomID, userID uuid.UUID) (*domain.RoomParticipant, error) {
	var part domain.RoomParticipant
	err := r.db.WithContext(ctx).
		Where("room_id = ? AND user_id = ?", roomID, userID).
		First(&part).Error
	return &part, err
}

func (r *roomRepository) UpdateParticipantStatus(ctx context.Context, roomID, userID uuid.UUID, status domain.JoinStatus) error {
	return r.db.WithContext(ctx).
		Model(&domain.RoomParticipant{}).
		Where("room_id = ? AND user_id = ?", roomID, userID).
		Update("status", status).Error
}

func (r *roomRepository) GetParticipants(ctx context.Context, roomID uuid.UUID) ([]domain.RoomParticipant, error) {
	var participants []domain.RoomParticipant
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("room_id = ?", roomID).
		Find(&participants).Error
	return participants, err
}

func (r *roomRepository) GetMaster(ctx context.Context, roomID uuid.UUID) (*domain.RoomParticipant, error) {
	var part domain.RoomParticipant
	err := r.db.WithContext(ctx).
		Where("room_id = ? AND role = ?", roomID, domain.ParticipantRoleGroupMaster).
		First(&part).Error
	return &part, err
}

func (r *roomRepository) CountParticipants(ctx context.Context, roomID uuid.UUID) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.RoomParticipant{}).
		Where("room_id = ? AND status = ?", roomID, domain.JoinStatusAccepted).
		Count(&count).Error
	return int(count), err
}

func (r *roomRepository) IsUserInRoom(ctx context.Context, roomID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.RoomParticipant{}).
		Where("room_id = ? AND user_id = ?", roomID, userID).
		Count(&count).Error
	return count > 0, err
}
