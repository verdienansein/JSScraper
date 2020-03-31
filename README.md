# JSScraper
Take a list of URLs, parses them for Javascript and scrape them for secrets

## Install

```
▶ go get -u github.com/verdienansein/jsscraper
```

## Basic Usage

JSScraper accepts line-delimited URLs on `stdin`:

```
▶ cat httpsites
https://example.com
https://example.edu
▶ cat httpsites | jsscraper
Scraping URL: http://example.com
Scraping URL: http://example.edu
https://example.edu/gtm.js

http://example.edu/_next/static/2A70LO1buIm0HnWqymD1l/app.js
)return M;switch(!0){case N&&"password"===O:return"new-password"
.concat(m.a.stringify({status_token:e})),options:{method:"GET"}}}
```


