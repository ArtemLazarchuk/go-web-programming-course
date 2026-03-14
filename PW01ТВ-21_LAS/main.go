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

func calculateTask1(hp, c, sp, np, op, ap, w string) string {
	hVal := parseFloat(hp)
	cVal := parseFloat(c)
	sVal := parseFloat(sp)
	nVal := parseFloat(np)
	oVal := parseFloat(op)
	aVal := parseFloat(ap)
	wVal := parseFloat(w)

	kpc := 100 / (100 - wVal)
	krg := 100 / (100 - wVal - aVal)

	qpn := (339*cVal) + (1030*hVal) - (108.8*(oVal-sVal)) - (25*wVal)
	qcn := ((qpn/1000) + 0.025*wVal) * (100 / (100 - wVal))
	qgn := ((qpn/1000) + 0.025*wVal) * (100 / (100 - wVal - aVal))

	return fmt.Sprintf(`1.1 Коефіцієнт переходу від робочої до сухої маси 
Кпс = %.2f

1.2 Коефіцієнт переходу від робочої до горючої маси 
Кпг = %.2f

1.3 Склад сухої маси палива
Hc = %.2f %%
Cc = %.2f %%
Sc = %.2f %%
Nc = %.2f %%
Oc = %.2f %%
Ac = %.2f %%

1.4 Склад горючої маси палива
Hг = %.2f %%
Cг = %.2f %%
Sг = %.2f %%
Nг = %.2f %%
Oг = %.2f %%

1.5 Нижча теплота згорання для робочої маси 
QpН = %.2f [кДж/кг] = %.2f [МДж/кг]

1.6 Нижча теплота згоряння для сухої маси
%.2f [МДж/кг]

1.7 Нижча теплота згоряння для горючої маси
%.2f [МДж/кг]`,
		kpc, krg,
		hVal*kpc, cVal*kpc, sVal*kpc, nVal*kpc, oVal*kpc, aVal*kpc,
		hVal*krg, cVal*krg, sVal*krg, nVal*krg, oVal*krg,
		qpn, qpn/1000,
		qcn, qgn)
}

func calculateTask2(hg, cg, sg, og, vg, wg, ag, Qdaf string) string {
	hVal := parseFloat(hg)
	cVal := parseFloat(cg)
	sVal := parseFloat(sg)
	oVal := parseFloat(og)
	vVal := parseFloat(vg)
	wVal := parseFloat(wg)
	aVal := parseFloat(ag)
	qVal := parseFloat(Qdaf)

	Qri := qVal*((100-wVal-aVal)/100) - (0.025 * wVal)

	return fmt.Sprintf(`2.1 Склад робочої маси мазуту становитеме
H = %.2f %%
C = %.2f %%
S = %.2f %%
O = %.2f %%
V = %.2f мг/кг
A = %.2f %%

2.2 Нижча теплота згорання мазуту
%.2f МДж/кг`,
		hVal*((100-wVal-aVal)/100),
		cVal*((100-wVal-aVal)/100),
		sVal*((100-wVal-aVal)/100),
		oVal*((100-wVal-aVal)/100),
		vVal*((100-wVal)/100),
		aVal*((100-wVal)/100),
		Qri)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	data := struct {
		Task1Result string
		Task2Result string
		Task1Values map[string]string
		Task2Values map[string]string
	}{
		Task1Values: map[string]string{
			"hp": "", "c": "", "sp": "", "np": "", "op": "", "ap": "", "w": "",
		},
		Task2Values: map[string]string{
			"hg": "", "cg": "", "sg": "", "og": "", "vg": "", "wg": "", "ag": "", "Qdaf": "",
		},
	}

	if r.Method == http.MethodPost {
		r.ParseForm()

		if r.FormValue("task") == "1" {
			hp := r.FormValue("hp")
			c := r.FormValue("c")
			sp := r.FormValue("sp")
			np := r.FormValue("np")
			op := r.FormValue("op")
			ap := r.FormValue("ap")
			w := r.FormValue("w")
			data.Task1Result = calculateTask1(hp, c, sp, np, op, ap, w)
			data.Task1Values = map[string]string{"hp": hp, "c": c, "sp": sp, "np": np, "op": op, "ap": ap, "w": w}
		} else if r.FormValue("task") == "2" {
			hg := r.FormValue("hg")
			cg := r.FormValue("cg")
			sg := r.FormValue("sg")
			og := r.FormValue("og")
			vg := r.FormValue("vg")
			wg := r.FormValue("wg")
			ag := r.FormValue("ag")
			Qdaf := r.FormValue("Qdaf")
			data.Task2Result = calculateTask2(hg, cg, sg, og, vg, wg, ag, Qdaf)
			data.Task2Values = map[string]string{"hg": hg, "cg": cg, "sg": sg, "og": og, "vg": vg, "wg": wg, "ag": ag, "Qdaf": Qdaf}
		}
	}

	tmpl.Execute(w, data)
}

func main() {
	http.HandleFunc("/", indexHandler)
	fmt.Println("Сервер запущено: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
