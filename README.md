# golang_IniConfigLoader
a smart tool class for loading configs from the dosini format file, written by golang.

`Author. chenyuebin `
`Date. 2016-5-5 02:03:55`

## eg.  

- a config.ini is like this:

```ini
[session-name]
key=value
              <--- support empty line
# this is a comment
key2=value2 # this is a comment too

bad thing here <--- this will arise an parse error when loading config file
[session2] # this is another comment

key=value

key=value2 <---- danger! here will cover the key=value, loader will tell you: "line 11: config key(key) confict value(value != value2)"

key3 =    value3     <--- we can also write some space in line

```

- after we load the config, the obj data structure in the memory is like this:

```go
map['session-name']['key']=value
```

- you can print the config set to log too
