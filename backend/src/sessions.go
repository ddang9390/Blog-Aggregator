package main

import (
	"blog-aggregator/backend/internal/database"
	"crypto/rand"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

type Session struct {
	sessionID string
	userId    string
	createdAt time.Time
	expiresAt time.Time
}

func createSession(cfg *apiConfig, userID string, r *http.Request) (*Session, error) {
	sessionID := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, sessionID)
	if err != nil {
		return nil, err
	}

	session := Session{
		sessionID: string(sessionID),
		userId:    userID,
		createdAt: time.Now().UTC(),
		expiresAt: time.Now().UTC().Add(time.Hour),
	}

	sessionParams := database.CreateSessionParams{
		SessionID: session.sessionID,
		UserID:    session.userId,
		CreatedAt: session.createdAt,
		ExpiresAt: session.expiresAt,
	}

	ctx := r.Context()
	_, err = cfg.DB.CreateSession(ctx, sessionParams)
	if err != nil {
		//reattempt to create session in case of duplicate session id
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			for i := 0; i < 5; i++ {
				sessionID = make([]byte, 32)
				_, err := io.ReadFull(rand.Reader, sessionID)
				if err != nil {
					return nil, err
				}
				sessionParams.SessionID = string(sessionID)
				_, err = cfg.DB.CreateSession(ctx, sessionParams)
				if err == nil {
					break
				}
				//for reducing load on system
				time.Sleep(time.Duration(1<<uint(i)) * time.Second)
			}
		}
		return nil, err
	}

	return &session, nil
}
