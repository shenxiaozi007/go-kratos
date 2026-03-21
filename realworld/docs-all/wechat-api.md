# Wechat 接口示例

以下示例默认服务地址为 `http://127.0.0.1:8000`。

## 1) 群发文本

```bash
curl -X POST "http://127.0.0.1:8000/api/wechat/v1/group/text" \
  -H "Content-Type: application/json" \
  -d '{
    "roomId": "123456@chatroom",
    "text": "大家好，这是一条群发文本消息"
  }'
```

## 2) 群发图片（URL 模式）

```bash
curl -X POST "http://127.0.0.1:8000/api/wechat/v1/group/image" \
  -H "Content-Type: application/json" \
  -d '{
    "roomId": "123456@chatroom",
    "imageUrl": "https://example.com/demo.png"
  }'
```

## 3) 群发图片（multipart 上传模式）

```bash
curl -X POST "http://127.0.0.1:8000/api/wechat/v1/group/image/upload" \
  -F "room_id=123456@chatroom" \
  -F "file=@/tmp/demo.png"
```

## 4) 群内 @ 人发文本

```bash
curl -X POST "http://127.0.0.1:8000/api/wechat/v1/group/mention" \
  -H "Content-Type: application/json" \
  -d '{
    "roomId": "123456@chatroom",
    "mentionIds": ["wxid_user_1", "wxid_user_2"],
    "text": "请两位确认一下今天的排期"
  }'
```

## 5) 获取登录二维码

```bash
curl "http://127.0.0.1:8000/api/wechat/v1/login/qrcode"
```

## 6) 查询登录状态

```bash
curl "http://127.0.0.1:8000/api/wechat/v1/login/status"
```

## 7) 诊断当前运行状态

```bash
curl "http://127.0.0.1:8000/api/wechat/v1/diag"
```

返回字段：

- `mode`：`service`（有 token）或 `fallback`（无 token）
- `has_token`：是否检测到 `WECHATY_PUPPET_SERVICE_TOKEN`
- `scan_received` / `qrcode_ready`：是否收到过扫码事件
- `logged_in` / `self_id` / `self_name`：当前登录态快照
- `startup_error`：最近启动错误（无则为空）

## Wechaty 运行约定（Go 内嵌模式）

当前后端已改为 **Go 进程内直接启动 Wechaty**，不再依赖独立 bridge 服务。

可配置环境变量：

- `WECHATY_PUPPET_SERVICE_TOKEN`：配置后自动走 `puppet-service`（推荐）
- `WECHATY_PUPPET_SERVICE_ENDPOINT`：可选，自定义 service endpoint

说明：

- 有 `WECHATY_PUPPET_SERVICE_TOKEN`：自动启用 `puppet-service`。
- 无 token：自动回退默认本地 puppet（Web 免费协议尝试模式，稳定性较差）。
- 若 `GET /api/wechat/v1/login/qrcode` 返回 `WECHATY_QRCODE_NOT_READY`，表示尚未收到 scan 事件。
- 建议优先在启动日志中确认 `scan/login/logout` 事件是否触发，再联调发送接口。
