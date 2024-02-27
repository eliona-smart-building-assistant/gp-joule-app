//  This file is part of the eliona project.
//  Copyright © 2022 LEICOM iTEC AG. All Rights Reserved.
//  ______ _ _
// |  ____| (_)
// | |__  | |_  ___  _ __   __ _
// |  __| | | |/ _ \| '_ \ / _` |
// | |____| | | (_) | | | | (_| |
// |______|_|_|\___/|_| |_|\__,_|
//
//  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING
//  BUT NOT LIMITED  TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
//  NON INFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
//  DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package gp_joule

import (
	"fmt"
	utilshttp "github.com/eliona-smart-building-assistant/go-utils/http"
	"gp-joule/apiserver"
	"gp-joule/model"
	"net/http"
	"strings"
	"time"
)

func GetClusters(config *apiserver.Configuration) ([]*model.Cluster, error) {

	// create request
	fullUrl := config.RootUrl + "/clusters"
	request, err := utilshttp.NewRequestWithApiKey(fullUrl, "x-api-key", config.ApiKey)
	if err != nil {
		return nil, fmt.Errorf("error creating request for %s: %w", fullUrl, err)
	}

	// lowercase, because http package converts x-api-key to X-Api-Key
	lowerCaseHeader := make(http.Header)
	for key, value := range request.Header {
		lowerCaseHeader[strings.ToLower(key)] = value
	}
	request.Header = lowerCaseHeader

	// read clusters
	clusters, statusCode, err := utilshttp.ReadWithStatusCode[[]*model.Cluster](request, time.Duration(*config.RequestTimeout)*time.Second, true)
	if err != nil || statusCode != http.StatusOK {
		return nil, fmt.Errorf("error reading request for %s: %d %w", fullUrl, statusCode, err)
	}

	return clusters, nil
}
