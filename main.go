package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/lib/pq"
)

type Country struct {
	Name       string   `json:"name"`
	Population int64    `json:"population"`
	FlagColors []string `json:"flag_colors"`
	Language   string   `json:"language"`
}

var countries = []Country{
	{Name: "Afghanistan", Population: 41454761, FlagColors: []string{"black", "red", "green", "white"}, Language: "Pashto"},
	{Name: "Algeria", Population: 44903225, FlagColors: []string{"green", "white", "red"}, Language: "Arabic"},
	{Name: "Argentina", Population: 46234830, FlagColors: []string{"light blue", "white", "yellow"}, Language: "Spanish"},
	{Name: "Australia", Population: 26658948, FlagColors: []string{"blue", "red", "white"}, Language: "English"},
	{Name: "Austria", Population: 9132383, FlagColors: []string{"red", "white"}, Language: "German"},
	{Name: "Bangladesh", Population: 171186372, FlagColors: []string{"green", "red"}, Language: "Bengali"},
	{Name: "Belgium", Population: 11763930, FlagColors: []string{"black", "yellow", "red"}, Language: "Dutch"},
	{Name: "Brazil", Population: 216422446, FlagColors: []string{"green", "yellow", "blue", "white"}, Language: "Portuguese"},
	{Name: "Canada", Population: 40097761, FlagColors: []string{"red", "white"}, Language: "English"},
	{Name: "Chile", Population: 19629590, FlagColors: []string{"red", "white", "blue"}, Language: "Spanish"},
	{Name: "China", Population: 1411750000, FlagColors: []string{"red", "yellow"}, Language: "Mandarin Chinese"},
	{Name: "Colombia", Population: 52215503, FlagColors: []string{"yellow", "blue", "red"}, Language: "Spanish"},
	{Name: "Denmark", Population: 5961249, FlagColors: []string{"red", "white"}, Language: "Danish"},
	{Name: "Egypt", Population: 112716598, FlagColors: []string{"red", "white", "black", "gold"}, Language: "Arabic"},
	{Name: "Ethiopia", Population: 126527060, FlagColors: []string{"green", "yellow", "red", "blue"}, Language: "Amharic"},
	{Name: "Finland", Population: 5545475, FlagColors: []string{"white", "blue"}, Language: "Finnish"},
	{Name: "France", Population: 68170228, FlagColors: []string{"blue", "white", "red"}, Language: "French"},
	{Name: "Germany", Population: 84482267, FlagColors: []string{"black", "red", "gold"}, Language: "German"},
	{Name: "Ghana", Population: 34121985, FlagColors: []string{"red", "yellow", "green", "black"}, Language: "English"},
	{Name: "Greece", Population: 10341277, FlagColors: []string{"blue", "white"}, Language: "Greek"},
	{Name: "India", Population: 1428627663, FlagColors: []string{"saffron", "white", "green", "navy blue"}, Language: "Hindi"},
	{Name: "Indonesia", Population: 277534122, FlagColors: []string{"red", "white"}, Language: "Indonesian"},
	{Name: "Iran", Population: 89172767, FlagColors: []string{"green", "white", "red"}, Language: "Persian"},
	{Name: "Iraq", Population: 45504560, FlagColors: []string{"red", "white", "black", "green"}, Language: "Arabic"},
	{Name: "Ireland", Population: 5262382, FlagColors: []string{"green", "white", "orange"}, Language: "English"},
	{Name: "Israel", Population: 9756700, FlagColors: []string{"white", "blue"}, Language: "Hebrew"},
	{Name: "Italy", Population: 58870762, FlagColors: []string{"green", "white", "red"}, Language: "Italian"},
	{Name: "Japan", Population: 123294513, FlagColors: []string{"white", "red"}, Language: "Japanese"},
	{Name: "Kenya", Population: 55100586, FlagColors: []string{"black", "red", "green", "white"}, Language: "Swahili"},
	{Name: "Malaysia", Population: 34308525, FlagColors: []string{"red", "white", "blue", "yellow"}, Language: "Malay"},
	{Name: "Mexico", Population: 128455567, FlagColors: []string{"green", "white", "red"}, Language: "Spanish"},
	{Name: "Morocco", Population: 37840044, FlagColors: []string{"red", "green"}, Language: "Arabic"},
	{Name: "Netherlands", Population: 17618299, FlagColors: []string{"red", "white", "blue"}, Language: "Dutch"},
	{Name: "New Zealand", Population: 5228100, FlagColors: []string{"blue", "red", "white"}, Language: "English"},
	{Name: "Nigeria", Population: 223804632, FlagColors: []string{"green", "white"}, Language: "English"},
	{Name: "Norway", Population: 5474360, FlagColors: []string{"red", "white", "blue"}, Language: "Norwegian"},
	{Name: "Pakistan", Population: 240485658, FlagColors: []string{"green", "white"}, Language: "Urdu"},
	{Name: "Peru", Population: 34352719, FlagColors: []string{"red", "white"}, Language: "Spanish"},
	{Name: "Philippines", Population: 117337368, FlagColors: []string{"blue", "red", "white", "yellow"}, Language: "Filipino"},
	{Name: "Poland", Population: 41026067, FlagColors: []string{"white", "red"}, Language: "Polish"},
	{Name: "Portugal", Population: 10247605, FlagColors: []string{"green", "red", "yellow"}, Language: "Portuguese"},
	{Name: "Russia", Population: 144444359, FlagColors: []string{"white", "blue", "red"}, Language: "Russian"},
	{Name: "Saudi Arabia", Population: 36947025, FlagColors: []string{"green", "white"}, Language: "Arabic"},
	{Name: "South Africa", Population: 60414495, FlagColors: []string{"black", "yellow", "green", "white", "red", "blue"}, Language: "Zulu"},
	{Name: "South Korea", Population: 51784059, FlagColors: []string{"white", "red", "blue", "black"}, Language: "Korean"},
	{Name: "Spain", Population: 47519628, FlagColors: []string{"red", "yellow"}, Language: "Spanish"},
	{Name: "Sweden", Population: 10612086, FlagColors: []string{"blue", "yellow"}, Language: "Swedish"},
	{Name: "Thailand", Population: 71801279, FlagColors: []string{"red", "white", "blue"}, Language: "Thai"},
	{Name: "Turkey", Population: 85816199, FlagColors: []string{"red", "white"}, Language: "Turkish"},
	{Name: "United States", Population: 339996563, FlagColors: []string{"red", "white", "blue"}, Language: "English"},
}

