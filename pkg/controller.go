package rewards

import (
    "encoding/json"
    "fmt"
    "net/http"
    "log"
)

type JSONResponse struct {
    Success bool `json:"success"`
    Data interface{} `json:"data,omitempty"`
    Error interface{} `json:"error,omitempty"`
}

type RewardAPI struct {
    *http.ServeMux
    store *RewardStore
}

func NewRewardAPI() *RewardAPI {
    var api RewardAPI = RewardAPI{
        ServeMux: http.NewServeMux(),
        store: NewRewardStore(),
    }

    api.HandleFunc("/add-points", api.HandleAddPoints)
    api.HandleFunc("/use-points", api.HandleUsePoints)
    api.HandleFunc("/check-balance", api.HandleCheckBalance)
    return &api
}

func http405(w http.ResponseWriter) {
    http.Error(w, "405 - Method Not Allowed!", http.StatusMethodNotAllowed)
}

func httpJSONResponse(w http.ResponseWriter, j JSONResponse) {
    err := json.NewEncoder(w).Encode(j)
	if err != nil {
        http.Error(w, fmt.Sprintf("Error encoding json request: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
}

func (a *RewardAPI) HandleAddPoints (w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http405(w)
        return
    }

    var rewards []interface{}
    err := json.NewDecoder(r.Body).Decode(&rewards)
    if err != nil {
        log.Printf("Failed to decode reward: %s", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    for _, ireward := range rewards  {
        reward := ireward.(map[string]interface{})
        ts, okTs := reward["timestamp"]
        payer, okPayer := reward["payer"]
        points, okPoints := reward["points"]
        if !okTs || !okPayer || !okPoints {
            httpJSONResponse(w, JSONResponse{false, nil, "Missing required parameters"})
            return
        }

        err = a.store.AddReward(ts.(string), int64(points.(float64)), payer.(string))
        if err != nil {
            httpJSONResponse(w, JSONResponse{false, nil, err.Error()})
            return
        }
    }

    httpJSONResponse(w, JSONResponse{true, nil, nil})
}

func (a *RewardAPI) HandleUsePoints (w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http405(w)
        return
    }

    var use map[string]interface{}
    err := json.NewDecoder(r.Body).Decode(&use)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    points, okPoints := use["points"]
    if !okPoints {
        httpJSONResponse(w, JSONResponse{false, nil, "Missing required parameters"})
        return
    }

    deducted, err := a.store.UsePoints(int64(points.(float64)))
    if err != nil {
        httpJSONResponse(w, JSONResponse{false, nil, err.Error()})
    } else {
        httpJSONResponse(w, JSONResponse{true, deducted, nil})
    }
}

func (a *RewardAPI) HandleCheckBalance (w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
        http405(w)
        return
    }

    balance := a.store.CheckBalance()
    httpJSONResponse(w, JSONResponse{true, balance, nil})
}
