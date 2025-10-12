package auth

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"shedoo-backend/internal/repositories"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type SessionData struct {
	JTI          string                            `json:"jti"`
	Sub          string                            `json:"sub"` // cmuitaccount
	BasicInfo    *repositories.CmuEntraIDBasicInfo `json:"basic_info"`
	AccessToken  string                            `json:"access_token"`
	RefreshToken string                            `json:"refresh_token"`
	ExpiresAt    time.Time                         `json:"expires_at"` // when the access_token expires (from provider)
}

type AuthService struct {
	repo        *repositories.AuthRepository
	JwtSecret   string
	redisClient *redis.Client
	sessionTTL  time.Duration
	jwtTTL      time.Duration
	issuer      string
	audience    string
}

func NewAuthService(repo *repositories.AuthRepository, jwtSecret string, rdb *redis.Client, sessionTTL, jwtTTL time.Duration, issuer, audience string) *AuthService {
	return &AuthService{
		repo:        repo,
		JwtSecret:   jwtSecret,
		redisClient: rdb,
		sessionTTL:  sessionTTL,
		jwtTTL:      jwtTTL,
		issuer:      issuer,
		audience:    audience,
	}
}

func (s *AuthService) SignIn(ctx context.Context, code string) (signedJWT string, cookieExpiry time.Time, err error) {
	tr, err := s.repo.ExchangeCode(ctx, code)
	if err != nil {
		return "", time.Time{}, err
	}

	basicInfo, err := s.repo.GetBasicInfo(ctx, tr.AccessToken)
	if err != nil {
		return "", time.Time{}, err
	}

	// create jti and session
	jti := uuid.NewString()
	now := time.Now().UTC()
	jwtExp := now.Add(s.jwtTTL)
	session := SessionData{
		JTI:          jti,
		Sub:          basicInfo.Cmuitaccount,
		BasicInfo:    basicInfo,
		AccessToken:  tr.AccessToken,
		RefreshToken: tr.RefreshToken,
		ExpiresAt:    now.Add(time.Duration(tr.ExpiresIn) * time.Second),
	}

	// save session in redis keyed by jti
	key := "session:" + jti
	bs, _ := json.Marshal(session)
	if err := s.redisClient.Set(ctx, key, bs, s.sessionTTL).Err(); err != nil {
		return "", time.Time{}, err
	}

	claims := jwt.MapClaims{
		"sub": basicInfo.Cmuitaccount,
		"jti": jti,
		"iss": s.issuer,
		"aud": s.audience,
		"iat": now.Unix(),
		"exp": jwtExp.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.JwtSecret))
	if err != nil {
		// best effort: clean up session if signing fails
		_ = s.redisClient.Del(ctx, key).Err()
		return "", time.Time{}, err
	}

	// cookie expiry: choose how long cookie should persist; here match sessionTTL or keep short
	cookieExpiry = now.Add(s.sessionTTL)
	return signed, cookieExpiry, nil
}

// helper to get session by jti
func (s *AuthService) GetSessionByJTI(ctx context.Context, jti string) (*SessionData, error) {
	key := "session:" + jti
	val, err := s.redisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, errors.New("session not found")
		}
		return nil, err
	}
	var sess SessionData
	if err := json.Unmarshal([]byte(val), &sess); err != nil {
		return nil, err
	}
	return &sess, nil
}

// revoke session
func (s *AuthService) RevokeSession(ctx context.Context, jti string) error {
	key := "session:" + jti
	return s.redisClient.Del(ctx, key).Err()
}
