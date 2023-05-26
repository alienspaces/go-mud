package salesforce

import (
	"fmt"
	"net/url"

	"github.com/nimajalali/go-force/force"
)

// The following functions are all wrappers around the internal Service Cloud
// package functions. These essentially exist to provide some level of
// abstraction from the underlying library that is used.

// Query -
func (scc *Client) Query(queryString string, results interface{}) error {
	l := scc.Log.WithFunctionContext("Query")

	if scc.Client == nil {
		msg := "failed, client is not connected"
		l.Warn(msg)
		return fmt.Errorf(msg)
	}

	err := scc.Client.Query(queryString, &results)
	if err != nil {
		l.Warn("failed Query >%v<", err)
		return err
	}

	return nil
}

// Get -
func (scc *Client) Get(path string, params url.Values, results interface{}) error {
	l := scc.Log.WithFunctionContext("Get")

	if scc.Client == nil {
		msg := "failed, client is not connected"
		l.Warn(msg)
		return fmt.Errorf(msg)
	}

	err := scc.Client.Get(path, params, &results)
	if err != nil {
		l.Warn("failed Get >%v<", err)
		return err
	}

	return nil
}

// Patch -
func (scc *Client) Patch(path string, params url.Values, payload, results interface{}) error {
	l := scc.Log.WithFunctionContext("Patch")

	if scc.Client == nil {
		msg := "failed, client is not connected"
		l.Warn(msg)
		return fmt.Errorf(msg)
	}

	err := scc.Client.Patch(path, params, payload, &results)
	if err != nil {
		l.Warn("failed Patch >%v<", err)
		return err
	}

	return nil
}

// Delete -
func (scc *Client) Delete(path string, params url.Values) error {
	l := scc.Log.WithFunctionContext("Delete")

	if scc.Client == nil {
		msg := "failed, client is not connected"
		l.Warn(msg)
		return fmt.Errorf(msg)
	}

	err := scc.Client.Delete(path, params)
	if err != nil {
		l.Warn("failed Delete >%v<", err)
		return err
	}

	return nil
}

// The following functions are less abstracted as the above as they expose elements
// of the underlying library that is used. Use the above functions when possible.

// GetSObject -
func (scc *Client) GetSObject(id string, fields []string, results force.SObject) error {
	l := scc.Log.WithFunctionContext("GetSObject")

	if scc.Client == nil {
		msg := "failed, client is not connected"
		l.Warn(msg)
		return fmt.Errorf(msg)
	}

	err := scc.Client.GetSObject(id, fields, results)
	if err != nil {
		l.Warn("failed GetSObject >%v<", err)
		return err
	}

	return nil
}

// GetSObject -
func (scc *Client) GetSObjectByExternalID(id string, fields []string, results force.SObject) error {
	l := scc.Log.WithFunctionContext("GetSObjectByExternalID")

	if scc.Client == nil {
		msg := "failed, client is not connected"
		l.Warn(msg)
		return fmt.Errorf(msg)
	}

	err := scc.Client.GetSObjectByExternalId(id, fields, results)
	if err != nil {
		l.Warn("failed GetSObject >%v<", err)
		return err
	}

	return nil
}

// InsertSObject -
func (scc *Client) InsertSObject(rec force.SObject) (*force.SObjectResponse, error) {
	l := scc.Log.WithFunctionContext("InsertSObject")

	if scc.Client == nil {
		msg := "failed, client is not connected"
		l.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	response, err := scc.Client.InsertSObject(rec)
	if err != nil {
		l.Warn("failed InsertSObject >%v<", err)
		return nil, err
	}

	return response, nil
}

func (scc *Client) UpsertByExternalID(id string, in force.SObject) (*force.SObjectResponse, error) {
	l := scc.Log.WithFunctionContext("UpsertByExternalID")

	if scc.Client == nil {
		msg := "failed, client is not connected"
		l.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	response, err := scc.Client.UpsertSObjectByExternalId(id, in)
	if err != nil {
		l.Warn("failed UpsertSObjectByExternalId >%v<", err)
		return nil, err
	}

	return response, nil
}

func (scc *Client) DeleteByExternalID(id string, in force.SObject) error {
	l := scc.Log.WithFunctionContext("DeleteByExternalID")

	if scc.Client == nil {
		msg := "failed, client is not connected"
		l.Warn(msg)
		return fmt.Errorf(msg)
	}

	err := scc.Client.DeleteSObjectByExternalId(id, in)
	if err != nil {
		l.Warn("failed DeleteSObject >%v<", err)
		return err
	}

	return nil
}

func (scc *Client) DeleteSObject(id string, in force.SObject) error {
	l := scc.Log.WithFunctionContext("DeleteSObject")

	if scc.Client == nil {
		msg := "failed, client is not connected"
		l.Warn(msg)
		return fmt.Errorf(msg)
	}

	err := scc.Client.DeleteSObject(id, in)
	if err != nil {
		l.Warn("failed DeleteSObject >%v<", err)
		return err
	}

	return nil
}
