package constants

import (
	"errors"
)

var (
	ErrCantParseBody                  = errors.New("cant parse body bad request")
	ErrPasswordMismatch               = errors.New("username or password not found")
	ErrOldPasswordNotMatch            = errors.New("old password is incorrect")
	ErrBadHostPortRegisterFormat      = errors.New("bad host:port format. Eg: localhost:8080")
	ErrBadPortType                    = errors.New("bad port type should be integer. Eg: 8080")
	ErrServiceUnavailable             = "%s unavailable"
	ErrChatRoomAlreadyExistInUserUuid = "chatroom with uuid %s already exist in user %s chatroom"
	ErrUserDeleted                    = "user was deleted at %s"
)
