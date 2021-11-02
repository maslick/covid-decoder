# covid-decoder

This is a super simple REST API that parses a COVID-19 certificate provided in textual form.

![overview](overview.png)

## Installation
```zsh
go build -ldflags="-s -w" && upx covid-decoder && mv covid-decoder decoder 
./decoder
```

## Usage
```zsh
CODE=`cat examples/cert-1.txt`

curl -sX POST localhost:8080/parse --data "{\"code\": \"${CODE}\"}" | jq
curl -sX POST localhost:8080/validate --data "{\"code\": \"${CODE}\"}" | jq
```

Output:
```json
{
    "dob": "2001-12-31",
    "expires_at": "2023-10-14 00:00:00 CEST",
    "issued_at": "2021-10-27 13:20:48 CEST",
    "issued_by": "CN=CSCA-FRANCE,O=Gouv,C=FR",
    "kid": "53FOjX/4aJs=",
    "name": "MICKEY MOUSE",
    "signed_by": "CN=DSC_FR_023,OU=180035024,O=CNAM,C=FR",
    "vaccinations": [{
        "target": "840539006",
        "vaccine": "1119349007",
        "product": "EU/1/20/1528",
        "manufacturer": "ORG-100030215",
        "doses": 1,
        "dose_series": 2,
        "date": "2021-10-27",
        "country": "FR",
        "issuer": "CNAM",
        "certificate_id": "URN:UVCI:01:FR:W7V2BE46QSBJ#L"
    }],
    "tests": [{
        "target": "840539006",
        "test_type": "LP217198-3",
        "test_name": "",
        "manufacturer": "1232",
        "sample_datetime": "2021-07-12T07:38:33Z",
        "test_result": "260415000",
        "testing_centre": "Health Institute, Paris",
        "country": "FR",
        "issuer": "CNAM",
        "certificate_id": "URN:UVCI:01:FR:W7V2BE46QSBJ#L"
        }],
    "recoveries": [{
        "target": "840539006",
        "country": "FR",
        "issuer": "CNAM",
        "first_positive_test": "2021-05-07",
        "valid_from": "2021-10-27 13:20:48 CEST",
        "valid_until": "2023-10-14 00:00:00 CEST",
        "certificate_id": "URN:UVCI:01:FR:W7V2BE46QSBJ#L"
    }]
}
```
