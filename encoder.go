package goat

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type EncoderFunc func(http.ResponseWriter, any) // the any part will receive a value, not a pointer

func JSONEncoder(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func HTMLEncoder(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "text/html")
	s, ok := v.(string)
	if !ok {
		fmt.Fprint(w, "value provided to HTMLEncoder is not a string")
		return
	}

	fmt.Fprint(w, s)
}
