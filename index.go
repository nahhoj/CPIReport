package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type IntegrationRuntimeArtifactsResponse struct {
	D struct {
		Results []struct {
			Metadata struct {
				ID          string `json:"id"`
				URI         string `json:"uri"`
				Type        string `json:"type"`
				ContentType string `json:"content_type"`
				MediaSrc    string `json:"media_src"`
				EditMedia   string `json:"edit_media"`
			} `json:"__metadata"`
			ID               string `json:"Id"`
			Version          string `json:"Version"`
			Name             string `json:"Name"`
			Type             string `json:"Type"`
			Deployedby       string `json:"DeployedBy"`
			Deployedon       string `json:"DeployedOn"`
			Status           string `json:"Status"`
			Delete           bool   `json:"Delete"` //Created for control
			Errorinformation struct {
				Deferred struct {
					URI string `json:"uri"`
				} `json:"__deferred"`
			} `json:"ErrorInformation"`
		} `json:"results"`
	} `json:"d"`
}

type IntegrationPackagesResponse struct {
	D struct {
		Results []struct {
			Metadata struct {
				ID          string `json:"id"`
				URI         string `json:"uri"`
				Type        string `json:"type"`
				ContentType string `json:"content_type"`
				MediaSrc    string `json:"media_src"`
				EditMedia   string `json:"edit_media"`
			} `json:"__metadata"`
			ID                             string      `json:"Id"`
			Name                           string      `json:"Name"`
			Description                    string      `json:"Description"`
			Shorttext                      string      `json:"ShortText"`
			Version                        string      `json:"Version"`
			Vendor                         string      `json:"Vendor"`
			Partnercontent                 bool        `json:"PartnerContent"`
			Updateavailable                bool        `json:"UpdateAvailable"`
			Mode                           string      `json:"Mode"`
			Supportedplatform              string      `json:"SupportedPlatform"`
			Modifiedby                     string      `json:"ModifiedBy"`
			Creationdate                   string      `json:"CreationDate"`
			Modifieddate                   string      `json:"ModifiedDate"`
			Createdby                      string      `json:"CreatedBy"`
			Products                       string      `json:"Products"`
			Keywords                       string      `json:"Keywords"`
			Countries                      string      `json:"Countries"`
			Industries                     string      `json:"Industries"`
			Lineofbusiness                 string      `json:"LineOfBusiness"`
			Packagecontent                 interface{} `json:"PackageContent"`
			Integrationdesigntimeartifacts struct {
				Deferred struct {
					URI string `json:"uri"`
				} `json:"__deferred"`
			} `json:"IntegrationDesigntimeArtifacts"`
			Valuemappingdesigntimeartifacts struct {
				Deferred struct {
					URI string `json:"uri"`
				} `json:"__deferred"`
			} `json:"ValueMappingDesigntimeArtifacts"`
			Customtags struct {
				Deferred struct {
					URI string `json:"uri"`
				} `json:"__deferred"`
			} `json:"CustomTags"`
		} `json:"results"`
	} `json:"d"`
}

type IntegrationDesigntimeArtifactsResponse struct {
	D struct {
		Results []struct {
			Metadata struct {
				ID          string `json:"id"`
				URI         string `json:"uri"`
				Type        string `json:"type"`
				ContentType string `json:"content_type"`
				MediaSrc    string `json:"media_src"`
				EditMedia   string `json:"edit_media"`
			} `json:"__metadata"`
			ID              string      `json:"Id"`
			Version         string      `json:"Version"`
			Packageid       string      `json:"PackageId"`
			Name            string      `json:"Name"`
			Description     string      `json:"Description"`
			Sender          string      `json:"Sender"`
			Receiver        string      `json:"Receiver"`
			Artifactcontent interface{} `json:"ArtifactContent"`
			Configurations  struct {
				Deferred struct {
					URI string `json:"uri"`
				} `json:"__deferred"`
			} `json:"Configurations"`
			Resources struct {
				Deferred struct {
					URI string `json:"uri"`
				} `json:"__deferred"`
			} `json:"Resources"`
		} `json:"results"`
	} `json:"d"`
}

