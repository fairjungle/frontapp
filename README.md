# frontapp

Build:

```
$ go build ./cmd/frontapp
```

Post custom message:

```
$ FRONTAPP_APITOKEN='xxxx' ./frontapp receive_custom_message cha_3zh2k "Sender Handle" "This is the body" --subject="This is the subject"
```
