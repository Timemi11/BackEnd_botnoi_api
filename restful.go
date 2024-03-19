package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Pokemon struct {
	Stats   []Stat  `json:"stats"`
	Name    string  `json:"name"`
	Sprites Sprites `json:"sprites"`
}

type Stat struct {
	BaseStat int `json:"base_stat"`
	Effort   int `json:"effort"`
	Stat     struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"stat"`
}
type Sprites struct {
	BackDefault      string  `json:"back_default"`
	BackFemale       *string `json:"back_female"`
	BackShiny        string  `json:"back_shiny"`
	BackShinyFemale  *string `json:"back_shiny_female"`
	FrontDefault     string  `json:"front_default"`
	FrontFemale      *string `json:"front_female"`
	FrontShiny       string  `json:"front_shiny"`
	FrontShinyFemale *string `json:"front_shiny_female"`
}

//1.รับข้อมูลจาก api ที่กำหนด
//2.สร้าง port และลงท้ายendpointด้วย id
//3.จากนั่นใช้ method post เพื่อสร้างข้อมูลที่เก็บไว้ลงไปใน id ในรูปแบบ json
//4.ลองใช้ postman ตรวจว่าข้อมูลตรงไหม

func getPokemonData(w http.ResponseWriter, r *http.Request) {
	// รับข้อมูลตัว api ตัวแรก
	resp1, err := http.Get("https://pokeapi.co/api/v2/pokemon/1/")
	if err != nil {
		log.Fatalf("Error making request to first endpoint: %v", err)
	}
	defer resp1.Body.Close()

	//แปลงข้อมูลjsonของ api มาเป็น ข้อมูลแบบ struct ที่กำนดไว้
	var pokemon Pokemon
	if err := json.NewDecoder(resp1.Body).Decode(&pokemon); err != nil {
		log.Fatalf("Error decoding JSON from first endpoint: %v", err)
	}

	// รับข้อมูลตัว api ที่สอง
	resp2, err := http.Get("https://pokeapi.co/api/v2/pokemon-form/1/")
	if err != nil {
		log.Fatalf("Error making request to second endpoint: %v", err)
	}
	defer resp2.Body.Close()

	// กำหนดรูปแบบข้อมูล
	var form struct {
		Name    string  `json:"name"`
		Sprites Sprites `json:"sprites"`
	}
	//แปลงข้อมูลjsonของ api มาเป็น ข้อมูลแบบ struct ที่กำนดไว้
	if err := json.NewDecoder(resp2.Body).Decode(&form); err != nil {
		log.Fatalf("Error decoding JSON from second endpoint: %v", err)
	}

	// รวมข้อมูลเข้าไปใน stuct pokemon
	pokemon.Name = form.Name
	pokemon.Sprites = form.Sprites

	//แปลงข้อมูลที่รวมกันเป็นรูปแบบ json และส่งกลับไป
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pokemon)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	// แสดงข้อความเมื่อเข้าไปในหน้า localhost:8000/
	fmt.Fprintf(w, "Welcome to the Pokemon API")
}
func status(w http.ResponseWriter, r *http.Request) {
	// แสดงข้อความเมื่อเข้าไปในหน้า localhost:8000/
	fmt.Fprintf(w, "Status 200")
}

// 1.POST localhost:8000/id เพิ่มข้อมูลลงไปใน id
// 2.GET localhost:8000/id ดูข้อมูล id หลัง post ไปแล้ว

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage).Methods("GET")
	myRouter.HandleFunc("/", status).Methods("POST")
	myRouter.HandleFunc("/id", getPokemonData).Methods("GET")
	myRouter.HandleFunc("/id", getPokemonData).Methods("POST")

	//สร้าง server localhost:8000 และรอ request จาก client ต่อ
	log.Fatal(http.ListenAndServe(":8000", myRouter))
}

func main() {
	handleRequests()
}
