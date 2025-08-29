package generator

import (
    "fmt"
    "math/rand"
    "time"
)

func init() {
    rand.Seed(time.Now().UnixNano())
}

type GeneratedEquation struct {
    ID       int
    Equation string
    Root1    float64
    Root2    float64
}

func GenerateEquation(id int) GeneratedEquation {
    // Генерация корней в диапазоне -100 до 100
    root1 := rand.Float64()*200 - 100
    root2 := rand.Float64()*200 - 100
    
    // Округляем до 2 знаков после запятой
    root1 = float64(int(root1*100)) / 100
    root2 = float64(int(root2*100)) / 100
    
    // По теореме Виета: x² + bx + c = 0
    b := -(root1 + root2)
    c := root1 * root2
    
    // Форматируем уравнение
    equationText := formatEquation(b, c)
    
    return GeneratedEquation{
        ID:       id,
        Equation: equationText,
        Root1:    root1,
        Root2:    root2,
    }
}

func formatEquation(b, c float64) string {
    var equation string
    
    if b == 0 {
        equation = "x²"
    } else if b > 0 {
        equation = fmt.Sprintf("x² + %.2fx", b)
    } else {
        equation = fmt.Sprintf("x² - %.2fx", -b)
    }
    
    if c == 0 {
        equation += " = 0"
    } else if c > 0 {
        equation += fmt.Sprintf(" + %.2f = 0", c)
    } else {
        equation += fmt.Sprintf(" - %.2f = 0", -c)
    }
    
    return equation
}

func GenerateEquations(count int) []GeneratedEquation {
    equations := make([]GeneratedEquation, count)
    for i := 0; i < count; i++ {
        equations[i] = GenerateEquation(i + 1)
    }
    return equations
}
