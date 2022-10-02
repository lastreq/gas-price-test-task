package controlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/lastreq/gas-price-test-task/internal/app/model"
)

type Controller struct {
	svc service
}
type service interface {
	GetGasInfo() (model.ProcessedGasInfo, error)
}

func New(svc service) Controller {
	ctr := Controller{svc: svc}
	return ctr
}

func (ctr *Controller) GetGasInfo(w http.ResponseWriter, r *http.Request) {
	gasInfo, err := ctr.svc.GetGasInfo()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	jsonResp, err := json.Marshal(gasInfo)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(jsonResp)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
