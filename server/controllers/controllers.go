package controllers

import (
	"C2-D2/server/auth"
	"C2-D2/server/database"
	"C2-D2/server/models"
	"C2-D2/server/obfuscation"
	"encoding/json"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Set up the routes for the API calls usint the mux Router
func InitializeRoutes(router *mux.Router) {
	router.HandleFunc("/api/agents", GetAgents).Methods("GET")
	router.HandleFunc("/api/register", RegisterAgent).Methods("POST")
	router.HandleFunc("/api/agent", GetAgentById).Methods("GET")
	router.HandleFunc("/api/agent", DeleteAgent).Methods("DELETE")
	router.HandleFunc("/api/checkin", AgentBeacon).Methods("PUT")
	//TO-DO: Tasking Handlers
	//router.HandleFunc("/api/agents/{id}/task", a.CreateTask).Methods("POST")
	//router.HandleFunc("/api/agents/{id}/task/{id}", a.PutResults).Methods("PUT")
	//router.HandleFunc("/api/agents/{id}/tasks", a.GetTasks).Methods("GET")
}

// Initial agent registration
func RegisterAgent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var agent models.Agent
	json.NewDecoder(r.Body).Decode(&agent)

	// Set these to empty so they don't overwrite the generated ones in the DB and cause shenanigans
	agent.UUID = ""
	agent.Token = ""
	// Setting IP to the IP from the request
	IP, _, _ := net.SplitHostPort(r.RemoteAddr)
	agent.IP = string(IP)
	agent.Created = time.Now().Format(time.RFC1123)
	database.DB.Create(&agent)

	response := auth.Login(agent)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Pass agent ID in the URL - may want to change this to another field since UUIDs are massive
func GetAgentById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// TO-DO: Set Auth for CLI/Authenticated GUI session to only hit this endpoint
	var agent models.Agent
	json.NewDecoder(r.Body).Decode(&agent)
	if !checkIfAgentExists(agent.UUID) {
		w.WriteHeader(http.StatusNotFound)
		response := models.Response{Type: "ERROR", Message: "Agent Not Found!", Data: "ID: " + agent.UUID}
		json.NewEncoder(w).Encode(response)
		return
	}
	database.DB.First(&agent, agent.UUID)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(agent)
}

// List all agents in DB
func GetAgents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// TO-DO: Set Auth for CLI/Authenticated GUI session to only hit this endpoint
	var agents []models.Agent
	database.DB.Find(&agents)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(agents)
}

// Update checkin time for agent based on ID
func AgentBeacon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var agent models.Agent
	json.NewDecoder(r.Body).Decode(&agent)

	// Check body for authenticated agent token
	result := auth.Auth(agent)
	if !result {
		w.WriteHeader(http.StatusUnauthorized)
		response := models.Response{Type: "ERROR", Message: "Unauthorized!"}
		json.NewEncoder(w).Encode(response)
		return
	}
	// Check DB to see if agent exists - otherwise what're ya even doin here?
	if !checkIfAgentExists(agent.UUID) {
		w.WriteHeader(http.StatusNotFound)
		response := models.Response{Type: "ERROR", Message: "Agent Not Found!", Data: "ID: " + agent.UUID}
		json.NewEncoder(w).Encode(response)
		return
	}
	// Query DB if checks pass and return
	database.DB.First(&agent, agent.UUID)
	// Grab IP from the request and update it in the DB record
	IP, _, _ := net.SplitHostPort(r.RemoteAddr)
	agent.IP = string(IP)
	// Post the latest checkin, IP, and JWT to the DB
	agent.Checkin = time.Now().Format(time.RFC1123)
	database.DB.Model(&agent).Updates(&models.Agent{IP: agent.IP, Checkin: agent.Checkin, Token: agent.Token})
	w.WriteHeader(http.StatusOK)
	response := models.Response{Type: "SUCCESS", Message: "Time: " + agent.Checkin, Data: "ID: " + agent.UUID}
	byt := obfuscation.EncryptResponse(response, agent.UUID)
	json.NewEncoder(w).Encode(byt)
}

// Delete the agent from the DB
func DeleteAgent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// TO-DO: Set Auth for CLI/Authenticated GUI session to only hit this endpoint
	var agent models.Agent
	json.NewDecoder(r.Body).Decode(&agent)
	if !checkIfAgentExists(agent.UUID) {
		w.WriteHeader(http.StatusNotFound)
		response := models.Response{Type: "ERROR", Message: "Agent Not Found!", Data: "UUID: " + agent.UUID}
		json.NewEncoder(w).Encode(response)
		return
	}
	result := auth.Auth(agent)
	if !result {
		w.WriteHeader(http.StatusUnauthorized)
		response := models.Response{Type: "ERROR", Message: "Unauthorized!"}
		json.NewEncoder(w).Encode(response)
		return
	}
	database.DB.Delete(&agent, agent.UUID)
	w.WriteHeader(http.StatusOK)
	response := models.Response{Type: "SUCCESS", Message: "Agent Deleted Successfully!", Data: "UUID: " + agent.UUID}
	json.NewEncoder(w).Encode(response)
}

// Verify we're not already creating an agent that exists
func checkIfAgentExists(agentUUID string) bool {
	var agent models.Agent
	database.DB.First(&agent, "uuid = ?", agentUUID)
	if agent.UUID == "" {
		return false
	}
	return true
}
