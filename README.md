# GoHead
Get interesting http headers from target(s)

## Run
```
$ gohead 

              ______      __  __               __ 
             / ____/___  / / / /__  ____ _____/ / 
            / / __/ __ \/ /_/ / _ \/ __ \/ __  /  
        __ / /_/ / /_/ / __  /  __/ /_/ / /_/ / __
      _/_/ \____/\____/_/ /_/\___/\__,_/\__,_/_/_/
    _/_/___________________________________ _/_/  
  _/_//_____/_____/_____/_____/_____/_____//_/    
 /_/                                     /_/      			  
		
Usage of gohead:
  -exclude string
    	Supply a file of headers to exclude seperated by newlines.
  -outdir string
    	Supply a directory to output the result to. Writes 1 file per supplied target.
  -secrets
    	Search JavaScript files for keys, passwords or secrets (default false)
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
Server: gws
X-Xss-Protection: 0
Set-Cookie: 1P_JAR=2020-04-12-14; expires=Tue, 12-May-2020 14:19:59 GMT; path=/; domain=.google.com; SecureNID=202=DYKO_wTCy_8JUtR9d5W5ljM01awsD46qaPhDM-0p2nce5YfJ6x3yldgk6OBTQ2bxCkC-ccEuj8zLTJD0D370D7m5ANMIMMTv4jg913PruAxLnjdb0vC6y8oi5XS1UhhbVbBb4eY5YQRIpKyeOB6py6yrqSminsckjMFh53CFVGI; expires=Mon, 12-Oct-2020 14:19:59 GMT; path=/; domain=.google.com; HttpOnly
Alt-Svc: quic=":443"; ma=2592000; v="46,43",h3-Q050=":443"; ma=2592000,h3-Q049=":443"; ma=2592000,h3-Q048=":443"; ma=2592000,h3-Q046=":443"; ma=2592000,h3-Q043=":443"; ma=2592000,h3-T050=":443"; ma=2592000
Content-Type: text/html; charset=ISO-8859-1
P3p: CP="This is not a P3P policy! See g.co/p3phelp for more info."
X-Frame-Options: SAMEORIGIN
```
