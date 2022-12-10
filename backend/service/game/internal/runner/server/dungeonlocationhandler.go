package runner

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	coreerror "gitlab.com/alienspaces/go-mud/server/core/error"
	"gitlab.com/alienspaces/go-mud/server/core/jsonschema"
	"gitlab.com/alienspaces/go-mud/server/core/server"
	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
	"gitlab.com/alienspaces/go-mud/server/core/type/modeller"
	"gitlab.com/alienspaces/go-mud/server/schema"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/model"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

const (
	getDungeonLocations server.HandlerConfigKey = "get-dungeon-locations"
	getDungeonLocation  server.HandlerConfigKey = "get-dungeon-location"
)

func (rnr *Runner) DungeonLocationHandlerConfig(hc map[server.HandlerConfigKey]server.HandlerConfig) map[server.HandlerConfigKey]server.HandlerConfig {

	return mergeHandlerConfigs(hc, map[server.HandlerConfigKey]server.HandlerConfig{
		getDungeonLocations: {
			Method:      http.MethodGet,
			Path:        "/api/v1/dungeons/:dungeon_id/locations",
			HandlerFunc: rnr.GetDungeonLocationsHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthenTypes: []server.AuthenticationType{
					server.AuthenTypePublic,
				},
				ValidateRequestSchema: jsonschema.SchemaWithReferences{
					Main: jsonschema.Schema{
						Location: "schema/docs/location",
						Name:     "response.schema.json",
					},
					References: []jsonschema.Schema{
						{
							Location: "schema/docs/location",
							Name:     "data.schema.json",
						},
					},
				},
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "List dungeon locations.",
			},
		},
		getDungeonLocation: {
			Method:      http.MethodGet,
			Path:        "/api/v1/dungeons/:dungeon_id/locations/:location_id",
			HandlerFunc: rnr.GetDungeonLocationHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthenTypes: []server.AuthenticationType{
					server.AuthenTypePublic,
				},
				ValidateResponseSchema: jsonschema.SchemaWithReferences{
					Main: jsonschema.Schema{
						Location: "schema/docs/location",
						Name:     "response.schema.json",
					},
					References: []jsonschema.Schema{
						{
							Location: "schema/docs/location",
							Name:     "data.schema.json",
						},
					},
				},
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Get a dungeon location.",
			},
		},
	})
}

// GetDungeonLocationHandler -
func (rnr *Runner) GetDungeonLocationHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) error {

	var recs []*record.Location
	var err error

	// Path parameters
	dungeonID := pp.ByName("dungeon_id")
	if dungeonID == "" {
		err := coreerror.NewNotFoundError("dungeon", dungeonID)
		server.WriteError(l, w, err)
		return err
	} else if !m.(*model.Model).IsUUID(dungeonID) {
		err := coreerror.NewPathParamInvalidTypeError("dungeon_id", dungeonID)
		server.WriteError(l, w, err)
		return err

	}

	locationID := pp.ByName("location_id")
	if locationID == "" {
		err := coreerror.NewNotFoundError("location", locationID)
		server.WriteError(l, w, err)
		return err
	} else if !m.(*model.Model).IsUUID(locationID) {
		err := coreerror.NewPathParamInvalidTypeError("location_id", locationID)
		server.WriteError(l, w, err)
		return err

	}

	l.Info("Getting dungeon ID >%s< location ID >%s<", dungeonID, locationID)

	dungeonRec, err := m.(*model.Model).GetDungeonRec(dungeonID, false)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	// Resource not found
	if dungeonRec == nil {
		err := coreerror.NewNotFoundError("dungeon", dungeonID)
		server.WriteError(l, w, err)
		return err
	}

	locationRec, err := m.(*model.Model).GetLocationRec(locationID, false)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	// Resource not found
	if locationRec == nil {
		err := coreerror.NewNotFoundError("location", locationID)
		server.WriteError(l, w, err)
		return err
	}

	if locationRec.DungeonID != dungeonRec.ID {
		err := coreerror.NewNotFoundError("location", locationID)
		server.WriteError(l, w, err)
		return err
	}

	recs = append(recs, locationRec)

	// Assign response properties
	data := []schema.LocationData{}
	for _, rec := range recs {

		// Response data
		responseData, err := rnr.RecordToLocationResponseData(*rec)
		if err != nil {
			server.WriteError(l, w, err)
			return err
		}

		data = append(data, responseData)
	}

	res := schema.LocationResponse{
		Data: data,
	}

	err = server.WriteResponse(l, w, res)
	if err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// GetDungeonLocationsHandler -
func (rnr *Runner) GetDungeonLocationsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) error {
	l.Info("** Get locations handler **")

	// Path parameters
	dungeonID := pp.ByName("dungeon_id")
	if dungeonID == "" {
		err := coreerror.NewNotFoundError("dungeon", dungeonID)
		server.WriteError(l, w, err)
		return err
	} else if !m.(*model.Model).IsUUID(dungeonID) {
		err := coreerror.NewPathParamInvalidTypeError("dungeon_id", dungeonID)
		server.WriteError(l, w, err)
		return err
	}

	l.Info("Getting dungeon ID >%s<", dungeonID)

	dungeonRec, err := m.(*model.Model).GetDungeonRec(dungeonID, false)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	// Resource not found
	if dungeonRec == nil {
		err := coreerror.NewNotFoundError("dungeon", dungeonID)
		server.WriteError(l, w, err)
		return err
	}

	params := make(map[string]interface{})
	params["dungeon_id"] = dungeonID

	for paramName, paramValue := range qp {
		l.Info("Querying location records with param name >%s< value >%v<", paramName, paramValue)
		params[paramName] = paramValue
	}

	l.Info("Querying dungeon location records with params >%#v<", params)

	locationRecs, err := m.(*model.Model).GetLocationRecs(
		map[string]interface{}{
			"dungeon_id": dungeonID,
		},
		nil, false,
	)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	// Assign response properties
	data := []schema.LocationData{}
	for _, locationRecs := range locationRecs {

		// Response data
		responseData, err := rnr.RecordToLocationResponseData(*locationRecs)
		if err != nil {
			server.WriteError(l, w, err)
			return err
		}

		data = append(data, responseData)
	}

	res := schema.LocationResponse{
		Data: data,
	}

	l.Info("Responding with >%#v<", res)

	err = server.WriteResponse(l, w, res)
	if err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// LocationRequestDataToRecord -
func (rnr *Runner) LocationRequestDataToRecord(data schema.LocationData, rec *record.Location) error {

	return nil
}

// RecordToLocationResponseData -
func (rnr *Runner) RecordToLocationResponseData(locationRec record.Location) (schema.LocationData, error) {

	data := schema.LocationData{
		LocationID:          locationRec.ID,
		LocationName:        locationRec.Name,
		LocationDescription: locationRec.Description,
		LocationDefault:     locationRec.IsDefault,
		LocationCreatedAt:   locationRec.CreatedAt,
		LocationUpdatedAt:   locationRec.UpdatedAt.Time,
	}

	return data, nil
}
