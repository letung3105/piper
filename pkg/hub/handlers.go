package hub

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

const clientBufSize = 1024

// ServeHTTP handles upgrading and maintaining websocket connection with client
func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// update client connection to websocket
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorf("could not upgrade websocket connection; got %v", err)
		return
	}

	// create and subscribe new client
	wsClient := &WSClient{
		username: r.RemoteAddr,
		nMsgRead: 0,
		free:     true,
		h:        h,
		wsConn:   wsConn,
		send:     make(chan *packet, clientBufSize),
	}
	h.subscribe <- wsClient

	// start reading messages from client and send to broadcast
	go wsClient.readPipe()

	// start writing messages from broadcast channel to client
	go wsClient.writePipe()
}

// Control starts and stops script from running
func (h *Hub) Control() http.HandlerFunc {
	type request struct {
		Action string `json:"action"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			log.Errorf("could not parse request; got %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		h.Lock()
		switch req.Action {
		case "start":
			h.isBroadcasting = true
			log.Info("Start broadcasting")
			w.WriteHeader(http.StatusOK)
		case "stop":
			h.isBroadcasting = false
			log.Info("Stop broadcasting")
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusBadRequest)
			log.Errorf("invalid action")
		}
		h.Unlock()
	}
}
