package apollo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/soyacen/easyhttp"
	easyhttpheader "github.com/soyacen/easyhttp/interceptor/header"
	easyhttprespbody "github.com/soyacen/easyhttp/interceptor/respbody"
	"github.com/soyacen/goutils/backoffutils"
	"github.com/soyacen/goutils/fileutils"
	"github.com/soyacen/goutils/retryutils"
	"github.com/soyacen/goutils/stringutils"

	"github.com/soyacen/easyconfmgr"
)

const InitNotificationID = -1

type EventDescription struct {
	appID         string
	cluster       string
	namespaceName string
}

func (ed *EventDescription) String() string {
	return fmt.Sprintf("%s.%s.%s", ed.appID, ed.cluster, ed.namespaceName)
}

type Watcher struct {
	scheme         string
	host           string
	port           int
	appID          string
	cluster        string
	namespaceName  string
	secret         string
	client         *easyhttp.Client
	events         chan *easyconfmgr.Event
	watchOnce      sync.Once
	stopCtx        context.Context
	stopFunc       context.CancelFunc
	stopOnce       sync.Once
	log            easyconfmgr.Logger
	notificationID int64
	maxAttempts    int
	backoffFunc    backoffutils.BackoffFunc
}

func (watcher *Watcher) Watch() error {
	err := errors.New("watcher had started watch")
	watcher.watchOnce.Do(func() {
		err = nil
		go watcher.watch()
	})
	return err
}

func (watcher *Watcher) watch() {
	defer close(watcher.events)
	for {
		select {
		case <-watcher.stopCtx.Done():
			watcher.log.Infof("stop watch apollo config, appID: %s, cluster: %s, namespaceName: %s", watcher.appID, watcher.cluster, watcher.namespaceName)
			return
		default:
			err := retryutils.Call(watcher.stopCtx, watcher.call, watcher.maxAttempts, watcher.backoffFunc)
			// 如果重试N次后还是有error，中断watch
			if err != nil {
				watcher.log.Error(err)
			}
		}
	}
}

func (watcher *Watcher) call() error {
	// 获取远程配置notificationID
	notificationID, err := watcher.notifyRemoteConfig(watcher.stopCtx, watcher.notificationID)
	if err != nil {
		return err
	}

	// 如果是初始化，则忽略
	if watcher.notificationID == InitNotificationID {
		watcher.notificationID = notificationID
		return nil
	}

	// 如果远程notificationID 与本地 notificationID 相同，则不需要更新
	if notificationID == watcher.notificationID {
		return nil
	}

	// 不同，有变化，同步
	uri := fmt.Sprintf("%s://%s:%d/configs/%s/%s/%s", watcher.scheme, watcher.host, watcher.port, watcher.appID, watcher.cluster, watcher.namespaceName)
	watcher.log.Info("reading apollo config:", uri)

	ctx, cancelFunc := context.WithTimeout(watcher.stopCtx, time.Second)
	defer cancelFunc()
	content, err := getConfigContent(ctx, uri, watcher.appID, watcher.secret, watcher.ContentType(), watcher.client)
	if err != nil {
		return err
	}
	watcher.events <- easyconfmgr.NewEvent(&EventDescription{}, []byte(content))
	watcher.notificationID = notificationID
	return nil
}

func (watcher *Watcher) notifyRemoteConfig(ctx context.Context, notificationID int64) (int64, error) {
	notifications := []*Notification{{
		NamespaceName:  watcher.namespaceName,
		NotificationID: notificationID,
	}}
	notificationsJson, err := json.Marshal(notifications)
	if err != nil {
		return 0, err
	}
	uri := fmt.Sprintf(
		"%s://%s:%d/notifications/v2?appId=%s&cluster=%s&notifications=%s",
		watcher.scheme, watcher.host, watcher.port,
		url.QueryEscape(watcher.appID), url.QueryEscape(watcher.cluster),
		url.QueryEscape(string(notificationsJson)))
	watcher.log.Info("reading apollo config:", uri)

	ctx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	var resp []*Notification
	chains := []easyhttp.Interceptor{easyhttprespbody.JSON(&resp)}

	if stringutils.IsNotBlank(watcher.secret) {
		chains = append(chains, easyhttpheader.SetMap(AuthSignatureHeaders(uri, watcher.appID, watcher.secret)))
	}

	reply, err := watcher.client.Get(ctx, uri, easyhttp.ChainInterceptor(chains...))
	if err != nil {
		return 0, err
	}

	httpResp := reply.RawResponse()
	switch httpResp.StatusCode {
	case http.StatusOK:
		return resp[0].NotificationID, nil
	case http.StatusNotModified:
		return notificationID, nil
	default:
		return 0, fmt.Errorf("failed to get config from apollo, URI: %s, StatusCode: %d, Status: %s", uri, httpResp.StatusCode, httpResp.Status)
	}
}

func (watcher *Watcher) Events() <-chan *easyconfmgr.Event {
	return watcher.events
}

func (watcher *Watcher) Stop() error {
	err := errors.New("watcher had stopped")
	watcher.stopOnce.Do(func() {
		err = watcher.stop()
	})
	return err
}

func (watcher *Watcher) stop() error {
	watcher.log.Info("stop watching apollo config...")
	watcher.stopFunc()
	return nil
}

func (watcher *Watcher) ContentType() string {
	contentType := fileutils.ExtName(watcher.namespaceName)
	if stringutils.IsBlank(contentType) {
		return "properties"
	}
	return contentType
}

type WatcherOption func(watcher *Watcher)

func WithScheme(scheme string) WatcherOption {
	return func(watcher *Watcher) {
		watcher.scheme = scheme
	}
}

func WithSecret(secret string) WatcherOption {
	return func(watcher *Watcher) {
		watcher.secret = secret
	}
}

func WithLogger(log easyconfmgr.Logger) WatcherOption {
	return func(watcher *Watcher) {
		watcher.log = log
	}
}

func WithRetry(maxAttempts int, backoffFunc backoffutils.BackoffFunc) WatcherOption {
	return func(watcher *Watcher) {
		watcher.maxAttempts = maxAttempts
		watcher.backoffFunc = backoffFunc
	}
}

func NewWatcher(host string, port int, appID string, cluster string, namespaceName string, opts ...WatcherOption) *Watcher {
	w := &Watcher{
		scheme:         "http",
		host:           host,
		port:           port,
		appID:          appID,
		cluster:        cluster,
		namespaceName:  namespaceName,
		client:         easyhttp.NewClient(easyhttp.WithTimeout(120 * time.Second)),
		events:         make(chan *easyconfmgr.Event),
		watchOnce:      sync.Once{},
		stopOnce:       sync.Once{},
		log:            easyconfmgr.DiscardLogger,
		notificationID: InitNotificationID,
		maxAttempts:    0,
		backoffFunc:    backoffutils.Linear(time.Second),
	}
	for _, opt := range opts {
		opt(w)
	}

	w.stopCtx, w.stopFunc = context.WithCancel(context.Background())

	return w
}
