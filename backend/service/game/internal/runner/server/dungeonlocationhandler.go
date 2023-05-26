package runner

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	coreerror "gitlab.com/alienspaces/go-mud/backend/core/error"
	"gitlab.com/alienspaces/go-mud/backend/core/jsonschema"
	"gitlab.com/alienspaces/go-mud/backend/core/queryparam"
	"gitlab.com/alienspaces/go-mud/backend/core/server"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/modeller"
	schema "gitlab.com/alienspaces/go-mud/backend/schema/game"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/model"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

const (
	getDungeonLocations string = "get-dungeon-locations"
	getDungeonLocation  string = "get-dungeon-location"
)

func (rnr *Runner) DungeonLocationHandlerConfig(hc map[string]server.HandlerConfig) map[string]server.HandlerConfig {

	return mergeHandlerConfigs(hc, map[string]server.HandlerConfig{
		getDungeonLocations: {
			Method:      http.MethodGet,
			Path:        "/api/v1/dungeons/:dungeon_id/locations",
			HandlerFunc: rnr.GetDungeonLocationsHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthenTypes: []server.AuthenticationType{
					server.AuthenticationTypePublic,
				},
				ValidateRequestSchema: jsonschema.SchemaWithReferences{
					Main: jsonschema.Schema{
						Location: "schema/game/location",
						Name:     "response.schema.json",
					},
					References: []jsonschema.Schema{
						{
							Location: "schema/game/location",
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
					server.AuthenticationTypePublic,
				},
				ValidateResponseSchema: jsonschema.SchemaWithReferences{
					Main: jsonschema.Schema{
						Location: "schema/game/location",
						Name:     "response.schema.json",
					},
					References: []jsonschema.Schema{
						{
							Location: "schema/game/location",
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
func (rnr *Runner) GetDungeonLocationHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m modeller.Modeller) error {

	var recs []*record.Location
	var err error

	// Path parameters
	dungeonID := pp.ByName("dungeon_id")
	locationID := pp.ByName("location_id")

	l.Info("Getting dungeon ID >%s< location ID >%s<", dungeonID, locationID)

	dungeonRec, err := m.(*model.Model).GetDungeonRec(dungeonID, nil)
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

	locationRec, err := m.(*model.Model).GetLocationRec(locationID, nil)
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

	err = server.WriteResponse(l, w, http.StatusOK, res)
	if err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// GetDungeonLocationsHandler -
func (rnr *Runner) GetDungeonLocationsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m modeller.Modeller) error {
	l.Info("** Get locations handler **")

	// Path parameters
	dungeonID := pp.ByName("dungeon_id")

	// Query parameters
	opts := queryparam.ToSQLOptions(qp)

	l.Info("Querying dungeon ID >%s< location records with opts >%#v<", dungeonID, opts)

	dungeonRec, err := m.(*model.Model).GetDungeonRec(dungeonID, nil)
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

	locationRecs, err := m.(*model.Model).GetLocationRecs(opts)
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

	err = server.WriteResponse(l, w, http.StatusOK, res)
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
