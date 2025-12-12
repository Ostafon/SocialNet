package repos

import (
	"errors"
	"gorm.io/gorm"
	"socialnet/services/chat/internal/model"
)

type ChatRepo struct {
	db *gorm.DB
}

func NewChatRepo(db *gorm.DB) *ChatRepo {
	return &ChatRepo{db: db}
}

//
// ─── CREATE ────────────────────────────────────────────────────────────────
//

func (r *ChatRepo) CreateChat(chat *model.Chat) error {
	return r.db.Create(chat).Error
}

func (r *ChatRepo) AddParticipant(chatID uint, userID string) error {
	p := &model.Participant{
		ChatID: chatID,
		UserID: userID,
	}
	return r.db.Create(p).Error
}

func (r *ChatRepo) SaveMessage(m *model.Message) error {
	return r.db.Create(m).Error
}

//
// ─── LIST MESSAGES ─────────────────────────────────────────────────────────
//

func (r *ChatRepo) GetMessages(chatID uint, limit, offset int) ([]model.Message, error) {
	var msgs []model.Message

	err := r.db.
		Where("chat_id = ?", chatID).
		Order("created_at ASC").
		Offset(offset).
		Limit(limit).
		Find(&msgs).Error

	return msgs, err
}

func (r *ChatRepo) GetLastMessage(chatID uint) (*model.Message, error) {
	var msg model.Message

	err := r.db.
		Where("chat_id = ?", chatID).
		Order("created_at DESC").
		Limit(1).
		First(&msg).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &msg, err
}

//
// ─── LIST CHATS BY USER ────────────────────────────────────────────────────
//

func (r *ChatRepo) GetChatsByUser(userID string) ([]model.Chat, error) {
	var chats []model.Chat

	err := r.db.
		Joins("JOIN participants ON participants.chat_id = chats.id").
		Where("participants.user_id = ?", userID).
		Preload("Participants").
		Find(&chats).Error

	return chats, err
}

//
// ─── GET PARTICIPANTS ─────────────────────────────────────────────────────
//

func (r *ChatRepo) GetChatParticipants(chatID uint) ([]model.Participant, error) {
	var participants []model.Participant

	err := r.db.
		Where("chat_id = ?", chatID).
		Find(&participants).Error

	return participants, err
}

//
// ─── FIND PRIVATE CHAT ─────────────────────────────────────────────────────
//

func (r *ChatRepo) FindPrivateChat(user1, user2 string) (*model.Chat, error) {
	var chat model.Chat

	err := r.db.
		Joins("JOIN participants p1 ON p1.chat_id = chats.id").
		Joins("JOIN participants p2 ON p2.chat_id = chats.id").
		Where("p1.user_id = ? AND p2.user_id = ?", user1, user2).
		Where("is_group = ?", false).
		First(&chat).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // приватный чат не найден
	}

	return &chat, err
}
