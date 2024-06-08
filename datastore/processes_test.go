package datastore_test

import (
	"testing"

	"github.com/ImDuong/vola-auto/datastore"
)

func TestParseFullPathByArgs(t *testing.T) {
	type testCase struct {
		input    string
		expected string
	}

	testCases := []testCase{
		{input: `%SystemRoot%\system32\csrss.exe ObjectDirectory=\Windows`, expected: `%SystemRoot%\system32\csrss.exe`},
		{input: `C:\Windows\system32\svchost.exe -k LocalServiceAndNoImpersonation`, expected: `C:\Windows\system32\svchost.exe`},
		{input: `"C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe" --profile-directory=Default`, expected: `C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe`},
		{input: `C:\simplepath.exe`, expected: `C:\simplepath.exe`},
		{input: `"C:\Path With Spaces\app.exe" argument`, expected: `C:\Path With Spaces\app.exe`},
		{input: `  C:\LeadingAndTrailingSpaces.exe  argument`, expected: `C:\LeadingAndTrailingSpaces.exe`},
		{input: `""`, expected: ``},
		{input: ` `, expected: ``},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			p := &datastore.Process{}
			p.Args = tc.input
			p.ParseFullPathByArgs()
			if p.FullPath != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, p.FullPath)
			}
		})
	}
}
