# REST API 

# Table of contents

<!-- TOC depthFrom:2 depthTo:6 insertAnchor:false orderedList:false updateOnSave:true withLinks:true -->

- [1. One-off requests](#1-one-off-requests)
	- [1.1 Read coil/register (**mbtcp.once.read**)](#11-read-coilregister-mbtcponceread)
	- [1.2 Write coil/register (**mbtcp.once.write**)](#12-write-coilregister-mbtcponcewrite)
	- [1.3 Get TCP connection timeout (**mbtcp.timeout.read**)](#13-get-tcp-connection-timeout-mbtcptimeoutread)
	- [1.4 Set TCP connection timeout (**mbtcp.timeout.update**)](#14-set-tcp-connection-timeout-mbtcptimeoutupdate)
- [2. Polling requests](#2-polling-requests)
	- [2.1 Add poll request (**mbtcp.poll.create**)](#21-add-poll-request-mbtcppollcreate)
	- [2.2 Update poll request interval (**mbtcp.poll.update**)](#22-update-poll-request-interval-mbtcppollupdate)
	- [2.3 Read poll request status (**mbtcp.poll.read**)](#23-read-poll-request-status-mbtcppollread)
	- [2.4 Delete poll request (**mbtcp.poll.delete**)](#24-delete-poll-request-mbtcppolldelete)
	- [2.5 Enable/Disable poll request (**mbtcp.poll.toggle**)](#25-enabledisable-poll-request-mbtcppolltoggle)
	- [2.6 Read all poll requests status (**mbtcp.polls.read**)](#26-read-all-poll-requests-status-mbtcppollsread)
	- [2.7 Delete all poll requests (**mbtcp.polls.delete**)](#27-delete-all-poll-requests-mbtcppollsdelete)
	- [2.8 Enable/Disable all poll requests (**mbtcp.polls.toggle**)](#28-enabledisable-all-poll-requests-mbtcppollstoggle)
	- [2.9 Import poll requests (**mbtcp.polls.import**)](#29-import-poll-requests-mbtcppollsimport)
	- [2.10 Export poll requests (**mbtcp.polls.export**)](#210-export-poll-requests-mbtcppollsexport)
	- [2.11 Read history (**mbtcp.poll.history**)](#211-read-history-mbtcppollhistory)
- [3. Filter requests](#3-filter-requests)
	- [3.1 Add filter request (**mbtcp.filter.create**)](#31-add-filter-request-mbtcpfiltercreate)
	- [3.2 Update filter request (**mbtcp.filter.update**)](#32-update-filter-request-mbtcpfilterupdate)
	- [3.3 Read filter request status (**mbtcp.filter.read**)](#33-read-filter-request-status-mbtcpfilterread)
	- [3.4 Delete filter request (**mbtcp.filter.delete**)](#34-delete-filter-request-mbtcpfilterdelete)
	- [3.5 Enable/Disable filter request (**mbtcp.filter.toggle**)](#35-enabledisable-filter-request-mbtcpfiltertoggle)
	- [3.6 Read all filter requests (**mbtcp.filters.read**)](#36-read-all-filter-requests-mbtcpfiltersread)
	- [3.7 Delete all filter requests (**mbtcp.filters.delete**)](#37-delete-all-filter-requests-mbtcpfiltersdelete)
	- [3.8 Enable/Disable all filter requests (**mbtcp.filters.toggle**)](#38-enabledisable-all-filter-requests-mbtcpfilterstoggle)
	- [3.9 Import filter requests (**mbtcp.filters.import**)](#39-import-filter-requests-mbtcpfiltersimport)
	- [3.10 Export filter requests (**mbtcp.filters.export**)](#310-export-filter-requests-mbtcpfiltersexport)
- [4. Authentication](#4-authentication)
	- [4.1 Login](#41-login)
	- [4.2 Logout](#42-logout)
	- [4.3 Update username](#43-update-username)
	- [4.4 Update password](#44-update-password)
	- [4.5 Reset authetication setting](#45-reset-authetication-setting)

<!-- /TOC -->

## 1. One-off requests

### 1.1 Read coil/register (**mbtcp.once.read**)

>|params |description            |In            |type          |range        |example     |required                                 |
>|:------|:----------------------|:-------------|:-------------|:------------|:-----------|:----------------------------------------|
>|fc     |function code          |path          |integer       |[1,4]        |1           |:heavy_check_mark:                       |
>|ip     |ip address             |query         |string        |-            |127.0.0.1   |:heavy_check_mark:                       |
>|port   |port number            |query         |string        |[1,65535]    |502         |default: 502                             |
>|slave  |slave id               |query         |integer       |[1, 253]     |1           |:heavy_check_mark:                       |
>|addr   |register start address |query         |integer       |-            |23          |:heavy_check_mark:                       |
>|len    |register length        |query         |integer       |-            |20          |default: 1                               |
>|type   |Data type              |query         |integer       |[1,8]        | see below  |default: 1, **fc 3, 4 only**             |
>|order  |Endian                 |query         |integer       |[1,4]        | see below  |default: 1, **fc 3, 4 and type 4~8 only**|
>|range  |Scale range            |query         |number        |-            | see below  |fc 3, 4 and type 3 only                  |
>|bytes  |response byte array    |response body |-             |[0XAB, 0X12] | see below  |fc 3, 4 and type 2~8 only                |
>|status |response status        |response body |string        |-            |"ok"        |:heavy_check_mark:                       |
>|data   |response value         |response body |integer array |-            |[1,0,24,1]  |if success                               |


- Verb: **GET**
- URI: /api/mb/tcp/fc/**{fc}**
- Query: ?ip=**{ip}**&port=**{port}**&slave=**{slave}**&addr=**{addr}**&len=**{len}**&type=**{type}**
- Example: Bits read (FC1, FC2)
    - Request
        - port: 502
        - len: 1
        - endpoint:
        ```Bash
        /api/mb/tcp/fc/1?ip=192.168.3.2&port=503&slave=1&addr=10&len=4
        ```

    - Response
        - Success
        ```JavaScript
        {
            "status": "ok",
            "data": [0,1,0,1,0,1]
        }
        ```

        - Fail
        ```JavaScript
        {
            "status": "timeout"
        }
        ```

- Examples: Register read (FC3, FC4) - type 1, 2 (raw)
    - Request
        - endpoint:

        ```Bash
        /api/mb/tcp/fc/3?ip=192.168.3.2&port=503&slave=1&addr=10&len=10&type=1
        ```

    - Response
        - Success - type 1 (RegisterArray):
        ```JavaScript
        {
            "status": "ok",
            "type": 1,
            "bytes": [0XFF, 0X34, 0XAB],
            "data": [255, 1234, 789]
        }
        ```

        - Success - type 2 (Hex String):
        ```JavaScript
        {
            "status": "ok",
            "type": 2,
            "bytes": [0XFF, 0X34, 0XAB],
            "data": "112C004F12345678"
        }
        ```

        - Fail:
        ```JavaScript
        {
            "type": 2,
            "status": "timeout"
        }
        ```

- Examples: Register read (FC3, FC4) - type 3 (scale)
    - Request
        - endpoint:

        ```Bash
        /api/mb/tcp/fc/3?ip=192.168.3.2&port=503&slave=1&addr=10&len=4&type=3&a=0&b=65535&c=100&d=500
        ```

    - Response
        - Success:
        ```JavaScript
        {
            "status": "ok",
            "type": 3,
            "bytes": [0XAB, 0X12, 0XCD, 0XED, 0X12, 0X34],
            "data": [22.34, 33.12, 44.56]
        }
        ```

        - Fail:
        ```JavaScript
        {
            "type": 3,
            "bytes": null,
            "status": "timeout"
        }
        ```

- Examples: Register read (FC3, FC4) - type 4, 5 (16-bit)
    - Request
        - endpoint:

        ```Bash
        /api/mb/tcp/fc/3?ip=192.168.3.2&port=503&slave=1&addr=10&len=4&type=4&order=1
        ```

    - Response
        - Success:
        ```JavaScript
        {
            "status": "ok",
            "type": 4,
            "bytes": [0XAB, 0X12, 0XCD, 0XED, 0X12, 0X34],
            "data": [255, 1234, 789]
        }
        ```

        - Fail:
        ```JavaScript
        {
            "type": 4,
            "bytes": null,
            "status": "timeout"
        }
        ```

- Examples: Register read (FC3, FC4) - type 6, 7, 8 (32-bit)
    - Request
        - endpoint:

        ```Bash
        /api/mb/tcp/fc/3?ip=192.168.3.2&port=503&slave=1&addr=10&len=4&type=6&order=3
        ```
    - Response
        - Success - type 6, 7 (UInt32, Int32):
        ```JavaScript
        {
            "status": "ok",
            "type": 6,
            "bytes": [0XAB, 0X12, 0XCD, 0XED, 0X12, 0X34],
            "data": [255, 1234, 789]
        }
        ```

        - Success - type 8 (Float32):
        ```JavaScript
        {
            "status": "ok",
            "type": 8,
            "bytes": [0XAB, 0X12, 0XCD, 0XED, 0X12, 0X34],
            "data": [22.34, 33.12, 44.56]
        }
        ```

        - Fail:
        ```JavaScript
        {
            "type": 8,
            "bytes": null,
            "status": "timeout"
        }
        ```
---

### 1.2 Write coil/register (**mbtcp.once.write**)

|params |description            |In            |type          |range              |example     |required          |
|:------|:----------------------|:-------------|:-------------|:------------------|:-----------|:-----------------|
|fc     |function code          |path          |integer       |**[5, 6, 15, 16]** |5           |:heavy_check_mark:|
|ip     |ip address             |request body  |string        |-                  |127.0.0.1   |:heavy_check_mark:|
|port   |port number            |request body  |string        |[1,65535]          |502         | default: 502     |
|slave  |slave id               |request body  |integer       |[1, 253]           |1           |:heavy_check_mark:|
|addr   |register start address |request body  |integer       |-                  |23          |:heavy_check_mark:|
|len    |register length        |request body  |integer       |-                  |20          | default: 1       |
|data   |data to be write       |request body  |integer array |                   |[1,0,24,1]  |if success        |
|status |response status        |response body |string        |-                  |"ok"        |:heavy_check_mark:|


- Verb: **POST**
- URI: /api/mb/tcp/fc/**{fc}**
- Example: write **single** coil/register

    - **Request**
        - port: 502
        - endpoint:
        ```Bash
        /api/mb/tcp/fc/5
        ```
        - body:
        ```JavaScript
        {
            "ip": "192.168.3.2",
            "slave": 22,
            "addr": 80,
            "data": 1
        }
        ```

    - **Response**

        - success:
        ```JavaScript
        {
            "status": "ok",
        }
        ```

        - fail:
        ```JavaScript
        {
            "status": "timeout"
        }
        ```

- Example: write **multiple** coils/registers

    - **Request**
        - endpoint:
        ```Bash
        /api/mb/tcp/fc/5
        ```
        - body:
        ```JavaScript
        {
            "ip": "192.168.3.2",
            "port": "503",
            "slave": 22,
            "addr": 80,
            "len": 4,
            "data": [1, 2, 3, 5]
        }
        ```

    - **Response**
        - success:
        ```JavaScript
        {
            "status": "ok",
        }
        ```

        - fail:
        ```JavaScript
        {
            "status": "timeout"
        }
        ```

---

### 1.3 Get TCP connection timeout (**mbtcp.timeout.read**)

|params |description            |In            |type          |range     |example     |required          |
|:------|:----------------------|:-------------|:-------------|:---------|:-----------|:-----------------|
|status |response status        |response body |string        |-         |"ok"        |:heavy_check_mark:|
|timeout|timeout in usec        |response body |integer       |[200000,~)|210000      |if success        |

- Verb: **GET**
- URI: /api/mb/tcp/timeout
- Example:

    - **Request**
        - endpoint:
        ```Bash
        /api/mb/tcp/timeout
        ```

    - **Response**
        - success:
        ```JavaScript
        {
            "status": "ok",
            "timeout": 210000
        }
        ```

        - fail:
        ```JavaScript
        {
            "status": "timeout"
        }
        ```
---

### 1.4 Set TCP connection timeout (**mbtcp.timeout.update**)

|params |description            |In            |type          |range     |example     |required          |
|:------|:----------------------|:-------------|:-------------|:---------|:-----------|:-----------------|
|timeout|timeout in usec        |request  body |integer       |[200000,~)|210000      |:heavy_check_mark:|
|status |response status        |response body |string        |-         |"ok"        |:heavy_check_mark:|

- Verb: **POST**
- URI: /api/mb/tcp/timeout
- Example:

    - **Request**
        - endpoint:
        ```Bash
        /api/mb/tcp/timeout
        ```
        
        - body:
        ```JavaScript
        {
            "timeout": 210000
        }
        ```

    - **Response**
        - success:
        ```JavaScript
        {
            "status": "ok"
        }
        ```

        - fail:
        ```JavaScript
        {
            "status": "timeout"
        }
        ```
---

## 2. Polling requests

### 2.1 Add poll request (**mbtcp.poll.create**)

|params       |description            |In            |type          |range                  |example     |required          |
|:------------|:----------------------|:-------------|:-------------|:----------------------|:-----------|:-----------------|
|**name**     |request/sensor name    |path          |string        |no space and **unique**|led_1       |:heavy_check_mark:|
|fc           |function code          |request body  |integer       |**[1,4]**              |1           |:heavy_check_mark:|
|ip           |ip address             |request body  |string        |-                      |127.0.0.1   |:heavy_check_mark:|
|port         |port number            |request body  |string        |[1,65535]              |502         | default: 502     |
|slave        |slave id               |request body  |integer       |[1, 253]               |1           |:heavy_check_mark:|
|addr         |register start address |request body  |integer       |-                      |23          |:heavy_check_mark:|
|len          |register length        |request body  |integer       |-                      |20          |default: 1        |
|**interval** |polling interval in sec|request body  |integer       |[1~)                   |[1,0,24,1]  |:heavy_check_mark:|
|status       |response status        |response body |string        |-                      |"ok"        |:heavy_check_mark:|

- Verb: **POST**
- URI: /api/mb/tcp/poll/**{name}**
- Example:

    - **Request**
        - endpoint:
        ```Bash
        /api/mb/tcp/poll/led_1
        ```
        
        - body:
        ```JavaScript
        {
            "fc": 1,
            "ip": "192.168.3.2",
            "port": "502",
            "slave": 22,
            "addr": 250,
            "len": 10,
            "interval" : 3
        }
        ```

    - **Response**
        - success:
        ```JavaScript
        {
            "status": "ok"
        }
        ```

        - fail:
        ```JavaScript
        {
            "status": "timeout"
        }
        ```

---

### 2.2 Update poll request interval (**mbtcp.poll.update**)

|params       |description            |In            |type          |range                  |example     |required          |
|:------------|:----------------------|:-------------|:-------------|:----------------------|:-----------|:-----------------|
|**name**     |request/sensor name    |path          |string        |no space and **unique**|led_1       |:heavy_check_mark:|
|**interval** |polling interval in sec|request body  |integer       |[1~)                   |[1,0,24,1]  |:heavy_check_mark:|
|status       |response status        |response body |string        |-                      |"ok"        |:heavy_check_mark:|

- Verb: **PUT**
- URI: /api/mb/tcp/poll/**{name}**
- Example:

    - **Request**
        - endpoint:
        ```Bash
        /api/mb/tcp/poll/led_1

        ```
        
        - body:
        ```JavaScript
        {
            "interval" : 3
        }
        ```

    - **Response**
        - success:
        ```JavaScript
        {
            "status": "ok"
        }
        ```

        - fail:
        ```JavaScript
        {
            "status": "timeout"
        }
        ```
---

### 2.3 Read poll request status (**mbtcp.poll.read**)

|params       |description            |In            |type          |range                  |example     |required          |
|:------------|:----------------------|:-------------|:-------------|:----------------------|:-----------|:-----------------|
|**name**     |request/sensor name    |path          |string        |no space and **unique**|led_1       |:heavy_check_mark:|
|fc           |function code          |response body |integer       |**[1,4]**              |1           |if success        |
|ip           |ip address             |response body |string        |-                      |127.0.0.1   |if success        |
|port         |port number            |response body |string        |[1,65535]              |502         |if success        |
|slave        |slave id               |response body |integer       |[1, 253]               |1           |if success        |
|addr         |register start address |response body |integer       |-                      |23          |if success        |
|len          |register length        |response body |integer       |-                      |20          |if success        |
|**interval** |polling interval in sec|response body |integer       |[1~)                   |[1,0,24,1]  |if success        |
|**enabled**  |polling enabled flag   |response body |boolean       |true, false            |true        |if success        |
|status       |response status        |response body |string        |-                      |"ok"        |:heavy_check_mark:|

- Verb: **GET**
- URI: /api/mb/tcp/poll/**{name}**
- Example:

    - **Request**
        - endpoint:
        ```Bash
        /api/mb/tcp/poll/led_1
        ```

    - **Response**
        - success:
        ```JavaScript
        {
            "fc": 1,
            "ip": "192.168.3.2",
            "port": "502",
            "slave": 22,
            "addr": 250,
            "len": 10,
            "interval" : 3,
            "status": "ok",
            "enabled": true
        }
        ```

        - fail:
        ```JavaScript
        {
            "status": "not exist"
        }
        ```
---

### 2.4 Delete poll request (**mbtcp.poll.delete**)

|params       |description            |In            |type          |range                  |example     |required          |
|:------------|:----------------------|:-------------|:-------------|:----------------------|:-----------|:-----------------|
|**name**     |request/sensor name    |path          |string        |no space and **unique**|led_1       |:heavy_check_mark:|
|status       |response status        |response body |string        |-                      |"ok"        |:heavy_check_mark:|

- Verb: **DELETE**
- URI: /api/mb/tcp/poll/**{name}**
- Example:

    - **Request**
        - endpoint:
        ```Bash
        /api/mb/tcp/poll/led_1
        ```
        - body: **No payload**

    - **Response**
        - success:
        ```JavaScript
        {
            "status": "ok"
        }
        ```

        - fail:
        ```JavaScript
        {
            "status": "timeout"
        }
        ```
---

### 2.5 Enable/Disable poll request (**mbtcp.poll.toggle**)

|params       |description            |In            |type          |range                  |example     |required          |
|:------------|:----------------------|:-------------|:-------------|:----------------------|:-----------|:-----------------|
|**name**     |request/sensor name    |path          |string        |no space and **unique**|led_1       |:heavy_check_mark:|
|**enabled**  |polling enabled flag   |request  body |boolean       |true, false            |true        |:heavy_check_mark:|
|status       |response status        |response body |string        |-                      |"ok"        |:heavy_check_mark:|


- Verb: **POST**
- URI: /api/mb/tcp/poll/**{name}**/toggle
- Example:

    - **Request**
        - endpoint:
        ```bash
        /api/mb/tcp/poll/led_1/toggle
        ```
        - body:
        ```JavaScript
        {
            "enabled": true
        }
        ```

    - **Response**
        - success:
        ```JavaScript
        {
            "status": "ok"
        }
        ```

        - fail:
        ```JavaScript
        {
            "status": "timeout"
        }
        ```
---

### 2.6 Read all poll requests status (**mbtcp.polls.read**)

|params       |description            |In            |type          |range                  |example     |required          |
|:------------|:----------------------|:-------------|:-------------|:----------------------|:-----------|:-----------------|
|polls        |request object array   |response body |object array  |-                      |-           |if success        |
|**name**     |request/sensor name    |response body |string        |no space and **unique**|led_1       |if success        |
|fc           |function code          |response body |integer       |**[1,4]**              |1           |if success        |
|ip           |ip address             |response body |string        |-                      |127.0.0.1   |if success        |
|port         |port number            |response body |string        |[1,65535]              |502         |if success        |
|slave        |slave id               |response body |integer       |[1, 253]               |1           |if success        |
|addr         |register start address |response body |integer       |-                      |23          |if success        |
|len          |register length        |response body |integer       |-                      |20          |if success        |
|**interval** |polling interval in sec|response body |integer       |[1~)                   |[1,0,24,1]  |if success        |
|**enabled**  |polling enabled flag   |response body |boolean       |true, false            |true        |if success        |
|status       |response status        |response body |string        |-                      |"ok"        |:heavy_check_mark:|

- Verb: **GET**
- URI: /api/mb/tcp/polls
- Example:

    - **Request**
        - endpoint:
        ```Bash
        /api/mb/tcp/polls
        ```

    - **Response**
        - success:
            ```JavaScript
            {
                "status": "ok",
                "polls": [
                    {
                        "name": "led_1",
                        "fc": 1,
                        "ip": "192.168.3.2",
                        "port": "502",
                        "slave": 22,
                        "addr": 250,
                        "len": 10,
                        "interval" : 3,
                        "status": "ok",
                        "enabled": true
                    },
                    {
                        "name": "led_2",
                        "fc": 1,
                        "ip": "192.168.3.2",
                        "port": "502",
                        "slave": 22,
                        "addr": 250,
                        "len": 10,
                        "interval" : 3,
                        "status": "ok",
                        "enabled": true
                    }]
            }
            ```

        - fail:
            ```JavaScript
            {
                "status": "timeout"
            }
            ```
---

### 2.7 Delete all poll requests (**mbtcp.polls.delete**)

|params       |description            |In            |type          |range      |example     |required          |
|:------------|:----------------------|:-------------|:-------------|:----------|:-----------|:-----------------|
|status       |response status        |response body |string        |-          |"ok"        |:heavy_check_mark:|

- Verb: **DELETE**
- URI: /api/mb/tcp/polls
- Example:
    - **Request**
        - endpoint:
        ```Bash
        /api/mb/tcp/polls
        ```

    - **Response**
        - success:
        ```JavaScript
        {
            "status": "ok"
        }
        ```

        - fail:
        ```JavaScript
        {
            "status": "timeout"
        }
        ```
---

### 2.8 Enable/Disable all poll requests (**mbtcp.polls.toggle**)

|params       |description            |In            |type          |range      |example     |required          |
|:------------|:----------------------|:-------------|:-------------|:----------|:-----------|:-----------------|
|**enabled**  |polling enabled flag   |request  body |boolean       |true, false|true        |:heavy_check_mark:|
|status       |response status        |response body |string        |-          |"ok"        |:heavy_check_mark:|


- Verb: **POST**
- URI: /api/mb/tcp/polls/toggle
- Example:

    - **Request**
        - endpoint:
        ```Bash
        /api/mb/tcp/polls/toggle
        ```
        - body:
        ```JavaScript
        {
            "enabled": true
        }
        ```
    - **Response**
        - success:
        ```JavaScript
        {
            "status": "ok"
        }
        ```

        - fail:
        ```JavaScript
        {
            "status": "timeout"
        }
        ```
---

### 2.9 Import poll requests (**mbtcp.polls.import**)

**TODO**
:heavy_exclamation_mark: This API should be modified to load file (filename required).

|params       |description            |In            |type          |range                  |example     |required          |
|:------------|:----------------------|:-------------|:-------------|:----------------------|:-----------|:-----------------|
|polls        |request object array   |request body  |object array  |-                      |-           |:heavy_check_mark:|
|**name**     |request/sensor name    |request body  |string        |no space and **unique**|led_1       |:heavy_check_mark:|
|fc           |function code          |request body  |integer       |**[1,4]**              |1           |:heavy_check_mark:|
|ip           |ip address             |request body  |string        |-                      |127.0.0.1   |:heavy_check_mark:|
|port         |port number            |request body  |string        |[1,65535]              |502         |:heavy_check_mark:|
|slave        |slave id               |request body  |integer       |[1, 253]               |1           |:heavy_check_mark:|
|addr         |register start address |request body  |integer       |-                      |23          |:heavy_check_mark:|
|len          |register length        |request body  |integer       |-                      |20          |:heavy_check_mark:|
|**interval** |polling interval in sec|request body  |integer       |[1~)                   |[1,0,24,1]  |:heavy_check_mark:|
|**enabled**  |polling enabled flag   |request body  |boolean       |true, false            |true        |:heavy_check_mark:|
|status       |response status        |response body |string        |-                      |"ok"        |:heavy_check_mark:|

- Verb: **POST**
- URI: /api/mb/tcp/polls/import
- Example:

    - **Request**
        - endpoint:
        ```Bash
        /api/mb/tcp/polls/import
        ```
        - body:
        ```JavaScript
        {
            "polls": [
                {
                    "name": "led_1",
                    "fc": 1,
                    "ip": "192.168.3.2",
                    "port": "502",
                    "slave": 22,
                    "addr": 250,
                    "len": 10,
                    "interval" : 3,
                    "status": "ok",
                    "enabled": true
                },
                {
                    "name": "led_2",
                    "fc": 1,
                    "ip": "192.168.3.2",
                    "port": "502",
                    "slave": 22,
                    "addr": 250,
                    "len": 10,
                    "interval" : 3,
                    "status": "ok",
                    "enabled": true
                }]
        }
        ```

    - **Response**
        - success:
        ```JavaScript
        {
            "status": "ok"
        }
        ```

        - fail:
        ```JavaScript
        {
            "status": "timeout"
        }
        ```

---

### 2.10 Export poll requests (**mbtcp.polls.export**)

**TODO**
:heavy_exclamation_mark: This API should be modified to save file (filename required).

|params       |description            |In             |type          |range                  |example     |required          |
|:------------|:----------------------|:--------------|:-------------|:----------------------|:-----------|:-----------------|
|polls        |request object array   |response body  |object array  |-                      |-           |if success        |
|**name**     |request/sensor name    |response body  |string        |no space and **unique**|led_1       |if success        |
|fc           |function code          |response body  |integer       |**[1,4]**              |1           |if success        |
|ip           |ip address             |response body  |string        |-                      |127.0.0.1   |if success        |
|port         |port number            |response body  |string        |[1,65535]              |502         |if success        |
|slave        |slave id               |response body  |integer       |[1, 253]               |1           |if success        |
|addr         |register start address |response body  |integer       |-                      |23          |if success        |
|len          |register length        |response body  |integer       |-                      |20          |if success        |
|**interval** |polling interval in sec|response body  |integer       |[1~)                   |[1,0,24,1]  |if success        |
|**enabled**  |polling enabled flag   |response body  |boolean       |true, false            |true        |if success        |
|status       |response status        |response body  |string        |-                      |"ok"        |:heavy_check_mark:|

- Verb: **GET**
- URI: /api/mb/tcp/polls/export
- Example:

    - **Request**
        - endpoint:
        ```Bash
        /api/mb/tcp/polls/export
        ```
    - **Response**
        - success:
        ```JavaScript
        {
            "status": "ok",
            "polls": [
                {
                    "name": "led_1",
                    "fc": 1,
                    "ip": "192.168.3.2",
                    "port": "502",
                    "slave": 22,
                    "addr": 250,
                    "len": 10,
                    "interval" : 3,
                    "status": "ok",
                    "enabled": true
                },
                {
                    "name": "led_2",
                    "fc": 1,
                    "ip": "192.168.3.2",
                    "port": "502",
                    "slave": 22,
                    "addr": 250,
                    "len": 10,
                    "interval" : 3,
                    "status": "ok",
                    "enabled": true
                }]
        }
        ```

        - fail:
        ```JavaScript
        {
            "status": "timeout"
        }
        ```

---

### 2.11 Read history (**mbtcp.poll.history**)

|params   |description            |In            |type          |range                  |example     |required          |
|:--------|:----------------------|:-------------|:-------------|:----------------------|:-----------|:-----------------|
|**name** |request/sensor name    |path          |string        |no space and **unique**|led_1       |:heavy_check_mark:|
|status   |response status        |response body |string        |-                      |"ok"        |:heavy_check_mark:|
|data(1)  |outer data             |response body |object array  |-                      |-           |if success        |
|data(2)  |inner data             |response body |integer array |                       |[1,0,24,1]  |if success        |
|ts       |time stamp             |response body |integer       |-                      |-           |if success        |

- Verb: **GET**
- URI: /api/mb/tcp/poll/**{name}**/logs
- Example:

    - **Request**
        - endpoint:
        ```Bash
        /api/mb/tcp/poll/led_1/logs
        ```

    - **Response**
        - success (len=1):
        ```JavaScript
        {
            "status": "ok",
            "data":[{"data": [1], "ts": 2012031203},
                    {"data": [0], "ts": 2012031205},
                    {"data": [1], "ts": 2012031207}]        
        }
        ```
        
        - success (len=n):
        ```JavaScript
        {
            "status": "ok",
            "data":[{"data": [1,0,1], "ts": 2012031203},
                    {"data": [1,1,1], "ts": 2012031205},
                    {"data": [0,0,1], "ts": 2012031207}]        
        }
        ```

        - fail:
        ```JavaScript
        {
            "status": "not exist"
        }
        ```

---

## 3. Filter requests

### 3.1 Add filter request (**mbtcp.filter.create**)

- Verb: **POST**
- URI: /api/mb/tcp/filter/**{name}**

---

### 3.2 Update filter request (**mbtcp.filter.update**)

- Verb: **PUT**
- URI: /api/mb/tcp/filter/**{name}**

---

### 3.3 Read filter request status (**mbtcp.filter.read**)

- Verb: **GET**
- URI: /api/mb/tcp/filter/**{name}**

---

### 3.4 Delete filter request (**mbtcp.filter.delete**)

- Verb: **DELETE**
- URI: /api/mb/tcp/filter/**{name}**

---

### 3.5 Enable/Disable filter request (**mbtcp.filter.toggle**)

- Verb: **POST**
- URI: /api/mb/tcp/filter/**{name}**/toggle

---

### 3.6 Read all filter requests (**mbtcp.filters.read**)

- Verb: **GET**
- URI: /api/mb/tcp/filters

---

### 3.7 Delete all filter requests (**mbtcp.filters.delete**)

- Verb: **DELETE**
- URI: /api/mb/tcp/filters

---

### 3.8 Enable/Disable all filter requests (**mbtcp.filters.toggle**)

- Verb: **POST**
- URI: /api/mb/tcp/filters/toggle

---

### 3.9 Import filter requests (**mbtcp.filters.import**)

- Verb: **POST**
- URI: /api/mb/tcp/filters/import

---

### 3.10 Export filter requests (**mbtcp.filters.export**)

- Verb: **GET**
- URI: /api/mb/tcp/filters/export

---

## 4. Authentication

### 4.1 Login

- Verb: **POST**
- URI: /api/auth/login

**TODO**

---

### 4.2 Logout

- Verb: **GET**
- URI: /api/auth/logout

---

### 4.3 Update username

---

### 4.4 Update password

---

### 4.5 Reset authetication setting

**TODO**

---