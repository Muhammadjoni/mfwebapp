package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/muhammadjoni/mfwebapp/internal/config"
	"github.com/muhammadjoni/mfwebapp/internal/domain/entity"
	"github.com/muhammadjoni/mfwebapp/internal/domain/repository"
	"github.com/muhammadjoni/mfwebapp/pkg/hash"
	"github.com/muhammadjoni/mfwebapp/pkg/jwt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserBanned         = errors.New("account is banned")
	ErrEmailTaken         = errors.New("email already in use")
	ErrTokenExpired       = errors.New("token expired or invalid")
)

type AuthService struct {
	userRepo repository.UserRepository
	jwtMgr   *jwt.Manager
	hasher   *hash.Hasher
	cfg      *config.JWTConfig
}

func NewAuthService(
	userRepo repository.UserRepository,
	jwtMgr *jwt.Manager,
	hasher *hash.Hasher,
	cfg *config.JWTConfig,
) *AuthService {
	return &AuthService{userRepo: userRepo, jwtMgr: jwtMgr, hasher: hasher, cfg: cfg}
}

type RegisterInput struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
	Phone     string
	Language  entity.Language
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

func (s *AuthService) Register(ctx context.Context, in RegisterInput) (*entity.User, *TokenPair, error) {
	if existing, _ := s.userRepo.GetByEmail(ctx, in.Email); existing != nil {
		return nil, nil, ErrEmailTaken
	}

	passwordHash, err := s.hasher.Hash(in.Password)
	if err != nil {
		return nil, nil, err
	}

	if in.Language == "" {
		in.Language = entity.LangRU
	}

	user := &entity.User{
		ID:           uuid.New(),
		Email:        in.Email,
		PasswordHash: passwordHash,
		FirstName:    in.FirstName,
		LastName:     in.LastName,
		Phone:        in.Phone,
		Role:         entity.RoleUser,
		Status:       entity.UserStatusActive,
		Language:     in.Language,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, nil, err
	}

	tokens, err := s.issueTokens(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*entity.User, *TokenPair, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, nil, ErrInvalidCredentials
	}
	if user.Status == entity.UserStatusBanned {
		return nil, nil, ErrUserBanned
	}
	if !s.hasher.Verify(password, user.PasswordHash) {
		return nil, nil, ErrInvalidCredentials
	}

	tokens, err := s.issueTokens(ctx, user)
	if err != nil {
		return nil, nil, err
	}
	return user, tokens, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*TokenPair, error) {
	claims, err := s.jwtMgr.ParseRefresh(refreshToken)
	if err != nil {
		return nil, ErrTokenExpired
	}

	stored, err := s.userRepo.GetRefreshToken(ctx, refreshToken)
	if err != nil || stored == nil || stored.ExpiresAt.Before(time.Now()) {
		return nil, ErrTokenExpired
	}

	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, ErrTokenExpired
	}

	if err := s.userRepo.DeleteRefreshToken(ctx, refreshToken); err != nil {
		return nil, err
	}

	return s.issueTokens(ctx, user)
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	return s.userRepo.DeleteRefreshToken(ctx, refreshToken)
}

func (s *AuthService) issueTokens(ctx context.Context, user *entity.User) (*TokenPair, error) {
	accessToken, err := s.jwtMgr.GenerateAccess(user.ID, string(user.Role))
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtMgr.GenerateRefresh(user.ID)
	if err != nil {
		return nil, err
	}

	rt := &entity.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(s.cfg.RefreshTTL),
		CreatedAt: time.Now(),
	}

	if err := s.userRepo.SaveRefreshToken(ctx, rt); err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.cfg.AccessTTL.Seconds()),
	}, nil
}
