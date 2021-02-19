# WechatPay GO(v3)

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/gunsluo/wechatpay-go/blob/master/LICENSE)
[![CI](https://github.com/gunsluo/wechatpay-go/workflows/ci/badge.svg)](https://github.com/gunsluo/wechatpay-go/actions?query=branch%3Amaster)
[![codecov](https://codecov.io/gh/gunsluo/wechatpay-go/branch/master/graph/badge.svg?token=VFZKUPNGXN)](https://codecov.io/gh/gunsluo/wechatpay-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/gunsluo/wechatpay-go)](https://goreportcard.com/report/github.com/gunsluo/wechatpay-go)

## Introduction

Wechat Pay SDK(V3) Write by Go. API V3 of Office document is [here](https://pay.weixin.qq.com/wiki/doc/apiv3/index.shtml).

## Features
* Signature/Verify messages
* Encrypt/Decrypt cert
* APIv3 Endpoints
* None third-party dependency package

When developing, you can use the `Makefile` for doing the following operations:

| Endpoint           | Description                                                      |        supported       |
| ------------------:| -----------------------------------------------------------------|:----------------------:|
| `pay`              | Merchant send the payment transaction                            |   :heavy_check_mark:   |
| `query`            | Merchant query payment transactions                              |   :heavy_check_mark:   |
| `close`            | Merchant close the payment transaction                           |   :heavy_check_mark:   |
| `notify`           | WeChat Pay notifies the merchant of the user's payment status    |   :heavy_check_mark:   |
| `certificate`      | obtain the platform cert and decrypt it to public key            |   :heavy_check_mark:   |
| `tradebill`        | obtain the download url of trade bill                            |   :heavy_check_mark:   |
| `fundflowbill`     | obtain the download url of trade bill                            |   :heavy_check_mark:   |
| `refund`           | Merchant send the refund transaction                             |   :heavy_check_mark:   |
| `refundquery`      | Merchant query payment transactions                              |   :heavy_check_mark:   |
| `refundnotify`     | WeChat Pay notifies the merchant of the refund status            |   :heavy_check_mark:   |


## Getting Started

Prepare your wechatp pay information, it includes App Id/Mech Id/Apiv3 Secret/Serial Number/Private Key Cert. You can find a getting started guide as shown below: 

1. *import package*
```
import "github.com/gunsluo/wechatpay-go/v3"
```

2. *use wechatpay-go package*
```Go
// create a client of wechat pay
client, err := wechatpay.NewClient(
    wechatpay.Config{
       ...
    })

// create a pay request
req := &wechatpay.PayRequest{
    AppId:       appId,
    MchId:       mchId,
    Description: "for testing",
        ...
    TradeType: wechatpay.Native,
}

resp, err := req.Do(r.Context(), client)
if err != nil {
    // do something
}
codeUrl := resp.CodeUrl
```

#### Config

Click [Wechat Pay](https://pay.weixin.qq.com/) and apply your account and configuration.
```
wechatpay.Config{
    AppId:       appId,
    MchId:       mchId,
    Apiv3Secret: apiv3Secret,
    Cert: wechatpay.CertSuite{
        SerialNo:       serialNo,
        PrivateKeyPath: privateKeyPath,
    },
}
```

#### Payment

Create a pay request and send it to wechat pay service.
```
req := &wechatpay.PayRequest{
    Description: "for testing",
    OutTradeNo:  tradeNo,
    TimeExpire:  time.Now().Add(10 * time.Minute),
    Attach:      "cipher code",
    NotifyUrl:   notifyURL,
    Amount: wechatpay.PayAmount{
        Total:    int(amount * 100),
        Currency: "CNY",
    },
    TradeType: wechatpay.Native,
}

resp, err := req.Do(r.Context(), payClient)
if err != nil {
    e := &wechatpay.Error{}
    if errors.As(err, &e) {
        fmt.Println("status", e.Status, "code:", e.Code, "message:", e.Message)
    }
    return
}
codeUrl := resp.CodeUrl
// use this code url to generate qr code
```

#### Notify

Receive the notification from wechat pay, use `ParseHttpRequest` or `Parse` to get notification information.
```
func notifyForPay(w http.ResponseWriter, r *http.Request) {
    notification := &wechatpay.PayNotification{}
    trans, err := notification.ParseHttpRequest(payClient, r)

    ...
}

func notifyForRefund(w http.ResponseWriter, r *http.Request) {
    notification := &wechatpay.RefundNotification{}
    trans, err := notification.ParseHttpRequest(payClient, r)

    ...
}
```

There is [a full example](https://github.com/gunsluo/wechatpay-example) for wechatpay-go.

#### Download

download bill file, Method `Download` get the decrypted byte array and and `UnmarshalDownload` to get a struct data.
```
req := wechatpay.TradeBillRequest{
    BillDate: billDate,
    BillType: wechatpay.AllBill,
    TarType:  wechatpay.GZIP,
}

ctx := context.Background()
data, err := req.Download(ctx, payClient)
//resp, err := req.UnmarshalDownload(ctx, payClient)
```


## TODO

* Pay combine transactions
* Close combine transactions
* Query combine transactions

## Contributing

See the [contributing documentation](CONTRIBUTING.md).

