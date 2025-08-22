package ext

import (
	"fmt"
	"net/http"
)

func (cfg *ApiConfig) ShowCountHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	template := "<html> <body> <h1>Welcome, Chirpy Admin</h1> <p>Chirpy has been visited %d times!</p> </body> </html>"

	w.Write([]byte(fmt.Sprintf(template, cfg.FileserverHits.Load())))
	fmt.Printf("Total Hits: %v\n", cfg.FileserverHits.Load())
}
