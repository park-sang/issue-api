package store

import (
	"issue-api/models"
	"sync"
)

var (
	Issues      = make(map[uint]*models.Issue)
	NextID uint = 1
	Mutex  sync.Mutex
)
