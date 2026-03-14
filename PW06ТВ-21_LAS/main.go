package main

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
)

type EPType struct {
	Name   string
	Eta    float64
	CosPhi float64
	U      float64
	N      int
}

type ParsedEPData struct {
	Type   EPType
	Pn     float64
	Kv     float64
	TgPhi  float64
}

var epTypes = []EPType{
	{"Шліфувальний верстат", 0.92, 0.9, 0.38, 4},
	{"Свердлильний верстат", 0.92, 0.9, 0.38, 2},
	{"Фугувальний верстат", 0.92, 0.9, 0.38, 4},
	{"Циркулярна пила", 0.92, 0.9, 0.38, 1},
	{"Прес", 0.92, 0.9, 0.38, 1},
	{"Полірувальний верстат", 0.92, 0.9, 0.38, 1},
	{"Фрезерний верстат", 0.92, 0.9, 0.38, 2},
	{"Вентилятор", 0.92, 0.9, 0.38, 1},
	{"Зварювальний трансформатор", 0.92, 0.9, 0.38, 2},
	{"Сушильна шафа", 0.92, 0.9, 0.38, 2},
}

var exampleValues = [][3]string{
	{"20", "0.15", "1.33"},
	{"14", "0.12", "1.0"},
	{"42", "0.15", "1.33"},
	{"36", "0.3", "1.52"},
	{"20", "0.5", "0.75"},
	{"40", "0.2", "1.0"},
	{"32", "0.2", "1.0"},
	{"20", "0.65", "0.75"},
	{"100", "0.2", "3.0"},
	{"120", "0.8", "0.0"},
}

