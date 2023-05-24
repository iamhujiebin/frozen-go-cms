package chatgpt_m

import (
	"git.hilo.cn/hilo-common/domain"
	"git.hilo.cn/hilo-common/resource/mysql"
	"gorm.io/gorm/clause"
)

type ChatgptSession struct {
	mysql.Entity
	UserId    mysql.ID
	SessionId mysql.ID
	Message   string
}

// 获取用户会话列表
// 一个都没有则初始化一个
func GetUserSessionsInit(model *domain.Model, userId mysql.ID) ([]ChatgptSession, error) {
	var sessions []ChatgptSession
	if err := model.DB().Model(ChatgptSession{}).Where("user_id = ?", userId).Find(&sessions).Error; err != nil {
		return sessions, err
	}
	if len(sessions) <= 0 {
		if err := model.DB().Model(ChatgptSession{}).Create(&ChatgptSession{
			UserId:    userId,
			SessionId: 0,
		}).Error; err != nil {
			return sessions, err
		}
		if err := model.DB().Model(ChatgptSession{}).Where("user_id = ?", userId).Find(&sessions).Error; err != nil {
			return sessions, err
		}
	}
	return sessions, nil
}

// 获取用户指定会话
func GetUserSession(model *domain.Model, userId, sessionId mysql.ID) (ChatgptSession, error) {
	var session ChatgptSession
	if err := model.DB().Model(ChatgptSession{}).
		Where("user_id = ? AND session_id = ?", userId, sessionId).First(&session).Error; err != nil {
		return session, err
	}
	return session, nil
}

// 更新用户会话
func UpdateSessionInit(model *domain.Model, session ChatgptSession) error {
	return model.DB().Model(ChatgptSession{}).
		Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "user_id"}, {Name: "session_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"message": session.Message,
			}),
		}).Create(&session).Error
}

// 更新用户会话
func CreateSessionInit(model *domain.Model, userId mysql.ID) (mysql.ID, error) {
	var maxSession ChatgptSession
	if err := model.DB().Model(ChatgptSession{}).Where("user_id = ?", userId).Order("session_id DESC").First(&maxSession).Error; err != nil {
		if err := model.DB().Model(ChatgptSession{}).Create(&maxSession).Error; err != nil {
			return 0, err
		}
		return maxSession.SessionId, nil
	}
	maxSession = ChatgptSession{UserId: maxSession.UserId, SessionId: maxSession.SessionId + 1}
	if err := model.DB().Model(ChatgptSession{}).Create(&maxSession).Error; err != nil {
		return 0, err
	}
	return maxSession.SessionId, nil
}

// 删除一个会话
func DeleteSession(model *domain.Model, userId, sessionId mysql.ID) error {
	if sessionId == 0 {
		return model.DB().Model(ChatgptSession{}).Where("user_id = ? AND session_id = ?", userId, sessionId).UpdateColumn("message", "").Error
	}
	if err := model.DB().Model(ChatgptSession{}).Where("user_id = ? AND session_id = ?", userId, sessionId).Delete(&ChatgptSession{}).Error; err != nil {
		return err
	}
	return nil
}
