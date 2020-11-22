<p align="center"><img src="https://user-images.githubusercontent.com/1092882/99123668-a3951000-2626-11eb-8cb3-7e25da2129b5.png" alt="htmlhost.live gopher" width="256px"/></p>

# v0.2.17

Currently this app is hosted under 3 domains, `htmlhost.live` was the first one I thought of, grabbed it quickly and deployed. After a few days, figured it could be named a tad bit better, so now it's available under:

1. ~~[htmlhost.live](https://htmlhost.live)~~
2. [hosthtml.live](https://hosthtml.live)
3. ~~[makehtml.live](https://makehtml.live)~~

Currently all the above domains are live, but based on some internal polls I conducted, `hosthtml.live` won! So I won't be renewing the others later.

## Why?

If you've tried hosting a simple HTML file, be it for sharing with someone or just for the kicks. You know it's a time consuming task. There are some solutions out there but requires registration, or is just tedious and some are on the way of becoming a great experience. At least for me, none of the existing solutions worked, and I had a very specific usecase as well!

So I was out setting up custom error pages on [Cloudflare](http://cloudflare.com). If you've used their service and have setup custom erorr pages, you know they require a URL where your template is hosted, and won't accept the HTML as is. And that's where my journey started, after searching around and trying to use some of the services, finally I hosted it as a [Github page](https://pages.github.com/). For me it turned out to be a time consuming endeavour and I didn't think it was worth it, and thus <a href="/#home">HTMLHost.live</a>!


## How to run?

Have latest stable version of Docker & Docker-compose installed. After that:

```bash
$ docker-compose up -d
$ curl http://localhost
```

Once docker-compose is up, you can open [localhost](http://localhost) in your browser


## Contributing

Open to suggestions and contribution. Please open a [Github issue](https://github.com/bnkamalesh/htmlhost/issues/new) if you'd like to see something changed. Even better if you submit a pull request!
Though please make sure if your suggestion, feedback, bug etc. is already raised in the [issue list](https://github.com/bnkamalesh/htmlhost/issues?q=), in which case just add a comment in the respective issue.

## The Gopher

The gopher used here was created using Gopherize.me. Bring HTML to life like our gopher here!