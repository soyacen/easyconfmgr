# easyconfmgr
Easy to use、extensible configuration manager for golang

# example
## file 
```go
// read config
manager := easyconfmgr.NewManager(
    easyconfmgr.WithLoader(file.NewLoader(configFile)),
    easyconfmgr.WithParser(parser.NewYamlParser()),
    easyconfmgr.WithValuer(valuer.NewTrieTreeValuer()),
)
err := manager.ReadConfig()
if err != nil {
    t.Fatal(err)
}

stringVal, err := manager.GetString("string")
assert.Nil(t, err)
assert.Equal(t, config.String, stringVal)

// watch
watcher, err := file.NewWatcher(fp)
assert.Nil(t, err)
manager := easyconfmgr.NewManager(easyconfmgr.WithWatcher(watcher))

err = manager.StartWatch()
assert.Nil(t, err)
events := manager.Events()
tmpConfContent := confContent
for event := range events {
    tmpConfContent += "float_32: 0.3\n"
    assert.Equal(t, tmpConfContent, string(event.Data()))
}

```

## 