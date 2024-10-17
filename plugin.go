/*
	Package
*/

package traefik_api_token_middleware

/*
	Imports
*/

import (
	"regexp"
	"io"
	"context"
	"fmt"
	"net"
	"net/http"
    "time"
	"strings"
	"log"
	"strconv"
	"encoding/json"
    "os"
)

/*
	Define > Color Codes
*/

var Reset = "\033[0m" 
var Red = "\033[31m" 
var Green = "\033[32m" 
var Yellow = "\033[33m" 
var Blue = "\033[34m" 
var Magenta = "\033[35m" 
var Cyan = "\033[36m" 
var Gray = "\033[37m" 
var White = "\033[97m"

/*
	Define > Header Values
*/

const (
	xForwardedFor                      = "X-Forwarded-For"
	xRealIP                            = "X-Real-IP"
	countryHeader                      = "X-IPCountry"
)

/*
	Construct Configurations

	OTP Secret can be generated at		https://it-tools.tech/otp-generator
*/

type Config struct {
	AuthenticationHeader     	bool     	`json:"authenticationHeader,omitempty"`
	AuthenticationHeaderName 	string   	`json:"authenticationHeaderName,omitempty"`
	AuthenticationErrorMsg 		string   	`json:"authenticationErrorMsg,omitempty"`
	BearerHeader             	bool     	`json:"bearerHeader,omitempty"`
	BearerHeaderName         	string   	`json:"bearerHeaderName,omitempty"`
	Tokens                     	[]string 	`json:"tokens,omitempty"`
	RemoveHeadersOnSuccess   	bool     	`json:"removeHeadersOnSuccess,omitempty"`
	RemoveTokenNameOnFailure	bool     	`json:"removeTokenNameOnError,omitempty"`
	TimestampUnix     			bool     	`json:"timestampUnix,omitempty"`
	AllowedIPAddresses			[]string 	`yaml:"allowedIPAddresses,omitempty"`
}

/*
	Construct Response
*/

type Response struct {
	Message    	string 	`json:"message"`
	StatusCode 	int    	`json:"status_code"`
	Timestamp 	string	`json:"timestamp"`
}

/*
	Create Config
*/

func CreateConfig() *Config {
	return &Config{
		AuthenticationHeader:     	true,
		AuthenticationHeaderName: 	"X-API-TOKEN",
		AuthenticationErrorMsg: 	"Access Denied",
		BearerHeader:             	true,
		BearerHeaderName:         	"Authorization",
		Tokens:                  	make([]string, 0),
		RemoveHeadersOnSuccess:   	true,
		RemoveTokenNameOnFailure:	false,
		TimestampUnix:				false,
		AllowedIPAddresses:			make([]string, 0),
	}
}

type KeyAuth struct {
	next                     	http.Handler
	authenticationHeader     	bool
	authenticationHeaderName 	string
	authenticationErrorMsg   	string
	bearerHeader             	bool
	bearerHeaderName         	string
	tokens                     	[]string
	removeHeadersOnSuccess   	bool
	removeTokenNameOnFailure	bool
	timestampUnix				bool
	allowedIPAddresses    		[]net.IP
}

func sliceString(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}

	return false
}

func sliceIp(a net.IP, list []net.IP) bool {
	for _, b := range list {
		if b.Equal(a) {
			return true
		}
	}
	return false
}

