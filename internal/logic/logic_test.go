package logic

import (
	"context"
	"sparky-back/internal/loader"
	"sparky-back/internal/models"
	"testing"
)

func TestLogic_SetReaction(t *testing.T) {
	logic := NewLogic(loader.New("localhost", 5432, "postgres", "123456", "sparky"))
	err := logic.SetReaction(context.TODO(), &models.Reaction{
		UserID: 2,
		ToID:   1,
		Like:   true,
	})
	if err != nil {
		panic(err)
	}
}
