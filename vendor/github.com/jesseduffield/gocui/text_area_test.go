package gocui

import (
	"reflect"
	"testing"
)

func Test_AutoWrapContent(t *testing.T) {
	tests := []struct {
		name                   string
		content                string
		autoWrapWidth          int
		expectedWrappedContent string
		expectedCursorMapping  []CursorMapping
	}{
		{
			"empty content",
			"",
			7,
			"",
			[]CursorMapping{},
		},
		{
			"auto-wrap off",
			"abcdefghijklm",
			-1,
			"abcdefghijklm",
			[]CursorMapping{},
		},
		{
			"no wrapping necessary",
			"abcde",
			7,
			"abcde",
			[]CursorMapping{},
		},
		{
			"wrap at whitespace",
			"abcde xyz",
			7,
			"abcde \nxyz",
			[]CursorMapping{{6, 7}},
		},
		{
			"lots of whitespace is preserved at end of line",
			"abcde      xyz",
			7,
			"abcde      \nxyz",
			[]CursorMapping{{11, 12}},
		},
		{
			"wrap inside word because there's no whitespace",
			"abcdefghijk",
			7,
			"abcdefg\nhijk",
			[]CursorMapping{{7, 8}},
		},
		{
			"hard line breaks",
			"abc\ndef\n",
			7,
			"abc\ndef\n",
			[]CursorMapping{},
		},
		{
			"mixture of hard and soft line breaks",
			"abc def ghi jkl mno\npqr stu vwx yz\n",
			7,
			"abc def \nghi jkl \nmno\npqr stu \nvwx yz\n",
			[]CursorMapping{{8, 9}, {16, 18}, {28, 31}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrappedContent, cursorMapping := AutoWrapContent([]rune(tt.content), tt.autoWrapWidth)
			if !reflect.DeepEqual(wrappedContent, []rune(tt.expectedWrappedContent)) {
				t.Errorf("autoWrapContentImpl() wrappedContent = %v, expected %v", string(wrappedContent), tt.expectedWrappedContent)
			}
			if !reflect.DeepEqual(cursorMapping, tt.expectedCursorMapping) {
				t.Errorf("autoWrapContentImpl() cursorMapping = %v, expected %v", cursorMapping, tt.expectedCursorMapping)
			}

			// As a sanity check, run through all runes of the original content,
			// convert the cursor to the wrapped cursor, and check that the rune
			// in the wrapped content at that position is the same:
			for i, r := range tt.content {
				wrappedIndex := origCursorToWrappedCursor(i, cursorMapping)
				if r != wrappedContent[wrappedIndex] {
					t.Errorf("Runes in orig content and wrapped content don't match at %d: expected %v, got %v", i, r, wrappedContent[wrappedIndex])
				}

				// Also, check that converting the wrapped position back to the
				// orig position yields the original value again:
				origIndexAgain := wrappedCursorToOrigCursor(wrappedIndex, cursorMapping)
				if i != origIndexAgain {
					t.Errorf("wrappedCursorToOrigCursor doesn't yield original position: expected %d, got %d", i, origIndexAgain)
				}
			}
		})
	}
}
