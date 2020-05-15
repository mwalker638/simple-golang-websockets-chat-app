# simple-golang-websockets-chat-app

I was looking for a exmaple of deploying and managing a [AWS Lambda](https://aws.amazon.com/lambda/) based app written in [Go](http://golang.org) that is deployed via [AWS SAM](https://aws.amazon.com/serverless/sam/) and uses WebSockets.  The closest I found was [simple-websockets-chat-app](https://github.com/aws-samples/simple-websockets-chat-app), which did all of the above but in Node.js.

This repository represents a fork of [simple-websockets-chat-app](https://github.com/aws-samples/simple-websockets-chat-app) but with the Lambda functions written as a single Go package.

I share incase it might be useful to anyone else looking to go down this path.

There is a single golang function which handles the the logic of the simple WebSocket which uses DynamoDB for its persistent storage.  A SAM template is used for ease of deployment and management.

```
.
├── README.md                <-- This instructions file
├── chatapp                  <-- Golang source code for websocket chatapp
└── template.yaml            <-- SAM template for Lambda Functions and DDB
```


# Deploying to your account

## AWS CLI commands

If you prefer, you can install the [AWS SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html) and use it to package, deploy, and describe your application.  Following are the commands I use:

```
sam build

sam deploy --guided
```

## Testing the chat API

To test the WebSocket API, you can use [wscat](https://github.com/websockets/wscat), an open-source command line tool.


1. [Install NPM](https://www.npmjs.com/get-npm).
2. Install wscat:
``` bash
$ npm install -g wscat
```
3. On the console, connect to your published API endpoint by executing the following command:
``` bash
$ wscat -c wss://{YOUR-API-ID}.execute-api.{YOUR-REGION}.amazonaws.com/{STAGE}
```

4. To test the sendMessage function, send a JSON message like the following example. The Lambda function sends it back using the callback URL: 
``` bash
$ wscat -c wss://{YOUR-API-ID}.execute-api.{YOUR-REGION}.amazonaws.com/prod
connected (press CTRL+C to quit)
> {"message":"sendmessage", "data":"hello world"}
< hello world
```

## License Summary

This sample code is made available under a MIT license. See the LICENSE file.
