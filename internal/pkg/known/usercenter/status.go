package usercenter

const (
	// UserStatusRegistered User has submitted registration information, the account is in a pending.
	// The user needs to complete email/phone number verification steps to transition to an active state.
	// The streaming project does not currently use this.
	UserStatusRegistered = "registered"
	// UserStatusActive The user has registered and been verified, and can use the system normally.
	// Most user operations are performed in this state.
	UserStatusActive = "active"
	// UserStatusLocked The user has entered the incorrect password too many times, and the account has been locked by the system.
	// The user needs to recover the password or contact the administrator to unlock the account.
	UserStatusLocked = "locked"
	// UserStatusBlacklisted The user has been added to the system blacklist due to serious misconduct.
	// Blacklisted users cannot register new accounts or use the system.
	UserStatusBlacklisted = "blacklisted"
	// UserStatusDisabled The administrator has manually disabled the user account, and the user cannot log in after being disabled.
	// This may be due to user misconduct or other reasons.
	UserStatusDisabled = "disabled"
	// UserStatusDeleted The user has actively deleted their own account, or the administrator has deleted the user account.
	// The deleted account can be chosen to be soft-deleted (with some data retained) or completely deleted.
	UserStatusDeleted = "deleted"
)
