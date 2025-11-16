package core

import "errors"

var ErrNotEnoughCoworkers = errors.New("not enough coworkers to assign as reviewer. Old reviewer was cleared")
var PullRequestAlreadyMerged = errors.New("pull request is already merged")
