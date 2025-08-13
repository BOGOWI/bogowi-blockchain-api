package storage

import (
	"context"
	"sync"
	"time"

	"bogowi-blockchain-go/internal/models"
)

// RewardsStorage interface defines methods for storing and retrieving reward data
type RewardsStorage interface {
	// Reward Claims
	CreateRewardClaim(ctx context.Context, claim *models.RewardClaim) error
	GetRewardClaim(ctx context.Context, id uint) (*models.RewardClaim, error)
	GetRewardClaimsByWallet(ctx context.Context, wallet string, limit int) ([]*models.RewardClaim, error)
	UpdateRewardClaimStatus(ctx context.Context, id uint, status string, txHash string) error

	// Referral Claims
	CreateReferralClaim(ctx context.Context, claim *models.ReferralClaim) error
	GetReferralClaimsByWallet(ctx context.Context, wallet string, limit int) ([]*models.ReferralClaim, error)
	UpdateReferralClaimStatus(ctx context.Context, id uint, status string, txHash string) error

	// Templates
	SaveRewardTemplate(ctx context.Context, template *models.RewardTemplate) error
	GetRewardTemplate(ctx context.Context, id, network string) (*models.RewardTemplate, error)
	GetAllRewardTemplates(ctx context.Context, network string, activeOnly bool) ([]*models.RewardTemplate, error)

	// Eligibility
	SaveUserEligibility(ctx context.Context, eligibility *models.UserRewardEligibility) error
	GetUserEligibility(ctx context.Context, userID, templateID, network string) (*models.UserRewardEligibility, error)
}

// InMemoryRewardsStorage is an in-memory implementation of RewardsStorage
// This can be replaced with a proper database implementation later
type InMemoryRewardsStorage struct {
	mu                sync.RWMutex
	rewardClaims      map[uint]*models.RewardClaim
	referralClaims    map[uint]*models.ReferralClaim
	templates         map[string]*models.RewardTemplate
	eligibilities     map[string]*models.UserRewardEligibility
	claimsByWallet    map[string][]uint
	referralsByWallet map[string][]uint
	nextID            uint
}

// NewInMemoryRewardsStorage creates a new in-memory storage instance
func NewInMemoryRewardsStorage() *InMemoryRewardsStorage {
	storage := &InMemoryRewardsStorage{
		rewardClaims:      make(map[uint]*models.RewardClaim),
		referralClaims:    make(map[uint]*models.ReferralClaim),
		templates:         make(map[string]*models.RewardTemplate),
		eligibilities:     make(map[string]*models.UserRewardEligibility),
		claimsByWallet:    make(map[string][]uint),
		referralsByWallet: make(map[string][]uint),
		nextID:            1,
	}

	// Initialize with default templates
	storage.initializeDefaultTemplates()

	return storage
}

func (s *InMemoryRewardsStorage) initializeDefaultTemplates() {
	now := time.Now()

	templates := []*models.RewardTemplate{
		{
			ID:                 "welcome_bonus",
			Name:               "Welcome Bonus",
			Description:        "One-time welcome bonus for new users",
			FixedAmount:        "10000000000000000000", // 10 BOGO
			MaxAmount:          "10000000000000000000",
			CooldownPeriod:     0,
			MaxClaimsPerWallet: 1,
			RequiresWhitelist:  false,
			Active:             true,
			Network:            "camino",
			CreatedAt:          now,
			UpdatedAt:          now,
		},
		{
			ID:                 "founder_bonus",
			Name:               "Founder Bonus",
			Description:        "Exclusive bonus for whitelisted founders",
			FixedAmount:        "100000000000000000000", // 100 BOGO
			MaxAmount:          "100000000000000000000",
			CooldownPeriod:     0,
			MaxClaimsPerWallet: 1,
			RequiresWhitelist:  true,
			Active:             true,
			Network:            "camino",
			CreatedAt:          now,
			UpdatedAt:          now,
		},
		{
			ID:                 "first_nft_mint",
			Name:               "First NFT Mint Reward",
			Description:        "Reward for minting your first NFT",
			FixedAmount:        "25000000000000000000", // 25 BOGO
			MaxAmount:          "25000000000000000000",
			CooldownPeriod:     86400, // 24 hours
			MaxClaimsPerWallet: 1,
			RequiresWhitelist:  false,
			Active:             true,
			Network:            "camino",
			CreatedAt:          now,
			UpdatedAt:          now,
		},
		{
			ID:                 "dao_participation",
			Name:               "DAO Participation Reward",
			Description:        "Reward for participating in DAO governance",
			FixedAmount:        "5000000000000000000",  // 5 BOGO
			MaxAmount:          "50000000000000000000", // Max 50 BOGO
			CooldownPeriod:     604800,                 // 7 days
			MaxClaimsPerWallet: 10,
			RequiresWhitelist:  false,
			Active:             true,
			Network:            "camino",
			CreatedAt:          now,
			UpdatedAt:          now,
		},
	}

	for _, template := range templates {
		key := template.ID + ":" + template.Network
		s.templates[key] = template
	}
}

