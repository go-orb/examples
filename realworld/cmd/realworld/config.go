package main

//nolint:gochecknoglobals
var config = `
---
realworld:
  service:
    httpgateway:
      client:
        middlewares:
          - name: retry
          - name: log
      registry:
        plugin: memory
    lobby:
      client:
        middlewares:
          - name: retry
          - name: log
      registry:
        plugin: memory
`
