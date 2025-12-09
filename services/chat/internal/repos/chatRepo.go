package repos

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
	"socialnet/services/chat/internal/model"
)

type ChatRepo struct {
	db *gorm.DB
}

func NewChatRepo(db *gorm.DB) *ChatRepo {
	return &ChatRepo{db: db}
}

func (r *ChatRepo) CreateChat(chat *model.Chat) error {
	return r.db.Create(chat).Error
}

func (r *ChatRepo) AddParticipant(chatID uint, userID string) error {
	return r.db.Create(&model.Participant{ChatID: chatID, UserID: userID}).Error
}

func (r *ChatRepo) GetChatByID(id uint) (*model.Chat, error) {
	var chat model.Chat
	err := r.db.Preload("Participants").First(&chat, id).Error
	return &chat, err
}

func (r *ChatRepo) ListUserChats(userID string) ([]model.Chat, error) {
	var chats []model.Chat
	err := r.db.
		Joins("JOIN participants ON participants.chat_id = chats.id").
		Where("participants.user_id = ?", userID).
		Preload("Participants").
		Find(&chats).Error
	return chats, err
}

func (r *ChatRepo) SaveMessage(msg *model.Message) error {
	return r.db.Create(msg).Error
}

func (r *ChatRepo) GetMessages(chatID uint, limit, offset int) ([]model.Message, error) {
	var msgs []model.Message
	err := r.db.Where("chat_id = ?", chatID).
		Order("created_at ASC").
		Limit(limit).
		Offset(offset).
		Find(&msgs).Error
	return msgs, err
}

func (r *ChatRepo) MarkAsRead(msgID uint) error {
	return r.db.Model(&model.Message{}).
		Where("id = ?", msgID).
		Update("read", true).Error
}

func (r *ChatRepo) GetChatsByUser(userID string) ([]model.Chat, error) {
	var chats []model.Chat
	err := r.db.Where("participants @> ?", pq.StringArray{userID}).Find(&chats).Error
	return chats, err
}

func (r *ChatRepo) GetLastMessage(chatID uint) (*model.Message, error) {
	var msg model.Message
	err := r.db.Where("chat_id = ?", chatID).Order("created_at desc").First(&msg).Error
	return &msg, err
}
