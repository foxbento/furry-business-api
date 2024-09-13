package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/foxbento/furry-business-api/db"
	"github.com/foxbento/furry-business-api/models"
)

func GetBusinesses(w http.ResponseWriter, r *http.Request) {
	page, pageSize := getPaginationParams(r)
	offset := (page - 1) * pageSize

	query := `
		SELECT id, "Name", "Store/Socials Link", "Type of Clothing", "Country/Continent", "Country/State", 
		"NSFW?", "General Overview/Personal Notes", "Gendered?", "Convention appearances?", "Notes"
		FROM businesses 
		ORDER BY id 
		LIMIT $1 OFFSET $2
	`
	rows, err := db.DB.Query(query, pageSize, offset)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve businesses")
		return
	}
	defer rows.Close()

	var businesses []models.Business
	for rows.Next() {
		var b models.Business
		var nsfw, gendered, conventions string
		err := rows.Scan(
			&b.ID, &b.Name, &b.Link, &b.Type, &b.Country, &b.State,
			&nsfw, &b.Overview, &gendered, &conventions, &b.Notes,
		)
		if err != nil {
			log.Printf("Error scanning business row: %v", err)
			continue
		}
		b.NSFW = strings.ToLower(nsfw) == "yes" || strings.ToLower(nsfw) == "true"
		b.Gendered = gendered
		b.Conventions = strings.ToLower(conventions) == "yes" || strings.ToLower(conventions) == "true"
		businesses = append(businesses, b)
	}

	if err = rows.Err(); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error processing business data")
		return
	}

	totalCount, err := getTotalCount()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve total count")
		return
	}

	response := struct {
		Businesses  []models.Business `json:"businesses"`
		TotalCount  int               `json:"totalCount"`
		CurrentPage int               `json:"currentPage"`
		PageSize    int               `json:"pageSize"`
	}{
		Businesses:  businesses,
		TotalCount:  totalCount,
		CurrentPage: page,
		PageSize:    pageSize,
	}

	respondWithJSON(w, http.StatusOK, response)
}

func getPaginationParams(r *http.Request) (int, int) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return page, pageSize
}

func getTotalCount() (int, error) {
	var count int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM businesses").Scan(&count)
	return count, err
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}