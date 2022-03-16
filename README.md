# Backend

## Prerequisites
1. Go version >= 1.14

## How to run this project

### Local
1. Install the requirements
2. Clone this repo
3. Run `Go mod tidy`
4. Run `make run-http`

because of this repo using aws database (dynamoDB), we need to specify the aws config and credentials.
1. create file `~/.aws/credentials`
```
[default]
aws_access_key_id = AKIART56DJS6KDC3BQ5W
aws_secret_access_key = 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
```
2. create file `~/.aws/config`
```
[default]
region=ap-southeast-1
output=json
```


### Cloud
We could access this running env from this url:
http://ec2-13-250-36-182.ap-southeast-1.compute.amazonaws.com:8000/

## Notes
For front end repo, kindly check from this link (https://github.com/papannn/fe-bgp-hackaton)


