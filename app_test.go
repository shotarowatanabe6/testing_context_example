package main

// mock/mock_db.go を使用してテストを作成してください

import (
	"testing"

	"go.uber.org/mock/gomock"

	"testing_context_example/db"
	mock "testing_context_example/mock"
)

func isExistValueInDB(db db.DB, key string) bool {
	value := db.Get(key)
	return value != ""
}

func Test_isExistValueInDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock.NewMockDB(ctrl)

	tests := []struct {
		name       string
		mockReturn string
		want       bool
	}{
		{
			name:       "DBから値を取得できた場合",
			mockReturn: "some_value",
			want:       true,
		},
		{
			name:       "DBから空文字を取得した場合",
			mockReturn: "",
			want:       false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockDB.EXPECT().Get("test").Return(test.mockReturn)
			got := isExistValueInDB(mockDB, "test")
			if got != test.want {
				t.Errorf("isExistValueInDB() = %v, want %v", got, test.want)
			}
		})
	}
}
