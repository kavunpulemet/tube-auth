package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"auth/internal/database"
	dbmodels "auth/internal/database/models"
	"auth/internal/models"
	"auth/internal/service/mappers"
	"auth/internal/utils"
)

type AuthorizationService interface {
	CreateUser(ctx utils.MyContext, user models.User) (string, error)
	Login(ctx utils.MyContext, input models.LoginInput, userIP string) (models.TokenPair, error)
	RefreshTokens(ctx utils.MyContext, accessToken, refreshToken, userIP string) (models.TokenPair, error)
	generateTokens(ctx utils.MyContext, userID, userIP string) (models.TokenPair, error)
}

type ImplAuthorizationService struct {
	repo            database.AuthorizationRepository
	jwtSecret       []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewAuthorizationService(repo database.AuthorizationRepository, jwtSecret []byte, accessTokenTTL, refreshTokenTTL time.Duration) *ImplAuthorizationService {
	return &ImplAuthorizationService{
		repo:            repo,
		jwtSecret:       jwtSecret,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (s *ImplAuthorizationService) CreateUser(ctx utils.MyContext, user models.User) (string, error) {
	user.Password, _ = generatePasswordHash(user.Password)
	return s.repo.CreateUser(ctx, mappers.MapToDBUser(user))
}

func (s *ImplAuthorizationService) Login(ctx utils.MyContext, input models.LoginInput, userIP string) (models.TokenPair, error) {
	userID, passwordHash, err := s.repo.GetUser(ctx, input.Email)
	if err != nil {
		return models.TokenPair{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(input.Password))
	if err != nil {
		return models.TokenPair{}, errors.New("invalid password")
	}

	return s.generateTokens(ctx, userID, userIP)
}

func (s *ImplAuthorizationService) RefreshTokens(ctx utils.MyContext, accessToken, refreshToken, userIP string) (models.TokenPair, error) {
	claims, err := parseAccessToken(accessToken, s.jwtSecret)
	if err != nil {
		return models.TokenPair{}, errors.New("invalid access token")
	}

	refreshTokenID := claims.RefreshTokenID
	storedIP := claims.UserIP
	storedUserID := claims.UserID

	if storedIP != userIP {
		go utils.SendWarningEmail(ctx, storedUserID, storedIP, userIP)
	}

	storedTokenHash, err := s.repo.GetRefreshToken(ctx, refreshTokenID)
	if err != nil {
		return models.TokenPair{}, err
	}

	if err = compareHashes(storedTokenHash, refreshToken); err != nil {
		return models.TokenPair{}, err
	}

	return s.generateTokens(ctx, storedUserID, userIP)
}

func (s *ImplAuthorizationService) generateTokens(ctx utils.MyContext, userID, userIP string) (models.TokenPair, error) {
	refreshTokenStr, err := generateRefreshToken()
	if err != nil {
		return models.TokenPair{}, err
	}

	refreshTokenHashByte, err := bcrypt.GenerateFromPassword([]byte(refreshTokenStr), bcrypt.DefaultCost)
	if err != nil {
		return models.TokenPair{}, errors.New("failed to hash refresh token")
	}
	refreshTokenHash := string(refreshTokenHashByte)

	refreshTokenID := uuid.New().String()

	err = s.repo.SaveRefreshToken(
		ctx,
		dbmodels.RefreshToken{
			ID:        refreshTokenID,
			UserID:    userID,
			TokenHash: refreshTokenHash,
			ExpiresAt: time.Now().Add(s.refreshTokenTTL),
		},
	)
	if err != nil {
		return models.TokenPair{}, err
	}

	accessToken, err := generateAccessToken(refreshTokenID, userID, userIP, s.jwtSecret, s.accessTokenTTL)
	if err != nil {
		return models.TokenPair{}, errors.New("failed to generate access token")
	}

	return models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
	}, nil
}

func generateAccessToken(refreshTokenID, userID, userIP string, jwtSecret []byte, tokenTTL time.Duration) (string, error) {
	claims := models.Claims{
		RefreshTokenID: refreshTokenID,
		UserID:         userID,
		UserIP:         userIP,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	return token.SignedString(jwtSecret)
}

func generateRefreshToken() (string, error) {
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", errors.New("failed to generate random bytes for refresh token")
	}

	refreshToken := base64.StdEncoding.EncodeToString(randomBytes)

	return refreshToken, nil
}

func parseAccessToken(accessToken string, jwtSecret []byte) (*models.Claims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("error parsing token")
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, errors.New("error parsing token")
	}

	claims, ok := token.Claims.(*models.Claims)
	if !ok || !token.Valid {
		return nil, errors.New("error parsing token")
	}

	return claims, nil
}

func compareHashes(hash string, token string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(token))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return errors.New("hash and token do not match")
		}

		return err
	}

	return nil
}

func generatePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}
