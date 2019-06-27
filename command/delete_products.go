package command

import (
	"strconv"

	"github.com/salaleser/vk-api/method/market"
)

// DeleteProducts удаляет указанное число товаров
func deleteProducts(ownerID string, limit int) {
	for i := 0; i < limit; i++ {
		o := market.Get(ownerID, "", "1", "", "")
		if o.R.Count == 0 {
			return
		}
		ID := o.R.Items[0].ID
		itemID := strconv.Itoa(ID)
		market.Delete(ownerID, itemID)
	}
}
