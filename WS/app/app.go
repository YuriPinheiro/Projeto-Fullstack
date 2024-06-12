package app

import (
	"database/sql"
	"fmt"
	"log"

	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/VPompeu/AgendaAstrologica/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func (a *App) Initialize(user, password, dbname string) {
	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) getUsers(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	users, err := models.GetUsers(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, users)
}

func authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Obtém o token do cabeçalho Authorization
		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == "" {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		// O token está no formato "Bearer {token}", então o dividimos
		tokenParts := strings.Split(tokenHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			respondWithError(w, http.StatusUnauthorized, "Invalid Token")
			return
		}

		tokenStr := tokenParts[1]

		claims := &Claims{}

		jwtKey := os.Getenv("JWT_KEY")

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				respondWithError(w, http.StatusUnauthorized, "Unauthorized")
				return
			}
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		if !token.Valid {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		// Se o token for válido, chama o próximo manipulador
		next.ServeHTTP(w, r)
	}
}

func (a *App) getUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	u := models.User{ID: id}
	if err := u.GetUser(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "User not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, u)
}

func (a *App) createUser(w http.ResponseWriter, r *http.Request) {
	var u models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if u.Name == "" {
		respondWithError(w, http.StatusBadRequest, "Empty Name")
		return
	}

	if u.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Empty E-mail")
		return
	}

	if len(u.Password) < 6 {
		respondWithError(w, http.StatusBadRequest, "Short Password")
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	u.Password = string(passwordHash)

	if err := u.CreateUser(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, u)
}

func (a *App) updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var u models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()
	u.ID = id

	if err := u.UpdateUser(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, u)
}

func (a *App) deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	u := models.User{ID: id}
	if err := u.DeleteUser(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "200"})
}

func (a *App) login(w http.ResponseWriter, r *http.Request) {
	var u models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	authenticated, err := u.Login(a.DB)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
	} else if authenticated {
		expirationTime := time.Now().Add(24 * time.Hour)

		expirationUnix := jwt.NewNumericDate(expirationTime)
		claims := &Claims{
			Email: u.Email,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: expirationUnix,
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		jwtKey := os.Getenv("JWT_KEY")
		tokenString, err := token.SignedString([]byte(jwtKey))
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		respondWithJSON(w, http.StatusOK, tokenString)
	} else {
		respondWithError(w, http.StatusUnauthorized, "Falha no login")
	}

}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/users", authenticate(a.getUsers)).Methods("GET")
	a.Router.HandleFunc("/user/{id:[0-9]+}", a.getUser).Methods("GET")
	a.Router.HandleFunc("/user", a.createUser).Methods("POST")
	a.Router.HandleFunc("/user/{id:[0-9]+}", a.updateUser).Methods("PUT")
	a.Router.HandleFunc("/user/{id:[0-9]+}", a.deleteUser).Methods("DELETE")

	a.Router.HandleFunc("/login", a.login).Methods("POST")
}

func (a *App) Run(addr string) {
	log.Printf("Conectando com banco de dados!")
	defer a.DB.Close()
	log.Printf("Iniciando serviço em: %s ", addr)
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
