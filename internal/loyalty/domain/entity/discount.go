package entity

var discountTable = map[LoyaltyStatus]int{
	Bronze: 5,
	Silver: 7,
	Gold:   10,
}

func getDiscount(loyaltyStatus LoyaltyStatus) int {
	return discountTable[loyaltyStatus]
}
