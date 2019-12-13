package memberservice

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/stodioo/roast/pkg/candi"
	"github.com/stodioo/roast/pkg/mailcore"
	"github.com/stodioo/roast/pkg/rest"
)

type MemberService struct {
	router      *mux.Router
	candiClient *candi.Client
	mailSvc     *mailcore.MailCore
}

func NewService(router *mux.Router, candiClient *candi.Client, mailSvc *mailcore.MailCore) *MemberService {
	return &MemberService{
		router:      router,
		candiClient: candiClient,
		mailSvc:     mailSvc,
	}
}

var httpClient = &http.Client{}

func (memberSvc *MemberService) SetupRouter(pathPrefix string) {
	r := memberSvc.router.PathPrefix(pathPrefix).Subrouter()
	// r.Use(rest.AuthAPIKeyMW)

	r.HandleFunc("/personal/data", memberSvc.postPersonalData).
		Methods(http.MethodPost).
		Name("Post personal data member registration")

	r.HandleFunc("/personal/terminal/register", memberSvc.terminalRegister).
		Methods(http.MethodPost).
		Name("Post terminal register")
}

func (memberSvc *MemberService) postPersonalData(resp http.ResponseWriter, req *http.Request) {
	memberSvc.candiClient.WithBearerToken("eyJhbGciOiJSUzI1NiIsImtpZCI6IjllNTUyYThmMjg1ZWM0NGY2MWNmNzU3ZjA5YjQ2ZWVkNTdiNzA5ZjgiLCJ0eXAiOiJKV1QifQ.eyJhenAiOiJjbC0weDI2NmQ0ZDAwIiwiaWF0IjoxNTc1MjczMTkyLCJqdGkiOiJhei0weDI2NmQ0ZDAwNTViYjJmN2QtMDAwMDAwMDA1ZGU0YzJlOCIsInN1YiI6ImktMHgwMDAwNjM5NjJjMjE1NTE0IiwidGVybWluYWxfaWQiOiJ0bC0weDI2NmQ0ZDAwNTViYjJmN2QifQ.rIndDk8wr8A1QemiEACjgXWm30VsBp7ShEfFjWJRuajBizyoC8nQuSnXLIUj5_mW0qEV7VUV5juN1xa-lh-ilxmjqJluE-s3Pi8xulnmazUdD3oxaB4wqVV47RrfjDwJghOZM7dGkDuoX0vkSmNK-cYpj4Gqdcnf_wdNcI_rrlQ")
	memberSvc.candiClient.TokenType(candi.BEARER_AUTH_SCHEME)
	user, err := memberSvc.candiClient.GetUserByPhoneNumber("+6282135770774,+6285640427774")
	if err != nil {

		http.Error(resp, "Unable to fetch user info", http.StatusBadRequest)
		return
	}

	rest.WriteHeaderAndJson(resp, http.StatusOK, user, rest.MIME_JSON)
}

type TerminalRegisterPostRequest struct {
	DisplayName              string `json:"display_name"`
	VerificationResourceType string `json:"verification_resource_type"`
	VerificationResourceName string `json:"verification_resource_name"`
	PlatformType             string `json:"platform_type,omitempty"`
}

func (memberSvc *MemberService) terminalRegister(resp http.ResponseWriter, req *http.Request) {
	memberSvc.candiClient.WithBearerToken("eyJhbGciOiJSUzI1NiIsImtpZCI6IjllNTUyYThmMjg1ZWM0NGY2MWNmNzU3ZjA5YjQ2ZWVkNTdiNzA5ZjgiLCJ0eXAiOiJKV1QifQ.eyJhenAiOiJjbC0weDI2NmQ0ZDAwIiwiaWF0IjoxNTc1MjczMTkyLCJqdGkiOiJhei0weDI2NmQ0ZDAwNTViYjJmN2QtMDAwMDAwMDA1ZGU0YzJlOCIsInN1YiI6ImktMHgwMDAwNjM5NjJjMjE1NTE0IiwidGVybWluYWxfaWQiOiJ0bC0weDI2NmQ0ZDAwNTViYjJmN2QifQ.rIndDk8wr8A1QemiEACjgXWm30VsBp7ShEfFjWJRuajBizyoC8nQuSnXLIUj5_mW0qEV7VUV5juN1xa-lh-ilxmjqJluE-s3Pi8xulnmazUdD3oxaB4wqVV47RrfjDwJghOZM7dGkDuoX0vkSmNK-cYpj4Gqdcnf_wdNcI_rrlQ")
	memberSvc.candiClient.TokenType(candi.BASIC_AUTH_SCHEME)
	var postReq TerminalRegisterPostRequest
	err := json.NewDecoder(req.Body).Decode(&postReq)

	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	terminal, err := memberSvc.candiClient.Register(postReq)

	if err != nil {
		http.Error(resp, "Unable to register", http.StatusBadRequest)
		return
	}

	rest.WriteHeaderAndJson(resp, http.StatusOK, terminal, rest.MIME_JSON)
}

type PersonalDataBodyRequest struct {
	Code        string    `json:"code"`
	DateOfBirth time.Time `json:"date_of_birth,omitempty"`
}
