package apollo

import (
	"context"
	"fmt"
	"net/http"

	"github.com/soyacen/easyhttp"
	easyhttpheader "github.com/soyacen/easyhttp/interceptor/header"
	easyhttprespbody "github.com/soyacen/easyhttp/interceptor/respbody"
	"github.com/soyacen/goutils/stringutils"
)

// ConfigResponse apollo配置响应
type ConfigResponse struct {
	AppID          string                 `json:"appId"`
	Cluster        string                 `json:"cluster"`
	NamespaceName  string                 `json:"namespaceName"`
	ReleaseKey     string                 `json:"releaseKey"`
	Configurations map[string]interface{} `json:"configurations"`
}

// Notification 用于保存 apollo Notification 信息
type Notification struct {
	NamespaceName  string `json:"namespaceName"`
	NotificationID int64  `json:"notificationId"`
}

func getConfigContent(ctx context.Context, uri string, appID, secret, contentType string, client *easyhttp.Client) (string, error) {
	var resp ConfigResponse
	itcptrs := []easyhttp.Interceptor{easyhttprespbody.JSON(&resp)}

	if stringutils.IsNotBlank(secret) {
		itcptrs = append(itcptrs, easyhttpheader.SetMap(AuthSignatureHeaders(uri, appID, secret)))
	}

	reply, err := client.Get(ctx, uri, itcptrs...)
	if err != nil {
		return "", err
	}

	httpResp := reply.RawResponse()
	if httpResp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed get config from apollo, StatusCode: %d, Status: %s", httpResp.StatusCode, httpResp.Status)
	}

	var content string
	switch contentType {
	case "yaml", "yml":
		if val, ok := resp.Configurations["content"]; ok {
			content, ok = val.(string)
			if !ok {
				return "", fmt.Errorf("content is not string, %v, %T", val, val)
			}
		}
	}
	return content, nil
}
