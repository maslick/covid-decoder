package src

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/stapelberg/coronaqr"
	"github.com/stapelberg/coronaqr/trustlist/trustlistmirror"
	"log"
	"net/http"
	"os"
	"strings"
)

///////////////////////
// Model
//////////////////////

type Request struct {
	Code string `json:"code"`
}

type Message struct {
	Name         string     `json:"name"`
	DOB          string     `json:"dob"`
	KID          string     `json:"kid"`
	IssuedOn     string     `json:"issued_on"`
	ExpiresOn    string     `json:"expires_on"`
	IssuedBy     string     `json:"issued_by"`
	SignedBy     string     `json:"signed_by"`
	Vaccinations []Vaccine  `json:"vaccinations"`
	Tests        []Test     `json:"tests"`
	Recoveries   []Recovery `json:"recoveries"`
}

type Vaccine struct {
	Target        string  `json:"target"`
	Vaccine       string  `json:"vaccine"`
	Product       string  `json:"product"`
	Manufacturer  string  `json:"manufacturer"`
	Doses         float64 `json:"doses"`
	DoseSeries    float64 `json:"dose_series"`
	Date          string  `json:"date"`
	Country       string  `json:"country"`
	Issuer        string  `json:"issuer"`
	CertificateID string  `json:"certificate_id"`
}

type Test struct {
	Target         string `json:"target"`
	TestType       string `json:"test_type"`
	Name           string `json:"test_name"`
	Manufacturer   string `json:"manufacturer"`
	SampleDatetime string `json:"sample_datetime"`
	TestResult     string `json:"test_result"`
	TestingCentre  string `json:"testing_centre"`
	Country        string `json:"country"`
	Issuer         string `json:"issuer"`
	CertificateID  string `json:"certificate_id"`
}

type Recovery struct {
	Target            string `json:"target"`
	Country           string `json:"country"`
	Issuer            string `json:"issuer"`
	FirstPositiveTest string `json:"first_positive_test"`
	ValidFrom         string `json:"valid_from"`
	ValidUntil        string `json:"valid_until"`
	CertificateID     string `json:"certificate_id"`
}

type IService interface {
	Validate(code string) (Message, error)
	Parse(code string) (Message, error)
}

type Service struct {
	provider coronaqr.PublicKeyProvider
}

func (t *Service) Validate(code string) (Message, error) {
	decoded, err := t.parseNverify(code)
	if err != nil {
		return Message{}, err
	}
	return createMessage(decoded), err
}

func (t *Service) Parse(code string) (Message, error) {
	decoded, err := t.parse(code)
	if err != nil {
		return Message{}, err
	}
	return createMessage(decoded), err
}

func (t *Service) parse(code string) (*coronaqr.Decoded, error) {
	cert, err := coronaqr.Decode(code)
	if err != nil {
		return nil, err
	}
	return cert.SkipVerification(), nil
}

func (t *Service) parseNverify(code string) (*coronaqr.Decoded, error) {
	cert, err := coronaqr.Decode(code)
	if err != nil {
		return nil, err
	}

	if t.provider == nil {
		t.provider, err = trustlistmirror.NewCertificateProvider(context.Background(), trustlistmirror.TrustlistAT)
		if err != nil {
			return nil, err
		}
	}

	return cert.Verify(t.provider)
}

///////////////////////
// Rest Controller
//////////////////////

type RestController struct {
	Service IService
}

func (it *RestController) Init() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/parse", it.ParseHandler)
	mux.HandleFunc("/validate", it.ValidateHandler)
	mux.HandleFunc("/", it.HealthHandler)
	return mux
}

func (it *RestController) Start() {
	port := getPort()
	fmt.Println("Starting server on port:", strings.Split(port, ":")[1])
	log.Fatal(http.ListenAndServe(port, it.Init()))
}

func (it *RestController) ParseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", 400)
		return
	}

	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	mess, err := it.Service.Parse(req.Code)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(messToBytes(mess))
}

func (it *RestController) ValidateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", 400)
		return
	}

	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	mess, err := it.Service.Validate(req.Code)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(messToBytes(mess))
}

func (it *RestController) HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", 400)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("UP"))
}

///////////////////////
// Helper functions
//////////////////////

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getPort() string {
	var port = getEnv("PORT", "8080")
	return ":" + port
}

func createMessage(decoded *coronaqr.Decoded) Message {
	const dateFormat = "2006-01-02 15:04:05 MST"
	message := Message{}
	for _, record := range decoded.Cert.VaccineRecords {
		message.Vaccinations = append(message.Vaccinations, vaccineRecordToVaccine(record))
	}
	for _, record := range decoded.Cert.TestRecords {
		message.Tests = append(message.Tests, testRecordToTest(record))
	}
	for _, record := range decoded.Cert.RecoveryRecords {
		message.Recoveries = append(message.Recoveries, recoveryRecordToRecovery(record))
	}

	kid := decoded.Kid
	if len(kid) == 0 && decoded.SignedBy != nil {
		hash := sha256.Sum256(decoded.SignedBy.Raw)
		kid = hash[:8]
	}

	message.KID = base64.StdEncoding.EncodeToString(kid)
	message.IssuedOn = fmt.Sprintf("%+v", decoded.IssuedAt.Format(dateFormat))
	if decoded.SignedBy != nil {
		message.SignedBy = fmt.Sprintf("%s", decoded.SignedBy.Subject)
		message.IssuedBy = fmt.Sprintf("%s", decoded.SignedBy.Issuer)
	}

	message.ExpiresOn = fmt.Sprintf("%+v", decoded.Expiration.Format(dateFormat))
	message.Name = fmt.Sprintf(
		"%s %s",
		decoded.Cert.PersonalName.GivenNameStd,
		decoded.Cert.PersonalName.FamilyNameStd)

	message.DOB = fmt.Sprintf("%+v", decoded.Cert.DateOfBirth)
	return message
}

func messToBytes(testStruct Message) []byte {
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(testStruct)
	return reqBodyBytes.Bytes()
}

func vaccineRecordToVaccine(record coronaqr.VaccineRecord) Vaccine {
	return Vaccine{
		Target:        record.Target,
		Vaccine:       record.Vaccine,
		Product:       record.Product,
		Manufacturer:  record.Manufacturer,
		Doses:         record.Doses,
		DoseSeries:    record.DoseSeries,
		Date:          record.Date,
		Country:       record.Country,
		Issuer:        record.Issuer,
		CertificateID: record.CertificateID,
	}
}

func testRecordToTest(record coronaqr.TestRecord) Test {
	return Test{
		Target:         record.Target,
		TestType:       record.TestType,
		Name:           record.Name,
		Manufacturer:   record.Manufacturer,
		SampleDatetime: record.SampleDatetime,
		TestResult:     record.TestResult,
		TestingCentre:  record.TestingCentre,
		Country:        record.Country,
		Issuer:         record.Issuer,
		CertificateID:  record.CertificateID,
	}
}

func recoveryRecordToRecovery(record coronaqr.RecoveryRecord) Recovery {
	return Recovery{
		Target:            record.Target,
		Country:           record.Country,
		Issuer:            record.Issuer,
		FirstPositiveTest: record.FirstPositiveTest,
		ValidFrom:         record.ValidFrom,
		ValidUntil:        record.ValidUntil,
		CertificateID:     record.CertificateID,
	}
}
