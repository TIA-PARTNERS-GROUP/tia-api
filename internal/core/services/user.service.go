package services

import (
	"context"
	"errors"
	"strings"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"github.com/TIA-PARTNERS-GROUP/tia-api/pkg/utils"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) CreateUser(ctx context.Context, data ports.UserCreationSchema) (*models.User, error) {
	if err := utils.ValidatePasswordComplexity(data.Password); err != nil {
		return nil, ports.ErrPasswordComplexity
	}

	hashedPassword, err := utils.HashPassword(data.Password)
	if err != nil {
		return nil, ports.ErrDatabase
	}

	user := models.User{
		FirstName:      data.FirstName,
		LastName:       &data.LastName,
		LoginEmail:     data.LoginEmail,
		PasswordHash:   &hashedPassword,
		ContactEmail:   &data.ContactEmail,
		ContactPhoneNo: &data.ContactPhoneNo,
		AdkSessionID:   &data.AdkSessionID,
		Active:         true,
		EmailVerified:  false,
	}

	if err := s.db.WithContext(ctx).Create(&user).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return nil, ports.ErrUserAlreadyExists
		}
		return nil, ports.ErrDatabase
	}

	return &user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id uint, data ports.UserUpdateSchema) (*models.User, error) {
	user, err := s.FindUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	updateData := make(map[string]interface{})

	if data.FirstName != nil {
		updateData["first_name"] = *data.FirstName
	}
	if data.LastName != nil {
		updateData["last_name"] = *data.LastName
	}
	if data.LoginEmail != nil {
		updateData["login_email"] = *data.LoginEmail
	}
	if data.ContactEmail != nil {
		updateData["contact_email"] = *data.ContactEmail
	}
	if data.ContactPhoneNo != nil {
		updateData["contact_phone_no"] = *data.ContactPhoneNo
	}
	if data.AdkSessionID != nil {
		updateData["adk_session_id"] = *data.AdkSessionID
	}
	if data.EmailVerified != nil {
		updateData["email_verified"] = *data.EmailVerified
	}
	if data.Active != nil {
		updateData["active"] = *data.Active
	}

	if data.Password != nil && *data.Password != "" {
		if err := utils.ValidatePasswordComplexity(*data.Password); err != nil {
			return nil, ports.ErrPasswordComplexity
		}
		hashedPassword, err := utils.HashPassword(*data.Password)
		if err != nil {
			return nil, ports.ErrDatabase
		}
		updateData["password_hash"] = hashedPassword
	}

	if len(updateData) == 0 {
		return nil, ports.ErrNoUpdateData
	}

	if err := s.db.WithContext(ctx).Model(user).Updates(updateData).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return nil, ports.ErrUserAlreadyExists
		}
		return nil, ports.ErrDatabase
	}

	return user, nil
}

func (s *UserService) ChangePassword(ctx context.Context, userID uint, currentPassword, newPassword string) error {
	var user models.User
	if err := s.db.WithContext(ctx).Select("password_hash").First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ports.ErrUserNotFound
		}
		return ports.ErrDatabase
	}

	if user.PasswordHash == nil || *user.PasswordHash == "" {
		return ports.ErrInvalidCredentials
	}

	if err := utils.VerifyPassword(currentPassword, *user.PasswordHash); err != nil {
		return ports.ErrIncorrectPassword
	}

	if err := utils.ValidatePasswordComplexity(newPassword); err != nil {
		return ports.ErrPasswordComplexity
	}

	newHashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return ports.ErrDatabase
	}

	if err := s.db.WithContext(ctx).Model(&models.User{ID: userID}).Update("password_hash", newHashedPassword).Error; err != nil {
		return ports.ErrDatabase
	}

	return nil
}

func (s *UserService) FindUserByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	if err := s.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ports.ErrUserNotFound
		}
		return nil, ports.ErrDatabase
	}
	return &user, nil
}

func (s *UserService) FindAllUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User
	if err := s.db.WithContext(ctx).Order("created_at desc").Find(&users).Error; err != nil {
		return nil, ports.ErrDatabase
	}
	return users, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", id).Delete(&models.UserSession{}).Error; err != nil {
			return err
		}

		result := tx.Delete(&models.User{}, id)
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return ports.ErrUserNotFound
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, ports.ErrUserNotFound) {
			return ports.ErrUserNotFound
		}
		return ports.ErrDatabase
	}

	return nil
}

func (s *UserService) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := s.db.WithContext(ctx).Where("login_email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ports.ErrUserNotFound
		}
		return nil, ports.ErrDatabase
	}
	return &user, nil
}

func (s *UserService) DeactivateUser(ctx context.Context, id uint) (*models.User, error) {
	user, err := s.FindUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !user.Active {
		return user, nil // Already inactive, no-op
	}

	if err := s.db.WithContext(ctx).Model(user).Update("active", false).Error; err != nil {
		return nil, ports.ErrDatabase
	}

	user.Active = false
	return user, nil
}
