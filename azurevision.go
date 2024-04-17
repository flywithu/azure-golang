package vision

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var apiVerstion string = "/vision/v3.2/"

type Tag struct {
	Name       string  `json:"name"`
	Confidence float64 `json:"confidence"`
}


type AzureTagResponse struct {
	Tags[] Tag
}




type AzureDescriptionResponse struct {
	Description  struct {
	Tags     []string  `json:"tags"`
	Captions []Caption `json:"captions"`
	}
}

type Caption struct {
	Text       string  `json:"text"`
	Confidence float64 `json:"confidence"`
}

type Client struct {
	Endpoint        string
	SubscriptionKey string
}

func ComputerVisionClient(endpoint, subscriptionKey string) *Client {
	return &Client{
		Endpoint:        endpoint,
		SubscriptionKey: subscriptionKey,
	}
}



func (c *Client) analyzeImage(inputPath string,apiPath string) ([]byte, error) {
	var reqBody []byte
	var err error

	if strings.HasPrefix(inputPath, "http") {
		// Handle URL
		urlPayload := map[string]string{"url": inputPath}
		reqBody, err = json.Marshal(urlPayload)
		if err != nil {
			return nil, fmt.Errorf("error marshaling URL: %w", err)
		}
	} else {
		// Handle local file
		reqBody, err = ioutil.ReadFile(inputPath)
		if err != nil {
			return nil, fmt.Errorf("error reading file: %w", err)
		}
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", c.Endpoint+apiPath, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	contentType := "application/octet-stream"
	if strings.HasPrefix(inputPath, "http") {
		contentType = "application/json"
	}
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Ocp-Apim-Subscription-Key", c.SubscriptionKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request to Azure Vision API: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}


	return body, nil
}


func (c *Client) GetImageTags(inputPath string) (AzureTagResponse, error) {

	data, err :=c.analyzeImage(inputPath,apiVerstion+"tag")
    if err != nil {
        return AzureTagResponse{}, fmt.Errorf("error from analyzeImage: %w", err)
    }

    var visionTag AzureTagResponse
    if err := json.Unmarshal(data, &visionTag); err != nil {
        return AzureTagResponse{}, fmt.Errorf("error parsing JSON response: %w", err)
    }
    return visionTag, nil
}


func (c *Client) GetImageDesc(inputPath string) (AzureDescriptionResponse, error) {

	data, err :=c.analyzeImage(inputPath,apiVerstion+"describe")
    if err != nil {
        return AzureDescriptionResponse{}, fmt.Errorf("error from analyzeImage: %w", err)
    }


    var visionDesc AzureDescriptionResponse
    if err := json.Unmarshal(data, &visionDesc); err != nil {
        return AzureDescriptionResponse{}, fmt.Errorf("error parsing JSON response: %w", err)
    }
    return visionDesc, nil
}

// Parent struct{
// 	child ObjectDetection
// } `json:"parent"`


type AzureDetectResponse struct {
	Objects []struct {
		Rectangle  Rectangle `json:"rectangle"`
		ObjectInfo `json:",inline"`
		Parent     *ObjectInfo `json:"parent,omitempty"`  
	} `json:"objects"`
	RequestId string `json:"requestId"`
	MetaData struct{} `json:"metadata"`
	ModelVersion string `json:"modelVersion"`
}

type Rectangle struct {
	Top    int `json:"x"`   
	Left   int `json:"y"`  
	Width  int `json:"w"`  
	Height int `json:"h"`
}

type ObjectInfo struct {
	ObjectName string  `json:"object"`
	Confidence float64 `json:"confidence"`
}


func (c *Client) GetImageObject(inputPath string) (AzureDetectResponse, error) {

	data, err :=c.analyzeImage(inputPath,apiVerstion+"detect")
    if err != nil {
        return AzureDetectResponse{}, fmt.Errorf("error from analyzeImage: %w", err)
    }

    var visionDetect AzureDetectResponse
    if err := json.Unmarshal(data, &visionDetect); err != nil {
        return AzureDetectResponse{}, fmt.Errorf("error parsing JSON response: %w", err)
    }


    return visionDetect, nil
}



type AzureAnalyzeResponse struct {
	Categories  []Category `json:"categories"`
	Adult       AdultContent `json:"adult"`
	Tags        []Tag `json:"tags"`
	AzureDescriptionResponse `json:",inline"`
	Faces       []struct{} `json:"faces"`
	Objects     []struct{
		Rectangle  Rectangle `json:"rectangle"`
		ObjectInfo `json:",inline"`
	} `json:"objects"`
	Brands     []struct {
		Name string `json:"name"`
		Confidence float64 `json:"confidence"`
		Rectangle  Rectangle `json:"rectangle"`

	}
}

type Category struct {
	Name   string  `json:"name"`
	Score  float64 `json:"score"`
	Detail *struct {
		Landmarks   []Landmark `json:"landmarks"`
	} `json:"detail,omitempty"`
}

type Landmark struct {
	Name       string  `json:"name"`
	Confidence float64 `json:"confidence"`
}

type AdultContent struct {
	IsAdultContent bool    `json:"isAdultContent"`
	IsRacyContent  bool    `json:"isRacyContent"`
	IsGoryContent  bool    `json:"isGoryContent"`
	AdultScore     float64 `json:"adultScore"`
	RacyScore      float64 `json:"racyScore"`
	GoreScore      float64 `json:"goreScore"`
}






func (c *Client) GetImageAnalyze(inputPath string) (AzureAnalyzeResponse, error) {

	data, err :=c.analyzeImage(inputPath,apiVerstion+"analyze?visualFeatures=Categories,Adult,Tags,Description,Faces,Color,ImageType,Objects,Brands&details=Landmarks")
    if err != nil {
        return AzureAnalyzeResponse{}, fmt.Errorf("error from analyzeImage: %w", err)
    }

    var visionAnalyze AzureAnalyzeResponse

    if err := json.Unmarshal(data, &visionAnalyze); err != nil {
        return AzureAnalyzeResponse{}, fmt.Errorf("error parsing JSON response: %w", err)
    }


    return visionAnalyze, nil
}


