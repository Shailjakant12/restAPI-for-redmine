package generator

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type Issue struct {
	Issue IssueDetails `json:"issue,omitempty"`
}

type IssueDetails struct {
	ProjectID    string      `json:"project_id,omitempty"`
	Subject      string      `json:"subject,omitempty"`
	CustomFields []FieldData `json:"custom_fields"`
	Description  string      `json:"description,omitempty"`
}

type FieldData struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Value []string `json:"value"`
}

type IncidentGenerator struct {
	APIClient  APIClientInterface
	FieldsData FieldData
}

type APIClientInterface interface {
	Publish(data []byte) (*http.Response, error)
}

//This function is use to receive the data from url and assign some manual values to other fields
func (gen *IncidentGenerator) createData(res http.ResponseWriter, req *http.Request) Issue {
	var IssueData IssueDetails
	IssueData.Description = "MonitorId : " + strings.Replace(req.FormValue("monitorID"), "*", "", -1) +
		"\nMonitorURL : " + strings.Replace(req.FormValue("monitorURL"), "*", "", -1) +
		"\nmonitorFriendlyName : " + strings.Replace(req.FormValue("monitorFriendlyName"), "*", "", -1) +
		"\nAlertType : " + strings.Replace(req.FormValue("alertType"), "*", "", -1) +
		"\nAlertTypeFriendlyName : " + strings.Replace(req.FormValue("alertTypeFriendlyName"), "*", "", -1) +
		"\nAlertDetails : " + strings.Replace(req.FormValue("alertDetails"), "*", "", -1) +
		"\nAlertDuration : " + strings.Replace(req.FormValue("alertDuration"), "*", "", -1) +
		"\nMonitorAlertContacts : " + strings.Replace(req.FormValue("monitorAlertContacts"), "*", "", -1)
	IssueData.Subject = strings.Replace(req.FormValue("alertDetails"), "*", "", -1)
	IssueData.ProjectID = "5"
	IssueData.CustomFields = make([]FieldData, 1)
	IssueData.CustomFields[0] = gen.FieldsData
	ParticularIssue := Issue{
		IssueData,
	}
	return ParticularIssue
}

func (gen *IncidentGenerator) CreateIncident(w http.ResponseWriter, req *http.Request) {
	data := gen.createData(w, req)
	//fmt.Println("created data ", data)
	response := gen.sendData(data)
	responseData, _ := ioutil.ReadAll(response.Body)
	//	fmt.Println("responseData ", string(responseData))
	w.Write(responseData)
}

func (gen *IncidentGenerator) sendData(issue Issue) *http.Response {
	issueEncoded, _ := json.Marshal(issue)
	response, _ := gen.APIClient.Publish(issueEncoded)
	return response
}
