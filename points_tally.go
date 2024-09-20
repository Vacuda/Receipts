package main

import (
	"fmt"
	"math"
	"strings"
	"unicode"
)

/* POINTS TALLY BY RULES */

//+1 = Every alphanumeric character in retailer
func PointsTally_Rule1(rec *Receipt) (int){

	var points int = 0

	//loop retailer name
	for _, char := range rec.Retailer{
		//if char is letter or digit
		if(unicode.IsLetter(char) || unicode.IsDigit(char)){
			points++
		}
	}

	/*DEBUG*/
	fmt.Print("Points - Rule1: ")
	fmt.Println(points)

	return points
}

//+50 = Total is round dollar amount
func PointsTally_Rule2(rec *Receipt) (int){

	var points int = 0

	//if round dolar amount
	if(int(rec.Total * 100) % 100 == 0){
		points += 50
	}

	/*DEBUG*/
	fmt.Print("Points - Rule2: ")
	fmt.Println(points)

	return points
}

//+25 = Total is multiple of .25
func PointsTally_Rule3(rec *Receipt) (int){

	var points int = 0

	//if multiple of .25
	if(int(rec.Total * 100) % 25 == 0){
		points += 25
	}

	/*DEBUG*/
	fmt.Print("Points - Rule3: ")
	fmt.Println(points)

	return points
}

//+5 = every two items in items
func PointsTally_Rule4(rec *Receipt) (int){
	
	var points int = 0

	//length of Items / 2, truncated down to an int, will give you the multiple of 5 points to add
	points += 5*(int(len(rec.Items)/2))

	/*DEBUG*/
	fmt.Print("Points - Rule4: ")
	fmt.Println(points)

	return points
}

//if trimmed length of desc multiple of 3, multiple price by 0.2, round up, result is points earned
func PointsTally_Rule5(rec *Receipt) (int){
	
	var points int = 0
	
	//loop Items
	for _,i := range rec.Items{

		//get trimmed length of ShortDescription chars
		var amount int = len(strings.TrimSpace(i.ShortDescription))
	
		//if amount multiple of 3
		if(amount % 3 == 0){

			//multiply price, round up, add to points
			points += int(math.Ceil(float64(i.Price * 0.2)))
		}
	}

	/*DEBUG*/
	fmt.Print("Points - Rule5: ")
	fmt.Println(points)

	return points
}

//+6 = day of purchase date is odd
func PointsTally_Rule6(rec *Receipt) (int){

		var points int = 0

		//if day is odd
		if(rec.PurchaseDateTime.Day() % 2 != 0){
			points += 6
		}

		/*DEBUG*/
		fmt.Print("Points - Rule6: ")
		fmt.Println(points)
	
		return points
}

//+10 = Time is between 2-4pm
func PointsTally_Rule7(rec *Receipt) (int){

	var points int = 0
	var hour int = rec.PurchaseDateTime.Hour()

	//if after 2pm(14) and before 4pm(16)
	if(hour == 14 || hour == 15){

		/* I'm assuming here that 2pm is inclusive and 4pm is exclusive */

		points += 10
	}

	/*DEBUG*/
	fmt.Print("Points - Rule7: ")
	fmt.Println(points)

	return points
}