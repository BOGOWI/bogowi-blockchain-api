package models

import (
	"time"
)

// RewardClaim represents a reward claim record in the database
type RewardClaim struct {
	ID            uint      `json:"id"`
	WalletAddress string    `json:"wallet_address"`
	TemplateID    string    `json:"template_id"`
	Amount        string    `json:"amount"`
	TxHash        string    `json:"tx_hash"`
	Status        string    `json:"status"` // pending, completed, failed
	ClaimedAt     time.Time `json:"claimed_at"`
	Network       string    `json:"network"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ReferralClaim represents a referral claim record
type ReferralClaim struct {
	ID               uint      `json:"id"`
	ReferrerAddress  string    `json:"referrer_address"`
	ReferredAddress  string    `json:"referred_address"`
	ReferralCode     string    `json:"referral_code"`
	BonusAmount      string    `json:"bonus_amount"`
	TxHash           string    `json:"tx_hash"`
	Status           string    `json:"status"`
	ClaimedAt        time.Time `json:"claimed_at"`
	Network          string    `json:"network"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// RewardTemplate represents a reward template configuration
type RewardTemplate struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	FixedAmount        string    `json:"fixed_amount"`
	MaxAmount          string    `json:"max_amount"`
	CooldownPeriod     uint64    `json:"cooldown_period"`
	MaxClaimsPerWallet uint64    `json:"max_claims_per_wallet"`
	RequiresWhitelist  bool      `json:"requires_whitelist"`
	Active             bool      `json:"active"`
	Network            string    `json:"network"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// UserRewardEligibility represents a user's eligibility for rewards
type UserRewardEligibility struct {
	ID             uint      `json:"id"`
	UserID         string    `json:"user_id"`
	WalletAddress  string    `json:"wallet_address"`
	TemplateID     string    `json:"template_id"`
	IsEligible     bool      `json:"is_eligible"`
	Reason         string    `json:"reason"`
	LastChecked    time.Time `json:"last_checked"`
	NextEligibleAt time.Time `json:"next_eligible_at"`
	ClaimCount     uint      `json:"claim_count"`
	Network        string    `json:"network"`
}