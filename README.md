# Receipt Processor - Adam Sikora

This is a web service for processor information from receipts.

* Infrastructure code is in the main.go file
* Points tallying code is in the points_tally.go file

## Acknowledgments

* All receipts in the examples folder will be processed and given an id.  This is why an array of ids is the initial response.
* I used an int as my id in order to better be able to test

## With Or Without You(Web Service)

By toggling comments in the main function, you'll be able to run the application as a web service or a regular program.  Instructions are in the main function.

## Instructions To Run

* Clone this repo onto your local machine
* I'm assuming you have Go installed
* In a terminal, go to the root directory of the project
* To test a unique receipt, just add one to the examples folder
* Type: go run .

### Web Service

This is the default.  The web service will run until you manually exit it.(Ctrl + C)

* Open a different terminal, again navigating to the root directory
* Type: curl http://localhost:8080/receipts/process
* This returns an array of int ids according to the files in the examples folder
* Type: curl http://localhost:8080/receipts/{id}/points
* Replace {id} with the int id that you want to know the points of

### No Web Service

Without a webservice, the program will just run and exit.  It will process the receipts and respond with a breakdown of points for the receipt id that it given in the main function(id_to_tally).
