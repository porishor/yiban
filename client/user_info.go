package client

type UserInfo struct {
	UniversityName string        `json:"UniversityName"`
	UniversityID   string        `json:"UniversityId"`  
	PersonID       string        `json:"PersonId"`      
	PersonName     string        `json:"PersonName"`    
	State          int64         `json:"State"`         
	PersonType     string        `json:"PersonType"`    
	UniversityIcon string        `json:"UniversityIcon"`
	Container      string        `json:"Container"`     
	HomeTheme      string        `json:"HomeTheme"`     
	CustomApps     CustomApps    `json:"CustomApps"`    
	Group          []interface{} `json:"Group"`         
	Apps           []App         `json:"Apps"`          
	WxState        int64         `json:"WxState"`       
	DingDingState  int64         `json:"DingDingState"` 
}

type App struct {
	ID         string `json:"Id"`        
	ServiceID  string `json:"ServiceId"` 
	AppIcon    string `json:"AppIcon"`   
	AppURL     string `json:"AppUrl"`    
	AppRuleURL string `json:"AppRuleUrl"`
	AuthCode   string `json:"AuthCode"`  
	AppName    string `json:"AppName"`   
}

type CustomApps struct {
	Up   []interface{} `json:"up"`  
	Down []Down        `json:"down"`
}

type Down struct {
	ID           string      `json:"Id"`          
	UniversityID string      `json:"UniversityId"`
	AppName      string      `json:"AppName"`     
	AppIcon      string      `json:"AppIcon"`     
	AppURL       string      `json:"AppUrl"`      
	Position     int64       `json:"Position"`    
	State        int64       `json:"State"`       
	Sort         int64       `json:"Sort"`        
	Remark       interface{} `json:"Remark"`      
	CreateTime   int64       `json:"CreateTime"`  
	UpdateTime   int64       `json:"UpdateTime"`  
}


