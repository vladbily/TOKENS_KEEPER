package main

import (
	"fmt"

	"gorm.io/gorm"
)

func totalSum(DB gorm.DB) float64 {
	var tS float64
	if err := DB.Model(&Infos{}).
		Select("SUM(Current_Price * Count)").
		Scan(&tS).Error; err != nil {
		fmt.Println("Error executing query:", err)
	}
	return tS
}

/*func totalDif() string {
	total_dif := 0.0
	for i := range tokens {
		count_res := tokens[i].Count
		price_res := tokens[i].Buy_Price
		curr_res := tokens[i].Current_Price
		total_dif += count_res * (curr_res - price_res)
	}
	return fmt.Sprintf("%.4f", total_dif)
}*/
