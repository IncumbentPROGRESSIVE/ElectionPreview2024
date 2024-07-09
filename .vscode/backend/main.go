package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const (
    dbDriver  = "mysql"
    dbUser    = "colinleary"
    dbPass    = "00005828"
    dbName    = "epdata"
    jwtSecret = "JSGWcKC9lXfuCvyDsaPAhNsvRs2hPiQHnh7La6PBEfw="
)

var db *sql.DB

func initDB() {
    var err error
    dsn := fmt.Sprintf("%s:%s@/%s", dbUser, dbPass, dbName)
    db, err = sql.Open(dbDriver, dsn)
    if err != nil {
        log.Fatal("Error opening database connection: ", err)
    }

    err = db.Ping()
    if err != nil {
        log.Fatal("Error pinging database: ", err)
    }

    log.Println("Database connection established")
}

type User struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

func generateToken(email string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "email": email,
        "exp":   time.Now().Add(time.Hour * 24).Unix(),
    })
    return token.SignedString([]byte(jwtSecret))
}

func registerUserHandler(w http.ResponseWriter, r *http.Request) {
    var user User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        log.Println("Error decoding request payload: ", err)
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    _, err = db.Exec("INSERT INTO users (email, password) VALUES (?, ?)", user.Email, user.Password)
    if err != nil {
        log.Println("Error inserting user into database: ", err)
        http.Error(w, fmt.Sprintf("Failed to register user: %v", err), http.StatusInternalServerError)
        return
    }

    tokenString, err := generateToken(user.Email)
    if err != nil {
        log.Println("Error generating token: ", err)
        http.Error(w, "Failed to generate token", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{
        "token": tokenString,
    })
}

func loginUserHandler(w http.ResponseWriter, r *http.Request) {
    var user User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        log.Println("Error decoding request payload: ", err)
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    var dbUser User
    err = db.QueryRow("SELECT email, password FROM users WHERE email = ?", user.Email).Scan(&dbUser.Email, &dbUser.Password)
    if err != nil {
        if err == sql.ErrNoRows {
            log.Println("User not found: ", err)
            http.Error(w, "User not found", http.StatusUnauthorized)
            return
        }
        log.Println("Error querying user from database: ", err)
        http.Error(w, "Failed to login user", http.StatusInternalServerError)
        return
    }

    if user.Password != dbUser.Password {
        log.Println("Invalid password")
        http.Error(w, "Invalid password", http.StatusUnauthorized)
        return
    }

    tokenString, err := generateToken(user.Email)
    if err != nil {
        log.Println("Error generating token: ", err)
        http.Error(w, "Failed to generate token", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "token": tokenString,
    })
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
    tokenString := r.Header.Get("Authorization")
    if tokenString == "" {
        log.Println("No token provided")
        http.Error(w, "No token provided", http.StatusUnauthorized)
        return
    }

    tokenString = tokenString[len("Bearer "):]
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(jwtSecret), nil
    })

    if err != nil {
        log.Println("Error parsing token: ", err)
        http.Error(w, "Invalid token", http.StatusUnauthorized)
        return
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        email := claims["email"].(string)
        json.NewEncoder(w).Encode(map[string]string{
            "email": email,
        })
    } else {
        log.Println("Invalid token")
        http.Error(w, "Invalid token", http.StatusUnauthorized)
    }
}

func main() {
    initDB()
    defer db.Close()

    r := mux.NewRouter()
    r.HandleFunc("/register", registerUserHandler).Methods("POST")
    r.HandleFunc("/login", loginUserHandler).Methods("POST")
    r.HandleFunc("/profile", profileHandler).Methods("GET")

    log.Println("Server listening on :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
