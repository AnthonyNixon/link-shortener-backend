package shortcode

import "math/rand"

var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")
var default_len = 6

type ShortCode struct {
	// required fields
	Len int
}

type Option func(f *ShortCode)

func Len(len int) Option {
	return func(f *ShortCode) {
		f.Len = len
	}
}

func New(opts ...Option) string {
	short := &ShortCode{Len: default_len}
	for _, opt := range opts {
		opt(short)
	}

	b := make([]rune, default_len)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
