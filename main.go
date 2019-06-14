//go:generate go run -tags=dev assets_generate.go

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"github.com/CyberGRX/api-connector-bulk/assets"
	"github.com/CyberGRX/api-connector-bulk/docs"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

var cyberGrxAPI string

func init() {
	hostname := strings.TrimSpace(os.Getenv("CYBERGRX_API"))
	if hostname == "" {
		hostname = "api.cybergrx.com"
	}

	hostname = strings.TrimRight(hostname, "/")

	if strings.TrimSpace(os.Getenv("GIN_MODE")) == "release" {
		cyberGrxAPI = "https://api.cybergrx.com"
	} else if strings.HasPrefix(hostname, "http://") {
		cyberGrxAPI = hostname
	} else {
		cyberGrxAPI = fmt.Sprintf("https://%s", strings.TrimLeft(hostname, "https://"))
	}

	log.Printf("[INFO] Configuring bulk query support for %s\n", cyberGrxAPI)
}

func getURL(uri string) string {
	return fmt.Sprintf("%s%s", cyberGrxAPI, uri)
}

func get(uri, authorization string) (response map[string]interface{}, err error) {
	url := getURL(uri)
	log.Printf("[INFO] Fetching %s\n", url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", authorization)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		return response, fmt.Errorf("Could not fetch %s status %d", url, resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &response)
	return
}

func allThirdParties(req *http.Request) []map[string]interface{} {
	resp := make([]map[string]interface{}, 0, 100)

	next := "/v1/third-parties?limit=50"
	for next != "" {
		tp, err := get(next, req.Header.Get("Authorization"))
		if err != nil {
			log.Printf("[ERROR] Issue looking up third parties, %s\n", err)
			break
		}

		if items, ok := tp["items"]; ok {
			if _items, ok := items.([]interface{}); ok {
				for _, thirdParty := range _items {
					resp = append(resp, thirdParty.(map[string]interface{}))
				}
			} else {
				log.Printf("[WARNING] Expected an array of items")
			}
		} else {
			log.Printf("[WARNING] 'items' was not present in the response")
		}

		if _next, ok := tp["next"]; !ok || _next == nil {
			next = ""
		} else {
			next = _next.(string)
			time.Sleep(1 * time.Second)
		}
	}

	return resp
}

func fieldAsList(key string, source map[string]interface{}) []interface{} {
	if list, ok := source[key]; !ok || list == nil {
		return []interface{}{}
	} else if _list, ok := list.([]interface{}); !ok {
		return []interface{}{}
	} else {
		return _list
	}
}

func latestRiskForThirdParty(thirdParty map[string]interface{}, req *http.Request) {
	if allResidualRisk, ok := thirdParty["residual_risk"]; !ok {
		return
	} else if _allResidualRisk, ok := allResidualRisk.([]interface{}); !ok || len(_allResidualRisk) == 0 {
		delete(thirdParty, "residual_risk")
	} else if residualRisk, ok := _allResidualRisk[0].(map[string]interface{}); !ok {
		delete(thirdParty, "residual_risk")
	} else if reportURI, ok := residualRisk["residual_risk_uri"]; !ok || reportURI == nil {
		delete(thirdParty, "residual_risk")
	} else if scoreURI, ok := residualRisk["scores_uri"]; !ok || scoreURI == nil {
		delete(thirdParty, "residual_risk")
	} else if report, err := get(reportURI.(string), req.Header.Get("Authorization")); err != nil {
		log.Printf("[ERROR] Issue looking up report for third party, %s\n", err)
		delete(thirdParty, "residual_risk")
	} else if allScores, err := get(scoreURI.(string), req.Header.Get("Authorization")); err != nil {
		log.Printf("[ERROR] Issue looking up scores for third party, %s\n", err)
		delete(thirdParty, "residual_risk")
	} else {
		findings := make([]interface{}, 0, 120)
		findings = append(findings, fieldAsList("high_risks", report)...)
		findings = append(findings, fieldAsList("medium_risks", report)...)
		findings = append(findings, fieldAsList("low_risks", report)...)

		scores := make([]interface{}, 0, 1200)
		scores = append(scores, fieldAsList("group_scores", allScores)...)
		scores = append(scores, fieldAsList("control_scores", allScores)...)

		normalizedReport := map[string]interface{}{
			"report_type": residualRisk["report_type"],
			"date":        report["date"],
			"tier":        report["tier"],
			"id":          report["id"],
			"findings":    findings,
			"scores":      scores,
		}

		thirdParty["residual_risk"] = normalizedReport
		time.Sleep(1 * time.Second)
	}
}

func thirdParties(w http.ResponseWriter, req *http.Request) {
	thirdParties := allThirdParties(req)

	for _, thirdParty := range thirdParties {
		latestRiskForThirdParty(thirdParty, req)
	}

	responseJSON, err := json.Marshal(thirdParties)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func main() {
	r := gin.Default()

	// Serve embedded assets, any matched paths will take precidence over mapped routes
	r.Use(static.Serve("", assets.NewEmbeddedFileSystem()))

	// Retrieve and transform the public swagger documentation
	r.GET("/v1/swagger.json", func(c *gin.Context) {
		director := func(req *http.Request) {
			req.URL.Scheme = "https"
			req.URL.Host = "api.cybergrx.com"
			req.Host = "api.cybergrx.com"
		}
		proxy := &httputil.ReverseProxy{Director: director, ModifyResponse: docs.TransformApiDocumentation}
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	r.GET("/v1/third-parties", func(c *gin.Context) {
		thirdParties(c.Writer, c.Request)
	})

	host := strings.TrimSpace(os.Getenv("HOST"))
	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "8080"
	}

	r.Run(fmt.Sprintf("%s:%s", host, port))
}