func parseIP(addr string) (net.IP, error) {
	ipAddress := net.ParseIP(addr)

	if ipAddress == nil {
		return nil, fmt.Errorf("cant parse IP address from address [%s]", addr)
	}

	return ipAddress, nil
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	fmt.Printf( Red + "[Aetherx-apikey]: " + Reset + "Starting Plugin " + Magenta + "%s" + Reset + "\n instance: " + Yellow + "%+v" + Reset + "\n ctx: " + Yellow + "%+v \n\n", name, *config, ctx)

	/*
		Ip whitelist
	*/

	var allowedIPAddresses []net.IP
	for _, ipAddressEntry := range config.AllowedIPAddresses {
		ip, ipBlock, err := net.ParseCIDR(ipAddressEntry)
		if err == nil {
			allowedIPAddresses = append(allowedIPAddresses, ip)
			continue
		}

		ipAddress := net.ParseIP(ipAddressEntry)
		if ipAddress == nil {
			fmt.Printf( Red + "[Aetherx-apikey]: " + Reset + "allowedIPAddresses whitelist contains %s" + Red + "%s" + Reset, "invalid ip address")
		}
		allowedIPAddresses = append(allowedIPAddresses, ipAddress)
	}

	/*
		Scan for empty tokens / keys
	*/

	if len(config.Tokens) == 0 {
		return nil, fmt.Errorf("Must specify at least one valid api token in plugin configurations")
	}

	/*
		At least one header must be specified
	*/

	if !config.AuthenticationHeader && !config.BearerHeader {
		return nil, fmt.Errorf("Must specify either authenticationHeader or bearerHeader in dynamic configuration")
	}

	/*
		return structure
	*/

	return &KeyAuth {
		next:                     	next,
		authenticationHeader:     	config.AuthenticationHeader,
		authenticationHeaderName: 	config.AuthenticationHeaderName,
		authenticationErrorMsg: 	config.AuthenticationErrorMsg,
		bearerHeader:             	config.BearerHeader,
		bearerHeaderName:         	config.BearerHeaderName,
		tokens:                  	config.Tokens,
		removeHeadersOnSuccess:   	config.RemoveHeadersOnSuccess,
		removeTokenNameOnFailure:	config.RemoveTokenNameOnFailure,
		timestampUnix:   			config.TimestampUnix,
		allowedIPAddresses:			allowedIPAddresses,
	}, nil
}

/*
	taken api tokens and compare to list of valid tokens.
	return if specified token is valid
*/

func contains(token string, validTokens []string) bool {
	for _, a := range validTokens {
		if a == token {
			return true
		}
	}

	return false
}

/*
	Bearer takes API token in `Authorization: Bearer $token` variant and compares it to list ov valid tokens.
	Token is extracted from header request value.
	Returns whether token is in list of valid tokens
*/

func bearer(token string, validTokens []string) bool {
	re, _ := regexp.Compile(`Bearer\s(?P<token>[^$]+)`)
	matches := re.FindStringSubmatch(token)

	/*
		No Match Found > Wrong form
	*/

	// If no match found the value is in the wrong form.
	if matches == nil {
		return false
	}

	/*
		Match Found > Compare to list of valid tokens
	*/

	tokenIndex := re.SubexpIndex("token")
	extractedToken := matches[tokenIndex]

	return contains(extractedToken, validTokens)
}

