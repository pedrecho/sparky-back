package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
	"io"
	"mime/multipart"
	"os"
	"slices"
	"sparky-back/internal/models"
	"time"
)

const (
	staticPath            = "static/"
	dbBufSize             = 1000
	clientBufSize         = 100
	defaultContextTimeout = 2 * time.Second
)

type Logic struct {
	db       *bun.DB
	dbCh     chan models.Message
	clientCh map[int64]chan models.Message
}

func NewLogic(db *bun.DB) *Logic {
	logic := &Logic{
		db:       db,
		clientCh: make(map[int64]chan models.Message),
		dbCh:     make(chan models.Message, dbBufSize),
	}
	return logic
}

func (l *Logic) AddUser(ctx context.Context, user *models.User) (int64, error) {
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

func (l *Logic) UpdateUser(ctx context.Context, user *models.User) (int64, error) {
	oldUser := new(models.User)
	err := l.db.NewSelect().Model(oldUser).Where("id = ?", user.ID).Scan(ctx)
	if err != nil {
		return 0, fmt.Errorf("select query: %w", err)
	}
	if user.ImgPath == "" {
		user.ImgPath = oldUser.ImgPath
	} else {
		l.DeleteImg(oldUser.ImgPath)
	}
	if user.Description == "" {
		user.Description = oldUser.Description
	}
	if user.Latitude == 0 {
		user.Latitude = oldUser.Latitude
	}
	if user.Longitude == 0 {
		user.Longitude = oldUser.Longitude
	}
	_, err = l.db.NewUpdate().
		Model(user).
		Set("description = ?, img_path = ?, latitude = ?, longitude = ?",
			user.Description, user.ImgPath, user.Latitude, user.Longitude).
		Where("id = ?", user.ID).
		Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("update query: %w", err)
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

func (l *Logic) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := l.db.NewSelect().Model(&user).Where("email = ?", email).Scan(ctx)
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
			message := &models.Message{
				UserID: reaction.UserID,
				ToID:   reaction.ToID,
				Time:   time.Now(),
				Text:   "",
			}
			return l.NewMessage(message)
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

func (l *Logic) NewMessage(message *models.Message) error {
	err := l.SaveMessage(context.TODO(), message)
	if err != nil {
		return fmt.Errorf("save message: %v", err)
	}
	if user, ok := l.clientCh[message.UserID]; ok {
		user <- *message
	}
	if to, ok := l.clientCh[message.ToID]; ok {
		to <- *message
	}
	return nil
}

func (l *Logic) SaveMessage(ctx context.Context, message *models.Message) error {
	_, err := l.db.NewInsert().Model(message).Exec(ctx)
	if err != nil {
		return fmt.Errorf("insert query: %v", err)
	}
	return nil
}

func (l *Logic) GetNewMessages(ctx context.Context, msg *models.Message) ([]models.Message, error) {
	messages := make([]models.Message, 0)
	err := l.db.NewSelect().
		Model(&messages).
		Where("user_id = ? OR to_id = ?", msg.UserID, msg.UserID).
		Where("time > ?", msg.Time).
		Order("time ASC").
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("select query: %w", err)
	}
	return messages, nil
}

func (l *Logic) SendMessages(ctx context.Context, send func([]byte), msg *models.Message) error {
	l.clientCh[msg.UserID] = make(chan models.Message, clientBufSize)
	messages, err := l.GetNewMessages(context.TODO(), msg)
	if err != nil {
		return fmt.Errorf("getting new messages: %w", err)
	}

	for i := range messages {
		jsonData, err := json.Marshal(messages[i])
		if err != nil {
			return fmt.Errorf("marshaling json: %w", err)
		}
		send(jsonData)
	}

	for {
		select {
		case message := <-l.clientCh[msg.UserID]:
			jsonData, err := json.Marshal(message)
			if err != nil {
				return fmt.Errorf("marshaling json: %w", err)
			}
			send(jsonData)
		case <-ctx.Done():
			delete(l.clientCh, msg.UserID)
			return nil
		}
	}
}

func (l *Logic) GetRecommendations(ctx context.Context, filter *models.Filter) ([]models.User, error) {
	user := new(models.User)
	err := l.db.NewSelect().
		Model(user).
		Relation("Reactions").
		Where("id = ?", filter.UserID).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("user select query: %w", err)
	}
	//это костыль
	reactedUserIDs := make([]int64, 1, len(user.Reactions)+1)
	for i := range user.Reactions {
		reactedUserIDs = append(reactedUserIDs, user.Reactions[i].ToID)
	}
	users := make([]models.User, 0)
	err = l.db.NewSelect().
		Model(&users).
		Where("id NOT IN (?)", bun.In(reactedUserIDs)).
		Where("sex = ?", filter.Sex).
		Where("EXTRACT(YEAR FROM AGE(CURRENT_TIMESTAMP, birthday)) BETWEEN ? AND ?", filter.MinAge, filter.MaxAge).
		Where("calculate_distance(latitude, longitude, ?, ?, 'K') < ?", user.Latitude, user.Longitude, filter.Distance).
		Limit(filter.Limit).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("users select query: %w", err)
	}
	return users, nil
}

func (l *Logic) DeleteImg(imagePath string) error {
	return os.Remove(imagePath)
}
