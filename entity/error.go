package entity

// Error ...
type Error struct {
	Message string `json:"message"`
}

const (
	WeakPasswordMessage              = "Password is weak"
	YourProfileHasChangedSuccusfully = "Your profile has changed succesfully"
	VerificationCodeSentYourEmail    = "Verification code sent  your email"
	SomethingWentWrong               = "Ooops something went wrong"
	ObjectNotFount                   = "Opject Not Fount"
)
