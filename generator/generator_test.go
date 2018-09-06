package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testData AlertDetails

type AlertDetails struct {
	MonitorID             string `json:"monitorID,omitempty"`
	MonitorURL            string `json:"monitorURL,omitempty"`
	MonitorFriendlyName   string `json:"monitorFriendlyName,omitempty"`
	AlertType             string `json:"alertType,omitempty"`
	AlertTypeFriendlyName string `json:"alertTypeFriendlyName,omitempty"`
	AlertDetails          string `json:"alertDetails,omitempty"`
	AlertDuration         string `json:"alertDuration,omitempty"`
	MonitorAlertContacts  string `json:"monitorAlertContacts,omitempty"`
}

var AlertDetailsData string

func TestCreatingIncident(t *testing.T) {
	gen := IncidentGenerator{
		APIClient: &FakeClient{},
	}
	gen.FieldsData = FieldData{}
	testData = AlertDetails{
		MonitorID:             "",
		MonitorURL:            "",
		MonitorFriendlyName:   "http://friendlyurl/",
		AlertType:             "RED",
		AlertTypeFriendlyName: "CRITICAL",
		AlertDetails:          "Critical alert fix it",
		AlertDuration:         "1234",
		MonitorAlertContacts:  "shailu",
	}

	AlertDetailsData = "MonitorId : " + testData.MonitorID +
		"\nMonitorURL : " + testData.MonitorURL +
		"\nmonitorFriendlyName : " + testData.MonitorFriendlyName +
		"\nAlertType : " + testData.AlertType +
		"\nAlertTypeFriendlyName : " + testData.AlertTypeFriendlyName +
		"\nAlertDetails : " + testData.AlertDetails +
		"\nAlertDuration : " + testData.AlertDuration +
		"\nMonitorAlertContacts : " + testData.MonitorAlertContacts

	url := fmt.Sprintf("/incident?monitorID=*%s*&monitorURL=*%s*&monitorFriendlyName=*%s*&alertType=*%s*&alertTypeFriendlyName=*%s*&alertDetails=*%s*&alertDuration=*%s*&monitorAlertContacts=*%s*",
		testData.MonitorID, testData.MonitorURL, testData.MonitorFriendlyName, testData.AlertType, testData.AlertTypeFriendlyName, testData.AlertDetails, testData.AlertDuration, testData.MonitorAlertContacts)
	req, _ := http.NewRequest("POST", url, nil)
	w := httptest.NewRecorder()

	gen.CreateIncident(w, req)
	fmt.Println(w.Body.String())
	// data := w.Body.String()

	// if expectedOutput != data {
	//  	t.Errorf("Expected success but received error")
	// }
}

/////////////////////////////////// Fake Implementation /////////////////////////////

type FakeClient struct {
}

func (fakeClient *FakeClient) Publish(data []byte) (resp *http.Response, err error) {
	issue := Issue{}
	errr := json.Unmarshal(data, &issue)
	if errr != nil {
		fmt.Println("error:", errr)
	}
	response := &http.Response{}
	var responseData string
	validiatedResponseData := validiateIssueDetailsElement(issue, responseData)

	if validiatedResponseData != "" {

		finalData := fmt.Sprintf("\n -\"errors\": [\n%s   ]\n", validiatedResponseData)
		response = &http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString(finalData)),
		}
		response.StatusCode = http.StatusUnprocessableEntity
	} else if issue.Issue.Description == AlertDetailsData {
		response = &http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString("Everything looks fine, data sended is received as such")),
		}
		response.StatusCode = http.StatusOK
	}
	//fmt.Println("response body ", response.Body)

	return response, nil
}

func validiateIssueDetailsElement(issue Issue, responseData string) string {

	if issue.Issue.Subject == "" && issue.Issue.ProjectID == "" {
		responseData = ""
		responseData = fmt.Sprintf("    \"Subject cannot be blank\"\n")
		responseData = responseData + fmt.Sprintf("    \"Project cannot be blank\"\n")
		responseData = responseData + fmt.Sprintf("    \"Tracker cannot be blank\"\n")
		responseData = responseData + fmt.Sprintf("    \"Status cannot be blank\"\n")
	} else if issue.Issue.Subject == "" && issue.Issue.ProjectID != "" && issue.Issue.CustomFields == nil {
		responseData = ""
		responseData = fmt.Sprintf("    \"Component cannot be blank\"\n")
		responseData = responseData + fmt.Sprintf("    \"Subject cannot be blank\"\n")
	} else if issue.Issue.ProjectID == "" && issue.Issue.CustomFields == nil || issue.Issue.ProjectID == "" {
		//	fmt.Println("I am inside proect id")
		responseData = ""
		responseData = fmt.Sprintf("    \"Project cannot be blank\"\n")
		responseData = responseData + fmt.Sprintf("    \"Tracker cannot be blank\"\n")
		responseData = responseData + fmt.Sprintf("    \"Status cannot be blank\"\n")
	} else if issue.Issue.Subject == "" {
		//	fmt.Println("i'm inside suject")
		responseData = ""
		responseData = fmt.Sprintf("    \"Subject cannot be blank\"\n")
	} else if issue.Issue.CustomFields == nil {
		responseData = ""
		responseData = fmt.Sprintf("    \"Component cannot be blank\"\n")
	}

	return responseData
}
