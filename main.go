package main

import (
	"fmt"
	"time"
	"os"
	"encoding/json"
	"log"
	"strconv"
	"github.com/gin-gonic/gin"
	"net/http"
)

/* STRUCTS */

type Item struct {
	ShortDescription string
	Price            float32
}

type Receipt struct{
	ID 					int
	Retailer 			string
	PurchaseDateTime	time.Time
	Total				float32
	Items        		[]*Item
}

type Receipt_raw struct{
	Retailer 		string 		`json:"retailer"`	
	PurchaseDate	string		`json:"purchaseDate"`	
	PurchaseTime	string		`json:"purchaseTime"`	
	Total			string		`json:"total"`	
	Items        	[]struct {
		ShortDescription string `json:"shortDescription"`
		Price            string `json:"price"`
	} 	
}

type response_ids struct {
	Ids			[]int
}

type response_points struct {
	Points		int
}

/* DATA */

var Receipts = make(map[int]*Receipt)
var UniqueID_Counter int = 1

//non web service id to run
var id_to_tally int

/* MAIN */

func main(){

	/* You can run this with or without the web service */
	/* Uncomment just Section A = With Web Service */
	/* Uncomment just Section B = Withou Web Service */

	/* SECTION A */
	ConfigureAndRun_LocalServer()

	/* SECTION B */

	// id_to_tally = 4    /* Change ID Here to change which receipt is run 1-4 */
	// GatherReceipts_FromExamplesFolder()
	// DEBUG_handleGetPointsTotal()

}

/* SERVER INFORMATION */

func ConfigureAndRun_LocalServer(){

	//use gin to create a router
	router := gin.Default()

	//ENDPOINT: Returns ID of all receipts in examples folder
	router.GET("/receipts/process", handleProcessReceipts)
	// (this needs to be a GET because we are not actually receiving the json through http.  It's just in memory)

	//ENDPOINT: Returns points of receipt at this ID
	router.GET("/receipts/:id/points", handleGetPointsTotal)

	//This starts a local server
	//Use command: curl http://localhost/{and the above routes} to access service
	router.Run("localhost:8080")

	/*
		COMMANDS TO COPY FOR SEPARATE TERMINAL
		--------------------------------------
		curl http://localhost:8080/receipts/process
		curl http://localhost:8080/receipts/       /points
		
	*/
}

/* HANDLER FUNCTIONS */

func handleProcessReceipts(c *gin.Context){

	GatherReceipts_FromExamplesFolder()

	//make slice of ids, empty but capacity set at the length of the map
	ids := make([]int, 0, len(Receipts))

	//loop Receipts map
	for key := range Receipts {
		//fill ids slice
		ids = append(ids, key)
	}

	//create response package
	res := new(response_ids)
	res.Ids = ids

	//send response to client
	c.JSON(http.StatusOK, res)
}

func handleGetPointsTotal(c *gin.Context){

	//get id to use
	id,_ := strconv.Atoi(c.Param("id"))

	//if id is in map
	if _, ok := Receipts[id]; ok{

		var receipt *Receipt = Receipts[id]
		var points int = 0
	
		//find points to add
		points += PointsTally_Rule1(receipt)
		points += PointsTally_Rule2(receipt)
		points += PointsTally_Rule3(receipt)
		points += PointsTally_Rule4(receipt)
		points += PointsTally_Rule5(receipt)
		points += PointsTally_Rule6(receipt)
		points += PointsTally_Rule7(receipt)

		//create response package
		res := new(response_points)
		res.Points = points

		//send response to client
		c.JSON(http.StatusOK, res)
		
		/*DEBUG*/
		fmt.Print("Total = = ")
		fmt.Println(points)

		
	}else{
		
		//send bad request
		c.JSON(http.StatusBadRequest, "Bad Request - ID Doesn't Exist")
	}
}

/* DEBUG */
func DEBUG_handleGetPointsTotal(){

	var receipt *Receipt = Receipts[id_to_tally]
	var points int = 0

	//find points to add
	points += PointsTally_Rule1(receipt)
	points += PointsTally_Rule2(receipt)
	points += PointsTally_Rule3(receipt)
	points += PointsTally_Rule4(receipt)
	points += PointsTally_Rule5(receipt)
	points += PointsTally_Rule6(receipt)
	points += PointsTally_Rule7(receipt)

	/*DEBUG*/
	fmt.Print("Total Points: ")
	fmt.Println(points)
}

/* UTILITIES */

func GatherReceipts_FromExamplesFolder(){

	/* This enables you to put new receipts in the examples folder and they will be processed */

	//read directory of examples
	file, err := os.Open("examples/")

	//err check
	if err != nil {
		log.Fatal(err)
	}

	//ensure order
	defer file.Close()

	//get filenames
	ListOf_ExampleReceipt_FileNames, err := file.Readdirnames(0)

	//err check
	if err != nil {
		log.Fatal(err)
	}

	//loop list of receipt file names
	for _, name := range ListOf_ExampleReceipt_FileNames {
		
		json_receipt, err := os.ReadFile("examples/" + name)
		
		//err check
		if err != nil {
			log.Fatal(err)
		}
		
		//get raw receipt
		var this_raw_receipt Receipt_raw
		json.Unmarshal(json_receipt, &this_raw_receipt)

		//process receipt
		var this_receipt *Receipt = Process_Raw_Receipt(&this_raw_receipt)

		//put into Receipts map
		Receipts[this_receipt.ID] = this_receipt
	}
}

func Process_Raw_Receipt(rec *Receipt_raw) (*Receipt) {

	//make new receipt
	var receipt *Receipt = new(Receipt)

	//create unique ID for this receipt
	receipt.ID = Get_UniqueID()

	//set DateTime - Merge Data and Time into one time object
	{
		//build layout strings to use in time.Parse
		var layout string = "2006-01-02T15:04"
		var datetime_string = rec.PurchaseDate + "T" + rec.PurchaseTime

		//get parsed datetime
		parsed_datetime, err := time.Parse(layout, datetime_string)

		//err check
		if err != nil {
			log.Fatal(err)
		}

		//set datetime value
		receipt.PurchaseDateTime = parsed_datetime
	}
	
	//set other values
	receipt.Retailer = rec.Retailer
	receipt.Total = Parse_String_ToFloat32(rec.Total)
	receipt.Items = make([]*Item, 0)

	//loop rec.Items
	for _, item := range rec.Items {

		//make new Item
		var new_item = new(Item)

		//alter values
		new_item.ShortDescription = item.ShortDescription
		new_item.Price = Parse_String_ToFloat32(item.Price)

		//add to .Items
		receipt.Items = append(receipt.Items, new_item)
	}

	return receipt
}

func Parse_String_ToFloat32(s string) (float32){

	//get float64
	num, err := strconv.ParseFloat(s, 32)

	//err check
	if err != nil {
		log.Fatal(err)
	}

	//return float32
	return float32(num)
}

func Get_UniqueID()(int){

	//store current ID available
	var int_to_return int = UniqueID_Counter

	//increment for next use
	UniqueID_Counter++

	//return stored value
	return int_to_return
}