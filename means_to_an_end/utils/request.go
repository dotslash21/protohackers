package utils

import (
	"encoding/binary"
	"fmt"

	"arunangshu.dev/protohackers/means_to_an_end/models"
)

type RequestProcessor struct {
	priceHistory models.PriceHistory
}

func NewRequestProcessor() *RequestProcessor {
	return &RequestProcessor{
		priceHistory: *models.NewPriceHistory(),
	}
}

func (rp *RequestProcessor) ProcessRequest(request *models.Request) models.Response {
	switch request.Method {
	case 'I':
		rp.priceHistory.AddPrice(request.Number1, request.Number2)

		return *models.NewResponse(request.Method, nil, nil)
	case 'Q':
		startTime := request.Number1
		endTime := request.Number2

		mean := rp.priceHistory.MeanPrice(startTime, endTime)

		return *models.NewResponse(request.Method, &mean, nil)
	default:
		return *models.NewResponse(request.Method, nil, fmt.Errorf("invalid method: %c", request.Method))
	}
}

func ParseRequest(requestBytes []byte) (*models.Request, error) {
	// Parse the request
	request := models.NewRequest(
		requestBytes[0],
		int32(binary.BigEndian.Uint32(requestBytes[1:5])),
		int32(binary.BigEndian.Uint32(requestBytes[5:9])))
	return request, nil
}