func parseFloat(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func getKpFromTable(kv float64, ne float64) float64 {
	neRounded := int(math.Round(ne))
	if neRounded < 1 {
		neRounded = 1
	}

	switch {
	case kv <= 0.1:
		return 1.0
	case kv <= 0.2:
		switch {
		case neRounded <= 2:
			return 1.0
		case neRounded <= 5:
			return 1.1
		case neRounded <= 10:
			return 1.15
		default:
			return 1.25
		}
	case kv <= 0.4:
		switch {
		case neRounded <= 30:
			return 1.0
		default:
			return 0.7
		}
	default:
		switch {
		case neRounded <= 30:
			return 1.0
		default:
			return 0.7
		}
	}
}

func calculateGroup(epList []ParsedEPData, useKpForQ bool) (kv, ne, kp, pp, qp, sp, ip float64) {
	var sumNPn, sumNPnKv, sumNPnKvTgPhi, sumNPn2 float64

	for _, ep := range epList {
		n := float64(ep.Type.N)
		sumNPn += n * ep.Pn
		sumNPnKv += n * ep.Pn * ep.Kv
		sumNPnKvTgPhi += n * ep.Pn * ep.Kv * ep.TgPhi
		sumNPn2 += n * ep.Pn * ep.Pn
	}

	if sumNPn > 0 {
		kv = sumNPnKv / sumNPn
	}
	if sumNPn2 > 0 {
		ne = (sumNPn * sumNPn) / sumNPn2
	}

	kp = getKpFromTable(kv, ne)
	pp = kp * sumNPnKv

	if useKpForQ {
		qp = kp * sumNPnKvTgPhi
	} else {
		qp = sumNPnKvTgPhi
	}

	sp = math.Sqrt(pp*pp + qp*qp)

	u := 0.38
	if len(epList) > 0 {
		u = epList[0].Type.U
	}
	ip = pp / u

	return kv, ne, kp, pp, qp, sp, ip
}

func calculateLoads(parsed []ParsedEPData) string {
	shrEPs := parsed[:8]
	shrKv, shrNe, shrKp, shrPp, shrQp, shrSp, shrIp := calculateGroup(shrEPs, false)

	wsKv, wsNe, wsKp, wsPp, wsQp, wsSp, wsIp := calculateGroup(parsed, true)

	return fmt.Sprintf(`1.1. Груповий коефіцієнт використання для ШР1=ШР2=ШР3: %.4f;
1.2. Ефективна кількість ЕП для ШР1=ШР2=ШР3: %d;
1.3. Розрахунковий коефіцієнт активної потужності для ШР1=ШР2=ШР3: %.2f;
1.4. Розрахункове активне навантаження для ШР1=ШР2=ШР3: %.2f кВт;
1.5. Розрахункове реактивне навантаження для ШР1=ШР2=ШР3: %.3f квар.;
1.6. Повна потужність для ШР1=ШР2=ШР3: %.4f кВ*А;
1.7. Розрахунковий груповий струм для ШР1=ШР2=ШР3: %.2f А;
1.8. Коефіцієнти використання цеху в цілому: %.2f;
1.9. Ефективна кількість ЕП цеху в цілому: %d;
1.10. Розрахунковий коефіцієнт активної потужності цеху в цілому: %.1f;
1.11. Розрахункове активне навантаження на шинах 0,38 кВ ТП: %.1f кВт;
1.12. Розрахункове реактивне навантаження на шинах 0,38 кВ ТП: %.1f квар;
1.13. Повна потужність на шинах 0,38 кВ ТП: %.0f кВ*А;
1.14. Розрахунковий груповий струм на шинах 0,38 кВ ТП: %.3f А.`,
		shrKv, int(shrNe), shrKp, shrPp, shrQp, shrSp, shrIp,
		wsKv, int(wsNe), wsKp, wsPp, wsQp, wsSp, wsIp)
}

func add(a, b int) int {
	return a + b
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("index.html").Funcs(template.FuncMap{"add": add}).ParseFiles("templates/index.html"))

	epValues := make([][3]string, 10)
	for i := 0; i < 10; i++ {
		epValues[i] = [3]string{"", "", ""}
	}

	data := struct {
		Result string
		Error  string
		EPList []EPFormData
	}{
		EPList: make([]EPFormData, 10),
	}

	for i := 0; i < 10; i++ {
		data.EPList[i] = EPFormData{
			Type:   epTypes[i],
			Pn:     "",
			Kv:     "",
			TgPhi:  "",
		}
	}

	if r.Method == http.MethodPost {
		r.ParseForm()

		if r.FormValue("action") == "example" {
			for i := 0; i < 10; i++ {
				data.EPList[i] = EPFormData{
					Type:   epTypes[i],
					Pn:     exampleValues[i][0],
					Kv:     exampleValues[i][1],
					TgPhi:  exampleValues[i][2],
				}
			}
		} else {
			var parsed []ParsedEPData
			for i := 0; i < 10; i++ {
				pn := r.FormValue(fmt.Sprintf("pn_%d", i))
				kv := r.FormValue(fmt.Sprintf("kv_%d", i))
				tgPhi := r.FormValue(fmt.Sprintf("tgPhi_%d", i))

				data.EPList[i] = EPFormData{Type: epTypes[i], Pn: pn, Kv: kv, TgPhi: tgPhi}

				pnVal := parseFloat(pn)
				kvVal := parseFloat(kv)
				tgPhiVal := parseFloat(tgPhi)

				if pnVal <= 0 {
					data.Error = "Рн для " + epTypes[i].Name
				} else if kvVal < 0 || kvVal > 1 {
					data.Error = "Кв для " + epTypes[i].Name
				} else if tgPhiVal < 0 {
					data.Error = "tgφ для " + epTypes[i].Name
				} else {
					parsed = append(parsed, ParsedEPData{Type: epTypes[i], Pn: pnVal, Kv: kvVal, TgPhi: tgPhiVal})
				}
			}

			if data.Error == "" && len(parsed) == 10 {
				data.Result = calculateLoads(parsed)
			}
		}
	}

	tmpl.Execute(w, data)
}

type EPFormData struct {
	Type  EPType
	Pn    string
	Kv    string
	TgPhi string
}

func main() {
	http.HandleFunc("/", indexHandler)
	fmt.Println("Сервер запущено: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
