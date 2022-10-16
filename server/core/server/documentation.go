package server

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gitlab.com/alienspaces/go-mud/server/core/jsonschema"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// GenerateHandlerDocumentation - generates documentation based on handler configuration
func (rnr *Runner) GenerateHandlerDocumentation(messageConfigs []MessageConfig, handlerConfigs []HandlerConfig) ([]byte, error) {

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
	h3 {
		color: DarkSlateGray;
	}
	th {
		text-align: left;
	}
	th, td {
		padding: 2px;
	}
	.header {
		padding-top: 10px;
		padding-bottom: 2px;
	}
	.path-method {
		color: #629153;
	}
	.summary {
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
	.params-value {
		display: inline-block;
	}
	.params-toggle-visibility {
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
		padding-right: 10px;
	}
	.schema-toggle-visibility {
	}
	.schema-data {
		background-color: #ffffff;
		padding: 10px;
		margin-left: 20px;
		margin-right: 20px;
	}
	.toggle-visibility {
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

	var hide = function (elem) {
		elem.style.display = 'none';
	};

	var toggleVisibility = function (elem) {

		// If the element is visible, hide it
		if (window.getComputedStyle(elem).display === 'block') {
			hide(elem);
			return;
		}

		show(elem);
	};

	document.addEventListener('click', function (event) {

		// Make sure clicked element is our toggle
		if (!event.target.classList.contains('toggle-visibility')) return;

		event.preventDefault();

		var content = document.querySelector(event.target.hash);
		if (!content) return;

		toggleVisibility(content);

	}, false);
</script>
</head>
<body>
	`)

	fmt.Fprintf(&b, "<div class='header'><h2>Schema Documentation</h2></div>")

	rnr.appendBuildInfo(&b)

	rnr.appendAPIDocumentationURL(&b)

	if len(messageConfigs) > 0 {
		fmt.Fprintf(&b, "<h3 class='header'>Messages</h3>")
		for count, cfg := range messageConfigs {
			appendMessageConfig(&b, cfg)

			schemaMain, schemaData, err := rnr.loadSchemaWithReferences(cfg.ValidateSchema)
			if err != nil {
				return nil, err
			}

			appendSchemaWithReferences(&b, count, schemaMain, schemaData, "Body Schema")
		}
	}

	if len(handlerConfigs) > 0 {
		fmt.Fprintf(&b, "<h3 class='header'>API</h3>")
		for count, config := range sortHandlerConfigs(handlerConfigs) {

			if !config.DocumentationConfig.Document {
				// skip documenting this endpoint
				continue
			}

			appendSummary(&b, config)

			querySchemaMain, querySchemaReferences, err := rnr.loadSchemaWithReferences(config.MiddlewareConfig.ValidateQueryParams)
			if err != nil {
				return nil, err
			}
			appendSchemaWithReferences(&b, count, querySchemaMain, querySchemaReferences, "Query Params Schema")

			requestSchemaMain, requestSchemaReferences, err := rnr.loadSchemaWithReferences(config.MiddlewareConfig.ValidateRequestSchema)
			if err != nil {
				return nil, err
			}
			appendSchemaWithReferences(&b, count, requestSchemaMain, requestSchemaReferences, "Request Body Schema")

			responseSchemaMain, responseSchemaReferences, err := rnr.loadSchemaWithReferences(config.MiddlewareConfig.ValidateResponseSchema)
			if err != nil {
				return nil, err
			}
			appendSchemaWithReferences(&b, count, responseSchemaMain, responseSchemaReferences, "Response Body Schema")
		}
	}

	fmt.Fprintf(&b, "<div class='footer'></div>")
	fmt.Fprintf(&b, `
	</body>
		`)

	return []byte(b.String()), nil
}

func (rnr *Runner) appendAPIDocumentationURL(b *strings.Builder) {
	if appHost := rnr.Config.Get("APP_HOST"); appHost != "" {
		fmt.Fprintf(b, "<div class='header'><a href='%s/'>View API Documentation</a></div>", appHost)
	}
}

func appendMessageConfig(b *strings.Builder, cfg MessageConfig) {
	fmt.Fprintf(b, "<h4 id='%s'>%s %s</h4>", strings.ToLower(string(cfg.Name)), cases.Title(language.Und).String(string(cfg.Subject)), cases.Title(language.Und).String(string(cfg.Event)))
	fmt.Fprintf(b, "<div class='params'>\n")
	fmt.Fprintf(b, "<div class='params-label'>Topic - </div><div class='params-value'>%s</div>", cfg.Topic)
	fmt.Fprintf(b, "</div>\n")
	fmt.Fprintf(b, "<div class='params'>\n")
	fmt.Fprintf(b, "<div class='params-label'>Subject - </div><div class='params-value'>%s</div>", cfg.Subject)
	fmt.Fprintf(b, "</div>\n")
	fmt.Fprintf(b, "<div class='params'>\n")
	fmt.Fprintf(b, "<div class='params-label'>Event - </div><div class='params-value'>%s</div>", cfg.Event)
	fmt.Fprintf(b, "</div>\n")
}

func appendSummary(b *strings.Builder, config HandlerConfig) {
	fmt.Fprintf(b, "<div id='%s' class='path'><h4><span class='path-method'>%s</span> - <span class='path=url'>%s</span></h4></div>", strings.ToLower(config.Name), config.Method, config.Path)
	if config.DocumentationConfig.Summary != "" {
		fmt.Fprintf(b, "<div class='summary'>%s</div>", config.DocumentationConfig.Summary)
	}
}

func (rnr *Runner) loadSchemaWithReferences(s jsonschema.SchemaWithReferences) (mainSchema []byte, referenceSchemas [][]byte, err error) {
	if s.IsEmpty() {
		return mainSchema, referenceSchemas, nil
	}

	mainSchemaPath := s.Main.GetFileName()

	rnr.Log.Debug("schema main content path >%s<", mainSchemaPath)

	mainSchema, err = ioutil.ReadFile(mainSchemaPath)
	if err != nil {
		return mainSchema, referenceSchemas, err
	}

	for _, schemaReference := range s.References {

		path := schemaReference.GetFileName()

		rnr.Log.Debug("schema reference content path >%s<", path)

		ds, err := ioutil.ReadFile(path)
		if err != nil {
			return mainSchema, referenceSchemas, err
		}
		referenceSchemas = append(referenceSchemas, ds)
	}

	return mainSchema, referenceSchemas, nil
}

func appendSchemaWithReferences(b *strings.Builder, count int, schemaMainContent []byte, schemaReferenceContents [][]byte, schemaLabel string) {
	if len(schemaMainContent) == 0 {
		return
	}

	schemaLabelID := strings.Join(strings.Split(strings.ToLower(schemaLabel), " "), "-")

	// fmt.Fprintf(b, "<div class='schema'>\n")
	// fmt.Fprintf(b, "<div class='schema-label'>%s</div>\n", schemaLabel)
	// fmt.Fprintf(b, "<div class='schema-toggle-visibility'><a href='#schema-%s-%d' class='toggle-visibility'>show / hide</a></div>", schemaLabelID, count)
	// fmt.Fprintf(b, "</span>\n</div>\n")
	fmt.Fprintf(b, "<div class='schema'>\n")
	fmt.Fprintf(b, "<span class='schema-label'>%s</span>\n", schemaLabel)
	fmt.Fprintf(b, "<span class='schema-toggle-visibility'><a href='#schema-%s-%d' class='toggle-visibility'>show / hide</a></span>", schemaLabelID, count)
	fmt.Fprintf(b, "</div>\n")
	fmt.Fprintf(b, "<div id='schema-%s-%d' style='display: none'>\n", schemaLabelID, count)
	if len(schemaMainContent) > 0 {
		fmt.Fprintf(b, "<pre class='schema-data'>%s</pre>\n", string(schemaMainContent))
	}
	for _, s := range schemaReferenceContents {
		fmt.Fprintf(b, "<pre class='schema-data'>%s</pre>\n", string(s))
	}
	fmt.Fprintf(b, "</div>\n")
}

func (rnr *Runner) appendBuildInfo(b *strings.Builder) {
	appImageBranch := rnr.Config.Get("APP_IMAGE_TAG_FEATURE_BRANCH")
	appImageSHA := rnr.Config.Get("APP_IMAGE_TAG_SHA")

	if appImageBranch != "" && appImageSHA != "" {
		fmt.Fprintf(b, "<div class='header'><h4>Branch: %s</h4><h4>SHA: %s</h4></div>", appImageBranch, appImageSHA)
	} else if appImageBranch != "" {
		fmt.Fprintf(b, "<div class='header'><h4>Branch: %s</h4></div>", appImageBranch)
	} else if appImageSHA != "" {
		fmt.Fprintf(b, "<div class='header'><h4>SHA: %s</h4></div>", appImageSHA)
	}
}
