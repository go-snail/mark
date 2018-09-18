package mark

import "testing"

func MartTest(t *testing.T) {
	Run("core")
	mes := make(map[string]interface{})
	mes["123"] = "qwe"
	mes["456"] = "asdf"
	Mark(mes)
}