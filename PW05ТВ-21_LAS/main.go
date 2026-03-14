package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func parseFloat(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func parseInt(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}

func formatValue(value float64, precision int, scientific bool) string {
	if scientific && value < 0.001 && value > 0 {
		return fmt.Sprintf("%.*e", precision, value)
	}
	return fmt.Sprintf("%.*f", precision, value)
}

func calculateReliability(connections int, accidentPrice, planedPrice float64) string {
	hoursPerYear := 8760.0

	omegaOc := 0.01 + 0.07 + 0.015 + 0.02 + 0.03*float64(connections)
	tV_oc := (0.01*30 + 0.07*10 + 0.015*100 + 0.02*15 + (0.03*float64(connections))*2) / omegaOc

	kA_oc := (omegaOc * tV_oc) / hoursPerYear
	kP_oc := 1.2 * (43.0 / hoursPerYear)

	omegaDk := 2 * omegaOc * (kA_oc + kP_oc)
	omegaDc := omegaDk + 0.02

	omega := 0.01
	tV := 45e-3
	Pm := 5.12e3
	Tm := 6451.0
	kP := 4e-3

	mathWnedA := omega * tV * Pm * Tm
	mathWnedP := kP * Pm * Tm
	mathLosses := accidentPrice*mathWnedA + planedPrice*mathWnedP

	return fmt.Sprintf(`═══════════════════════════════════
РЕЗУЛЬТАТИ РОЗРАХУНКУ НАДІЙНОСТІ
═══════════════════════════════════

1. Одноколова система
   Частота відмов ω_oc = %s рік⁻¹
   Середня тривалість відновлення t_в.ос = %s год
   Коефіцієнт аварійного простою k_a.oc = %s
   Коефіцієнт планового простою k_п.ос = %s

2. Двоколова система
   Частота відмов одночасно двох кіл ω_дк = %s рік⁻¹
   Частота відмов з урахуванням секційного вимикача ω_дс = %s рік⁻¹

   Висновок: Надійність двоколової системи є значно вищою ніж одноколової

3. Збитки від перерв електропостачання
   Математичне сподівання аварійного недовідпущення M(W_нед.а) = %.0f кВт·год
   Математичне сподівання планового недовідпущення M(W_нед.п) = %.0f кВт·год
   Математичне сподівання збитків M(З_пер) = %.2f грн

═══════════════════════════════════`,
		formatValue(omegaOc, 3, false),
		formatValue(tV_oc, 1, false),
		formatValue(kA_oc, 5, true),
		formatValue(kP_oc, 5, true),
		formatValue(omegaDk, 5, true),
		formatValue(omegaDc, 4, false),
		mathWnedA, mathWnedP, mathLosses)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	data := struct {
		Result string
		Error  string
		Values map[string]string
	}{
		Values: map[string]string{
			"connections":   "6",
			"accidentPrice": "23.6",
			"planedPrice":   "17.6",
		},
	}

	if r.Method == http.MethodPost {
		r.ParseForm()

		action := r.FormValue("action")

		if action == "example" {
			data.Values = map[string]string{
				"connections":   "6",
				"accidentPrice": "23.6",
				"planedPrice":   "17.6",
			}
		} else {
			connectionsStr := r.FormValue("connections")
			accidentPriceStr := r.FormValue("accidentPrice")
			planedPriceStr := r.FormValue("planedPrice")

			connections := parseInt(connectionsStr)
			accidentPrice := parseFloat(accidentPriceStr)
			planedPrice := parseFloat(planedPriceStr)

			if connections <= 0 {
				data.Error = "Введіть коректну кількість приєднань"
			} else if accidentPrice <= 0 {
				data.Error = "Введіть коректну ціну для аварійних перерв"
			} else if planedPrice <= 0 {
				data.Error = "Введіть коректну ціну для планових перерв"
			} else {
				data.Result = calculateReliability(connections, accidentPrice, planedPrice)
			}

			data.Values = map[string]string{
				"connections":   connectionsStr,
				"accidentPrice": accidentPriceStr,
				"planedPrice":   planedPriceStr,
			}
		}
	}

	tmpl.Execute(w, data)
}

func main() {
	http.HandleFunc("/", indexHandler)
	fmt.Println("Сервер запущено: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
