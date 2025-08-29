package api

import (
    "log"
    "net/http"

    "algebrawl/internal/database"
    "algebrawl/internal/generator"
    "algebrawl/internal/models"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

type Handler struct {
    repo *database.Repository
}

func NewHandler(repo *database.Repository) *Handler {
    return &Handler{repo: repo}
}

func (h *Handler) Register(c *gin.Context) {
    var req models.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Проверяем уникальность логина
    exists, err := h.repo.UserExists(req.Login)
    if err != nil {
        log.Printf("Error checking user existence: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }
    if exists {
        c.JSON(http.StatusConflict, gin.H{"error": "Login already exists"})
        return
    }

    userID, err := h.repo.CreateUser(req.FirstName, req.LastName, req.Login)
    if err != nil {
        log.Printf("Error creating user: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }

    c.JSON(http.StatusOK, models.RegisterResponse{UserID: userID})
}

func (h *Handler) NewEquationSet(c *gin.Context) {
    var req models.NewEquationRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    userID := c.Query("user_id")
    if userID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "user_id parameter is required"})
        return
    }

    // Создаем новый набор уравнений
    setID, err := h.repo.CreateEquationSet(userID)
    if err != nil {
        log.Printf("Error creating equation set: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }

    // Генерируем уравнения
    generatedEqs := generator.GenerateEquations(req.Count)
    var equations []models.EquationItem

    for _, genEq := range generatedEqs {
        equation := models.Equation{
            ID:           uuid.New().String(),
            SetID:        setID,
            EquationText: genEq.Equation,
            Root1:        genEq.Root1,
            Root2:        genEq.Root2,
        }

        if err := h.repo.SaveEquation(&equation); err != nil {
            log.Printf("Error saving equation: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
            return
        }

        equations = append(equations, models.EquationItem{
            ID:       genEq.ID,
            Equation: genEq.Equation,
        })
    }

    response := models.NewEquationResponse{
        Count: req.Count,
        SetID: setID,
        List:  equations,
    }

    c.JSON(http.StatusOK, response)
}

func (h *Handler) GetEquationSet(c *gin.Context) {
    userID := c.Query("user_id")
    if userID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "user_id parameter is required"})
        return
    }

    setID, err := h.repo.GetLastEquationSet(userID)
    if err != nil {
        log.Printf("Error getting last equation set: %v", err)
        c.JSON(http.StatusNotFound, gin.H{"error": "No equation sets found"})
        return
    }

    equations, err := h.repo.GetEquationsBySetID(setID)
    if err != nil {
        log.Printf("Error getting equations: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }

    var equationItems []models.EquationItem
    for i, eq := range equations {
        equationItems = append(equationItems, models.EquationItem{
            ID:       i + 1,
            Equation: eq.EquationText,
        })
    }

    response := models.NewEquationResponse{
        Count: len(equationItems),
        SetID: setID,
        List:  equationItems,
    }

    c.JSON(http.StatusOK, response)
}

func (h *Handler) Statistics(c *gin.Context) {
    var req models.StatisticsRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Сохраняем ответы
    var answers []models.Answer
    for _, ans := range req.Answers {
        answers = append(answers, models.Answer{
            EquationID: ans.EquationID,
            Root1: ans.Root1,
            Root2: ans.Root2,
        })
    }

    if err := h.repo.SaveAnswers(req.SetID, answers); err != nil {
        log.Printf("Error saving answers: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }

    // Получаем результаты
    results, err := h.repo.GetEquationResults(req.SetID)
    if err != nil {
        log.Printf("Error getting results: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }

    // Вычисляем процент правильных ответов
    correctCount := 0
    for _, res := range results {
        if res.Correct {
            correctCount++
        }
    }
    rate := float64(correctCount) / float64(len(results)) * 100

    response := models.StatisticsResponse{
        Count:  req.Count,
        Rate:   rate,
        Result: results,
    }

    c.JSON(http.StatusOK, response)
}
