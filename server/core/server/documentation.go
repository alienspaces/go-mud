package server

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// GenerateHandlerDocumentation - generates documentationbased on handler configuration
func (rnr *Runner) GenerateHandlerDocumentation() ([]byte, error) {

	rnr.Log.Info("** Generate Handler Documentation **")

	var b strings.Builder

	fmt.Fprintf(&b, `
<html>
<head>
<style>
	body {
		background-color: #efefef;
	}
	h2 {
		color: #504ecc;
	}
	.header {
		padding-top: 10px;
		padding-bottom: 2px;
	}
	.path-method {
		color: #629153;
	}
	.description {
		margin-left: 20px;
		margin-right: 20px;
		padding-bottom: 8px;
	}
	.params {
		margin-left: 20px;
		padding-top: 8px;
		padding-bottom: 8px;
	}
	.params-label {
		display: inline-block;
		width: 140px;
	}
	.params-toggle {
		display: inline-block;
		width: 140px;
	}
	.params-list {
		background-color: #ffffff;
		padding: 10px;
		margin-left: 20px;
		margin-right: 20px;
	}
	.param {
		background-color: #efefef;
		border: 1px solid #cdcdcd;
		padding: 2px;
		margin-bottom: 2px;
		display: inline-block;
	}
	.schema {
		margin-left: 20px;
		padding-top: 8px;
		padding-bottom: 8px;
	}
	.schema-label {
		display: inline-block;
		width: 140px;
	}
	.schema-toggle {
		display: inline-block;
		width: 140px;
	}
	.schema-data {
		background-color: #ffffff;
		padding: 10px;
		margin-left: 20px;
		margin-right: 20px;
	}
	.toggle {
		font-style: italic;
		font-size: small;
	}
	.footer {
		padding-top: 50px;
		padding-bottom: 50px;
	}
</style>
<script language="javascript">
	// Show an element
	var show = function (elem) {
		elem.style.display = 'block';
	};

	// Hide an element
	var hide = function (elem) {
		elem.style.display = 'none';
	};

	// Toggle element visibility
	var toggle = function (elem) {

		// If the element is visible, hide it
		if (window.getComputedStyle(elem).display === 'block') {
			hide(elem);
			return;
		}

		// Otherwise, show it
		show(elem);
	};

	// Listen for click events
	document.addEventListener('click', function (event) {

		// Make sure clicked element is our toggle
		if (!event.target.classList.contains('toggle')) return;

		// Prevent default link behaviour
		event.preventDefault();

		// Get the content
		var content = document.querySelector(event.target.hash);
		if (!content) return;

		// Toggle the content
		toggle(content);

	}, false);
</script>
</head>
<body>
	`)

	fmt.Fprintf(&b, "<div class='header'><h2>API Documentation</h2></div>")

	for count, config := range rnr.HandlerConfig {

		// Skip documenting this endpoint
		if !config.DocumentationConfig.Document {
			continue
		}

		rnr.Log.Info("Processing handler documentation path >%s: %s<", config.Method, config.Path)

		var requestSchemaMainContent []byte
		var requestSchemaDataContent []byte
		var responseSchemaMainContent []byte
		var responseSchemaDataContent []byte
		var err error

		schemaPath := rnr.Config.Get("APP_SERVER_SCHEMA_PATH")
		schemaLoc := config.MiddlewareConfig.ValidateSchemaLocation
		if schemaLoc != "" {

			// Request schema
			requestSchemaMain := config.MiddlewareConfig.ValidateSchemaRequestMain
			if requestSchemaMain != "" {

				filename := fmt.Sprintf("%s/%s/%s", schemaPath, schemaLoc, requestSchemaMain)

				rnr.Log.Info("Request schema main content filename >%s<", filename)

				requestSchemaMainContent, err = ioutil.ReadFile(filename)
				if err != nil {
					return nil, err
				}

				schemaReferences := config.MiddlewareConfig.ValidateSchemaRequestReferences
				for _, schemaReference := range schemaReferences {

					filename := fmt.Sprintf("%s/%s/%s", schemaPath, schemaLoc, schemaReference)

					rnr.Log.Info("Request schema reference content filename >%s<", filename)

					requestSchemaDataContent, err = ioutil.ReadFile(filename)
					if err != nil {
						return nil, err
					}
				}
			}

			// Response schema
			responseSchemaMain := config.MiddlewareConfig.ValidateSchemaResponseMain
			if responseSchemaMain != "" {

				filename := fmt.Sprintf("%s/%s/%s", schemaPath, schemaLoc, responseSchemaMain)

				rnr.Log.Info("Response schema main content filename >%s<", filename)

				responseSchemaMainContent, err = ioutil.ReadFile(filename)
				if err != nil {
					return nil, err
				}

				schemaReferences := config.MiddlewareConfig.ValidateSchemaResponseReferences
				for _, schemaReference := range schemaReferences {

					filename := fmt.Sprintf("%s/%s/%s", schemaPath, schemaLoc, schemaReference)

					rnr.Log.Info("Response schema reference content filename >%s<", filename)

					responseSchemaDataContent, err = ioutil.ReadFile(filename)
					if err != nil {
						return nil, err
					}
				}
			}
		}

		// Description
		fmt.Fprintf(&b, "<div class='path'><h4><span class='path-method'>%s</span> - <span class='path=url'>%s</span></h4></div>", config.Method, config.Path)
		if config.DocumentationConfig.Description != "" {
			fmt.Fprintf(&b, "<div class='description'>%s</div>", config.DocumentationConfig.Description)
		}

		// Query parameters
		queryParams := config.MiddlewareConfig.ValidateQueryParams
		if len(queryParams) != 0 {
			fmt.Fprintf(&b, "<div class='params'>\n")
			fmt.Fprintf(&b, "<div class='params-label'>Query Parameters -</div>\n")
			fmt.Fprintf(&b, "<div class='params-toggle'><a href='#params-%d' class='toggle'>show / hide</a></div>", count)
			fmt.Fprintf(&b, "</div>\n")
			fmt.Fprintf(&b, "<div id='params-%d' class='params-list' style='display: none'>\n", count)
			for _, param := range queryParams {
				fmt.Fprintf(&b, "<span class='param'>%s</span>\n", param)
			}
			fmt.Fprintf(&b, "</div>\n")
		}

		// Request schema
		if requestSchemaMainContent != nil || requestSchemaDataContent != nil {
			fmt.Fprintf(&b, "<div class='schema'>\n")
			fmt.Fprintf(&b, "<div class='schema-label'>Request schema -</div>\n")
			fmt.Fprintf(&b, "<div class='schema-toggle'><a href='#request-schema-%d' class='toggle'>show / hide</a></div>", count)
			fmt.Fprintf(&b, "</span>\n</div>\n")
			fmt.Fprintf(&b, "<div id='request-schema-%d' style='display: none'>\n", count)
			if requestSchemaMainContent != nil {
				fmt.Fprintf(&b, "<pre class='schema-data'>%s</pre>\n", string(requestSchemaMainContent))
			}
			if requestSchemaDataContent != nil {
				fmt.Fprintf(&b, "<pre class='schema-data'>%s</pre>\n", string(requestSchemaDataContent))
			}
			fmt.Fprintf(&b, "</div>\n")
		}

		// Response schema
		if responseSchemaMainContent != nil || responseSchemaDataContent != nil {
			fmt.Fprintf(&b, "<div class='schema'>\n")
			fmt.Fprintf(&b, "<div class='schema-label'>Response schema -</div>\n")
			fmt.Fprintf(&b, "<div class='schema-toggle'><a href='#response-schema-%d' class='toggle'>show / hide</a></div>", count)
			fmt.Fprintf(&b, "</span>\n</div>\n")
			fmt.Fprintf(&b, "<div id='response-schema-%d' style='display: none'>\n", count)
			if responseSchemaMainContent != nil {
				fmt.Fprintf(&b, "<pre class='schema-data'>%s</pre>\n", string(responseSchemaMainContent))
			}
			if responseSchemaDataContent != nil {
				fmt.Fprintf(&b, "<pre class='schema-data'>%s</pre>\n", string(responseSchemaDataContent))
			}
			fmt.Fprintf(&b, "</div>\n")
		}

	}

	fmt.Fprintf(&b, "<div class='footer'></div>")
	fmt.Fprintf(&b, `
	</body>
		`)

	return []byte(b.String()), nil
}
