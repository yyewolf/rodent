package mischief

import "errors"

var (
	ErrGettingBrowser           = errors.New("error getting browser from pool")
	ErrGettingPage              = errors.New("error when getting page")
	ErrNavigatingToPage         = errors.New("error when navigating to page")
	ErrWaitingForPageToBeStable = errors.New("error when waiting for page to be stable")
	ErrWhileTakingScreenshot    = errors.New("error while taking screenshot")
)
