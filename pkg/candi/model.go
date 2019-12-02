package candi

import "time"

type User struct {
	DisplayName     string  `json:"display_name"`
	EmailAddress    string  `json:"email_address,omitempty"`
	ID              *string `json:"id"`
	PhoneNumber     string  `json:"phone_number"`
	ProfileImageURL *string `json:"profile_image_url,omitempty"`
}

type UserPhoneNumber struct {
	UserID      string `json:"user_id"`
	PhoneNumber string `json:"phone_number"`
}

type UserPhoneNumberList struct {
	Items []UserPhoneNumber `json:"items"`
}

type TerminalRegisterPostResponse struct {
	TerminalID     string     `json:"terminal_id"`
	TerminalSecret string     `json:"terminal_secret,omitempty"`
	CodeExpiry     *time.Time `json:"code_expiry,omitempty"`
}
