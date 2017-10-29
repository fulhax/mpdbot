package ircFormat

import (
	"fmt"
	"testing"
)

func TestFgColor(t *testing.T) {
	text := "test"
	expect := fmt.Sprintf("\x038%s\x03", text)

	c := New(text).SetFg(Yellow).String()

	if c != expect {
		t.Fatalf("Got %q, expected %q", c, expect)
	}
}

func TestBgColor(t *testing.T) {
	text := "test"
	expect := fmt.Sprintf("\x03,6%s\x03", text)

	c := New(text).SetBg(Purple).String()

	if c != expect {
		t.Fatalf("Got %q, expected %q", c, expect)
	}
}

func TestBgAndFgColor(t *testing.T) {
	text := "test"
	expect := fmt.Sprintf("\x038,6%s\x03", text)

	c := New(text).SetFg(Yellow).SetBg(Purple).String()

	if c != expect {
		t.Fatalf("Got %q, expected %q", c, expect)
	}
}

func TestChainAllMethods(t *testing.T) {
	text := "test"
	expect := fmt.Sprintf("\x038,6\x1D\x02\x1F%s\x1F\x02\x1D\x03", text)
	c := New(text).SetFg(Yellow).SetBg(Purple).SetBold().SetItalic().SetUnderline().String()

	if c != expect {
		t.Fatalf("Got %q, expected %q", c, expect)
	}
}

func TestUnderlineText(t *testing.T) {
	text := "test"
	expect := fmt.Sprintf("\x1F%s\x1F", text)

	c := Underline(text)
	if c != expect {
		t.Fatalf("Got %q, expected %q", c, expect)
	}
}

func TestBold(t *testing.T) {
	text := "test"
	expect := fmt.Sprintf("\x02%s\x02", text)

	c := Bold(text)
	if c != expect {
		t.Fatalf("Got %q, expected %q", c, expect)
	}
}

func TestItalicText(t *testing.T) {
	text := "test"
	expect := fmt.Sprintf("\x1D%s\x1D", text)

	c := Italic(text)
	if c != expect {
		t.Fatalf("Got %q, expected %q", c, expect)
	}
}

func TestColorizeNone(t *testing.T) {
	text := "test"
	c := Colorize("test", None, None)
	if c != text {
		t.Fatalf("Got %q, expected %q", c, text)
	}
}
