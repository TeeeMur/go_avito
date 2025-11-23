package models

// ErrorResponseErrorCode defines model for ErrorResponse.Error.Code.
type ErrorResponseErrorCode string

// PullRequestEnumStatus defines model for PullRequestShort.Status.
type PullRequestEnumStatus string

// Defines values for ErrorResponseErrorCode.
const (
	NOCANDIDATE ErrorResponseErrorCode = "NO_CANDIDATE"
	NOTASSIGNED ErrorResponseErrorCode = "NOT_ASSIGNED"
	NOTFOUND    ErrorResponseErrorCode = "NOT_FOUND"
	PREXISTS    ErrorResponseErrorCode = "PR_EXISTS"
	PRMERGED    ErrorResponseErrorCode = "PR_MERGED"
	TEAMEXISTS  ErrorResponseErrorCode = "TEAM_EXISTS"
)

// Defines values for PullRequestStatus.
const (
	PullRequestStatusMERGED PullRequestEnumStatus = "MERGED"
	PullRequestStatusOPEN   PullRequestEnumStatus = "OPEN"
)
