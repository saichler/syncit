package handlers

import (
	"github.com/saichler/syncit/model"
	"sync"
)

type FetchJob struct {
	cond          *sync.Cond
	queue         []*model.Command
	last          int32
	hadOrderIssue bool
}
