package models

type TimestampedPrice struct {
	Timestamp int32
	Price     int32
}

type PriceHistory struct {
	TimestampedPrices []TimestampedPrice
}

func NewPriceHistory() *PriceHistory {
	return &PriceHistory{
		TimestampedPrices: make([]TimestampedPrice, 0),
	}
}

func (p *PriceHistory) AddPrice(timestamp int32, price int32) {
	p.TimestampedPrices = append(p.TimestampedPrices, TimestampedPrice{
		Timestamp: timestamp,
		Price:     price,
	})
}

func (p *PriceHistory) MeanPrice(startTime int32, endTime int32) int32 {
	var totalPrice int64
	var totalCount int64
	for _, timestampedPrice := range p.TimestampedPrices {
		if timestampedPrice.Timestamp < startTime || timestampedPrice.Timestamp > endTime {
			continue
		}

		totalPrice += int64(timestampedPrice.Price)
		totalCount++
	}

	if totalCount == 0 {
		return 0
	} else {
		return int32(totalPrice / totalCount)
	}
}