// CreateRewardClaim stores a new reward claim
func (s *InMemoryRewardsStorage) CreateRewardClaim(ctx context.Context, claim *models.RewardClaim) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	claim.ID = s.nextID
	s.nextID++
	claim.CreatedAt = time.Now()
	claim.UpdatedAt = time.Now()

	s.rewardClaims[claim.ID] = claim
	s.claimsByWallet[claim.WalletAddress] = append(s.claimsByWallet[claim.WalletAddress], claim.ID)

	return nil
}

// GetRewardClaim retrieves a reward claim by ID
func (s *InMemoryRewardsStorage) GetRewardClaim(ctx context.Context, id uint) (*models.RewardClaim, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	claim, exists := s.rewardClaims[id]
	if !exists {
		return nil, nil
	}

	return claim, nil
}

// GetRewardClaimsByWallet retrieves reward claims for a wallet
func (s *InMemoryRewardsStorage) GetRewardClaimsByWallet(ctx context.Context, wallet string, limit int) ([]*models.RewardClaim, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	claimIDs, exists := s.claimsByWallet[wallet]
	if !exists {
		return []*models.RewardClaim{}, nil
	}

	var claims []*models.RewardClaim
	count := 0

	// Get claims in reverse order (most recent first)
	for i := len(claimIDs) - 1; i >= 0 && (limit == 0 || count < limit); i-- {
		if claim, exists := s.rewardClaims[claimIDs[i]]; exists {
			claims = append(claims, claim)
			count++
		}
	}

	return claims, nil
}

// UpdateRewardClaimStatus updates the status of a reward claim
func (s *InMemoryRewardsStorage) UpdateRewardClaimStatus(ctx context.Context, id uint, status string, txHash string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	claim, exists := s.rewardClaims[id]
	if !exists {
		return nil
	}

	claim.Status = status
	if txHash != "" {
		claim.TxHash = txHash
	}
	claim.UpdatedAt = time.Now()

	return nil
}

// CreateReferralClaim stores a new referral claim
func (s *InMemoryRewardsStorage) CreateReferralClaim(ctx context.Context, claim *models.ReferralClaim) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	claim.ID = s.nextID
	s.nextID++
	claim.CreatedAt = time.Now()
	claim.UpdatedAt = time.Now()

	s.referralClaims[claim.ID] = claim
	s.referralsByWallet[claim.ReferredAddress] = append(s.referralsByWallet[claim.ReferredAddress], claim.ID)

	return nil
}

// GetReferralClaimsByWallet retrieves referral claims for a wallet
func (s *InMemoryRewardsStorage) GetReferralClaimsByWallet(ctx context.Context, wallet string, limit int) ([]*models.ReferralClaim, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	claimIDs, exists := s.referralsByWallet[wallet]
	if !exists {
		return []*models.ReferralClaim{}, nil
	}

	var claims []*models.ReferralClaim
	count := 0

	// Get claims in reverse order (most recent first)
	for i := len(claimIDs) - 1; i >= 0 && (limit == 0 || count < limit); i-- {
		if claim, exists := s.referralClaims[claimIDs[i]]; exists {
			claims = append(claims, claim)
			count++
		}
	}

	return claims, nil
}

// UpdateReferralClaimStatus updates the status of a referral claim
func (s *InMemoryRewardsStorage) UpdateReferralClaimStatus(ctx context.Context, id uint, status string, txHash string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	claim, exists := s.referralClaims[id]
	if !exists {
		return nil
	}

	claim.Status = status
	if txHash != "" {
		claim.TxHash = txHash
	}
	claim.UpdatedAt = time.Now()

	return nil
}

// SaveRewardTemplate saves or updates a reward template
func (s *InMemoryRewardsStorage) SaveRewardTemplate(ctx context.Context, template *models.RewardTemplate) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := template.ID + ":" + template.Network
	if existing, exists := s.templates[key]; exists {
		template.CreatedAt = existing.CreatedAt
	} else {
		template.CreatedAt = time.Now()
	}
	template.UpdatedAt = time.Now()

	s.templates[key] = template
	return nil
}

// GetRewardTemplate retrieves a reward template
func (s *InMemoryRewardsStorage) GetRewardTemplate(ctx context.Context, id, network string) (*models.RewardTemplate, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := id + ":" + network
	template, exists := s.templates[key]
	if !exists {
		return nil, nil
	}

	return template, nil
}

// GetAllRewardTemplates retrieves all reward templates
func (s *InMemoryRewardsStorage) GetAllRewardTemplates(ctx context.Context, network string, activeOnly bool) ([]*models.RewardTemplate, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var templates []*models.RewardTemplate

	for _, template := range s.templates {
		if template.Network == network {
			if !activeOnly || template.Active {
				templates = append(templates, template)
			}
		}
	}

	return templates, nil
}

// SaveUserEligibility saves or updates a user's eligibility
func (s *InMemoryRewardsStorage) SaveUserEligibility(ctx context.Context, eligibility *models.UserRewardEligibility) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := eligibility.UserID + ":" + eligibility.TemplateID + ":" + eligibility.Network
	s.eligibilities[key] = eligibility

	return nil
}

// GetUserEligibility retrieves a user's eligibility for a template
func (s *InMemoryRewardsStorage) GetUserEligibility(ctx context.Context, userID, templateID, network string) (*models.UserRewardEligibility, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := userID + ":" + templateID + ":" + network
	eligibility, exists := s.eligibilities[key]
	if !exists {
		return nil, nil
	}

	return eligibility, nil
}
