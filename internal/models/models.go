package models

import (
    "time"
)

type User struct {
    ID        string    `json:"id" db:"id"`
    FirstName string    `json:"first_name" db:"first_name"`
    LastName  string    `json:"last_name" db:"last_name"`
    Login     string    `json:"login" db:"login"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type EquationSet struct {
    ID        string    `json:"id" db:"id"`
    UserID    string    `json:"user_id" db:"user_id"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Equation struct {
    ID           string  `json:"id" db:"id"`
    SetID        string  `json:"set_id" db:"set_id"`
    EquationText string  `json:"equation_text" db:"equation_text"`
    Root1        float64 `json:"root1" db:"root1"`
    Root2        float64 `json:"root2" db:"root2"`
    UserAnswer1  *float64 `json:"user_answer1,omitempty" db:"user_answer1"`
    UserAnswer2  *float64 `json:"user_answer2,omitempty" db:"user_answer2"`
    SolvedCorrectly *bool `json:"solved_correctly,omitempty" db:"solved_correctly"`
}

type RegisterRequest struct {
    FirstName string `json:"first_name" binding:"required"`
    LastName  string `json:"last_name" binding:"required"`
    Login     string `json:"login" binding:"required"`
}

type RegisterResponse struct {
    UserID string `json:"user_id"`
}

type NewEquationRequest struct {
    Count int `json:"count" binding:"required,min=1,max=100"`
}

type EquationItem struct {
    ID       int    `json:"id"`
    Equation string `json:"equation"`
}

type NewEquationResponse struct {
    Count int            `json:"count"`
    SetID string         `json:"set_id"`
    List  []EquationItem `json:"list"`
}

type Answer struct {
    EquationID int     `json:"equation_id"`
    Root1      float64 `json:"root1"`
    Root2      float64 `json:"root2"`
}

type StatisticsRequest struct {
    Count   int      `json:"count" binding:"required"`
    SetID   string   `json:"set_id" binding:"required"`
    Answers []Answer `json:"answers" binding:"required"`
}

type EquationResult struct {
    ID           int     `json:"id"`
    Equation     string  `json:"equation"`
    UserAnswer   string  `json:"user_answer"`
    CorrectAnswer string `json:"correct_answer"`
    Correct      bool    `json:"correct"`
}

type StatisticsResponse struct {
    Count int             `json:"count"`
    Rate  float64         `json:"rate"`
    Result []EquationResult `json:"result"`
}
