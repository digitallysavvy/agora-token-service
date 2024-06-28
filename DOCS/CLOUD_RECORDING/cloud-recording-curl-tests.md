# CuRL scripts for testing cloud recording service using POST requests

Start (simple)

```bash
curl -X POST http://localhost:8080/cloud_recording/startRecording \
-H "Content-Type: application/json" \
-d '{
  "channelName": "testChannel",
  "sceneMode": "realtime",
  "recordingMode": "mix",
  "excludeResourceIds": []
}'

```

Start (full)

```bash
curl -X POST http://localhost:8080/cloud_recording/startRecording \
-H "Content-Type: application/json" \
-d '{
  "channelName": "testChannel",
  "sceneMode": "realtime",
  "recordingMode": "mix",
  "excludeResourceIds": [],
  "recordingConfig": {
    "channelType": 0,
    "decryptionMode": 1,
    "secret": "your_secret",
    "salt": "your_salt",
    "maxIdleTime": 120,
    "streamTypes": 2,
    "videoStreamType": 0,
    "subscribeAudioUids": ["#allstream#"],
    "unsubscribeAudioUids": [],
    "subscribeVideoUids": ["#allstream#"],
    "unsubscribeVideoUids": [],
    "subscribeUidGroup": 0,
    "streamMode": "individual",
    "audioProfile": 1,
    "transcodingConfig": {
      "width": 640,
      "height": 360,
      "fps": 15,
      "bitrate": 500,
      "maxResolutionUid": "1",
      "layoutConfig": [
        {
          "x_axis": 0,
          "y_axis": 0,
          "width": 640,
          "height": 360,
          "alpha": 1,
          "render_mode": 1
        }
      ]
    }
  }
}'

```

Stop

```bash
curl -X POST http://localhost:8080/cloud_recording/stopRecording \
-H "Content-Type: application/json" \
-d '{
  "cname": "testChannel",
  "UID":"123456789",
  "recordingId":"35BBaeb2a94l95f3a4bcfmkeab1d1fS0;",
  "resourceId":"GFc1iaSt1fBOIr0YcMm-goRnbSnylNl3sssqoZXDCJsZ78dT6UNy4vUwhj-JYra1wbjbWalOBDOZLKpEa0TLLjndRJ0b983mIXMINqm0fqUU_PBJgM3xhn2ip9KTKJs19QOkuMu7bdcvhAgataGbNAGm6_n-3IuhflwqeP06EXH5ZqE3HUQfdOknrHB_r_uuk9t5yw8RAT0oODAEHMaU2kODMcJdGf4oGC27k-8XM_IiJNmW3phzu2IOPvo5nV4YbSrUrxTghrSpi30iFFTdLg",
  "recordingMode": "mix",
  "async_stop": true
}'
```
