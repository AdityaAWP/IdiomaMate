package domain

import "errors"

// --- Authentication & User Errors ---
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrShadowBanned       = errors.New("account restricted")
)

// --- Room & Matchmaking Errors ---
var (
	ErrRoomNotFound     = errors.New("room not found")
	ErrRoomFull         = errors.New("room is full")
	ErrRoomClosed       = errors.New("room is already closed")
	ErrAlreadyInRoom    = errors.New("user is already in this room")
	ErrNotRoomMaster    = errors.New("only the group master can perform this action")
	ErrNoMatchAvailable = errors.New("no matching partner available")
)

// --- Friendship Errors ---
var (
	ErrFriendshipNotFound   = errors.New("friendship not found")
	ErrAlreadyFriends       = errors.New("users are already friends")
	ErrFriendRequestExists  = errors.New("a pending friend request already exists")
	ErrCannotFriendSelf     = errors.New("cannot send friend request to yourself")
	ErrNotFriends           = errors.New("users are not friends")
)

// --- General Errors ---
var (
	ErrNotFound        = errors.New("resource not found")
	ErrBadRequest      = errors.New("bad request")
	ErrInternalServer  = errors.New("internal server error")
	ErrForbidden       = errors.New("forbidden")
	ErrDuplicateRating = errors.New("you have already rated this user for this session")
)
