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
| `refund`           | Merchant send the refund transaction                             |:heavy_multiplication_x:|
| `refundquery`      | Merchant query payment transactions                              |:heavy_multiplication_x:|
| `refundnotify`     | WeChat Pay notifies the merchant of the refund status            |:heavy_multiplication_x:|

Note: *Endpoints about refund still uses v2 version, will update once wechat-pay upgrade*


## Getting Started

Prepare your wechatp pay information, it includes App Id/Mech Id/Apiv3 Secret/Serial Number/Private Key Cert. You can find a getting started guide as shown below: 

1. *import package*
```
import "github.com/gunsluo/wechatpay-go/v3"
```

2. *use this sdk*
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

There is [a full example](https://github.com/gunsluo/wechatpay-example) for wechatpay-go.

## Contributing

See the [contributing documentation](CONTRIBUTING.md).

