package player

type OrderType int64

const (
	OrderA OrderType = iota
	OrderB
	OrderC
)

func OrderKey(typ OrderType) string {
	switch typ {
	case OrderA:
		return "OrderA"
	case OrderB:
		return "OrderB"
	case OrderC:
		return "OrderC"
	}
	return "Order"
}
