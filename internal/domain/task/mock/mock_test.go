package mock

import (
	"context"
	"testing"

	"github.com/henriqueassiss/advanced-golang-api/internal/utils/errorMsg"
	"github.com/henriqueassiss/advanced-golang-api/third_party/database"

	sqlxmock "github.com/zhashkevych/go-sqlxmock"
)

func TestMock(t *testing.T) {
	db, mock := database.NewSqlxMock(t)
	defer db.Close()

	type args struct {
		ctx context.Context
	}

	type test struct {
		name       string
		args       args
		beforeTest func()
		wantErr    error
	}

	tests := []test{
		{
			name: "Success - Create",
			args: args{
				ctx: context.TODO(),
			},
			beforeTest: func() {
				rows := mock.NewRows([]string{"count"}).AddRow(0)
				mock.ExpectQuery("SELECT").WillReturnRows(rows)

				mock.ExpectExec("INSERT INTO tasks").WillReturnResult(sqlxmock.NewResult(0, 5))
			},
		},
		{
			name: "Fail - Already populated",
			args: args{
				ctx: context.TODO(),
			},
			beforeTest: func() {
				rows := mock.NewRows([]string{"count"}).AddRow(5)
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			wantErr: errorMsg.ErrTableIsPopulated,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.beforeTest()

			err := Mock(db)
			if err != nil && err != test.wantErr {
				t.Errorf("error %v, wantErr = %v", err, test.wantErr)
			}
		})
	}
}
