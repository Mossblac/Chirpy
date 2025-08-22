package ext

import (
	"fmt"
	"net/http"
)

func (cfg *ApiConfig) ResetHandler(w http.ResponseWriter, r *http.Request) {
	access, err := DevAccess()
	if err != nil {
		w.Write([]byte(fmt.Sprintln("requires Dev permission")))
		w.WriteHeader(403)
	}
	if access {
		cfg.FileserverHits.Store(0)
		w.Write([]byte(fmt.Sprintln("hit count reset")))
		fmt.Printf("hit counter reset\n")

		err := cfg.DB.DeleteAllUsers(r.Context())
		if err != nil {
			WriteError(w, err)
		}
		w.Write([]byte(fmt.Sprintln("Users Reset")))
		fmt.Printf("users reset\n")
	}
}
