package docs

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func TransformApiDocumentation(resp *http.Response) error {
	b, err := ioutil.ReadAll(resp.Body) //Read json response
	if err != nil {
		return err
	}

	err = resp.Body.Close()
	if err != nil {
		return err
	}

	var schema map[string]interface{}
	err = json.Unmarshal(b, &schema)
	if err != nil {
		return err
	}

	thirdParty := schema["definitions"].(map[string]interface{})["ThirdParty"].(map[string]interface{})
	tpResidualRiskProperties := thirdParty["properties"].(map[string]interface{})["residual_risk"].(map[string]interface{})["items"].(map[string]interface{})["properties"].(map[string]interface{})
	residualRiskProperties := schema["definitions"].(map[string]interface{})["ThirdPartyResidualRisk"].(map[string]interface{})["properties"].(map[string]interface{})
	detaildScoresProperties := schema["definitions"].(map[string]interface{})["ThirdPartyDetailedScores"].(map[string]interface{})["properties"].(map[string]interface{})
	thirdParty["properties"].(map[string]interface{})["residual_risk"] = map[string]interface{}{
		"type":        "object",
		"description": "The residual risk for a third party at a specific point in time",
		"properties": map[string]interface{}{
			"report_type": tpResidualRiskProperties["report_type"],
			"date":        residualRiskProperties["date"],
			"tier":        residualRiskProperties["tier"],
			"id":          residualRiskProperties["id"],
			"findings":    residualRiskProperties["high_risks"],
			"scores":      detaildScoresProperties["group_scores"],
		},
		"required": []string{
			"report_type",
			"date",
			"tier",
			"id",
		},
	}

	getThirdParties := schema["paths"].(map[string]interface{})["/v1/third-parties"].(map[string]interface{})["get"].(map[string]interface{})

	cleanResponse := getThirdParties["responses"].(map[string]interface{})["200"].(map[string]interface{})
	cleanResponse["description"] = "Successful response, a list of all your third parties."
	cleanResponse["schema"] = map[string]interface{}{
		"type": "array",
		"items": map[string]interface{}{
			"$ref": "#/definitions/ThirdParty",
		},
	}

	getThirdParties["responses"] = map[string]interface{}{
		"200": cleanResponse,
	}

	schemes := schema["schemes"]
	if strings.TrimSpace(os.Getenv("GIN_MODE")) != "release" {
		schemes = []string{"https", "http"}
	}

	updatedSchema := map[string]interface{}{
		"basePath": "/",
		"definitions": map[string]interface{}{
			"ThirdParty": thirdParty,
		},
		"info":    schema["info"],
		"schemes": schemes,
		"swagger": schema["swagger"],
		"paths": map[string]interface{}{
			"/v1/third-parties": map[string]interface{}{
				"get": map[string]interface{}{
					"description": getThirdParties["description"],
					"parameters":  []string{},
					"produces":    getThirdParties["produces"],
					"responses":   getThirdParties["responses"],
					"security":    getThirdParties["security"],
					"tags":        getThirdParties["tags"],
				},
			},
		},
		"securityDefinitions": map[string]interface{}{
			"ApiV1AuthToken": schema["securityDefinitions"].(map[string]interface{})["ApiV1AuthToken"],
		},
	}

	encodedSchema, err := json.Marshal(updatedSchema)
	if err != nil {
		return err
	}

	body := ioutil.NopCloser(bytes.NewReader(encodedSchema))
	resp.Body = body
	resp.ContentLength = int64(len(encodedSchema))
	resp.Header.Set("Content-Length", strconv.Itoa(len(encodedSchema)))
	return nil
}
