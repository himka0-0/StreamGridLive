package db

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	db, err := Init()
	require.NoError(t, err, "инициализация БД должна проходить без ошибок")
	require.NotNil(t, db, "должен получить не-nil объект *gorm.DB")

	err = db.Exec("SELECT 1 FROM rooms LIMIT 1;").Error
	require.NoError(t, err, "таблица rooms должна существовать после миграции")
}
