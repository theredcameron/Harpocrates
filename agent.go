package harpocrates

import (
    "encoding/json"
    "net/http"
    "bytes"
    "io"
    "fmt"
)

type LoginInfo struct {
    UserName	string	`json:"username"`
    Password	string	`json:"password"`
}

type LoginResult struct {
    Token	string	`json:"Token"`
}

type LinkAgent struct {
    linkStorePort   string
    userName                string
    password                string
    secretToken             string
}

type UserSearch struct {
	UserName	string	`json:"username"`
	FirstName	string	`json:"firstname"`
	LastName	string	`json:"lastname"`
	Active		string	`json:"active"`
}

type User struct {
	Id			int		`json:"UserId"`
	UserName		string		`json:"UserName"`
	Password		string		`json:"UserPassword"`
	CreatedDate 		string		`json:"CreatedDate"`
	FirstName		string		`json:"FirstName"`
	LastName		string		`json:"LastName"`
	Active			bool		`json:"Active"`
	LoginAttemptCount	int		`json:"LoginAttemptCount"`
}

func NewLinkAgent(port, username, password string) *LinkAgent {
    return &LinkAgent{
        linkStorePort:  fmt.Sprintf(":%s", port),
        userName:       username,
        password:       password,
    }
}

func (this *LinkAgent) authenticate() error {
    loginBody := &LoginInfo{
        UserName: this.userName,
        Password: this.password,
    }    

    bodyString, err := json.Marshal(loginBody)
    if err != nil {
        return err
    }
    
    resp, err := http.Post(fmt.Sprintf("http://127.0.0.1%s/api/User/Login", this.linkStorePort), "application/json", bytes.NewBuffer(bodyString))
    if err != nil {
        return err
    }

    result, err := io.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    
    if resp.StatusCode != 200 {
        return fmt.Errorf("Status code is (%v) when it should be (200 OK) Body: %v", resp.Status, string(result))
    }

    var token LoginResult
    err = json.Unmarshal(result, &token)
    if err != nil {
        return err
    }

    this.secretToken = token.Token

    return nil
}

func (this *LinkAgent) checkAuthenticationStatus() error {
    fullUrl := fmt.Sprintf("http://127.0.0.1%s/api/Ping", this.linkStorePort)

    request, err := http.NewRequest("GET", fullUrl, nil)
    if err != nil {
        return err
    }

    request.Header.Add("Content-Type", "application/json")
    request.Header.Add("X-Session-Token", this.secretToken)

    client := &http.Client{}
    response, err := client.Do(request)
    if err != nil {
        return err
    }

    if response.StatusCode == 401 {
        err = this.authenticate()
        if err != nil {
            return err
        }
        return nil
    } 
    
    stringReturn, err := io.ReadAll(response.Body)
    if err != nil {
        return err
    }

    if response.StatusCode != 200 {
        return fmt.Errorf("Status code is (%v) when it should be (200 OK). Body: %v", response.Status, string(stringReturn))
    }

    return nil
}

func (this *LinkAgent) makeRequest(method string, url string, body []byte) ([]byte, error) {
    err := this.checkAuthenticationStatus()
    if err != nil {
        return nil, err
    }

    fullUrl := fmt.Sprintf("http://127.0.0.1%s/api/%s", this.linkStorePort, url)

    request, err := http.NewRequest(method, fullUrl, bytes.NewBuffer(body))
    if err != nil {
        return nil, err
    }

    request.Header.Add("Content-Type", "application/json")
    request.Header.Add("X-Session-Token", this.secretToken)

    client := &http.Client{}
    response, err := client.Do(request)
    if err != nil {
        return nil, err
    }

    stringReturn, err := io.ReadAll(response.Body)
    if err != nil {
        return nil, err
    }

    if response.StatusCode != 200 {
        return nil, fmt.Errorf("Status code is (%v) when it should be (200 OK). Body: %v", response.Status, string(stringReturn))
    }

    return stringReturn, nil
}


func (this *LinkAgent) GetAllUsersCount() (int, error) {
    bodyContent := &UserSearch{
        UserName: "",
        FirstName: "",
        LastName: "",
        Active: "active",
    }

    userContent, err := jsgn.Marshal(bodyContent)
    if err != nil {
        return 0, err
    }

    users, err := this.makeRequest("POST", "User/_search", userContent)
    if err != nil {
        return 0, err
    }
   
    var userList []User

    err = json.Unmarshal(users, &userList)
    if err != nil {
        return 0, err
    }

    return len(userList), nil
}
