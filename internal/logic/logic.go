package logic

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
	"io"
	"mime/multipart"
	"os"
	"slices"
	"sparky-back/internal/models"
)

const staticPath = "static/"

type Logic struct {
	db *bun.DB
}

func NewLogic(db *bun.DB) *Logic {
	return &Logic{
		db: db,
	}
}

func (l *Logic) SaveUser(ctx context.Context, user *models.User) (int64, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("hashing password: %w", err)
	}
	user.Password = string(hashedPassword)
	_, err = l.db.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("insert query: %w", err)
	}
	return user.ID, nil
}

func (l *Logic) SaveImg(file multipart.File, filename string) (string, error) {
	filePath := staticPath + uuid.New().String() + filename
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("error creating image file: %v", err)
	}
	defer dst.Close()
	if _, err = io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("error copying image file: %v", err)
	}
	return filePath, nil
}

func (l *Logic) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	var user models.User
	err := l.db.NewSelect().Model(&user).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("select query: %w", err)
	}
	return &user, nil
}

func (l *Logic) GetFile(filename string) ([]byte, error) {
	file, err := os.Open(staticPath + filename)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}
	return data, nil
}

func (l *Logic) LogIn(ctx context.Context, email, password string) (int64, error) {
	var user models.User
	err := l.db.NewSelect().Model(&user).Where("email = ?", email).Scan(ctx)
	if err != nil {
		return 0, fmt.Errorf("select query: %w", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return 0, fmt.Errorf("comparing passwords: %w", err)
	}
	return user.ID, nil
}

func (l *Logic) SetReaction(ctx context.Context, reaction *models.Reaction) error {
	_, err := l.db.NewInsert().Model(reaction).Exec(ctx)
	if err != nil {
		return fmt.Errorf("insert reaction: %w", err)
	}
	user := new(models.User)
	err = l.db.NewSelect().Model(user).Relation("Reactions").Where("id = ?", reaction.ToID).Scan(ctx)
	if err != nil {
		return fmt.Errorf("select to user: %w", err)
	}
	if i := slices.IndexFunc(user.Reactions, func(r models.Reaction) bool {
		return reaction.UserID == r.ToID
	}); i != -1 {
		if user.Reactions[i].Like {
			chat := &models.Chat{}
			_, err = l.db.NewInsert().Model(chat).Exec(ctx)
			if err != nil {
				return fmt.Errorf("new chat: %w", err)
			}
			userChat := &models.UserChat{
				UserID: reaction.UserID,
				ChatID: chat.ID,
			}
			_, err = l.db.NewInsert().Model(userChat).Exec(ctx)
			if err != nil {
				return fmt.Errorf("user chat: %w", err)
			}
			toChat := &models.UserChat{
				UserID: reaction.ToID,
				ChatID: chat.ID,
			}
			_, err = l.db.NewInsert().Model(toChat).Exec(ctx)
			if err != nil {
				return fmt.Errorf("to chat: %w", err)
			}
		}
	} else {
		if !reaction.Like {
			_, err = l.db.NewInsert().Model(&models.Reaction{
				UserID: user.ID,
				ToID:   reaction.UserID,
				Like:   false,
			}).Exec(ctx)
			if err != nil {
				return fmt.Errorf("insert false reaction: %w", err)
			}
		}
	}
	return nil
}
