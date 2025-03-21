package mischief

import (
	"time"

	"github.com/go-rod/rod"
)

type Rat struct {
	CreatedAt time.Time

	*rod.Browser
}
