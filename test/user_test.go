package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	models "ecommers/internal/domain"
	H "ecommers/internal/handler/http"
	"ecommers/internal/service/auth/mocks"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetProfile_Simple(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserSvc := mocks.NewMockUserService(ctrl)
	h := H.NewUserHandler(mockUserSvc, nil) // Второй мок можно в nil, если он не нужен

	userID := uuid.New()

	t.Run("Успех", func(t *testing.T) {
		// 1. Ожидаем вызов сервиса
		mockUserSvc.EXPECT().
			GetProfile(gomock.Any(), userID).
			Return(&models.User{ID: userID, Email: "test@mail.com"}, nil)

		req, _ := http.NewRequest("GET", "/profile", nil)

		// 2. ПРОЩЕ НЕКУДА: Пихаем ID всеми способами, чтобы хендлер его нашел
		ctx := req.Context()
		ctx = context.WithValue(ctx, "user_id", userID) // Вариант 1 (строковый ключ)
		ctx = context.WithValue(ctx, "userID", userID)  // Вариант 2
		// Если в хендлере ID достается как строка, а не UUID, добавим и это:
		ctx = context.WithValue(ctx, "user_id", userID.String())

		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()

		// 3. Запуск
		h.GetProfile(rr, req)

		// 4. Проверка
		assert.Equal(t, http.StatusOK, rr.Code, "Если тут 401, значит в хендлере другой ключ контекста")

		var res models.User
		json.NewDecoder(rr.Body).Decode(&res)
		assert.Equal(t, "test@mail.com", res.Email)
	})
}
