package main

import (
    "github.com/eyedeekay/susc"
)

func main(){
    susc,err := susc.NewClient()
    if err != nil {
        panic(err)
    }
    err = susc.Hello()
    if err != nil {
        panic(err)
    }
}
