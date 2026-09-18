package issues

import "errors"

var (
	ErrIssueNotFound      = errors.New("issue not found")
	ErrIssueAlreadyExists = errors.New("issue already exists")

	ErrTitleTooShort      = errors.New("the title is too short")
	ErrTitleTooLong       = errors.New("the title is too long")
	ErrDescriptionTooLong = errors.New("the description is too long")
	ErrInvalidType        = errors.New("invalid issue type")
	ErrInvalidStatus      = errors.New("invalid issue status")
	ErrInvalidPriority    = errors.New("invalid issue priority")
	ErrInvalidParent      = errors.New("invalid parent issue")
	ErrIssueArchived      = errors.New("issue is archived")
	ErrCommentTooLong     = errors.New("the comment is too long")

	ErrCommentNotFound      = errors.New("comment not found")
	ErrCommentAlreadyExists = errors.New("comment already exists")
)
