package main

import (
    "net/http"
    "fmt"
    "log"
    "api/pkg"
    "flag"
)

func main() {
    apiPortFlag := flag.Int("port", 8080, "port to host api on")
    flag.Parse()

    api := rewards.NewRewardAPI()
    
    log.Printf("Rewards API listening on http://localhost:%d", *apiPortFlag)
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *apiPortFlag), api))
}
