package smdb2

import (
	"net/http"
)

func SMServer(w http.ResponseWriter, req *http.Request) {
	w.Write("hello world!")
}
func httpSrv(port int) error {
	http.HandleFunc("/sm", SMServer)
	return http.ListenAndServe(":12345", nil)
}