func (ka *KeyAuth) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	reqIPAddr, err := ka.collectRemoteIP(req)

	/*
		Authentication Header > check for valid token
	*/

	if ka.authenticationHeader {
		if contains(req.Header.Get(ka.authenticationHeaderName), ka.tokens) {

			/*
				Authentication Header > X-API-TOKEN header request contains valid token
			*/

			if ka.removeHeadersOnSuccess {
				req.Header.Del(ka.authenticationHeaderName)
			}
	
			ka.next.ServeHTTP(rw, req)
	
			return
		}
	}

	/*
		Bearer Header > check for valid token
	*/

	if ka.bearerHeader {
		if bearer(req.Header.Get(ka.bearerHeaderName), ka.tokens) {

			/*
				Bearer Header > Request header contains valid Bearer Token
			*/

			if ka.removeHeadersOnSuccess {
				req.Header.Del(ka.bearerHeaderName)
			}

			ka.next.ServeHTTP(rw, req)

			return
		}
	}

	/*
		IP Whitelist
	*/

	if len(ka.allowedIPAddresses) > 0 {
		fmt.Printf( Red + "[Aetherx-apikey]: " + Reset + "IPs specified in setting " + Magenta + "%s" + Reset, "allowedIPAddresses")

		for _, ipAddress := range reqIPAddr {

			fmt.Printf( Red + "[Aetherx-apikey]: " + Reset + "Checking IP for whitelist access " + Magenta + "%s" + Reset, ipAddress)

			if sliceIp(*ipAddress, ka.allowedIPAddresses) {
				fmt.Printf( Red + "[Aetherx-apikey]: " + Reset + "Allowing whitelisted IP " + Magenta + "%s" + Reset + " to bypass apikey", ipAddress)
				req.Header.Del(ka.authenticationHeaderName)
				req.Header.Del(ka.bearerHeaderName)
				ka.next.ServeHTTP(rw, req)
				return
			}
		}
	}

	/*
		Gather some settings and values

		- default output msg
		- timestamp (Unix Timestamp || UnixDate)
	*/

	output := "Access Denied"
	now := time.Now().Format(time.UnixDate) // UnixDate

    if len(ka.authenticationErrorMsg) > 0 {
        output = ka.authenticationErrorMsg
    }

	/*
		ANSIC 		"Mon Jan _2 15:04:05 2006"
		UnixDate 	"Mon Jan _2 15:04:05 PST 2006"
		RubyDate 	"Mon Jan 02 15:04:05 -0700 2006"
		RFC822 		"02 Jan 06 15:04 PST"
		RFC822Z 	"02 Jan 06 15:04 -0700"
		RFC850 		"Monday, 02-Jan-06 15:04:05 PST"
		RFC1123 	"Mon, 02 Jan 2006 15:04:05 PST"
		RFC1123Z 	"Mon, 02 Jan 2006 15:04:05 -0700"
		RFC3339 	"2006-01-02T15:04:05Z07:00"
		RFC3339Nano	"2006-01-02T15:04:05.999999999Z07:00"
	*/

    if ka.timestampUnix {
        var ts int = int(time.Now().Unix()) // Unix timestamp
		now = strconv.Itoa(ts)
		// int - 
    }

	/*
		Determine Auth Method & return response
	*/

	var response Response
	if ka.authenticationHeader && ka.bearerHeader {
		if !ka.removeTokenNameOnFailure {
			output = fmt.Sprintf(output + ". Provide a valid API Token header using either %s: $token or %s: Bearer $token", ka.authenticationHeaderName, ka.bearerHeaderName)
		}

		response = Response{
			Message:    	output,
			StatusCode: 	http.StatusForbidden,
			Timestamp: 		now,
		}
	} else if ka.authenticationHeader && !ka.bearerHeader {
		if !ka.removeTokenNameOnFailure {
			output = fmt.Sprintf(output + ". Provide a valid API Token header using %s: $token", ka.authenticationHeaderName)
		}

		response = Response{
			Message:    	output,
			StatusCode: 	http.StatusForbidden,
			Timestamp: 		now,
		}
	} else if !ka.authenticationHeader && ka.bearerHeader {
		if !ka.removeTokenNameOnFailure {
			output = fmt.Sprintf(output + ". Provide a valid API Token header using %s: Bearer $token", ka.bearerHeaderName)
		}

		response = Response{
			Message:    	output,
			StatusCode: 	http.StatusForbidden,
			Timestamp: 		now,
		}
	}

	/*
		Set Headers
	*/

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(response.StatusCode)

	/*
		No Headers or Invalid Key
		return 403 Forbidden
	*/

	if err := json.NewEncoder(rw).Encode(response); err != nil {

		/*
			Response can't be written > log error
		*/
		
		fmt.Printf( Red + "[Aetherx-apikey]: " + Reset + "Erroneous response due to invalid api token: " + Red + "%s" + Reset, err.Error())
	}
}

/*
	Collect Remote IP
*/

func (a *KeyAuth) collectRemoteIP(req *http.Request) ([]*net.IP, error) {
	var ipList []*net.IP

	splitFn := func(c rune) bool {
		return c == ','
	}

	xForwardedForValue := req.Header.Get(xForwardedFor)
	xForwardedForIPs := strings.FieldsFunc(xForwardedForValue, splitFn)

	xRealIPValue := req.Header.Get(xRealIP)
	xRealIPList := strings.FieldsFunc(xRealIPValue, splitFn)

	for _, value := range xForwardedForIPs {
		ipAddress, err := parseIP(value)
		if err != nil {
			return ipList, fmt.Errorf("parsing failed: %s", err)
		}

		ipList = append(ipList, &ipAddress)
	}

	for _, value := range xRealIPList {
		ipAddress, err := parseIP(value)
		if err != nil {
			return ipList, fmt.Errorf("parsing failed: %s", err)
		}

		ipList = append(ipList, &ipAddress)
	}

	return ipList, nil
}
