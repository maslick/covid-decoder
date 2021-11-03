# =covid-decoder=
A simple REST API that parses COVID-19 certificates (provided in textual form)

![overview](overview.png)

## :bulb: Motivation
There are plenty of Green Pass (COVID-19 certificate) validators out there (Github) written in different languages - Golang, Node.js, Python. Some just parse the QR-code, some verify the signature. I thought it might be a good idea to create a REST API which you can feed a QR code (in textual format) and get all the info contained in it. 

The server app can be compiled into a small binary (2-5 MB) and run on any server:
* Locally 
* Docker, Kubernetes, etc.
* Serverless on AWS
* Raspberry PI


## :white_check_mark: Installation
```zsh
# local
make serverfull 
make run

# docker
make docker-build

docker images | grep covid-decoder
covid-decoder        latest    4806d6ec9fef   1 minute ago   5.47MB

# AWS
aws configure
make deploy
```

## :rocket: Usage
```zsh
CODE=`cat my-covid-certificate.txt`

# Run locally
URL="localhost:8080"
curl -sX POST $URL/parse --data "{\"code\": \"${CODE}\"}" | jq
curl -sX POST $URL/validate --data "{\"code\": \"${CODE}\"}" | jq

# Run via docker
make docker-run
URL="localhost:8081"
curl -sX POST $URL/parse --data "{\"code\": \"${CODE}\"}" | jq
curl -sX POST $URL/validate --data "{\"code\": \"${CODE}\"}" | jq

# Run on AWS 
URL="https://f00wc4ugz6.execute-api.eu-central-1.amazonaws.com"
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
