package card_types

type CardType string

const (
	MasterCard      CardType = "MasterCard"
	VisaCard        CardType = "Visa"
	Mir             CardType = "Mir"
	UnionPay        CardType = "UnionPay"
	JCB             CardType = "JCB"
	AmericanExpress CardType = "AmericanExpress"
	DinersClub      CardType = "DinersClub"
	DiscoverCard    CardType = "DiscoverCard"
	InstaPayment    CardType = "InstaPayment"
	InstaPaymentTM  CardType = "InstaPaymentTM"
	Laser           CardType = "Laser"
	Dankort         CardType = "Dankort"
	Solomon         CardType = "Solo"
	Switch          CardType = "Switch"
	Unknown         CardType = "Unknown"
)
