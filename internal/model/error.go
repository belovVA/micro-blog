package model

import "errors"

var ErrUserNotFound = errors.New("user not found")
var ErrPostNotFound = errors.New("post not found")
var ErrLikeQueue = errors.New("likeQueue not attached")
