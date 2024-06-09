package utils_test

import (
	"testing"

	"github.com/ImDuong/vola-auto/utils"
)

func TestGetPathInCamelCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"c:\\windows\\system32\\lsass.exe", "C:\\Windows\\System32\\lsass.exe"},
		{"c:\\program files\\common files", "C:\\Program Files\\Common Files"},
		{"d:\\data\\my folder\\myfile.txt", "D:\\Data\\My Folder\\myfile.txt"},
		{"c:\\user\\documents\\sample file.txt", "C:\\User\\Documents\\sample file.txt"},
		{"c:\\user\\documents\\sample folder", "C:\\User\\Documents\\Sample Folder"},
	}

	for _, test := range tests {
		result := utils.GetPathInCamelCase(test.input)
		if result != test.expected {
			t.Errorf("For input '%s', expected '%s' but got '%s'", test.input, test.expected, result)
		}
	}
}
