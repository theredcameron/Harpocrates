package harpocrates

import (
    "testing"
    "flag"
)
    
var username = flag.String("username", "", "Username for testing")
var password = flag.String("password", "", "Password for testing")
var port     = flag.String("port", "", "Port number")

func setup() *LinkAgent {
    flag.Parse()

    linkAgent := NewLinkAgent(*port, *username, *password)

    return linkAgent
}

func TestGetAllUsers(t *testing.T) {
    linkAgent := setup()

    userCount, err := linkAgent.GetAllUsersCount()
    if err != nil {
        t.Fatal(err)
    }

    if userCount <= 0 {
        t.Error("The number of users is 0.")
    }
}
