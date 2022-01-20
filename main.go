package main

import "os"

func main(){
	a := App{}
	a.Initialize()
	    port, err := os.Getenv("PORT")
    if err != nil {
        port = "3000"
    } 
    a.Run(":"+port)
}