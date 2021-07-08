# easyconfig
Easy to use、extensible configuration manager for golang 

## Install

```shell
go get github.com/soyacen/easyconfmgr
```


# example
## file
```go
watcher, err := easyconfmgrfile.NewWatcher(fp, easyconfmgr.DiscardLogger)

manager := easyconfmgr.NewManager(
    easyconfmgr.WithLoader(easyconfmgrfile.NewLoader(configFile, easyconfmgr.DiscardLogger)),
    easyconfmgr.WithParser(easyconfmgrparser.NewYamlParser()),
    easyconfmgr.WithValuer(easyconfmgrvaluer.NewTrieTreeValuer()),
    easyconfmgr.WithWatcher(watcher),
)
err := manager.ReadConfig()
if err != nil {
    t.Fatal(err)
}

stringVal, err := manager.GetString("key")

var conf Config
err = manager.Unmarshal(&conf)

events := manager.Events()
tmpConfContent := confContent
for event := range events {
    
}

err = manager.StopWatch()
```

## 