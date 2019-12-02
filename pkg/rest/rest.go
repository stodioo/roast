package rest

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.WithFields(logrus.Fields{"module": "rest"})

const (
	MIME_XML   = "application/xml"          // Accept or Content-Type used in Consumes() and/or Produces()
	MIME_JSON  = "application/json"         // Accept or Content-Type used in Consumes() and/or Produces()
	MIME_OCTET = "application/octet-stream" // If Content-Type is not present in request, use the default

	HEADER_Allow                         = "Allow"
	HEADER_Accept                        = "Accept"
	HEADER_Origin                        = "Origin"
	HEADER_ContentType                   = "Content-Type"
	HEADER_LastModified                  = "Last-Modified"
	HEADER_AcceptEncoding                = "Accept-Encoding"
	HEADER_ContentEncoding               = "Content-Encoding"
	HEADER_AccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	HEADER_AccessControlRequestMethod    = "Access-Control-Request-Method"
	HEADER_AccessControlRequestHeaders   = "Access-Control-Request-Headers"
	HEADER_AccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HEADER_AccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HEADER_AccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	HEADER_AccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HEADER_AccessControlMaxAge           = "Access-Control-Max-Age"

	ENCODING_GZIP    = "gzip"
	ENCODING_DEFLATE = "deflate"
)

var (
	MarshalIndent = json.MarshalIndent
	NewDecoder    = json.NewDecoder
	NewEncoder    = json.NewEncoder

	ErrAuthorizationFormatInvalid             = &Error{Msg: "authorization format invalid"}
	ErrAuthorizationTypeUnsupported           = &Error{Msg: "authorization type unsupported"}
	ErrAuthorizationTokenEmpty                = &Error{Msg: "authorization token empty"}
	ErrAuthorizationTokenInvalid              = &Error{Msg: "authorization token invalid"}
	ErrAuthorizationIDInvalid                 = &Error{Msg: "authorization ID invalid"}
	ErrAuthorizationTokenSubjectFormatInvalid = &Error{Msg: "authorization token subject user ID invalid"}
)

type ErrorResponse struct {
	Code string `json:"code,omitempty"`

	// We use the term description because it describes the error
	// to the developer rather than a message for the end user.
	Description string `json:"description,omitempty"`

	Fields []ErrorResponseField `json:"fields,omitempty"`
	DocURL string               `json:"doc_url,omitempty"`
}

type ErrorResponseField struct {
	Field       string `json:"field"`
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
	DocURL      string `json:"doc_url,omitempty"`
}

type EmptyResponse struct{}

type EmptyRequest struct{}

/*
AuthAPIKeyMW : Check x-api-key header middleware
Middleware ini digunakan untuk melakukan pengecekan header: x-api-key
api-key yang di cek antara lain: EXT_AUTH, INT_AUTH, MOB_AUTH,
semua auth key didefinisikan dari environment variable
*/
func AuthAPIKeyMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		var xAPIKeys []string

		if extAUTH := os.Getenv("EXT_AUTH"); extAUTH != "" {
			xAPIKeys = append(xAPIKeys, extAUTH)
		}
		if intAuth := os.Getenv("INT_AUTH"); intAuth != "" {
			xAPIKeys = append(xAPIKeys, intAuth)
		}

		if mobAUTH := os.Getenv("MOB_AUTH"); mobAUTH != "" {
			xAPIKeys = append(xAPIKeys, mobAUTH)
		}

		xAPIKey := req.Header.Get("x-api-key")
		if xAPIKey == "" {
			http.Error(resp, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if len(xAPIKeys) > 0 && !contains(xAPIKeys, xAPIKey) {
			log.Warnf("Invalid x-api-key")
			http.Error(resp, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(resp, req)
	})
}

func WriteHeaderAndJson(w http.ResponseWriter, status int, value interface{}, contentType string) error {
	if value == nil {
		w.WriteHeader(status)
		return nil
	}

	output, err := MarshalIndent(value, "", " ")
	if err != nil {
		return err
	}

	w.Header().Set(HEADER_ContentType, contentType)
	w.WriteHeader(status)
	_, err = w.Write(output)
	return err

}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
