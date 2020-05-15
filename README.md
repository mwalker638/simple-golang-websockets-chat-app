# simple-golang-websockets-chat-app

This is the code and template for the simple-golang-websocket-chat-app.  There are three functions contained within the directories and a SAM template that wires them up to a DynamoDB table and provides the minimal set of permissions needed to run the app:

```
.
├── README.md                <-- This instructions file
├── chatapp                  <-- Golang source code for websocket chatapp
└── template.yaml            <-- SAM template for Lambda Functions and DDB
```


It should be noted that this is a modified version of the sample app [simple-golang-websockets-chat-app](https://github.com/aws-samples/simple-websockets-chat-app).  The primary changes were implementing the Lambda code in golang rather then the original Node.js code.


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
