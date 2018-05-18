package iapi

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func createKeyValuePairs(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=%s,", key, value)
	}
	return b.String()
}

// GetHost ...
func (server *Server) GetHost(hostname string) ([]HostStruct, error) {

	var hosts []HostStruct

	results, err := server.NewAPIRequest("GET", "/objects/hosts/"+hostname, nil)
	if err != nil {
		return nil, err
	}

	// Contents of the results is an interface object. Need to convert it to json first.
	jsonStr, marshalErr := json.Marshal(results.Results)
	if marshalErr != nil {
		return nil, marshalErr
	}

	// then the JSON can be pushed into the appropriate struct.
	// Note : Results is a slice so much push into a slice.

	if unmarshalErr := json.Unmarshal(jsonStr, &hosts); unmarshalErr != nil {
		return nil, unmarshalErr
	}

	return hosts, err

}

// CreateHost ...
func (server *Server) CreateHost(hostname, address, checkCommand string, variables map[string]string, templates []string) ([]HostStruct, error) {

	var newAttrs HostAttrs
	newAttrs.Address = address
	newAttrs.CheckCommand = "hostalive"
	if variables != nil {
		newAttrs.Vars = Flatten(variables)
		//newAttrs.Vars = variables
	}

	//newAttrs.Vars = variables
	//newAttrs.Vars = createKeyValuePairs(variables)

	nhatr, _ := json.Marshal(newAttrs)
	m := make(map[string]interface{})

	errd := json.Unmarshal(nhatr, &m)
	if errd != nil {

	}

	var newHost HostStruct
	newHost.Name = hostname
	newHost.Type = "Host"
	newHost.Templates = templates
	newHost.Attrs = Flatten(m)
	//newAttrs
	//newHost.Attrs = newAttrs

	// Create JSON from completed struct
	payloadJSON, marshalErr := json.Marshal(newHost)
	if marshalErr != nil {
		return nil, marshalErr
	}

	//fmt.Printf("<payload> %s\n", payloadJSON)

	// Make the API request to create the hosts.
	results, err := server.NewAPIRequest("PUT", "/objects/hosts/"+hostname, []byte(payloadJSON))
	if err != nil {
		return nil, err
	}

	if results.Code == 200 {
		hosts, err := server.GetHost(hostname)
		return hosts, err
	}

	return nil, fmt.Errorf("%s", results.ErrorString)

}

// DeleteHost ...
func (server *Server) DeleteHost(hostname string) error {

	results, err := server.NewAPIRequest("DELETE", "/objects/hosts/"+hostname+"?cascade=1", nil)
	if err != nil {
		return err
	}

	if results.Code == 200 {
		return nil
	} else {
		return fmt.Errorf("%s", results.ErrorString)
	}

}