func main() {
	db, err := openDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := setupDatabase(db); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/api/countries", countriesHandler(db))
	mux.HandleFunc("/countries", countriesHandler(db))

	port := getEnv("PORT", "8080")
	log.Printf("server running on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"message":  "Countries API is running",
		"endpoint": "/api/countries",
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func countriesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		countries, err := getCountries(db)
		if err != nil {
			log.Printf("failed to get countries: %v", err)
			http.Error(w, "failed to get countries", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(countries); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}

func openDatabase() (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_NAME", "postgres"),
		getEnv("DB_SSLMODE", "disable"),
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	for attempt := 1; attempt <= 10; attempt++ {
		if err := db.Ping(); err == nil {
			return db, nil
		} else if attempt == 10 {
			db.Close()
			return nil, fmt.Errorf("connect database: %w", err)
		}

		time.Sleep(time.Second)
	}

	return db, nil
}

func setupDatabase(db *sql.DB) error {
	createTable := `
		CREATE TABLE IF NOT EXISTS countries (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			population BIGINT NOT NULL,
			flag_colors TEXT[] NOT NULL,
			language TEXT NOT NULL
		);
	`
	if _, err := db.Exec(createTable); err != nil {
		return fmt.Errorf("create countries table: %w", err)
	}

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM countries").Scan(&count); err != nil {
		return fmt.Errorf("count countries: %w", err)
	}
	if count > 0 {
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin seed transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO countries (name, population, flag_colors, language)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (name) DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("prepare seed statement: %w", err)
	}
	defer stmt.Close()

	for _, country := range countries {
		if _, err := stmt.Exec(country.Name, country.Population, pq.Array(country.FlagColors), country.Language); err != nil {
			return fmt.Errorf("seed country %q: %w", country.Name, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit seed transaction: %w", err)
	}

	log.Printf("seeded %d countries", len(countries))
	return nil
}

func getCountries(db *sql.DB) ([]Country, error) {
	rows, err := db.Query(`
		SELECT name, population, flag_colors, language
		FROM countries
		ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]Country, 0)
	for rows.Next() {
		var country Country
		if err := rows.Scan(&country.Name, &country.Population, pq.Array(&country.FlagColors), &country.Language); err != nil {
			return nil, err
		}
		result = append(result, country)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
