package main

// // Get All Instruments
// func getInstruments(w http.ResponseWriter, r *http.Request){

// }
// // Get One Instrument
// func getInstruments(w http.ResponseWriter, r *http.Request){

// }

// // Get All Instruments
// func getInstruments(w http.ResponseWriter, r *http.Request){

// }
// // Get All Instruments
// func getInstruments(w http.ResponseWriter, r *http.Request){

// }

// func insert(collection *mongo.Collection, inst {}){
// 	instrument := Instrument{
// 		count(collection), inst.Name, inst.Family, inst.Price
// 	}
// 	insertResult, err := collection.InsertOne(context.TODO(), instrument)
// 	if err != nil {
//     	log.Fatal(err)
// 	}
// 	fmt.Println(insertResult)
// }


func main(){
	a := App{}
	a.Initialize()
	a.Run(":8080")
}