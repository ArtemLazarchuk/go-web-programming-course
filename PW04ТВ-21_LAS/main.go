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

func calculateTask1(smStr, unomStr string) string {
	Sm := parseFloat(smStr)
	U_nom := parseFloat(unomStr)

	if Sm == 0 || U_nom == 0 {
		return "Будь ласка, введіть коректні значення."
	}

	t_f := 2.5
	Jek := 1.4
	Ik := 2500.0
	Ct := 92.0

	Im := (Sm / 2) / (math.Sqrt(3) * U_nom)
	Im_pa := 2 * Im

	Sek := Im / Jek
	S_min := (Ik * math.Sqrt(t_f)) / Ct

	S_selected := math.Max(Sek, S_min)

	return fmt.Sprintf(`═══════════════════════════════════
РЕЗУЛЬТАТИ РОЗРАХУНКУ ЗАВДАННЯ 1
═══════════════════════════════════

1. Розрахунковий струм:
   Іm = (Sm/2) / (√3 × U_nom) = %.2f А
   Іm*па = 2 × Іm = %.2f А

2. Економічний переріз:
   Sек = Іm / Jек = %.2f мм²

3. Термічна стійкість:
   S_min = (Iк × √t_f) / Ct = %.2f мм²

4. Вибір перерізу кабеля:
   Sек = %.2f мм²
   S_min = %.2f мм²
   Обрано: S = %.2f мм²
═══════════════════════════════════`,
		Im, Im_pa, Sek, S_min, Sek, S_min, S_selected)
}

func calculateTask2(ucnStr, skStr, ukStr, snomTStr string) string {
	UcnVal := parseFloat(ucnStr)
	SkVal := parseFloat(skStr)
	UkVal := parseFloat(ukStr)
	Snom_tVal := parseFloat(snomTStr)

	if SkVal == 0 || Snom_tVal == 0 {
		return "Помилка: Sk та Snom_t не можуть бути нульовими!"
	}

	Xc := (UcnVal * UcnVal) / SkVal
	Xt := (UkVal * (UcnVal * UcnVal)) / (100 * Snom_tVal)
	X_all := Xc + Xt
	Ipo := UcnVal / (math.Sqrt(3) * X_all)

	return fmt.Sprintf(`═══════════════════════════════════
РЕЗУЛЬТАТИ РОЗРАХУНКУ ЗАВДАННЯ 2
═══════════════════════════════════

1. Опори елементів заступної схеми:
   Xc = Ucn² / Sk = %.3f Ом
   Xt = (Uk%% / 100) × (Ucn² / Snom_t) = %.3f Ом

2. Сумарний опір для точки К1:
   XΣ = Xc + Xt = %.3f Ом

3. Початкове діюче значення струму трифазного КЗ:
   Iп₀ = Ucn / (√3 × XΣ) = %.3f кА
═══════════════════════════════════`,
		Xc, Xt, X_all, Ipo)
}

