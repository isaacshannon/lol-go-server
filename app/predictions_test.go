package main

import (
	"strings"
	"testing"
)

func TestValidateImage(t *testing.T) {
	t.Run("the image is valid", func(t *testing.T) {
		img := "data:image/png;base64," + strings.Repeat("a", 20000)
		res := validateImage(img)
		if res != nil {
			t.Fail()
		}
	})

	t.Run("the image has no file type", func(t *testing.T) {
		img := strings.Repeat("a", 20000)
		res := validateImage(img)
		if res == nil {
			t.Fail()
		}
	})

	t.Run("the image is too small", func(t *testing.T) {
		img := "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAASwAAAEsCA"
		res := validateImage(img)
		if res == nil {
			t.Fail()
		}
	})

	t.Run("the image is too large", func(t *testing.T) {
		img := "data:image/png;base64,"+strings.Repeat("a", 40000)
		res := validateImage(img)
		if res == nil {
			t.Fail()
		}
	})
}
