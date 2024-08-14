package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
)

type Handler struct {
	secretKey []byte
	log       *log.Logger
}

func NewHandler(log *log.Logger, secretKey []byte) *Handler {
	return &Handler{
		secretKey: secretKey,
		log:       log,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		log.Info("Handling GET request")

		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Missing authorization header")
			log.Warn("Missing authorization header")
			return
		}
		tokenString = tokenString[len("Bearer "):]

		err := h.verifyToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, err.Error())
			log.Warn(err)
			return
		}

		log.Info("Token verified")
		w.WriteHeader(http.StatusOK)
	} else if r.Method == http.MethodPost {
		log.Info("Handling POST request")
		userName := r.Header.Get("X-User-Name")

		if userName == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing X-User-Name"))
			log.Warn("Missing X-User-Name")
			return
		}

		token, err := h.createToken(userName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			log.Error(err)
			return
		}

		log.Info("Token created")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(token))
	} else {
		log.Info("Handling unsupported request")

		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) createToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Minute * 5).Unix(),
		})

	tokenString, err := token.SignedString(h.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (h *Handler) verifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return h.secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
