package database

import (
    "database/sql"
    "fmt"
    "math"
    "time"

    "algebrawl/internal/models"
    "github.com/google/uuid"
    _ "github.com/lib/pq"
)

type Repository struct {
    db *sql.DB
}

func NewRepository(connStr string) (*Repository, error) {
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }
    
    if err := db.Ping(); err != nil {
        return nil, err
    }
    
    return &Repository{db: db}, nil
}

func (r *Repository) Close() {
    r.db.Close()
}

func (r *Repository) CreateUser(firstName, lastName, login string) (string, error) {
    id := uuid.New().String()
    _, err := r.db.Exec(
        "INSERT INTO users (id, first_name, last_name, login, created_at) VALUES ($1, $2, $3, $4, $5)",
        id, firstName, lastName, login, time.Now(),
    )
    if err != nil {
        return "", err
    }
    return id, nil
}

func (r *Repository) UserExists(login string) (bool, error) {
    var count int
    err := r.db.QueryRow("SELECT COUNT(*) FROM users WHERE login = $1", login).Scan(&count)
    if err != nil {
        return false, err
    }
    return count > 0, nil
}

func (r *Repository) CreateEquationSet(userID string) (string, error) {
    id := uuid.New().String()
    _, err := r.db.Exec(
        "INSERT INTO equation_sets (id, user_id, created_at) VALUES ($1, $2, $3)",
        id, userID, time.Now(),
    )
    return id, err
}

func (r *Repository) SaveEquation(equation *models.Equation) error {
    _, err := r.db.Exec(
        "INSERT INTO equations (id, set_id, equation_text, root1, root2) VALUES ($1, $2, $3, $4, $5)",
        equation.ID, equation.SetID, equation.EquationText, equation.Root1, equation.Root2,
    )
    return err
}

func (r *Repository) GetLastEquationSet(userID string) (string, error) {
    var setID string
    err := r.db.QueryRow(
        "SELECT id FROM equation_sets WHERE user_id = $1 ORDER BY created_at DESC LIMIT 1",
        userID,
    ).Scan(&setID)
    if err != nil {
        return "", err
    }
    return setID, nil
}

func (r *Repository) GetEquationsBySetID(setID string) ([]models.Equation, error) {
    rows, err := r.db.Query(
        "SELECT id, equation_text, root1, root2 FROM equations WHERE set_id = $1 ORDER BY created_at",
        setID,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var equations []models.Equation
    for rows.Next() {
        var eq models.Equation
        if err := rows.Scan(&eq.ID, &eq.EquationText, &eq.Root1, &eq.Root2); err != nil {
            return nil, err
        }
        equations = append(equations, eq)
    }

    return equations, nil
}

func (r *Repository) SaveAnswers(setID string, answers []models.Answer) error {
    tx, err := r.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    for _, answer := range answers {
        var root1, root2 float64
        err := tx.QueryRow(
            "SELECT root1, root2 FROM equations WHERE set_id = $1 AND id = $2",
            setID, answer.EquationID,
        ).Scan(&root1, &root2)
        if err != nil {
            return err
        }

        correct := isAnswerCorrect(root1, root2, answer.Root1, answer.Root2)
        
        _, err = tx.Exec(
            "UPDATE equations SET user_answer1 = $1, user_answer2 = $2, solved_correctly = $3 WHERE set_id = $4 AND id = $5",
            answer.Root1, answer.Root2, correct, setID, answer.EquationID,
        )
        if err != nil {
            return err
        }
    }

    return tx.Commit()
}

func isAnswerCorrect(root1, root2, userAnswer1, userAnswer2 float64) bool {
    // Проверяем оба варианта порядка корней с учетом погрешности
    tolerance := 0.01
    return (abs(root1-userAnswer1) < tolerance && abs(root2-userAnswer2) < tolerance) ||
           (abs(root1-userAnswer2) < tolerance && abs(root2-userAnswer1) < tolerance)
}

func abs(x float64) float64 {
    return math.Abs(x)
}

func (r *Repository) GetEquationResults(setID string) ([]models.EquationResult, error) {
    rows, err := r.db.Query(`
        SELECT equation_text, user_answer1, user_answer2, root1, root2, solved_correctly 
        FROM equations WHERE set_id = $1 ORDER BY created_at
    `, setID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var results []models.EquationResult
    for i := 1; rows.Next(); i++ {
        var eqText string
        var user1, user2, root1, root2 sql.NullFloat64
        var solved sql.NullBool
        
        if err := rows.Scan(&eqText, &user1, &user2, &root1, &root2, &solved); err != nil {
            return nil, err
        }

        userAnswer := formatAnswer(user1, user2)
        correctAnswer := formatAnswer(root1, root2)

        results = append(results, models.EquationResult{
            ID:           i,
            Equation:     eqText,
            UserAnswer:   userAnswer,
            CorrectAnswer: correctAnswer,
            Correct:      solved.Bool,
        })
    }

    return results, nil
}

func formatAnswer(a1, a2 sql.NullFloat64) string {
    if !a1.Valid || !a2.Valid {
        return "Нет ответа"
    }
    return fmt.Sprintf("%.2f, %.2f", a1.Float64, a2.Float64)
}