func calculateTask3(ukMaxStr, uvNStr, snomT3Str, unNStr, rcMinStr string) string {
	Uk_max := parseFloat(ukMaxStr)
	Uv_n := parseFloat(uvNStr)
	Snom_t := parseFloat(snomT3Str)
	Un_n := parseFloat(unNStr)
	Rc_min := parseFloat(rcMinStr)

	if Snom_t == 0 || Uv_n == 0 || Un_n == 0 {
		return "Помилка: Snom_t, Uv_n та Un_n не можуть бути нульовими!"
	}

	l := 12.37
	Rc_n := 10.65
	Ro := 0.64
	Xc_n := 24.02
	Xo := 0.363
	Xc_min := 65.68

	Xt := (Uk_max * Uv_n * Uv_n) / (100.0 * Snom_t)

	Rsh := Rc_n
	Xsh := Xc_n + Xt
	Zsh := math.Sqrt(Rsh*Rsh + Xsh*Xsh)

	Rsh_min := Rc_min
	Xsh_min := Xc_min + Xt
	Zsh_min := math.Sqrt(Rsh_min*Rsh_min + Xsh_min*Xsh_min)

	Ish_3 := (Uv_n * 1000) / (math.Sqrt(3) * Zsh)
	Ish_2 := Ish_3 * (math.Sqrt(3) / 2)
	Ish_3_min := (Uv_n * 1000) / (math.Sqrt(3) * Zsh_min)
	Ish_2_min := Ish_3_min * (math.Sqrt(3) / 2)

	Kpr := (Un_n * Un_n) / (Uv_n * Uv_n)

	Rsh_n := Rsh * Kpr
	Xsh_n := Xsh * Kpr
	Zsh_n := math.Sqrt(Rsh_n*Rsh_n + Xsh_n*Xsh_n)

	Rsh_n_min := Rsh_min * Kpr
	Xsh_n_min := Xsh_min * Kpr
	Zsh_n_min := math.Sqrt(Rsh_n_min*Rsh_n_min + Xsh_n_min*Xsh_n_min)

	Ish_n_3 := (Un_n * 1000) / (math.Sqrt(3) * Zsh_n)
	Ish_n_2 := Ish_n_3 * (math.Sqrt(3) / 2)
	Ish_n_3_min := (Un_n * 1000) / (math.Sqrt(3) * Zsh_n_min)
	Ish_n_2_min := Ish_n_3_min * (math.Sqrt(3) / 2)

	Rl := l * Ro
	Xl := l * Xo

	Rall_n := Rl + Rsh_n
	Xall_n := Xl + Xsh_n
	Zall_n := math.Sqrt(Rall_n*Rall_n + Xall_n*Xall_n)

	Rall_n_min := Rl + Rsh_n_min
	Xall_n_min := Xl + Xsh_n_min
	Zall_n_min := math.Sqrt(Rall_n_min*Rall_n_min + Xall_n_min*Xall_n_min)

	Il_n_3 := (Un_n * 1000) / (math.Sqrt(3) * Zall_n)
	Il_n_2 := Il_n_3 * (math.Sqrt(3) / 2)
	Il_n_3_min := (Un_n * 1000) / (math.Sqrt(3) * Zall_n_min)
	Il_n_2_min := Il_n_3_min * (math.Sqrt(3) / 2)

	return fmt.Sprintf(`═══════════════════════════════════
РЕЗУЛЬТАТИ РОЗРАХУНКУ ЗАВДАННЯ 3
═══════════════════════════════════

1. Реактивний опір трансформатора:
   Xt = (Uk_max × Uv_n²) / (100 × Snom_t) = %.3f Ом

2. Опори на шинах 10 кВ (приведені до 110 кВ):
   Rsh = Rc_n = %.3f Ом
   Xsh = Xc_n + Xt = %.3f Ом
   Zsh = √(Rsh² + Xsh²) = %.3f Ом
   Rsh_min = Rc_min = %.3f Ом
   Xsh_min = Xc_min + Xt = %.3f Ом
   Zsh_min = √(Rsh_min² + Xsh_min²) = %.3f Ом

3. Струми КЗ на шинах 10 кВ (приведені до 110 кВ):
   Ish_3 = (Uv_n × 10³) / (√3 × Zsh) = %.3f А
   Ish_2 = Ish_3 × (√3 / 2) = %.3f А
   Ish_3_min = (Uv_n × 10³) / (√3 × Zsh_min) = %.3f А
   Ish_2_min = Ish_3_min × (√3 / 2) = %.3f А

4. Коефіцієнт приведення:
   Kpr = Un_n² / Uv_n² = %.3f

5. Опори на шинах 10 кВ:
   Rsh_n = Rsh × Kpr = %.3f Ом
   Xsh_n = Xsh × Kpr = %.3f Ом
   Zsh_n = √(Rsh_n² + Xsh_n²) = %.3f Ом
   Rsh_n_min = Rsh_min × Kpr = %.3f Ом
   Xsh_n_min = Xsh_min × Kpr = %.3f Ом
   Zsh_n_min = √(Rsh_n_min² + Xsh_n_min²) = %.3f Ом

6. Дійсні струми КЗ на шинах 10 кВ:
   Ish_n_3 = (Un_n × 10³) / (√3 × Zsh_n) = %.3f А
   Ish_n_2 = Ish_n_3 × (√3 / 2) = %.3f А
   Ish_n_3_min = (Un_n × 10³) / (√3 × Zsh_n_min) = %.3f А
   Ish_n_2_min = Ish_n_3_min × (√3 / 2) = %.3f А

7. Опори лінії:
   l = %.2f км
   Rl = l × Ro = %.3f Ом
   Xl = l × Xo = %.3f Ом

8. Опори в точці 10:
   Rall_n = Rl + Rsh_n = %.3f Ом
   Xall_n = Xl + Xsh_n = %.3f Ом
   Zall_n = √(Rall_n² + Xall_n²) = %.3f Ом
   Rall_n_min = Rl + Rsh_n_min = %.3f Ом
   Xall_n_min = Xl + Xsh_n_min = %.3f Ом
   Zall_n_min = √(Rall_n_min² + Xall_n_min²) = %.3f Ом

9. Струми КЗ в точці 10:
   Il_n_3 = (Un_n × 10³) / (√3 × Zall_n) = %.3f А
   Il_n_2 = Il_n_3 × (√3 / 2) = %.3f А
   Il_n_3_min = (Un_n × 10³) / (√3 × Zall_n_min) = %.3f А
   Il_n_2_min = Il_n_3_min × (√3 / 2) = %.3f А
═══════════════════════════════════`,
		Xt,
		Rsh, Xsh, Zsh, Rsh_min, Xsh_min, Zsh_min,
		Ish_3, Ish_2, Ish_3_min, Ish_2_min,
		Kpr,
		Rsh_n, Xsh_n, Zsh_n, Rsh_n_min, Xsh_n_min, Zsh_n_min,
		Ish_n_3, Ish_n_2, Ish_n_3_min, Ish_n_2_min,
		l, Rl, Xl,
		Rall_n, Xall_n, Zall_n, Rall_n_min, Xall_n_min, Zall_n_min,
		Il_n_3, Il_n_2, Il_n_3_min, Il_n_2_min)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	data := struct {
		Task1Result string
		Task2Result string
		Task3Result string
		Task1Values map[string]string
		Task2Values map[string]string
		Task3Values map[string]string
	}{
		Task1Values: map[string]string{"sm": "", "unom": ""},
		Task2Values: map[string]string{"ucn": "", "sk": "", "uk": "", "snomT": ""},
		Task3Values: map[string]string{"ukMax": "", "uvN": "", "snomT3": "", "unN": "", "rcMin": ""},
	}

	if r.Method == http.MethodPost {
		r.ParseForm()

		if r.FormValue("task") == "1" {
			sm := r.FormValue("sm")
			unom := r.FormValue("unom")
			data.Task1Result = calculateTask1(sm, unom)
			data.Task1Values = map[string]string{"sm": sm, "unom": unom}
		} else if r.FormValue("task") == "2" {
			ucn := r.FormValue("ucn")
			sk := r.FormValue("sk")
			uk := r.FormValue("uk")
			snomT := r.FormValue("snomT")
			data.Task2Result = calculateTask2(ucn, sk, uk, snomT)
			data.Task2Values = map[string]string{"ucn": ucn, "sk": sk, "uk": uk, "snomT": snomT}
		} else if r.FormValue("task") == "3" {
			ukMax := r.FormValue("ukMax")
			uvN := r.FormValue("uvN")
			snomT3 := r.FormValue("snomT3")
			unN := r.FormValue("unN")
			rcMin := r.FormValue("rcMin")
			data.Task3Result = calculateTask3(ukMax, uvN, snomT3, unN, rcMin)
			data.Task3Values = map[string]string{"ukMax": ukMax, "uvN": uvN, "snomT3": snomT3, "unN": unN, "rcMin": rcMin}
		}
	}

	tmpl.Execute(w, data)
}

func main() {
	http.HandleFunc("/", indexHandler)
	fmt.Println("Сервер запущено: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
