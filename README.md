# Prayers

This is an application that has two parts: a server, and a client integration.

The server runs on a web host, and requires:
* Go
* CoffeeScript
  * Through Node.js and npm
* MySQL

The client integration is straight forward. At the minimum, the following will enable this service:

```
<div id="prayers" data-integration="first"></div>
<script src="//ajax.googleapis.com/ajax/libs/jquery/2.0.3/jquery.min.js"></script>
<script src="//prayers.bryankendall.com/js/prayers.js"></script>
<noscript>Please enable javascript for prayers to load.</noscript>
```

You can see a sample of the (current working master branch) code [here](http://prayers.bryankendall.com)
