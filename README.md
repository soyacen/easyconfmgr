# leoconfmgr
Easy to use、extensible configuration manager for golang

# example
## file
```go
watcher, err := leoconfmgrfile.NewWatcher(fp)

manager := leoconfmgr.NewManager(
    leoconfmgr.WithLoader(leoconffile.NewLoader(configFile)),
    leoconfmgr.WithParser(leoconfparser.NewYamlParser()),
    leoconfmgr.WithValuer(leoconfvaluer.NewTrieTreeValuer()),
    leoconfmgr.WithWatcher(watcher),
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