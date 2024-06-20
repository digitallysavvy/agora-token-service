# How to Build a Token Server for Agora Applications using GoLang #
1[](https://miro.medium.com/max/1400/1*kO_453SIm9v_a9UsVGpjWw.jpeg)

This is a guide on how to build a token server using Golang to generate a token for use with the Agora SDKs. Within the Agora Platform, one layer of security comes in the form of token authentication. A token, for those of you that don’t know, is a dynamic key that is generated using a set of given inputs. Agora’s Platform uses tokens to authenticate users.

Agora offers token security for both its RTC and RTM SDKs. This guide will explain how to build a simple microservice using Golang and the Gin framework to generate an Agora RTC and RTM tokens.

## Pre Requisites ##
- A basic understanding of [Golang](https://golang.org/) _(minimal knowledge needed)_
- An understanding of how web servers function _(minimal knowledge needed)_
- An Agora Developer Account (see: [How To Get Started with Agora](https://www.agora.io/en/blog/how-to-get-started-with-agora?utm_source=medium&utm_medium=blog&utm_campaign=How_to_Build_a_Token_Server_for_Agora_Applications_using_GoLang))

## Project Setup ##
To start, let’s open our terminal and create a new folder for our project and cd into it.
```sh
mkdir agora-token-service
cd agora-token-service
```
Now that the project has been created, let’s initialize the project’s Go module.
```
go mod init agora-token-service
```
Lastly, we’ll use `go get` to add our Gin and Agora dependencies.
```
go get github.com/gin-gonic/gin
go get github.com/AgoraIO-Community/go-tokenbuilder
```

## Build the Gin web server ##
Now that the project is set up, open the folder in your favorite code editor and create the `main.go` file.
![](https://miro.medium.com/max/1400/1*bedBLzVwVkALrJvk5AaJ4Q.png)
Within the `main.go` we’ll start by declaring our package and adding the `main` function.
```go
package main

func main() {

}
```
Next we’ll import the [Gin framework](https://github.com/gin-gonic/gin), create our Gin app, set up a simple `GET` endpoint and set it to listen and serve on `localhost` port `8080`. For the simple endpoint, we’ll set it to take the request context and return a JSON response with a `200` status header.
```go
package main

import (
  "github.com/gin-gonic/gin"
)

func main() {

  api := gin.Default()

  api.GET("/ping", func(c *gin.Context) {
    c.JSON(200, gin.H{
      "message": "pong",
    })
  })

  api.Run(":8080") // listen and serve on localhost:8080
}
```
We are ready to test our server. Go back to the terminal window and run:
```
go run main.go
```
![](https://miro.medium.com/max/1400/1*P9PiY7gPClMQhSzg4g4Atw.png)
To test the endpoint open your web browser and visit:
```
localhost:8080/ping
```
You’ll see the server respond as expected.
```JSON
{"message":"pong"}
```
After we confirm that our endpoint is working, return to the terminal window and use the keyboard command `ctrl c` to terminate the process.

## Generate the Agora Tokens ##
Now that we have our Gin server setup, we are ready to add the functionality to generate the RTC and RTM tokens.

Before we can generate our tokens we need to add our `AppID` and `AppCertificate`. We’ll declare the `appID` and `appCertificate` as _Strings_ in the `global` scope. For this guide, use environment variables to store the project credentials, so we’ll need to retrieve them. Within `main()` we’ll use `os.LookupEnv` to retrieve the environment variables. The `os.LookupEnv` returns a _String_ for the environment variable along with a _boolean_ for whether the variable existed. We’ll use the latter return values to check if the environment is configured correctly. If so we can assign the environment variable values to our global `appID` and `AppCertificate` variables, respectively.

```go
package main

import (
  "log"
  "os"

  "github.com/gin-gonic/gin"
)

var appID, appCertificate string

func main() {

  appIDEnv, appIDExists := os.LookupEnv("APP_ID")
  appCertEnv, appCertExists := os.LookupEnv("APP_CERTIFICATE")

  if !appIDExists || !appCertExists {
    log.Fatal("FATAL ERROR: ENV not properly configured, check APP_ID and APP_CERTIFICATE")
  } else {
    appID = appIDEnv
    appCertificate = appCertEnv
  }

  api := gin.Default()

  api.GET("/ping", func(c *gin.Context) {
    c.JSON(200, gin.H{
      "message": "pong",
    })
  })

  api.Run(":8080")
}
```
Next we will add 3 endpoints, one for RTC tokens, another for RTM tokens, and one that returns both tokens.

The RTC token will require a channel name, the UID, the user role, tokenType to distinguish between string and integer based UIDs, and lastly an expiration time. The RTM endpoint only requires a UID and an expiration time. The dual token endpoint will need to accept the same structure as the RTC token endpoint.

```go
api.GET("rtc/:channelName/:role/:tokenType/:uid/", getRtcToken)
api.GET("rtm/:uid/", getRtmToken)
api.GET("rte/:channelName/:role/:tokenType/:uid/", getBothTokens)
```

To minimize the amount of repeated code, the three functions `getRtcToken`, `getRtmToken`, and `getBothTokens` will call separate functions (`parseRtcParams`/`parseRtmParams`) to validate and extract the values passed to each endpoint. Then each function will use the returned values to generate the tokens and return them as JSON in the response `body`.

RTC tokens can be generated using two types of UIDs (`uint`/`string`), so we’ll use a function (`generateRtcToken`) to wrap the [Agora RTC Token Builder](https://github.com/AgoraIO-Community/go-tokenbuilder/blob/master/rtctokenbuilder/RtcTokenBuilder.go) functions `BuildTokenWithUserAccount`/`BuildTokenWithUID`.

Below is the base template for our token server. We’ll walk through each function and fill in the blanks.
```go
package main

import (
  "log"
  "os"

  "github.com/AgoraIO-Community/go-tokenbuilder/rtctokenbuilder"
  "github.com/gin-gonic/gin"
)

var appID, appCertificate string

func main() {

  appIDEnv, appIDExists := os.LookupEnv("APP_ID")
  appCertEnv, appCertExists := os.LookupEnv("APP_CERTIFICATE")

  if !appIDExists || !appCertExists {
    log.Fatal("FATAL ERROR: ENV not properly configured, check APP_ID and APP_CERTIFICATE")
  } else {
    appID = appIDEnv
    appCertificate = appCertEnv
  }

  api := gin.Default()

  api.GET("/ping", func(c *gin.Context) {
    c.JSON(200, gin.H{
      "message": "pong",
    })
  })

  api.GET("rtc/:channelName/:role/:tokenType/:uid/", getRtcToken)
  api.GET("rtm/:uid/", getRtmToken)
  api.GET("rte/:channelName/:role/:tokenType/:uid/", getBothTokens)

  api.Run(":8080")
}

func getRtcToken(c *gin.Context) {

}

func getRtmToken(c *gin.Context) {

}

func getBothTokens(c *gin.Context) {

}

func parseRtcParams(c *gin.Context) (channelName, tokenType, uidStr string, role rtctokenbuilder.Role, expireTimestamp uint32, err error) {

}

func parseRtmParams(c *gin.Context) (uidStr string, expireTimestamp uint32, err error) {

}

func generateRtcToken(channelName, uidStr, tokenType string, role rtctokenbuilder.Role, expireTimestamp uint32) (rtcToken string, err error) {

}
```

### Build the RTC Token ###
We’ll start with `getRtcToken`. This function takes a reference to the `gin.Context`, using it to call `parseRtcParams` which will extract the required values. Then using the returned values to call `generateRtcToken` to generate the token _String_. We’ll also include a few checks for errors to make sure there weren’t any issues along the way. Lastly we’ll build the response.
```go
func getRtcToken(c *gin.Context) {
  log.Printf("rtc token\n")
  // get param values
  channelName, tokenType, uidStr, role, expireTimestamp, err := parseRtcParams(c)

  if err != nil {
    c.Error(err)
    c.AbortWithStatusJSON(400, gin.H{
      "message": "Error Generating RTC token: " + err.Error(),
      "status":  400,
    })
    return
  }

  rtcToken, tokenErr := generateRtcToken(channelName, uidStr, tokenType, role, expireTimestamp)

  if tokenErr != nil {
    log.Println(tokenErr) // token failed to generate
    c.Error(tokenErr)
    errMsg := "Error Generating RTC token - " + tokenErr.Error()
    c.AbortWithStatusJSON(400, gin.H{
      "status": 400,
      "error":  errMsg,
    })
  } else {
    log.Println("RTC Token generated")
    c.JSON(200, gin.H{
      "rtcToken": rtcToken,
    })
  }
}
```
Next let’s fill in `parseRtcParams`. This function also takes a reference to the gin.Context, which we’ll use to extract the parameters and return them. You’ll notice `parseRtcParams` also returns an error in case we run into any issues we can return an error message.
```go

func parseRtcParams(c *gin.Context) (channelName, tokenType, uidStr string, role rtctokenbuilder.Role, expireTimestamp uint32, err error) {
  // get param values
  channelName = c.Param("channelName")
  roleStr := c.Param("role")
  tokenType = c.Param("tokenType")
  uidStr = c.Param("uid")
  expireTime := c.DefaultQuery("expiry", "3600")

  if roleStr == "publisher" {
    role = rtctokenbuilder.RolePublisher
  } else {
    role = rtctokenbuilder.RoleSubscriber
  }

  expireTime64, parseErr := strconv.ParseUint(expireTime, 10, 64)
  if parseErr != nil {
    // if string conversion fails return an error
    err = fmt.Errorf("failed to parse expiry: %s, causing error: %s", expireTime, parseErr)
  }

  // set timestamps
  expireTimeInSeconds := uint32(expireTime64)
  currentTimestamp := uint32(time.Now().UTC().Unix())
  expireTimestamp = currentTimestamp + expireTimeInSeconds

  return channelName, tokenType, uidStr, role, expireTimestamp, err
}
```
Lastly, we’ll fill in the `generateRtcToken` function. This function takes the channel name, the UID as a _String_, the type of token (`uid` or `userAccount`), the role, and the expire time.

Using these values, the function calls the appropriate [Agora RTC Token Builder](https://github.com/AgoraIO-Community/go-tokenbuilder/blob/master/rtctokenbuilder/RtcTokenBuilder.go) function (`BuildTokenWithUserAccount`/`BuildTokenWithUID`) to generate a token _String_. Once the token builder function returns we’ll first check for errors and if there aren’t any we’ll return the token _String_ value.
```go
func generateRtcToken(channelName, uidStr, tokenType string, role rtctokenbuilder.Role, expireTimestamp uint32) (rtcToken string, err error) {

  if tokenType == "userAccount" {
    log.Printf("Building Token with userAccount: %s\n", uidStr)
    rtcToken, err = rtctokenbuilder.BuildTokenWithUserAccount(appID, appCertificate, channelName, uidStr, role, expireTimestamp)
    return rtcToken, err

  } else if tokenType == "uid" {
    uid64, parseErr := strconv.ParseUint(uidStr, 10, 64)
    // check if conversion fails
    if parseErr != nil {
      err = fmt.Errorf("failed to parse uidStr: %s, to uint causing error: %s", uidStr, parseErr)
      return "", err
    }

    uid := uint32(uid64) // convert uid from uint64 to uint 32
    log.Printf("Building Token with uid: %d\n", uid)
    rtcToken, err = rtctokenbuilder.BuildTokenWithUID(appID, appCertificate, channelName, uid, role, expireTimestamp)
    return rtcToken, err

  } else {
    err = fmt.Errorf("failed to generate RTC token for Unknown Tokentype: %s", tokenType)
    log.Println(err)
    return "", err
  }
}
```
### Build the RTM token ###
Next, let’s move on to `getRtmToken`. Just like the code above, `getRtmToken` takes a reference to the `gin.Context`, uses it to call `parseRtmParams` to extract the required values, and uses the returned values to generate an RTM token. The difference here is that we call the Agora RTM Token builder directly to generate the token, String. We’ll include the error checks to make sure there weren’t any issues, and lastly we’ll build the response.
```go
func getRtmToken(c *gin.Context) {
  log.Printf("rtm token\n")
  // get param values
  uidStr, expireTimestamp, err := parseRtmParams(c)

  if err != nil {
    c.Error(err)
    c.AbortWithStatusJSON(400, gin.H{
      "message": "Error Generating RTC token: " + err.Error(),
      "status":  400,
    })
    return
  }

  rtmToken, tokenErr := rtmtokenbuilder.BuildToken(appID, appCertificate, uidStr, rtmtokenbuilder.RoleRtmUser, expireTimestamp)

  if tokenErr != nil {
    log.Println(err) // token failed to generate
    c.Error(err)
    errMsg := "Error Generating RTM token: " + tokenErr.Error()
    c.AbortWithStatusJSON(400, gin.H{
      "error":  errMsg,
      "status": 400,
    })
  } else {
    log.Println("RTM Token generated")
    c.JSON(200, gin.H{
      "rtmToken": rtmToken,
    })
  }
}
```
Next let’s fill in `parseRtmParams`. This function also takes a reference to the `gin.Context`, then extracts and returns the parameters.
```go
func parseRtmParams(c *gin.Context) (uidStr string, expireTimestamp uint32, err error) {
  // get param values
  uidStr = c.Param("uid")
  expireTime := c.DefaultQuery("expiry", "3600")

  expireTime64, parseErr := strconv.ParseUint(expireTime, 10, 64)
  if parseErr != nil {
    // if string conversion fails return an error
    err = fmt.Errorf("failed to parse expiry: %s, causing error: %s", expireTime, parseErr)
  }

  // set timestamps
  expireTimeInSeconds := uint32(expireTime64)
  currentTimestamp := uint32(time.Now().UTC().Unix())
  expireTimestamp = currentTimestamp + expireTimeInSeconds

  // check if string conversion fails
  return uidStr, expireTimestamp, err
}
```
### Build both RTC and RTM tokens ###
Now that we are able to generate both RTC and RTM tokens with individual server requests, we are going to fill in `getBothTokens` to allow for generating both tokens from a single request. We’ll use code very similar to the `getRtcToken`, except this time we’ll include the RTM token.
```go
func getBothTokens(c *gin.Context) {
  log.Printf("dual token\n")
  // get rtc param values
  channelName, tokenType, uidStr, role, expireTimestamp, rtcParamErr := parseRtcParams(c)

  if rtcParamErr != nil {
    c.Error(rtcParamErr)
    c.AbortWithStatusJSON(400, gin.H{
      "message": "Error Generating RTC token: " + rtcParamErr.Error(),
      "status":  400,
    })
    return
  }
  // generate the rtcToken
  rtcToken, rtcTokenErr := generateRtcToken(channelName, uidStr, tokenType, role, expireTimestamp)
  // generate rtmToken
  rtmToken, rtmTokenErr := rtmtokenbuilder.BuildToken(appID, appCertificate, uidStr, rtmtokenbuilder.RoleRtmUser, expireTimestamp)

  if rtcTokenErr != nil {
    log.Println(rtcTokenErr) // token failed to generate
    c.Error(rtcTokenErr)
    errMsg := "Error Generating RTC token - " + rtcTokenErr.Error()
    c.AbortWithStatusJSON(400, gin.H{
      "status": 400,
      "error":  errMsg,
    })
  } else if rtmTokenErr != nil {
    log.Println(rtmTokenErr) // token failed to generate
    c.Error(rtmTokenErr)
    errMsg := "Error Generating RTC token - " + rtmTokenErr.Error()
    c.AbortWithStatusJSON(400, gin.H{
      "status": 400,
      "error":  errMsg,
    })
  } else {
    log.Println("RTC Token generated")
    c.JSON(200, gin.H{
      "rtcToken": rtcToken,
      "rtmToken": rtmToken,
    })
  }

}
```
## Test the Token Server ##
Let’s go back to our terminal window and run our token server.
```
run main.go
```
Once the server instance is running we’ll see the list of endpoints and the message: `Listening and serving HTTP on :8080`.
![](https://miro.medium.com/max/1400/1*AHC01sizJAlQoq3lHpAY2w.png)
Now that our server instance is running, let’s open our web browser and test. For these tests we’ll try a few variations that omit various query params.

### Test the RTC endpoint ###
We’ll start with the RTC token:
```
http://localhost:8080/rtc/testing/publisher/userAccount/1234/
http://localhost:8080/rtc/testing/publisher/uid/1234/
```
The endpoints will generate a token that can be used in the channel: `testing` by a user with the role of `publisher` and the `UID` (_String_ and _uint_) of `1234`.
```JSON
{
  "rtcToken": "0062ec0d84c41c4442d88ba6f5a2beb828bIADJRwbbO8J93uIDi4J305xNXA0A+pVDTPLPavzwsLW3uAZa8+ij4OObIgDqFTEDoOMyXwQAAQAwoDFfAgAwoDFfAwAwoDFfBAAwoDFf"
}
```
To test this token we can use the [Agora 1:1 Web Demo](https://webdemo.agora.io/agora-web-showcase/examples/Agora-Web-Tutorial-1to1-Web).

### Test the RTM endpoint ###
Next, we’ll test the RTM token:
```
http://localhost:8080/rtm/1234/
```
The endpoints will generate a token that can be used by a user with the UID of `1234` to log into RTM with the given `AppID`.
```JSON
{
  "rtmToken": "0062ec0d84c41c4442d88ba6f5a2beb828bIABSMH0fzaqy7sa0erk8u4Bp6FJ4sO1kQ/o6HCRECBRrzKPg45sAAAAAEAAjAkAEO+cyXwEA6APLozFf"
}
```
To test this token we can use the [Agora RTM Tutorial Demo](https://webdemo.agora.io/agora-web-showcase/examples/Agora-RTM-Tutorial-Web).

### Test the Dual token endpoint ###
We’ll finish our testing with the Dual token endpoint:
```
http://localhost:8080/rte/testing/publisher/userAccount/1234/
http://localhost:8080/rte/testing/publisher/uid/1234/
```
The endpoints will generate both `RTC` and `RTM` tokens that can be used by a user with the UID (_String_ or _uint_) of `1234` and for the Video channel: `testing` with the role of `publisher`.
```JSON
{
  "rtcToken": "0062ec0d84c41c4442d88ba6f5a2beb828bIAD33wY6pO+xp6iBY8mbYz2YtOIiRoTTrzdIPF9DEFlSIwZa8+ij4OObIgAQ6e0EX+UyXwQAAQDvoTFfAgDvoTFfAwDvoTFfBADvoTFf",
  "rtmToken": "0062ec0d84c41c4442d88ba6f5a2beb828bIABbCwQgl2te3rk0MEDZ2xrPoalb37fFhTqmTIbGeWErWaPg45sAAAAAEAD1WwYBX+UyXwEA6APvoTFf"
}
```
To test the tokens we can use the [Agora 1:1 Web Demo](https://webdemo.agora.io/agora-web-showcase/examples/Agora-Web-Tutorial-1to1-Web) for the RTC token and [Agora RTM Tutorial Demo](https://webdemo.agora.io/agora-web-showcase/examples/Agora-RTM-Tutorial-Web) for the RTM token.
![](https://miro.medium.com/max/1400/1*JezQwHNFaHwJFshq50N4Vg.png)
After testing the endpoints, your terminal window will display all the requests.

## Done! ##
And just like that we are done! Thanks for taking the time to read my tutorial and if you have any questions please let me know. If you see any room for improvement feel free to fork the repo and make a pull request!

## Other Resources ##
For more information about the Tokens for Agora.io applications, please take a look at the [Set up Authentication](https://docs.agora.io/en/Agora%20Platform/token?platform=All%20Platforms) guide and [Agora Advanced Guide: How to build a Token]f(https://docs.agora.io/en/Video/token_server_go?platform=Go)(Go).

I also invite you to [join the Agoira.io Developer Slack community](http://bit.ly/2IWexJQ).