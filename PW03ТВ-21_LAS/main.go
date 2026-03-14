package main

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
)

func parseFloat(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

// calculateNormalProbability обчислює ймовірність того, що значення нормального розподілу
// знаходиться в діапазоні [lowerBound, upperBound]
func calculateNormalProbability(mean, stdDev, lowerBound, upperBound float64) float64 {
	z1 := (lowerBound - mean) / (stdDev * math.Sqrt(2))
	z2 := (upperBound - mean) / (stdDev * math.Sqrt(2))
	return 0.5 * (math.Erf(z2) - math.Erf(z1))
}

func calculateSolarProfit(pcStr, sigma1Str, sigma2Str, bStr string) string {
	pc := parseFloat(pcStr)
	sigma1 := parseFloat(sigma1Str)
	sigma2 := parseFloat(sigma2Str)
	b := parseFloat(bStr)

	if pc <= 0 || sigma1 <= 0 || sigma2 <= 0 || b <= 0 {
		return "Помилка: всі значення повинні бути більше нуля"
	}

	delta := pc * 0.05
	lowerBound := pc - delta
	upperBound := pc + delta

	// Розрахунок для σ1
	deltaW1 := calculateNormalProbability(pc, sigma1, lowerBound, upperBound)
	w1 := pc * 24 * deltaW1
	w2 := pc * 24 * (1 - deltaW1)
	profit1 := w1 * b * 1000
	penalty1 := w2 * b * 1000
	netResult1 := profit1 - penalty1

	// Розрахунок для σ2
	deltaW2 := calculateNormalProbability(pc, sigma2, lowerBound, upperBound)
	w3 := pc * 24 * deltaW2
	w4 := pc * 24 * (1 - deltaW2)
	profit2 := w3 * b * 1000
	penalty2 := w4 * b * 1000
	netResult2 := profit2 - penalty2

	var result string

	// Початкова система
	result += "═══════════════════════════════════\n"
	result += fmt.Sprintf("РОЗРАХУНКИ ДЛЯ ПОЧАТКОВОЇ СИСТЕМИ (σ₁ = %.2f МВт)\n", sigma1)
	result += "═══════════════════════════════════\n\n"
	result += fmt.Sprintf("Діапазон без штрафів: %.2f - %.2f МВт\n\n", lowerBound, upperBound)
	result += fmt.Sprintf("Частка енергії без небалансів δw₁ = %.1f%%\n\n", deltaW1*100)
	result += fmt.Sprintf("За %.1f%% електроенергії:\n", deltaW1*100)
	result += fmt.Sprintf("  W₁ = Pc × 24 × δw₁ = %.2f × 24 × %.2f = %.2f МВт-год\n", pc, deltaW1, w1)
	result += fmt.Sprintf("  Прибуток П₁ = W₁ × B = %.2f × %.2f × 1000 = %.2f тис. грн\n\n", w1, b, profit1)
	result += fmt.Sprintf("За %.1f%% енергії:\n", (1-deltaW1)*100)
	result += fmt.Sprintf("  W₂ = Pc × 24 × (1 - δw₁) = %.2f × 24 × %.2f = %.2f МВт-год\n", pc, 1-deltaW1, w2)
	result += fmt.Sprintf("  Штраф Ш₁ = W₂ × B = %.2f × %.2f × 1000 = %.2f тис. грн\n\n", w2, b, penalty1)
	result += "═══════════════════════════════════\n"
	result += "ВИСНОВОК ДЛЯ ПОЧАТКОВОЇ СИСТЕМИ:\n"
	if netResult1 >= 0 {
		result += fmt.Sprintf("Електростанція працює з прибутком: %.2f тис. грн\n", netResult1)
	} else {
		result += fmt.Sprintf("Електростанція є нерентабельною і працює в збиток: %.2f тис. грн\n", math.Abs(netResult1))
	}
	result += "═══════════════════════════════════\n\n\n"

	// Вдосконалена система
	result += "═══════════════════════════════════\n"
	result += fmt.Sprintf("РОЗРАХУНКИ ДЛЯ ВДОСКОНАЛЕНОЇ СИСТЕМИ (σ₂ = %.2f МВт)\n", sigma2)
	result += "═══════════════════════════════════\n\n"
	result += fmt.Sprintf("Діапазон без штрафів: %.2f - %.2f МВт\n\n", lowerBound, upperBound)
	result += fmt.Sprintf("Частка енергії без небалансів δw₂ = %.1f%%\n\n", deltaW2*100)
	result += fmt.Sprintf("За %.1f%% електроенергії:\n", deltaW2*100)
	result += fmt.Sprintf("  W₃ = Pc × 24 × δw₂ = %.2f × 24 × %.2f = %.2f МВт-год\n", pc, deltaW2, w3)
	result += fmt.Sprintf("  Прибуток П₂ = W₃ × B = %.2f × %.2f × 1000 = %.2f тис. грн\n\n", w3, b, profit2)
	result += fmt.Sprintf("За %.1f%% енергії:\n", (1-deltaW2)*100)
	result += fmt.Sprintf("  W₄ = Pc × 24 × (1 - δw₂) = %.2f × 24 × %.2f = %.2f МВт-год\n", pc, 1-deltaW2, w4)
	result += fmt.Sprintf("  Штраф Ш₂ = W₄ × B = %.2f × %.2f × 1000 = %.2f тис. грн\n\n", w4, b, penalty2)
	result += "═══════════════════════════════════\n"
	result += "ВИСНОВОК ДЛЯ ВДОСКОНАЛЕНОЇ СИСТЕМИ:\n"
	if netResult2 >= 0 {
		result += fmt.Sprintf("Електростанція працює з прибутком: %.2f тис. грн\n", netResult2)
	} else {
		result += fmt.Sprintf("Електростанція працює в збиток: %.2f тис. грн\n", math.Abs(netResult2))
	}
	result += "═══════════════════════════════════\n"

	return result
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	data := struct {
		Result string
		Values map[string]string
	}{
		Values: map[string]string{
			"pc": "5", "sigma1": "1", "sigma2": "0.25", "b": "7",
		},
	}

	if r.Method == http.MethodPost {
		r.ParseForm()

		pc := r.FormValue("pc")
		sigma1 := r.FormValue("sigma1")
		sigma2 := r.FormValue("sigma2")
		b := r.FormValue("b")

		data.Result = calculateSolarProfit(pc, sigma1, sigma2, b)
		data.Values = map[string]string{"pc": pc, "sigma1": sigma1, "sigma2": sigma2, "b": b}
	}

	tmpl.Execute(w, data)
}

func main() {
	http.HandleFunc("/", indexHandler)
	fmt.Println("Сервер запущено: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
