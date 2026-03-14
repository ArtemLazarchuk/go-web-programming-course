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

func calculateEmissions(bCoal, bOil, bGas, qCoal, qOil, qGas, kCoal, kOil, kGas string) string {
	bC := parseFloat(bCoal)
	bO := parseFloat(bOil)

	qC := parseFloat(qCoal)
	qO := parseFloat(qOil)

	kC := parseFloat(kCoal)
	kO := parseFloat(kOil)

	colaK := (1000000/qC)*0.8*(kC/(100-1.5))*(1-0.985)
	oilK := (1000000/qO)*1*(kO/(100-0))*(1-0.985)
	coalEmission := 0.000001 * colaK * qC * bC
	oilEmission := 0.000001 * oilK * qO * bO

	return fmt.Sprintf(`1.1 Показник емісії твердих частинок при спалюванні вугілля: 
%.2f г/ГДж

1.2 Валовий викид при спалюванні вугілля: 
%.2f т.

1.3 Показник емісії твердих частинок при спалюванні мазуту:
%.2f г/ГДж

1.4 Валовий викид при спалюванні мазуту:
%.2f т.

1.5 Показник емісії твердих частинок при спалюванні природного газу:
0 г/ГДж

1.6 Валовий викид при спалюванні природного газу:
0 т.`,
		colaK, coalEmission, oilK, oilEmission)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	data := struct {
		Result string
		Values map[string]string
	}{
		Values: map[string]string{
			"bCoal": "1096363", "bOil": "70945", "bGas": "84762",
			"qCoal": "20.47", "qOil": "39.48", "qGas": "33.08",
			"kCoal": "25.20", "kOil": "0.15", "kGas": "0.723",
		},
	}

	if r.Method == http.MethodPost {
		r.ParseForm()

		bCoal := r.FormValue("bCoal")
		bOil := r.FormValue("bOil")
		bGas := r.FormValue("bGas")
		qCoal := r.FormValue("qCoal")
		qOil := r.FormValue("qOil")
		qGas := r.FormValue("qGas")
		kCoal := r.FormValue("kCoal")
		kOil := r.FormValue("kOil")
		kGas := r.FormValue("kGas")

		data.Result = calculateEmissions(bCoal, bOil, bGas, qCoal, qOil, qGas, kCoal, kOil, kGas)
		data.Values = map[string]string{
			"bCoal": bCoal, "bOil": bOil, "bGas": bGas,
			"qCoal": qCoal, "qOil": qOil, "qGas": qGas,
			"kCoal": kCoal, "kOil": kOil, "kGas": kGas,
		}
	}

	tmpl.Execute(w, data)
}

func main() {
	http.HandleFunc("/", indexHandler)
	fmt.Println("Сервер запущено: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
