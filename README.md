# StatusOK

Monitor your Website and APIs from your computer. Get notified through Slack or E-mail when your server is down or response time is more than expected.

## Fork informations

This fork from [sanathp/statusok](https://github.com/sanathp/statusok) brings those changes:

* Shorter notifications, aiming for better readability in chat channels such as Slack
* Upgrade from Go 1.11 to 1.17 (add `go.mod`)
* Simplify the `Dockerfile` to build from source code (instead of building from a downloaded zip file)
* Add `docker-compose.yml` file to easily deploy the stack
* Add a GitHub Actions workflow to publish an image to [`ghcr.io/vemonet/statusok`](https://github.com/vemonet/statusok/pkgs/container/statusok)

[![Publish StatusOK image](https://github.com/vemonet/statusok/actions/workflows/publish-docker.yml/badge.svg)](https://github.com/vemonet/statusok/actions/workflows/publish-docker.yml)

## Simple Version

Simple Setup to monitor your website and recieve a notification to your Gmail when your website is down.

Step 1: Write a `config.json` with the url information 
```json
{
	"notifications":{
		"mail":{
			"smtpHost":"smtp.gmail.com",
			"port":587,
			"username":"yourmailid@gmail.com",
			"password":"your gmail password",
			"from":"yourmailid@gmail.com",
			"to":"destemailid@gmail.com"
		}
	},
	"requests":[
		{
			"url":"http://mywebsite.com",
			"requestType":"GET",
			"checkEvery":30,	
			"responseTime":800
		}
	]
}
```
Turn on access for your gmail https://www.google.com/settings/security/lesssecureapps .

Step 2: Download bin file from [here](https://github.com/vemonet/statusok/releases/) and run the below command from your terminal
```bash
./statusok --config config.json
```
Thats it !!!! You will receive a mail when your website is down or response time is more.

To run as background process add & at the end

```bash
./statusok --config config.json &	
```
to stop the process 
```bash
jobs
kill %jobnumber
```

## Complete Version using InfluxDb

![alt text](https://github.com/vemonet/statusok/raw/master/screenshots/graphana.png "Graphana Screenshot")

You can save data to influx db and view response times over a period of time as above using graphana.

[Guide to install influxdb and grafana](https://github.com/vemonet/statusok/blob/master/Config.md#database) 

With StatusOk you can monitor all your REST APIs by adding api details to config file as below.A Notification will be triggered when you api is down or response time is more than expected.

```json
{
	"url":"http://mywebsite.com/v1/data",
	"requestType":"POST",
	"headers":{
		"Authorization":"Bearer ac2168444f4de69c27d6384ea2ccf61a49669be5a2fb037ccc1f",
		"Content-Type":"application/json"
	},
	"formParams":{
		"description":"sanath test",
		"url":"http://google.com"
	},
	"checkEvery":30,
	"responseCode":200,		
	"responseTime":800
},

{
	"url":"http://mywebsite.com/v1/data",
	"requestType":"GET",
	"headers":{
		"Authorization":"Bearer ac2168444f4de69c27d6384ea2ccf61a49669be5a2fb037ccc1f",		
	},
	"urlParams":{
		"name":"statusok"
	},
	"checkEvery":300,
	"responseCode":200,		
	"responseTime":800
},

{
	"url":"http://something.com/v1/data",
	"requestType":"DELETE",
	"formParams":{
		"name":"statusok"
	},
	"checkEvery":300,
	"responseCode":200,		
	"responseTime":800
}

```
[Guide to write config.json file](https://github.com/vemonet/statusok/blob/master/Config.md#writing-a-config-file)

[Sample config.json file](https://github.com/vemonet/statusok/blob/master/sample_config.json)

To run the app

```bash
./statusok --config config.json &
```

## Database

Save Requests response time information and error information to your database by adding database details to config file. Currently only Influxdb 0.9.3+ is supported.

You can also add data to your own database.[view details](https://github.com/vemonet/statusok/blob/master/Config.md#save-data-to-any-other-database)

## Notifications

Notifications will be triggered when mean response time is below given response time for a request or when an error is occured . Currently the below clients are supported to receive notifications.For more information on setup [click here](https://github.com/vemonet/statusok/blob/master/Config.md#notifications)

1. [Slack](https://github.com/vemonet/statusok/blob/master/Config.md#slack)
2. [Smtp Email](https://github.com/vemonet/statusok/blob/master/Config.md#e-mail)
3. [Mailgun](https://github.com/vemonet/statusok/blob/master/Config.md#mailgun)
4. [Http EndPoint](https://github.com/vemonet/statusok/blob/master/Config.md#http-endpoint)
5. [Dingding](https://github.com/vemonet/statusok/blob/master/Config.md#dingding)

Adding support to other clients is simple: [view details](https://github.com/vemonet/statusok/blob/master/Config.md#write-your-own-notification-client)

## Running with plain Docker

```bash
docker run -d -v /path/to/config/folder:/config sanathp/statusok
```

> Config folder should contain config file with name `config.json`

## Running with Docker Compose

Change the `config.json` file in the `config/` folder.

Now run it:

```bash
docker-compose up -d
```

Stop all services:

```bash
docker-compose down
```

## Contribution

Contributions are welcomed and greatly appreciated. Create an issue if you find bugs.
Send a pull request if you have written a new feature or fixed an issue .Please make sure to write test cases.

## License
```
Copyright 2015 Sanath Kumar Pasumarthy

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License
```
