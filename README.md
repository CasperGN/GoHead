# GoHead
Get interesting http headers from target(s)

## Run
```
$ ./gohead
Usage of ./gohead:
  -silent
        Print header (default false).
  -target string
        Supply single target for probing.
  -targets string
        Supply a file of targets seperated by newlines.
  -threads int
        Number of threads (default 5)
$ ./gohead -target https://google.com

              ______      __  __               __
             / ____/___  / / / /__  ____ _____/ /
            / / __/ __ \/ /_/ / _ \/ __ \/ __  /
        __ / /_/ / /_/ / __  /  __/ /_/ / /_/ / __
      _/_/ \____/\____/_/ /_/\___/\__,_/\__,_/_/_/
    _/_/___________________________________ _/_/
  _/_//_____/_____/_____/_____/_____/_____//_/
 /_/                                     /_/

https://google.com
Date: [Sun, 12 Apr 2020 13:04:37 GMT]
Expires: [-1]
Cache-Control: [private, max-age=0]
Content-Type: [text/html; charset=ISO-8859-1]
X-Frame-Options: [SAMEORIGIN]
Alt-Svc: [quic=":443"; ma=2592000; v="46,43",h3-Q050=":443"; ma=2592000,h3-Q049=":443"; ma=2592000,h3-Q048=":443"; ma=2592000,h3-Q046=":443"; ma=2592000,h3-Q043=":443"; ma=2592000,h3-T050=":443"; ma=2592000]
P3p: [CP="This is not a P3P policy! See g.co/p3phelp for more info."]
Server: [gws]
X-Xss-Protection: [0]
Set-Cookie: [1P_JAR=2020-04-12-13; expires=Tue, 12-May-2020 13:04:37 GMT; path=/; domain=.google.com; Secure NID=202=rOEOReHfj0DD7NWt57liHrY0BXd5GVbn6BhLRqIHNTZ3c_aWg2MSOyv8DuDzaHHDM-cBXAKumtPeWcaINnsSRnuObJuG0bMIOf6sDFwdN--Y2aMibnZj8LDqF__YCie0R8OkOk5G-mSIuIkmrhmo76ag4dg2mdM4hXu1TKCCCrE; expires=Mon, 12-Oct-2020 13:04:37 GMT; path=/; domain=.google.com; HttpOnly]
```