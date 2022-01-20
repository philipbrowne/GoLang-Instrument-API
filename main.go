package main

import "os"

func main(){
	a := App{}
	a.Initialize()
	port, ok := os.LookupEnv("PORT")
    if !ok {
        port = "3000"
    } 
    a.Run(":"+port)
}