package mischief

import "errors"

var (
	ErrGettingBrowser           = errors.New("error getting browser from pool")
	ErrOpeningPage              = errors.New("error when opening page")
	ErrWaitingForPageToBeStable = errors.New("error when waiting for page to be stable")
	ErrWhileTakingScreenshot    = errors.New("error while taking screenshot")
)
