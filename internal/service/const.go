package service

const (
	defaultAssistantPassword = "assistant123"

	prefixRedisRefreshToken   = "refresh-token:"
	prefixRedisBlacklistToken = "blacklist-token:"
	prefixRedisPasswordChange = "password-change:"
)

const emailPasswordChange = `
Confirm to change your password!
You need to click <a href="%s">this link</a> to confirm.
`
