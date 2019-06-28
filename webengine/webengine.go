package webengine

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/morfien101/go-metrics-auth/config"
	"github.com/morfien101/go-metrics-auth/redisengine"
)

type WebEngine struct {
	config      config.WebServerConfig
	router      *mux.Router
	server      *http.Server
	redisEngine redisengine.Redis
}

// New creates a new WebEngine and returns it to you. You will still need to start it when you are
// ready.
func New(config config.WebServerConfig, redis *redisengine.RedisEngine) *WebEngine {
	we := &WebEngine{
		config:      config,
		redisEngine: redis,
	}
	we.loadRoutes()
	listenerAddress := we.config.ListenAddress + ":" + we.config.Port
	we.server = &http.Server{Addr: listenerAddress, Handler: we.router}
	return we
}

func (we *WebEngine) loadRoutes() {
	we.router = mux.NewRouter()
	we.router.HandleFunc("/auth", we.getAuth).Methods("GET")
	we.router.HandleFunc("/auth", we.validateAuth).Methods("POST")
	we.router.HandleFunc("/_status", we.getStatus).Methods("GET")
}

// ServeHTTP is used to allow the router to start accepting requests before the start is started up. This will help with testing.
func (we *WebEngine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	we.router.ServeHTTP(w, r)
}

func setContentJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
}

func jsonMarshal(x interface{}) ([]byte, error) {
	return json.MarshalIndent(x, "", "  ")
}

func printJSON(w http.ResponseWriter, jsonbytes []byte) (int, error) {
	return fmt.Fprint(w, string(jsonbytes), "\n")
}

// Start will start the web server using the configuration provided.
// It returns a channel that will give the error if there is one
func (we *WebEngine) Start() <-chan error {
	c := make(chan error, 1)
	startfunc := we.startClear
	if we.config.UseTLS {
		startfunc = we.startTLS
	}
	go func() {
		c <- startfunc()
	}()

	return c
}

func (we *WebEngine) startTLS() error {
	return we.server.ListenAndServeTLS(we.config.CertPath, we.config.KeyPath)
}

func (we *WebEngine) startClear() error {
	return we.server.ListenAndServe()
}

func (we *WebEngine) getAuth(w http.ResponseWriter, r *http.Request) {
	credentials, err := we.redisEngine.CreateCredentials()
	if err != nil {
		// Log here
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	endpoint, err := we.redisEngine.GetEndpoint()
	if err != nil {
		// Log here
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	retStruct := struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Endpoint string `json:"endpoint"`
	}{
		Username: credentials[0],
		Password: credentials[1],
		Endpoint: endpoint,
	}
	b, err := jsonMarshal(retStruct)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	setContentJSON(w)
	w.Write(b)
}

func (we *WebEngine) validateAuth(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	username, ok := q["username"]
	if !ok || len(username) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	password, ok := q["password"]
	if !ok || len(password) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ok, err := we.redisEngine.ValidateAuth(username[0], password[0])
	if err != nil {
		log.Println("Error while validating credentials. Error:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func (we *WebEngine) getStatus(w http.ResponseWriter, r *http.Request) {
	// Make this better.
	w.Write([]byte("OK"))
}
