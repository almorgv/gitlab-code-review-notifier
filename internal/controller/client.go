package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"gitlab-code-review-notifier/internal/database"
	client2 "gitlab-code-review-notifier/pkg/config"
)

type ClientController struct {
	repo *database.ClientRepository
}

func NewClientController(repo *database.ClientRepository) *ClientController {
	return &ClientController{repo: repo}
}

func (c *ClientController) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	val, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "Failed to parse :id from value %s: %v", vars["id"], err)
		return
	}

	id := int(val)
	client, err := c.repo.Get(id)

	if err == database.ErrNotFound {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprintf(w, "Client id %d not found", id)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "Failed to get client with id %d: %v", id, err)
		return
	}

	maskClient(client)

	if err := json.NewEncoder(w).Encode(&client); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "Failed to serialize client with id %d: %v", id, err)
		return
	}
}

func (c *ClientController) GetAll(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	clients, err := c.repo.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "Failed to get clients: %v", err)
		return
	}

	for _, client := range clients {
		maskClient(client)
	}

	if err := json.NewEncoder(w).Encode(&clients); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "Failed to serialize clients: %v", err)
		return
	}
}

func (c *ClientController) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var client client2.FiringConfig
	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "Failed to deserialize client from request body: %v", err)
		return
	}

	if err := c.repo.Create(&client); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "Failed to create client %d: %v", client.Id, err)
		return
	}
}

func (c *ClientController) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var client client2.FiringConfig
	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "Failed to deserialize client from request body: %v", err)
		return
	}

	vars := mux.Vars(r)
	val, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "Failed to parse :id from value %s: %v", vars["id"], err)
		return
	}

	id := int(val)
	client.Id = id

	if err := c.repo.Update(&client); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "Failed to save client %d: %v", client.Id, err)
		return
	}
}

func (c *ClientController) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	val, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "Failed to parse :id from value %s: %v", vars["id"], err)
		return
	}

	id := int(val)
	if err := c.repo.Delete(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "Failed to delete client with id %d: %v", id, err)
		return
	}
}

func maskClient(client *client2.FiringConfig) {
	client.GitlabToken = "<MASKED>"
	client.WebhookUrl = "<MASKED>"
}
