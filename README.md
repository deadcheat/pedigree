# Pedigree
![pedigree!!](https://cloud.githubusercontent.com/assets/2797681/26349267/33e4e418-3fea-11e7-89f8-26c0c17386be.gif)
```
__________             .___.__                                  ._.
\______   \  ____    __| _/|__|   ____  _______   ____    ____  | |
 |     ___/_/ __ \  / __ | |  |  / ___\ \_  __ \_/ __ \ _/ __ \ | |
 |    |    \  ___/ / /_/ | |  | / /_/  > |  | \/\  ___/ \  ___/  \|
 |____|     \___  >\____ | |__| \___  /  |__|    \___  > \___  > __
                \/      \/     /_____/               \/      \/  \/

pedigree is a simply request-logging server application.

Usage:
  pedigree [command]

Available Commands:
  help        Help about any command
  logger      startup request-logging server

Flags:
  -h, --help   help for pedigree

Use "pedigree [command] --help" for more information about a command.

```

## Pedigree-Logger

```
__________             .___.__                                           .____                                               ._.
\______   \  ____    __| _/|__|   ____  _______   ____    ____           |    |      ____     ____     ____    ____  _______ | |
 |     ___/_/ __ \  / __ | |  |  / ___\ \_  __ \_/ __ \ _/ __ \   ______ |    |     /  _ \   / ___\   / ___\ _/ __ \ \_  __ \| |
 |    |    \  ___/ / /_/ | |  | / /_/  > |  | \/\  ___/ \  ___/  /_____/ |    |___ (  <_> ) / /_/  > / /_/  >\  ___/  |  | \/ \|
 |____|     \___  >\____ | |__| \___  /  |__|    \___  > \___  >         |_______ \ \____/  \___  /  \___  /  \___  > |__|    __
                \/      \/     /_____/               \/      \/                  \/        /_____/  /_____/       \/          \/

pedigree-logger is one & only function of pedigree.

startup request-logging server

Usage:
  pedigree logger [flags]

Flags:
      --fluent-host string   specify fluentd host default is not set and never access fluentd
      --fluent-port int      specify fluentd port default is not set and never access fluentd
  -h, --help                 help for logger
  -H, --host string          specify hostname, default: localhost (default "localhost")
  -n, --name string          Top-Level object's name that will be logged, default: RequestData (default "RequestData")
  -p, --port int             specify portnum, default: 3000 (default 3000)
  -t, --tag string           Tag name that should be passed to fluentd, default: tracking.request (default "tracking.request")

```

