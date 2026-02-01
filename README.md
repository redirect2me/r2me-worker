# redirect2me worker node [<img alt="redirect2me Logo" src="https://www.redirect2.me/favicon.svg" height="90" align="right" />](https://www.redirect2.me/)

Server that handles the actual redirects.

## Usage

Build with:
```
go build -o server ./cmd/server
```

You can see the available command line options by running:
```
./server --help
```

See [`run.sh`](run.sh) for how I run the server during development.  You can test the server locally with:
```
curl --header "Host: redirect2me" --header "X-Redirect2me-Debug: 1" http://localhost:4000
```

See [pkg/README.md](pkg/README.md) for how I packaging and run it in production

## License

Copyright (c) 2018-2026 by Andrew Marcuse.  All Rights Reserved.

## Credits

[![bash](https://www.vectorlogo.zone/logos/gnu_bash/gnu_bash-ar21.svg)](https://www.gnu.org/software/bash/ "scripting")
[![certmagic](https://www.vectorlogo.zone/logos/github_mholt_certmagic/github_mholt_certmagic-ar21.svg)](https://github.com/mholt/certmagic "Certificate management")
[![Digital Ocean](https://www.vectorlogo.zone/logos/digitalocean/digitalocean-ar21.svg)](https://m.do.co/c/976f479b2317 "Hosting")
[![Git](https://www.vectorlogo.zone/logos/git-scm/git-scm-ar21.svg)](https://git-scm.com/ "Version control")
[![Github](https://www.vectorlogo.zone/logos/github/github-ar21.svg)](https://github.com/ "Code hosting")
[![golang](https://www.vectorlogo.zone/logos/golang/golang-ar21.svg)](https://golang.org/ "Programming language")
[![NodePing](https://www.vectorlogo.zone/logos/nodeping/nodeping-ar21.svg)](https://nodeping.com?rid=201109281250J5K3P "Uptime monitoring")
[![Ubuntu](https://www.vectorlogo.zone/logos/ubuntu/ubuntu-ar21.svg)](https://www.ubuntu.com/ "Server operating system")
[![water.css](https://www.vectorlogo.zone/logos/netlifyapp_watercss/netlifyapp_watercss-ar21.svg)](https://watercss.netlify.app/ "Classless CSS")

* [cosmtrek/air](https://github.com/cosmtrek/air)
* [spf13/viper](https://github.com/spf13/viper)