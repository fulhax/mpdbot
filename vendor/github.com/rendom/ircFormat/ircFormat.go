package ircFormat

import "strconv"

const (
	White = iota
	Black
	Blue
	Green
	Red
	Brown
	Purple
	Orange
	Yellow
	Lime
	Teal
	Cyan
	Darkblue
	Pink
	Grey
	Lightgrey
	None
)

const (
	ResetFormatDelim = "\x0F"
	BoldDelim        = "\x02"
	ItalicDelim      = "\x1D"
	UnderlineDelim   = "\x1F"
	SwapDelim        = "\x16"
	ColorDelim       = "\x03"
)

type IrcText struct {
	text      string
	bgColor   int
	fgColor   int
	bold      bool
	italic    bool
	underline bool
	swap      bool
}

// New creates IrcText object and take your text as argument.
func New(s string) *IrcText {
	return &IrcText{text: s, bgColor: None, fgColor: None}
}

// SetFg sets foreground color
func (i *IrcText) SetFg(c int) *IrcText {
	i.fgColor = c
	return i
}

// SetBg sets background color
func (i *IrcText) SetBg(c int) *IrcText {
	i.bgColor = c
	return i
}

// Bold text
func (i *IrcText) SetBold() *IrcText {
	i.bold = true
	return i
}

// Italic text
func (i *IrcText) SetItalic() *IrcText {
	i.italic = true
	return i
}

// Underline text
func (i *IrcText) SetUnderline() *IrcText {
	i.underline = true
	return i
}

func (i *IrcText) String() string {
	c := i.text

	if i.underline {
		c = Underline(c)
	}

	if i.bold {
		c = Bold(c)
	}

	if i.italic {
		c = Italic(c)
	}

	if i.bgColor != None || i.fgColor != None {
		c = Colorize(c, i.fgColor, i.bgColor)
	}

	return c
}

func Colorize(text string, fg int, bg int) string {
	if fg == None && bg == None {
		return text
	}

	c := ColorDelim

	if fg != None {
		c += strconv.Itoa(fg)
	}

	if bg != None {
		c += ","
		c += strconv.Itoa(bg)
	}

	c += text
	c += ColorDelim

	return c
}

func Underline(text string) string {
	return wrapText(text, UnderlineDelim)
}

func Bold(text string) string {
	return wrapText(text, BoldDelim)
}

func Italic(text string) string {
	return wrapText(text, ItalicDelim)
}

func wrapText(text string, delim string) string {
	return delim + text + delim
}
