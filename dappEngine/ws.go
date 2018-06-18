package dappEngine

import (
	"encoding/json"
	"fmt"
	"gx/ipfs/QmZH5VXfAJouGMyCCHTRPGCT3e5MG9Lu78Ln3YAYW1XTts/websocket"
	"net/http"
	"sync"
)

type stack struct {
	lock     sync.Mutex
	requests map[string]request
}

type request struct {
	ID     string                 `json:"id"`
	Body   map[string]interface{} `json:"body"`
	closer chan struct{}
}

type response struct {
	ID string `json:"id"`
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

	sta := stack{
		lock:     sync.Mutex{},
		requests: map[string]request{},
	}

	// write logic
	go func() {

		for {

			select {
			case request := <-sendChan:

				sta.lock.Lock()
				sta.requests[request.ID] = request
				sta.lock.Unlock()

				raw, err := json.Marshal(request)
				if err != nil {
					fmt.Println(err)
				}
				conn.WriteJSON(string(raw))

			}
		}
	}()

	// read logic
	go func() {

		for {

			req := response{}
			err := conn.ReadJSON(&req)
			if err != nil {
				fmt.Println("Error reading json.", err)
			}

			sta.lock.Lock()
			request := sta.requests[req.ID]
			request.closer <- struct{}{}
			sta.lock.Unlock()

		}

	}()

}
