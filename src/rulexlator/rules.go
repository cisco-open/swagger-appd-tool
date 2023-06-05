/*
	SPDX-License-Identifier: Apache-2.0

	Copyright (c) 2023 Cisco Systems, Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"

	appdrest "gitlab.com/Divis/appdynamics-rest-api"

	"sigs.k8s.io/yaml"
)

const SWAGGER_GENERATED_RULE_MARK = "SWAGGER-GENERATED"

var appdClient *appdrest.Client
var appdClientInitted = false

type SwaggerRule struct {
	Method      string
	UrlTemplate string
	Description string
	UriRegex    string
}

type AppDRule struct {
	Method      string
	UriRegex    string
	Description string
	Name        string
	Id          string
}

const (
	LANG_JAVA = iota
	LANG_DOTNET
	LANG_NODEJS
	LANG_PYTHON
	LANG_PHP
	LANG_APACHE
)

type LangConsts struct {
	AgentType        string
	TxEntrypointType string
}

var langConsts = map[int]LangConsts{
	LANG_JAVA:   {AgentType: "APPLICATION_SERVER", TxEntrypointType: "SERVLET"},
	LANG_DOTNET: {AgentType: "DOT_NET_APPLICATION_SERVER", TxEntrypointType: "ASP_DOTNET"},
	LANG_NODEJS: {AgentType: "NODE_JS_SERVER", TxEntrypointType: "NODEJS_WEB"},
	LANG_PYTHON: {AgentType: "PYTHON_SERVER", TxEntrypointType: "PYTHON_WEB"},
	LANG_PHP:    {AgentType: "PHP_APPLICATION_SERVER", TxEntrypointType: "PHP_WEB"},
	LANG_APACHE: {AgentType: "NATIVE_WEB_SERVER", TxEntrypointType: "WEB"},
}

func processSwaggerFile(spec string, appId string, scopeName string, langId int) error {

	fmt.Printf("Swagger for app/scope: %s/%s\n%s\n", appId, scopeName, spec)
	if scopeName == "undefined" {
		scopeName = "Default Scope"
	}

	fmt.Printf("Parsing swagger...")
	sr, err := parseSwagger(spec)
	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Printf("Getting AppD rules...")
	ar, err := getAppDRules(appId, scopeName)
	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Printf("Searching for rule diffs...")
	swaggerRulesToAdd, appdRulesToDelete := diffRules(sr, ar)

	fmt.Printf("ToDo: \nAdd: %v\nDelete: %v\n\n", swaggerRulesToAdd, appdRulesToDelete)

	mdsData, err := getAppdRulesFromSwaggerRules(scopeName, swaggerRulesToAdd, langId)
	if err != nil {
		return err
	}

	// Printf("Swagger To Add: %v\n", *mdsData)

	// Delete no longer used rules
	for _, r := range *appdRulesToDelete {
		err = appdClient.TxDetectionRule.DeleteTransactionDetectionRule(r.Id)
		if err != nil {
			fmt.Printf("Error deleting trx detection rules, %v\n", err)
		}
	}

	err = uploadRulesToAppd(appId, scopeName, mdsData)
	if err != nil {
		return err
	}

	return nil
}

func diffRules(swaggerRules []*SwaggerRule, appdRules *[]AppDRule) (*[]SwaggerRule, *[]AppDRule) {
	swaggerRulesToAdd := []SwaggerRule{}
	appdRulesToDelete := []AppDRule{}

	sr := swaggerRules
	ar := *appdRules
	// fmt.Printf("SR:\n%v\n\nAR:\n%v\n\n", sr, ar)

	sort.Slice(sr[:], func(i, j int) bool {
		if sr[i].Method < sr[j].Method {
			return true
		}
		if sr[i].Method == sr[j].Method {
			return sr[i].UriRegex < sr[j].UriRegex
		}
		return false
	})

	sort.SliceStable(ar[:], func(i, j int) bool {
		if ar[i].Method < ar[j].Method {
			return true
		}
		if ar[i].Method == ar[j].Method {
			return ar[i].UriRegex < ar[j].UriRegex
		}
		return false
	})

	//for _, r := range sr {
	//	fmt.Printf("SR: %v\n", r)
	//}
	//fmt.Printf("SR:\n%v\n\nAR:\n%v\n\n", sr, ar)

	n := len(sr)
	m := len(ar)

	// Traverse both arrays simultaneously.
	i := 0
	j := 0
	for i < n && j < m {

		// Print smaller element and move
		// ahead in array with smaller element
		if compareSwaggerToAppdRule(*sr[i], ar[j]) < 0 {
			// fmt.Printf("SR < AR: %s-%s < %s-%s\n", sr[i].Method, sr[i].UriRegex, ar[j].Method, ar[j].UriRegex)
			swaggerRulesToAdd = append(swaggerRulesToAdd, *sr[i])
			i++
		} else {
			if compareSwaggerToAppdRule(*sr[i], ar[j]) > 0 {
				// fmt.Printf("AR < SR: %s-%s < %s-%s\n", ar[j].Method, ar[j].UriRegex, sr[i].Method, sr[i].UriRegex)
				appdRulesToDelete = append(appdRulesToDelete, ar[j])
				j++
			} else {
				i++
				j++
			}
		}
	}

	for i < n {
		// fmt.Printf("SR left: %s-%s\n", sr[i].Method, sr[i].UriRegex)
		swaggerRulesToAdd = append(swaggerRulesToAdd, *sr[i])
		i++
	}
	for j < m {
		// fmt.Printf("AR left: %s-%s\n", ar[j].Method, ar[j].UriRegex)
		appdRulesToDelete = append(appdRulesToDelete, ar[j])
		j++
	}

	return &swaggerRulesToAdd, &appdRulesToDelete
}

func compareSwaggerToAppdRule(sr SwaggerRule, ar AppDRule) int {
	if sr.Method < ar.Method {
		return -1
	} else {
		if sr.Method > ar.Method {
			return 1
		}
	}

	if sr.UriRegex < ar.UriRegex {
		return -1
	} else {
		if sr.UriRegex > ar.UriRegex {
			return 1
		}
	}

	return 0
}

func parseSwagger(spec string) ([]*SwaggerRule, error) {
	result := []*SwaggerRule{}

	var parsedSwagger map[string]interface{}

	if err := yaml.Unmarshal([]byte(spec), &parsedSwagger); err != nil {
		log.Fatal(err)
	}

	// log.Printf("%+v\n\n", parsedSwagger)

	paths := parsedSwagger["paths"].(map[string]interface{})
	// log.Printf("%+v\n\n", paths)

	for path, pathSpec := range paths {
		ps := pathSpec.(map[string]interface{})
		// log.Printf("%+v\n->\n%+v\n\n", path, ps)
		for method, methodSpec := range ps {
			ms := methodSpec.(map[string]interface{})

			desc := ""
			if ms["summary"] != nil {
				desc = ms["summary"].(string)
			}

			result = append(result, &SwaggerRule{
				Method:      swaggerMethodToAppD(method),
				UrlTemplate: path,
				Description: desc,
				UriRegex:    swaggerUriToRegex(path),
			})
			// log.Printf("%s %+v -> %+v\n\n", method, path, ms["summary"])
		}
	}

	return result, nil
}

func getAppdRulesFromSwaggerRules(scope string, rules *[]SwaggerRule, langId int) (*appdrest.MdsData, error) {

	ruleList := appdrest.RuleList{
		Rule: []appdrest.Rule{},
	}

	for _, r := range *rules {
		matchString := getRuleMatchString("MATCHES_REGEX", r.UriRegex, r.Method, langId)
		ruleName := r.Method + " - " + r.UrlTemplate
		ruleList.Rule = append(ruleList.Rule, getRule(
			ruleName, SWAGGER_GENERATED_RULE_MARK+": "+ruleName, true, matchString, langId,
		))
	}

	scopeList := appdrest.ScopeList{
		Scope: []appdrest.Scope{
			{
				ScopeDescription: "",
				ScopeName:        scope,
				ScopeType:        "ALL_TIERS_IN_APP",
				ScopeVersion:     "1",
			},
		},
	}

	scopeRuleMappingList := appdrest.ScopeRuleMappingList{
		ScopeRuleMapping: []appdrest.ScopeRuleMapping{},
	}

	var scopeRules = []appdrest.ScopeRule{}
	for _, r := range ruleList.Rule {
		scopeRules = append(scopeRules, appdrest.ScopeRule{
			RuleDescription: r.RuleDescription,
			RuleName:        r.RuleName,
		})
	}

	scopeRuleMappingList.ScopeRuleMapping = append(scopeRuleMappingList.ScopeRuleMapping, appdrest.ScopeRuleMapping{
		ScopeName: scope,
		ScopeRule: scopeRules,
	})

	var mdsData = appdrest.MdsData{
		ControllerVersion:    "021-009-002-000",
		ScopeList:            scopeList,
		RuleList:             ruleList,
		ScopeRuleMappingList: scopeRuleMappingList,
	}
	/*
		encoded, err := xml.Marshal(mdsData)
		if err != nil {
			fmt.Printf("Error marhalling trx detection rules\n")
		}

		fmt.Printf("MdsData: %s\n", string(encoded))
	*/

	return &mdsData, nil
}

