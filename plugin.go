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
var RedL = "\033[91m"
var Green = "\033[32m"
var GreenL = "\033[92m"
var Orange = "\033[33m"
var Yellow = "\033[93m"
var Blue = "\033[34m"
var BlueL = "\033[94m"
var PurpleL = "\033[95m"
var Magenta = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var GrayD = "\033[90m"
var White = "\033[97m"

var BBlack = "\033[1;30m"
var BRed = "\033[1;31m"
var BGreen = "\033[1;32m"
var BYellow = "\033[1;33m"
var BBlue = "\033[1;34m"
var BPurple = "\033[1;35m"
var BCyan = "\033[1;36m"
var BWhite = "\033[1;37m"

type logWriter struct {
}

/*
	Logs > Writer
*/

func (writer logWriter) Write(bytes []byte) (int, error) {
	str := GrayD + time.Now().Format("2006-01-02T15:04:05") + Reset + " " + string(bytes)
	return io.WriteString(os.Stderr, str)
}

/*
	Logging
*/

var ( 
	logInfo = log.New(io.Discard, BPurple + "[ APITOKEN ] " + BlueL + "[INFO] " + Reset + ": ", log.Ldate|log.Ltime)
	logErr = log.New(io.Discard, BPurple + "[ APITOKEN ] " + RedL + "[ERROR] " + Reset + ": ", log.Ldate|log.Ltime)
	logWarn = log.New(io.Discard, BPurple + "[ APITOKEN ] " + Orange + "[WARN] " + Reset + ": ", log.Ldate|log.Ltime)
	logDebug = log.New(io.Discard, BPurple + "[ APITOKEN ] " + GrayD + "[Debug] " + Reset + ": ", log.Ldate|log.Ltime)
)

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
	AuthenticationHeader		bool		`json:"authenticationHeader,omitempty"`
	AuthenticationHeaderName	string		`json:"authenticationHeaderName,omitempty"`
	AuthenticationErrorMsg		string		`json:"authenticationErrorMsg,omitempty"`
	BearerHeader				bool		`json:"bearerHeader,omitempty"`
	BearerHeaderName			string		`json:"bearerHeaderName,omitempty"`
	PermissiveMode				bool		`json:"permissiveMode,omitempty"`
	RemoveHeadersOnSuccess		bool		`json:"removeHeadersOnSuccess,omitempty"`
	RemoveTokenNameOnFailure	bool		`json:"removeTokenNameOnError,omitempty"`
	TimestampUnix				bool		`json:"timestampUnix,omitempty"`
	DebugLogs					bool		`json:"debugLogs,omitempty"`
	InternalErrorRoute			string		`json:"internalErrorRoute,omitempty"`
	Tokens						[]string	`json:"tokens,omitempty"`
	WhitelistIPs				[]string	`yaml:"whitelistIPs,omitempty"`
	RegexAllow 					[]string 	`json:"regexAllow,omitempty"`
	RegexDeny					[]string 	`json:"regexDeny,omitempty"`
}

/*
	Construct Response
*/

type Response struct {
	Message    	string 	`json:"message"`
	StatusCode 	int    	`json:"status_code"`
	Timestamp 	string	`json:"timestamp"`
	UserAgent	string	`json:"user-agent"`
	RemoteAddr 	string 	`json:"ip"`
	Host       	string 	`json:"host"`
	RequestURI	string	`json:"uri"`
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
		PermissiveMode:				false,
		RemoveHeadersOnSuccess:   	true,
		RemoveTokenNameOnFailure:	false,
		TimestampUnix:				false,
		DebugLogs:					false,
		InternalErrorRoute:			"",
		Tokens:                  	make([]string, 0),
		WhitelistIPs:				make([]string, 0),
		RegexAllow: 				make([]string, 0),
		RegexDeny: 					make([]string, 0),
	}
}

type KeyAuth struct {
	next                     	http.Handler
	authenticationHeader     	bool
	authenticationHeaderName 	string
	authenticationErrorMsg   	string
	bearerHeader             	bool
	bearerHeaderName         	string
	permissiveMode				bool
	removeHeadersOnSuccess   	bool
	removeTokenNameOnFailure	bool
	timestampUnix				bool
	debugLogs           		bool
	internalErrorRoute			string
	tokens                     	[]string
	whitelistIPs    			[]net.IP
	regexpsAllow 				[]*regexp.Regexp
	regexpsDeny  				[]*regexp.Regexp
}

