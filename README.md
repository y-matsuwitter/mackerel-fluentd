# mackerel-fluentd

Install
-----

```
$ curl -L -o /path/to/mackerel-fluentd  https://github.com/y-matsuwitter/mackerel-fluentd/releases/download/v0.0.1/mackerel-fluentd.linux
$ chmod a+x /path/to/mackerel-fluentd
```

Usage
----

```
[plugin.metrics.fluentd]
command=/path/to/mackerel-fluentd -host=localhost -port=24220
```

Tips
----
This plugin uses fluentd plugin_id as a mackerel metrics name.
`plugin_id` is not a constant value if you don't specify it in a fluentd configuration.

Please specify a plugin id like this.

```
<match **>
  id stdout
  type stdout
</match>
```

License
----

Mantle is released under the MIT license. See [LICENSE.md](https://github.com/y-matsuwitter/mackerel-fluentd//blob/master/LICENSE.md).