type DesigntimeArtifacts struct {
	PackageID                  string
	PackageName                string
	Type                       string
	DesigntimeArtifactsID      string
	DesigntimeArtifactsName    string
	DesigntimeArtifactsVersion string
	Deployed                   bool
	Delete                     bool
}

type Response struct {
	StatusCode int
	Headers    map[string][]string
	Response   []byte
}

func main() {
	if len(os.Args) != 5 {
		fmt.Println("Not found parameters.")
		os.Exit(3)
	}
	dev := os.Args[1]
	prd := os.Args[2]
	credectial := os.Args[3] + ":" + os.Args[4]
	fmt.Println("Start...")
	fmt.Println("Getting data from API CPI for deploys artifacs...")
	authorization := "Basic " + base64.StdEncoding.EncodeToString([]byte(credectial))
	headers := map[string][]string{
		"Authorization": {authorization},
	}
	res := callService(dev+"/api/v1/IntegrationRuntimeArtifacts?$format=json", headers)
	if res.StatusCode != 200 {
		fmt.Println("There is an error: ", res.StatusCode, string(res.Response))
		os.Exit(3)
	}
	var res_IRA_DEV IntegrationRuntimeArtifactsResponse
	json.Unmarshal(res.Response, &res_IRA_DEV)
	res = callService(prd+"/api/v1/IntegrationRuntimeArtifacts?$format=json", headers)
	if res.StatusCode != 200 {
		fmt.Println("There is an error: ", res.StatusCode, string(res.Response))
		os.Exit(3)
	}
	var res_IRA_PRD IntegrationRuntimeArtifactsResponse
	json.Unmarshal(res.Response, &res_IRA_PRD)
	var fileDeploy []string
	fmt.Println("Procesing Deploys artifacs...")
	fileDeploy = append(fileDeploy, "Type|Name|ID DEV|Version Dev|ID PRD|Version PRD|Equal ID|Equal Version|Equal")
	for x, result_DEV := range res_IRA_DEV.D.Results {
		for y, result_PRD := range res_IRA_PRD.D.Results {
			line := result_DEV.Type + "|" + result_DEV.Name + "|" + result_DEV.ID + "|" + result_DEV.Version + "|" + result_PRD.ID + "|" + result_PRD.Version
			if result_PRD.Name == result_DEV.Name {
				if result_DEV.ID == result_PRD.ID {
					line += "|true"
				} else {
					line += "|false"
				}
				dev_v, _ := strconv.ParseInt(strings.ReplaceAll(result_DEV.Version, ".", ""), 10, 10)
				prd_v, _ := strconv.ParseInt(strings.ReplaceAll(result_PRD.Version, ".", ""), 10, 10)
				if dev_v >= prd_v {
					line += "|true"
				} else {
					line += "|false"
				}
				if strings.Contains(line, "|false") {
					line += "|false"
				} else {
					line += "|true"
				}
				//fmt.Println(line)
				fileDeploy = append(fileDeploy, line)
				res_IRA_DEV.D.Results[x].Delete = true
				res_IRA_PRD.D.Results[y].Delete = true
			}
		}
	}
	for _, result_DEV := range res_IRA_DEV.D.Results {
		if result_DEV.Delete != true {
			fileDeploy = append(fileDeploy, result_DEV.Type+"|"+result_DEV.Name+"|"+result_DEV.ID+"|"+result_DEV.Version+"|n/a|n/a|false|false|false")
		}
	}
	for _, result_PRD := range res_IRA_PRD.D.Results {
		if result_PRD.Delete != true {
			fileDeploy = append(fileDeploy, result_PRD.Type+"|"+result_PRD.Name+"|n/a|n/a|"+result_PRD.ID+"|"+result_PRD.Version+"|false|false|false")
		}
	}

	fmt.Println("Creating File for deploy artifcts...")
	file, _ := os.Create("/Users/johan/Desktop/deploys.csv")
	for _, line := range fileDeploy {
		file.WriteString(fmt.Sprintln(line))
	}
	file.Close()
	fmt.Println("It has created file deploys")
	fmt.Println("Getting data from API CPI for packages and artifacs...")
	res = callService(dev+"/api/v1/IntegrationPackages?$format=json", headers)
	if res.StatusCode != 200 {
		fmt.Println("There is an error: ", res.StatusCode, string(res.Response))
		os.Exit(3)
	}
	var res_IP_DEV IntegrationPackagesResponse
	json.Unmarshal(res.Response, &res_IP_DEV)
	res = callService(prd+"/api/v1/IntegrationPackages?$format=json", headers)
	if res.StatusCode != 200 {
		fmt.Println("There is an error: ", res.StatusCode, string(res.Response))
		os.Exit(3)
	}
	var res_IP_PRD IntegrationPackagesResponse
	json.Unmarshal(res.Response, &res_IP_PRD)
	var DesigntimeArtifacts_DEV []DesigntimeArtifacts
	for _, packages := range res_IP_DEV.D.Results {
		var res_IDA IntegrationDesigntimeArtifactsResponse
		res = callService(packages.Metadata.ID+"/IntegrationDesigntimeArtifacts?$format=json", headers)
		if res.StatusCode != 200 {
			fmt.Println("There is an error: ", res.StatusCode, string(res.Response))
			os.Exit(3)
		}
		json.Unmarshal(res.Response, &res_IDA)
		var DesigntimeArtifact_DEV DesigntimeArtifacts
		for _, DesigntimeArtifact := range res_IDA.D.Results {
			DesigntimeArtifact_DEV.PackageID = packages.ID
			DesigntimeArtifact_DEV.PackageName = packages.Name
			DesigntimeArtifact_DEV.DesigntimeArtifactsID = DesigntimeArtifact.ID
			DesigntimeArtifact_DEV.DesigntimeArtifactsName = DesigntimeArtifact.Name
			DesigntimeArtifact_DEV.DesigntimeArtifactsVersion = DesigntimeArtifact.Version
			DesigntimeArtifact_DEV.Type = DesigntimeArtifact.Metadata.Type
			for _, result_DEV := range res_IRA_DEV.D.Results {
				if DesigntimeArtifact.Name == result_DEV.Name && DesigntimeArtifact.ID == result_DEV.ID {
					DesigntimeArtifact_DEV.Deployed = true
				}
			}
			DesigntimeArtifacts_DEV = append(DesigntimeArtifacts_DEV, DesigntimeArtifact_DEV)
		}

		res = callService(packages.Metadata.ID+"/ValueMappingDesigntimeArtifacts?$format=json", headers)
		if res.StatusCode != 200 {
			fmt.Println("There is an error: ", res.StatusCode, string(res.Response))
			os.Exit(3)
		}
		json.Unmarshal(res.Response, &res_IDA)
		for _, DesigntimeArtifact := range res_IDA.D.Results {
			DesigntimeArtifact_DEV.PackageID = packages.ID
			DesigntimeArtifact_DEV.PackageName = packages.Name
			DesigntimeArtifact_DEV.DesigntimeArtifactsID = DesigntimeArtifact.ID
			DesigntimeArtifact_DEV.DesigntimeArtifactsName = DesigntimeArtifact.Name
			DesigntimeArtifact_DEV.DesigntimeArtifactsVersion = DesigntimeArtifact.Version
			DesigntimeArtifact_DEV.Type = DesigntimeArtifact.Metadata.Type
			for _, result_DEV := range res_IRA_DEV.D.Results {
				if DesigntimeArtifact.Name == result_DEV.Name && DesigntimeArtifact.ID == result_DEV.ID {
					DesigntimeArtifact_DEV.Deployed = true
				}
			}
			DesigntimeArtifacts_DEV = append(DesigntimeArtifacts_DEV, DesigntimeArtifact_DEV)
		}

	}
	var DesigntimeArtifacts_PRD []DesigntimeArtifacts
	for _, packages := range res_IP_PRD.D.Results {
		var res_IDA IntegrationDesigntimeArtifactsResponse
		res = callService(packages.Metadata.ID+"/IntegrationDesigntimeArtifacts?$format=json", headers)
		if res.StatusCode != 200 {
			fmt.Println("There is an error: ", res.StatusCode, string(res.Response))
			os.Exit(3)
		}
		json.Unmarshal(res.Response, &res_IDA)
		var DesigntimeArtifact_PRD DesigntimeArtifacts
		for _, DesigntimeArtifact := range res_IDA.D.Results {
			DesigntimeArtifact_PRD.PackageID = packages.ID
			DesigntimeArtifact_PRD.PackageName = packages.Name
			DesigntimeArtifact_PRD.DesigntimeArtifactsID = DesigntimeArtifact.ID
			DesigntimeArtifact_PRD.DesigntimeArtifactsName = DesigntimeArtifact.Name
			DesigntimeArtifact_PRD.DesigntimeArtifactsVersion = DesigntimeArtifact.Version
			DesigntimeArtifact_PRD.Type = DesigntimeArtifact.Metadata.Type
			for _, result_PRD := range res_IRA_PRD.D.Results {
				if DesigntimeArtifact.Name == result_PRD.Name && DesigntimeArtifact.ID == result_PRD.ID {
					DesigntimeArtifact_PRD.Deployed = true
				}
			}
			DesigntimeArtifacts_PRD = append(DesigntimeArtifacts_PRD, DesigntimeArtifact_PRD)
		}

		res = callService(packages.Metadata.ID+"/ValueMappingDesigntimeArtifacts?$format=json", headers)
		if res.StatusCode != 200 {
			fmt.Println("There is an error: ", res.StatusCode, string(res.Response))
			os.Exit(3)
		}
		json.Unmarshal(res.Response, &res_IDA)
		for _, DesigntimeArtifact := range res_IDA.D.Results {
			DesigntimeArtifact_PRD.PackageID = packages.ID
			DesigntimeArtifact_PRD.PackageName = packages.Name
			DesigntimeArtifact_PRD.DesigntimeArtifactsID = DesigntimeArtifact.ID
			DesigntimeArtifact_PRD.DesigntimeArtifactsName = DesigntimeArtifact.Name
			DesigntimeArtifact_PRD.DesigntimeArtifactsVersion = DesigntimeArtifact.Version
			DesigntimeArtifact_PRD.Type = DesigntimeArtifact.Metadata.Type
			for _, result_PRD := range res_IRA_PRD.D.Results {
				if DesigntimeArtifact.Name == result_PRD.Name && DesigntimeArtifact.ID == result_PRD.ID {
					DesigntimeArtifact_PRD.Deployed = true
				}
			}
			DesigntimeArtifacts_PRD = append(DesigntimeArtifacts_PRD, DesigntimeArtifact_PRD)
		}
	}
	var filePackages []string
	filePackages = append(filePackages, "Type|Packages Name|Packages ID DEV|Packages ID PRD|Artifacts Name|Artifacts ID DEV|Artifacts ID PRD|Artifacts Version DEV|Artifacts Version PRD|Deployed DEV|Deployed PRD|Equal")
	for x, DesigntimeArtifact_DEV := range DesigntimeArtifacts_DEV {
		for y, DesigntimeArtifact_PRD := range DesigntimeArtifacts_PRD {
			if DesigntimeArtifact_PRD.PackageName == DesigntimeArtifact_DEV.PackageName && DesigntimeArtifact_PRD.DesigntimeArtifactsName == DesigntimeArtifact_DEV.DesigntimeArtifactsName {
				line := DesigntimeArtifact_DEV.Type + "|" + DesigntimeArtifact_DEV.PackageName + "|" + DesigntimeArtifact_DEV.PackageID + "|" + DesigntimeArtifact_PRD.PackageID + "|" + DesigntimeArtifact_DEV.DesigntimeArtifactsName + "|" + DesigntimeArtifact_DEV.DesigntimeArtifactsID + "|" + DesigntimeArtifact_PRD.DesigntimeArtifactsID + "|" + DesigntimeArtifact_DEV.DesigntimeArtifactsVersion + "|" + DesigntimeArtifact_PRD.DesigntimeArtifactsVersion + "|" + strconv.FormatBool(DesigntimeArtifact_DEV.Deployed) + "|" + strconv.FormatBool(DesigntimeArtifact_PRD.Deployed)
				if DesigntimeArtifact_PRD.DesigntimeArtifactsVersion == "Active" || DesigntimeArtifact_PRD.DesigntimeArtifactsVersion == "Draft" {
					line += "|false"
				} else if DesigntimeArtifact_DEV.DesigntimeArtifactsVersion == "Active" || DesigntimeArtifact_DEV.DesigntimeArtifactsVersion == "Draft" {
					line += "|true"
				} else {
					dev_v, _ := strconv.ParseInt(strings.ReplaceAll(DesigntimeArtifact_DEV.DesigntimeArtifactsVersion, ".", ""), 10, 10)
					prd_v, _ := strconv.ParseInt(strings.ReplaceAll(DesigntimeArtifact_PRD.DesigntimeArtifactsVersion, ".", ""), 10, 10)
					if dev_v >= prd_v {
						line += "|true"
					} else {
						line += "|false"
					}
				}
				filePackages = append(filePackages, line)
				DesigntimeArtifacts_PRD[y].Delete = true
				DesigntimeArtifacts_DEV[x].Delete = true
			}
		}
	}

	for _, DesigntimeArtifact_DEV := range DesigntimeArtifacts_DEV {
		if DesigntimeArtifact_DEV.Delete == false {
			line := DesigntimeArtifact_DEV.Type + "|" + DesigntimeArtifact_DEV.PackageName + "|" + DesigntimeArtifact_DEV.PackageID + "|n/a|" + DesigntimeArtifact_DEV.DesigntimeArtifactsName + "|" + DesigntimeArtifact_DEV.DesigntimeArtifactsID + "|n/a|" + DesigntimeArtifact_DEV.DesigntimeArtifactsVersion + "|n/a|" + strconv.FormatBool(DesigntimeArtifact_DEV.Deployed) + "|n/a|TRUE"
			filePackages = append(filePackages, line)
		}
	}
	for _, DesigntimeArtifact_PRD := range DesigntimeArtifacts_PRD {
		if DesigntimeArtifact_PRD.Delete == false {
			line := DesigntimeArtifact_PRD.Type + "|" + DesigntimeArtifact_PRD.PackageName + "|n/a|" + DesigntimeArtifact_PRD.PackageID + "|" + DesigntimeArtifact_PRD.DesigntimeArtifactsName + "|n/a|" + DesigntimeArtifact_PRD.DesigntimeArtifactsID + "|n/a|" + DesigntimeArtifact_PRD.DesigntimeArtifactsVersion + "|n/a|" + strconv.FormatBool(DesigntimeArtifact_PRD.Deployed) + "|FALSE"
			filePackages = append(filePackages, line)
		}
	}

	fmt.Println("Creating File for packages and artifcts...")
	file, _ = os.Create("/Users/johan/Desktop/packages.csv")
	for _, line := range filePackages {
		file.WriteString(fmt.Sprintln(line))
	}
	file.Close()
	fmt.Println("It has created file packages and artifcts")
}

func callService(urlString string, headers map[string][]string) Response {
	urlParse, _ := url.Parse(urlString)
	req := http.Request{
		Method: "GET",
		URL:    urlParse,
		Header: http.Header(headers),
	}
	res, err := http.DefaultClient.Do(&req)
	if err != nil {
		return Response{
			StatusCode: 500,
			Response:   []byte(err.Error()),
		}
	}
	resBytes, _ := ioutil.ReadAll(res.Body)
	return Response{
		StatusCode: res.StatusCode,
		Headers:    map[string][]string(res.Header),
		Response:   resBytes,
	}
}
