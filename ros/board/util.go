package board

import (
	"reflect"
	"strings"
)

// updateBoard updates board from req.
// Only field available in Board will be updated.
// If there is field with tag `field:"ro"` in the req struct,
// that field value will not be updated
func updateBoard(board *Board, req *CreateBoardRequest) {
	boardValue := reflect.ValueOf(board).Elem()
	reqType := reflect.TypeOf(req).Elem()

	for i := 0; i < reqType.NumField(); i++ {
		reqField := reqType.Field(i)
		boardField := boardValue.FieldByName(reqField.Name)

		// Ignore fields that are not set in the Board
		if !boardField.IsValid() {
			continue
		}

		// don't update readonly field tag
		if f, _ := reqType.FieldByName(reqField.Name); strings.Contains(f.Tag.Get("field"), "ro") {
			continue
		}

		reqFieldValue := reflect.ValueOf(req).Elem().Field(i)
		boardField.Set(reqFieldValue)
	}
}
