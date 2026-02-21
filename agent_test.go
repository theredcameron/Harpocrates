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

    config := LinkAgentConfig{
        LinkStorePort: *port,
        UserName: *username,
        Password: *password,
    }

    linkAgent := NewLinkAgent(config)

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

func TestCreateLocation(t *testing.T) {
    linkAgent := setup()

    err := linkAgent.CreateLocationEntry(1, 0.0000000000, 1.111111111111)
    if err != nil {
        t.Fatal(err)
    }
}
