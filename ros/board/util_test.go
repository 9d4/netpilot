package board

import (
	"gorm.io/gorm"
	"reflect"
	"testing"
	"time"
)

func Test_updateBoard(t *testing.T) {
	type args struct {
		board       *Board
		wantedBoard *Board
		req         *CreateBoardRequest
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "update 1",
			args: args{
				board: &Board{
					ID:                 0,
					CreatedAt:          time.UnixMilli(1213131212),
					UpdatedAt:          time.UnixMilli(1213131212),
					DeletedAt:          gorm.DeletedAt{},
					UUID:               "this-is-uuid",
					Host:               "localhost",
					Port:               "8080",
					InsecureSkipVerify: true,
					User:               "",
					Password:           "",
				},
				wantedBoard: &Board{
					ID:                 0,
					CreatedAt:          time.UnixMilli(1213131212),
					UpdatedAt:          time.UnixMilli(1213131212),
					DeletedAt:          gorm.DeletedAt{},
					UUID:               "this-is-uuid",
					Host:               "127.0.1.1",
					Port:               "443",
					InsecureSkipVerify: false,
					User:               "admin",
					Password:           "admin",
				},
				req: &CreateBoardRequest{
					Host:               "127.0.1.1",
					Port:               "443",
					InsecureSkipVerify: false,
					User:               "admin",
					Password:           "admin",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateBoard(tt.args.board, tt.args.req)
			if reflect.DeepEqual(tt.args.board, tt.args.wantedBoard) == false {
				t.Fatalf("something not updated wanted: \n%+v \ngot:\n%+v", tt.args.wantedBoard, tt.args.board)
			}
		})
	}
}