func uploadRulesToAppd(applicationId string, scope string, ruleData *appdrest.MdsData) error {
	ensureConnected()

	err := appdClient.TxDetectionRule.UploadTransactionRules(applicationId, ruleData)
	if err != nil {
		fmt.Printf("Error uploading trx detection rules, %v\n", err)
		return err
	}

	return nil
}

func getAppDRules(appId string, scope string) (*[]AppDRule, error) {
	result := []AppDRule{}
	ensureConnected()

	data, err := appdClient.TxDetectionRule.GetTransactionDetectionRules(appId)
	if err != nil {
		fmt.Printf("Error getting trx detection rules, %v\n", err)
		return nil, err
	}

	for _, ruleMapping := range data.RuleScopeSummaryMappings {
		// fmt.Printf("%v %v %v\n", ruleMapping.Rule.Summary.ID, ruleMapping.Rule.Summary.Name,
		//	ruleMapping.Rule.Summary.Description)

		if strings.HasPrefix(ruleMapping.ScopeSummaries[0].Name, scope) {
			if strings.HasPrefix(ruleMapping.Rule.Summary.Description, SWAGGER_GENERATED_RULE_MARK) {
				//ruleId = ruleMapping.Rule.Summary.ID

				txMatchConditions := ruleMapping.Rule.TxMatchRule.TxCustomRule.MatchConditions
				// fmt.Printf("%v\n", txMatchConditions)
				if txMatchConditions[0].Type == "HTTP" {
					// fmt.Printf("%v\n", txMatchConditions[0].HttpMatch.Uri.MatchStrings[0])
					result = append(result, AppDRule{
						Method:      txMatchConditions[0].HttpMatch.HttpMethod,
						UriRegex:    txMatchConditions[0].HttpMatch.Uri.MatchStrings[0],
						Description: ruleMapping.Rule.Summary.Description,
						Name:        ruleMapping.Rule.Summary.Name,
						Id:          ruleMapping.Rule.Summary.ID,
					})
				}
			}
		}
	}

	return &result, nil
}

