# Publichost.me

## Protocol Example

```
> EXPOSE http 127.0.0.1 4000
< OK endpoint 1
< NOK some error message

< CONN {

    local:
}

< CONN route 1
> OK stream 1 
```