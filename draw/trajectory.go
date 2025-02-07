package draw

import (
	_ "embed"
	"encoding/json"
	"net/http"
	"os"
)

//go:embed index.html
var drawHtml []byte

type Point struct {
	X    float32 `json:"x"`
	Y    float32 `json:"y"`
	Time int64   `json:"time"`
}

type Stroke []*Point

type Trajectories struct {
	Strokes []Stroke `json:"trajectories"`
}

func LoadTrajectories(path string) (*Trajectories, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var trajectories Trajectories
	if err := json.Unmarshal(data, &trajectories); err != nil {
		return nil, err
	}
	return &trajectories, nil
}

func ListTrajectories(addr string, robot *Robot) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(drawHtml)
	})
	http.HandleFunc("/trajectories", func(w http.ResponseWriter, r *http.Request) {
		var trajectories Trajectories
		if err := json.NewDecoder(r.Body).Decode(&trajectories); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := robot.Draw(&trajectories, -50, 10, true); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{status: 0}"))
	})
	go http.ListenAndServe(addr, nil)
}