/*
	Strings > Slice
*/

func sliceString(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}

	return false
}

/*
	iP > Slice

	returns true of ips match
*/

func sliceIp(a net.IP, list []net.IP) bool {
	for _, b := range list {
		if b.Equal(a) {
			return true
		}
	}

	return false
}

/*
	IP > Parse
*/

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
		@TODO		merge logs
	*/

	logInfo.SetFlags(0)
	logInfo.SetOutput(new(logWriter))
	// logInfo.SetOutput(os.Stdout)

	logErr.SetFlags(0)
	logErr.SetOutput(new(logWriter))

	logWarn.SetFlags(0)
	logWarn.SetOutput(new(logWriter))

	logDebug.SetFlags(0)
	logDebug.SetOutput(new(logWriter))

	/*
		Regex

		allows you to filter out user-agents based on configured regex rules
	*/

	regexpsAllow := make([]*regexp.Regexp, len(config.RegexAllow))
	regexpsDeny := make([]*regexp.Regexp, len(config.RegexDeny))

	for i, regex := range config.RegexAllow {
		re, err := regexp.Compile(regex)
		if err != nil {
			return nil, fmt.Errorf("error compiling regexAllow %q: %w", regex, err)
		}

		regexpsAllow[i] = re
	}

	for i, regex := range config.RegexDeny {
		re, err := regexp.Compile(regex)
		if err != nil {
			return nil, fmt.Errorf("error compiling regex %q: %w", regex, err)
		}

		regexpsDeny[i] = re
	}

	/*
		Ip Whitelist
	*/

	var whitelistIPs []net.IP
	for _, ipAddressEntry := range config.WhitelistIPs {
		ip, ipBlock, err := net.ParseCIDR(ipAddressEntry)
		if err == nil {
			whitelistIPs = append(whitelistIPs, ip)
			continue
		}

		ipAddress := net.ParseIP(ipAddressEntry)
		if ipAddress == nil {
			logInfo.Printf(Reset + "whitelistIPs whitelist contains %s" + Red + "%s" + Reset, "invalid ip address")
		}
		whitelistIPs = append(whitelistIPs, ipAddress)
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
		permissiveMode:				config.PermissiveMode,
		removeHeadersOnSuccess:   	config.RemoveHeadersOnSuccess,
		removeTokenNameOnFailure:	config.RemoveTokenNameOnFailure,
		timestampUnix:   			config.TimestampUnix,
		whitelistIPs:				whitelistIPs,
		debugLogs:					config.DebugLogs,
		internalErrorRoute:			config.InternalErrorRoute,
		regexpsAllow: 				regexpsAllow,
		regexpsDeny:  				regexpsDeny,
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

/*
	Permissive Mode

	allows the request to pass even if no valid key is provided
*/

func (ka *KeyAuth) permissiveOk(rw http.ResponseWriter, req *http.Request) {
	logWarn.Printf(Reset + "Permissive Mode enabled, no valid credentials found. Allowing request anyway: " + Red + "\"%s\"" + Reset, req.URL)

	req.RequestURI = req.URL.RequestURI()
	ka.next.ServeHTTP(rw, req)
}

/*
	Serve
*/

func (ka *KeyAuth) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	reqIPAddr, err := ka.collectRemoteIP(req)
	now := time.Now().Format(time.UnixDate) // UnixDate
	bRegexWhitelist := false
	bRegexBlacklist := false
	userAgent := req.UserAgent()
	userIp := "Unknown"
	output := "Access Denied"
	var response Response

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
		logs > output user
		assign user ip to string
	*/

	for _, ipAddress := range reqIPAddr {
		userIp = ipAddress.String()
		logInfo.Printf(Reset + "Initializing new request " + Yellow + "%s" + Reset + " for url " + Yellow + "%s" + Reset, userIp, req.URL)
	}

	/*
		User-agent > Whitelist
	*/

	for _, re := range ka.regexpsAllow {
		if re.MatchString(userAgent) {
			if ka.debugLogs {
				logDebug.Printf(Reset + "Client " + Green + "%s" + Reset + " has whitelisted useragent " + Green + "%s" + Reset, userIp, userAgent)
			}
			bRegexWhitelist = true
		}
	}

	/*
		User-agent > Blacklist
	*/

	for _, re := range ka.regexpsDeny {
		if re.MatchString(userAgent) {
			if ka.debugLogs {
				logDebug.Printf(Reset + "Client " + Red + "%s" + Reset + " has blacklisted useragent " + Green + "%s" + Reset, userIp, userAgent)
			}
			bRegexBlacklist = true
		}
	}

	/*
		Debug Logs
		All users should pass this step before being directed to their proper destinations
	*/

	if ka.debugLogs {
		logDebug.Printf(Yellow + "%s :" + Reset + " %s ", "UserAgent", userAgent)
		logDebug.Printf(Yellow + "%s :" + Reset + " %s ", "RemoteAddr", req.RemoteAddr)
		logDebug.Printf(Yellow + "%s :" + Reset + " %s ", "ipAddress", userIp)
		logDebug.Printf(Yellow + "%s :" + Reset + " %s ", "Host", req.Host)
		logDebug.Printf(Yellow + "%s :" + Reset + " %s ", "RequestURI", req.RequestURI)
		logDebug.Printf(Yellow + "%s :" + Reset + " %s ", "req.URL", req.URL)
		logDebug.Printf(Yellow + "%s :" + Reset + " %s ", "Timestamp", now)
	}

	/*
		Authentication Error Message
	*/

    if len(ka.authenticationErrorMsg) > 0 {
        output = ka.authenticationErrorMsg
    }

	/*
		Allows you to use a "Dry-run" to pass no matter what even if no valid token was provided
	*/

	if ka.permissiveMode {
		ka.permissiveOk(rw, req)
		return
	}

	/*
		Override > User Agent
	*/

	if bRegexWhitelist {
		logInfo.Printf(Reset + "Client " + Green + "%s" + Reset + " passed useragent whitelist " + Green + "%s" + Reset, userIp, userAgent)

		req.Header.Del(ka.authenticationHeaderName)
		req.Header.Del(ka.bearerHeaderName)
		ka.next.ServeHTTP(rw, req)
		return
	}

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

	if len(ka.whitelistIPs) > 0 {
		if ka.debugLogs {
			logDebug.Printf(Reset + "Client IP whitelist check " + Green + "%s" + Reset, "enabled")
		}

		for _, ipAddress := range reqIPAddr {
			if ka.debugLogs {
				logDebug.Printf(Reset + "IP whitelist check initiated on client" + Green + "%s" + Reset + " with useragent " + Green + "%s" + Reset, userIp, userAgent)
			}

			if sliceIp(*ipAddress, ka.whitelistIPs) {
				logInfo.Printf(Reset + "Client " + Green + "%s" + Reset + " passed IP whitelist override " + Green + "%s" + Reset, userIp, userAgent)

				req.Header.Del(ka.authenticationHeaderName)
				req.Header.Del(ka.bearerHeaderName)
				ka.next.ServeHTTP(rw, req)
				return
			}
		}
	}

	/*
		Determine Auth Method & return response
	*/

	if bRegexBlacklist {
		if !ka.removeTokenNameOnFailure {
			output = fmt.Sprintf("Blacklisted useragent detected")
		}
	} else if ka.authenticationHeader && ka.bearerHeader {
		if !ka.removeTokenNameOnFailure {
			output = fmt.Sprintf(output + ". Provide a valid API Token header using either %s: $token or %s: Bearer $token", ka.authenticationHeaderName, ka.bearerHeaderName)
		}
	} else if ka.authenticationHeader && !ka.bearerHeader {
		if !ka.removeTokenNameOnFailure {
			output = fmt.Sprintf(output + ". Provide a valid API Token header using %s: $token", ka.authenticationHeaderName)
		}
	} else if !ka.authenticationHeader && ka.bearerHeader {
		if !ka.removeTokenNameOnFailure {
			output = fmt.Sprintf(output + ". Provide a valid API Token header using %s: Bearer $token", ka.bearerHeaderName)
		}
	}

	response = Response{
		Message:    	output,
		StatusCode: 	http.StatusForbidden,
		UserAgent:  	userAgent,
		RemoteAddr: 	req.RemoteAddr,
		Host:       	req.Host,
		RequestURI: 	req.RequestURI,
		Timestamp: 		now,
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
		logErr.Printf(Reset + "Erroneous response due to invalid api token: " + Red + "%s" + Reset, err.Error())
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