func getRule(name string, description string, enabled bool, match string, langId int) appdrest.Rule {
	rule := appdrest.Rule{
		AgentType:       langConsts[langId].AgentType, // "APPLICATION_SERVER",
		Enabled:         enabled,
		Priority:        "2",
		RuleDescription: description,
		RuleName:        name,
		RuleType:        "TX_MATCH_RULE",
		Version:         "1",
		TxMatchRule: appdrest.TxMatchRuleStr{
			Text: match,
		},
	}

	return rule
}

func getRuleMatchString(matchType string, matchUri string, method string, langId int) string {

	var httpMatch = appdrest.RuleMatchHttpmatch{
		URI: appdrest.RuleMatchURI{
			Type:         matchType,
			Matchstrings: []string{matchUri},
		},
		Httpmethod: method,
		Parameters: []interface{}{},
		Headers:    []interface{}{},
		Cookies:    []interface{}{},
	}
	var matchcondition = appdrest.RuleMatchMatchconditions{
		Type:      "HTTP",
		Httpmatch: httpMatch,
	}
	var txcustomrule = appdrest.RuleMatchTxcustomrule{
		Type:             "INCLUDE",
		Txentrypointtype: langConsts[langId].TxEntrypointType, // "SERVLET",
		Matchconditions:  []appdrest.RuleMatchMatchconditions{matchcondition},
		Actions:          []interface{}{},
		Properties:       []interface{}{},
	}
	var matchRule = appdrest.RuleMatchTxMatchRule{
		Type: "CUSTOM",
		Txautodiscoveryrule: appdrest.RuleMatchTxautodiscoveryrule{
			Autodiscoveryconfigs: []interface{}{},
		},
		Txcustomrule: txcustomrule,
		Agenttype:    langConsts[langId].AgentType, // "APPLICATION_SERVER",
	}

	encoded, err := json.Marshal(matchRule)
	if err != nil {
		return ""
	}

	return string(encoded)

}

func swaggerMethodToAppD(method string) string {
	appdMethod := "GET"
	switch method {
	case "get":
		appdMethod = "GET"
	case "put":
		appdMethod = "PUT"
	case "post":
		appdMethod = "POST"
	case "delete":
		appdMethod = "DELETE"
	default:
		appdMethod = "GET"
	}
	return appdMethod
}

func swaggerUriToRegex(uri string) string {

	swaggerParamRegex := "{[a-z,A-Z,0-9]*}"
	re := regexp.MustCompile(swaggerParamRegex)
	appdRegex := re.ReplaceAllLiteralString(uri, ".*")

	return appdRegex
}
