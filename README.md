# covid-decoder
A simple REST API that parses COVID-19 certificates (provided in textual form)

![overview](overview.png)

## Installation
```zsh
# local
make serverfull 
make run

# AWS
aws configure
make deploy
```

## Usage
```zsh
CODE=`cat examples/cert-1.txt`

# Run locally
URL="localhost:8080"
curl -sX POST $URL/parse --data "{\"code\": \"${CODE}\"}" | jq
curl -sX POST $URL/validate --data "{\"code\": \"${CODE}\"}" | jq

# Run on AWS 
URL="https://f00wc4ugz6.execute-api.eu-central-1.amazonaws.com/dev"
curl -sX POST $URL/parse --data "{\"code\": \"${CODE}\"}" | jq
curl -sX POST $URL/validate --data "{\"code\": \"${CODE}\"}" | jq

# Scan QR code and decode
CODE=`zbarimg --raw qr.jpeg 2>&1 | head -n 1`
curl -sX POST $URL/parse --data "{\"code\": \"${CODE}\"}" | jq
curl -sX POST $URL/validate --data "{\"code\": \"${CODE}\"}" | jq
```

Output:
```json
{
    "dob": "2001-12-31",
    "expires_on": "2023-10-14 00:00:00 CEST",
    "issued_on": "2021-10-27 13:20:48 CEST",
    "issued_by": "CN=CSCA-FRANCE,O=Gouv,C=FR",
    "kid": "53FOjX/4aJs=",
    "name": "John Doe",
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

Fields description:
```
dob: Date of birth
name: Certificate owner's full name 
kid: Public key identifier

target: Disease or agent targeted
vaccine: Vaccine/prophylaxis
product: Vaccine medical product
manufacturer: Vaccine marketing authorization holder or manufacturer
doses: Number in a series of vaccinations/doses
dose_series: Total number of vaccination per series
date: Date of vaccination, indicating the date of the latest dose recieved
country: Member State or third country in which the vaccine was administered
issuer: Certificate issuer
certificate_id: Unique certificate identifier
```
