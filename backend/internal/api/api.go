package api

import (
	"backend/internal/influxdb"
	"encoding/json"
	"fmt"
	"net/http"
)

// StartServer démarre le serveur HTTP
func StartServer() {
	http.HandleFunc("/", handleGetFlows)
	fmt.Println("API Server running ...!")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Erreur lors du démarrage du serveur HTTP : %v\n", err)
	}
}

// handleGetFlows gère les requêtes pour récupérer les flux NetFlow
func handleGetFlows(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Récupération des flux depuis InfluxDB
	flows, err := influxdb.QueryNetFlowData()
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de la récupération des données : %v", err), http.StatusInternalServerError)
		return
	}

	// Conversion des données en JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(flows); err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de l'encodage JSON : %v", err), http.StatusInternalServerError)
	}
}
