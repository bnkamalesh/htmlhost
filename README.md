# HTMLHost.live

If you've tried hosting a simple HTML file, be it for sharing with someone or just for the kicks. You know it's a time consuming task. There are some solutions out there but requires registration, or is just tedious and some are on the way of becoming a great experience. At least for me, none of the existing solutions worked, and I had a very specific usecase as well!

So I was out setting up custom error pages on [Cloudflare](http://cloudflare.com). If you've used their service and have setup custom erorr pages, you know they require a URL where your template is hosted, and won't accept the HTML as is. And that's where my journey started, after searching around and trying to use some of the services, finally I hosted it as a [Github page](https://pages.github.com/). For me it turned out to be a time consuming endeavour and I didn't think it was worth it, and thus <a href="/#home">HTMLHost.live</a>!


## How to run?

Have latest stable version of Docker & Docker-compose installed. After that:

```bash
$ docker-compose up -d
$ curl http://localhost
```

Once docker-compose is up, you can open [localhost](http://localhost) in your browser