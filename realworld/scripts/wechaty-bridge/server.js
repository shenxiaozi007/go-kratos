const express = require('express')
const QRCode = require('qrcode')
const { WechatyBuilder } = require('wechaty')
const { FileBox } = require('file-box')

const port = Number(process.env.PORT || 8787)
const apiKey = (process.env.WECHATY_BRIDGE_API_KEY || '').trim()
const puppet = (process.env.WECHATY_PUPPET || 'wechaty-puppet-wechat4u').trim()
const puppetToken = (process.env.WECHATY_PUPPET_SERVICE_TOKEN || '').trim()

const app = express()
app.use(express.json({ limit: '15mb' }))

let latestQrcode = ''
let latestQrcodeDataUrl = ''
let selfContact = null

const botOptions = {
  name: process.env.WECHATY_NAME || 'kratos-wechaty-bridge',
  puppet,
}
if (puppetToken) {
  botOptions.puppetOptions = { token: puppetToken }
}
const bot = WechatyBuilder.build(botOptions)

function unauthorized(res) {
  return res.status(401).json({ error: 'unauthorized' })
}

function checkApiKey(req, res) {
  if (!apiKey) {
    return true
  }
  const auth = req.headers.authorization || ''
  if (auth === `Bearer ${apiKey}`) {
    return true
  }
  unauthorized(res)
  return false
}

async function resolveRoom(roomId, roomTopic) {
  if (roomId) {
    const room = await bot.Room.find({ id: roomId })
    if (room) return room
  }
  if (roomTopic) {
    const rooms = await bot.Room.findAll({ topic: roomTopic })
    if (rooms && rooms.length > 0) return rooms[0]
  }
  return null
}

bot.on('scan', async (qrcode) => {
  latestQrcode = qrcode || ''
  // #region agent log
  fetch('http://127.0.0.1:7484/ingest/c35e479e-0584-4487-8e6b-3dd35b9a1ce2',{method:'POST',headers:{'Content-Type':'application/json','X-Debug-Session-Id':'6b26ca'},body:JSON.stringify({sessionId:'6b26ca',runId:'pre-fix',hypothesisId:'H1',location:'scripts/wechaty-bridge/server.js:35',message:'scan event fired',data:{error:'scan_event',hasQrcode:Boolean(latestQrcode),qrcodeLen:latestQrcode.length},timestamp:Date.now()})}).catch(()=>{});
  // #endregion
  if (latestQrcode) {
    try {
      latestQrcodeDataUrl = await QRCode.toDataURL(latestQrcode)
      // #region agent log
      fetch('http://127.0.0.1:7484/ingest/c35e479e-0584-4487-8e6b-3dd35b9a1ce2',{method:'POST',headers:{'Content-Type':'application/json','X-Debug-Session-Id':'6b26ca'},body:JSON.stringify({sessionId:'6b26ca',runId:'pre-fix',hypothesisId:'H4',location:'scripts/wechaty-bridge/server.js:40',message:'qrcode data url generated',data:{error:'qrcode_dataurl_generated',dataUrlLen:latestQrcodeDataUrl.length},timestamp:Date.now()})}).catch(()=>{});
      // #endregion
    } catch (err) {
      // #region agent log
      fetch('http://127.0.0.1:7484/ingest/c35e479e-0584-4487-8e6b-3dd35b9a1ce2',{method:'POST',headers:{'Content-Type':'application/json','X-Debug-Session-Id':'6b26ca'},body:JSON.stringify({sessionId:'6b26ca',runId:'pre-fix',hypothesisId:'H4',location:'scripts/wechaty-bridge/server.js:43',message:'qrcode data url generation failed',data:{error:String(err?.message||err)},timestamp:Date.now()})}).catch(()=>{});
      // #endregion
    }
  }
  console.log('[bridge] scan qrcode updated')
})

bot.on('login', (user) => {
  selfContact = user
  console.log('[bridge] login:', user?.name())
})

bot.on('logout', (user) => {
  selfContact = null
  console.log('[bridge] logout:', user?.name())
})

bot.on('error', (err) => {
  // #region agent log
  const axiosUrl = err?.response?.config?.url || err?.config?.url || ''
  const axiosStatus = err?.response?.status || ''
  const axiosCode = err?.code || ''
  fetch('http://127.0.0.1:7484/ingest/c35e479e-0584-4487-8e6b-3dd35b9a1ce2',{method:'POST',headers:{'Content-Type':'application/json','X-Debug-Session-Id':'6b26ca'},body:JSON.stringify({sessionId:'6b26ca',runId:'pre-fix',hypothesisId:'H1',location:'scripts/wechaty-bridge/server.js:64',message:'wechaty error event (url='+String(axiosUrl)+' status='+String(axiosStatus)+')',data:{error:String(err?.message||err)},timestamp:Date.now()})}).catch(()=>{});
  // #endregion
  console.error('[bridge] wechaty error:', err)
})

app.get('/login/qrcode', async (req, res) => {
  if (!checkApiKey(req, res)) return
  // #region agent log
  fetch('http://127.0.0.1:7484/ingest/c35e479e-0584-4487-8e6b-3dd35b9a1ce2',{method:'POST',headers:{'Content-Type':'application/json','X-Debug-Session-Id':'6b26ca'},body:JSON.stringify({sessionId:'6b26ca',runId:'pre-fix',hypothesisId:'H2',location:'scripts/wechaty-bridge/server.js:72',message:'login qrcode endpoint called',data:{error:'endpoint_called',hasQrcode:Boolean(latestQrcode),qrcodeLen:latestQrcode.length,hasDataUrl:Boolean(latestQrcodeDataUrl),dataUrlLen:latestQrcodeDataUrl.length,loggedIn:Boolean(selfContact)},timestamp:Date.now()})}).catch(()=>{});
  // #endregion
  return res.json({
    qrcode: latestQrcode,
    qrcode_data_url: latestQrcodeDataUrl,
  })
})

app.get('/login/status', async (req, res) => {
  if (!checkApiKey(req, res)) return
  return res.json({
    logged_in: Boolean(selfContact),
    self_id: selfContact ? selfContact.id : '',
    self_name: selfContact ? selfContact.name() : '',
  })
})

app.post('/group/text', async (req, res) => {
  if (!checkApiKey(req, res)) return
  try {
    const { room_id, room_topic, text } = req.body || {}
    const room = await resolveRoom(room_id, room_topic)
    if (!room) return res.status(404).json({ success: false, error: 'room not found' })
    await room.say(text || '')
    return res.json({ success: true, message_id: '' })
  } catch (err) {
    return res.status(500).json({ success: false, error: String(err?.message || err) })
  }
})

app.post('/group/image', async (req, res) => {
  if (!checkApiKey(req, res)) return
  try {
    const { room_id, room_topic, image_url, image_base64, filename } = req.body || {}
    const room = await resolveRoom(room_id, room_topic)
    if (!room) return res.status(404).json({ success: false, error: 'room not found' })

    let file
    if (image_url) {
      file = FileBox.fromUrl(image_url)
    } else if (image_base64) {
      const buffer = Buffer.from(image_base64, 'base64')
      file = FileBox.fromBuffer(buffer, filename || 'upload.png')
    } else {
      return res.status(400).json({ success: false, error: 'image_url or image_base64 required' })
    }

    await room.say(file)
    return res.json({ success: true, message_id: '' })
  } catch (err) {
    return res.status(500).json({ success: false, error: String(err?.message || err) })
  }
})

app.post('/group/mention', async (req, res) => {
  if (!checkApiKey(req, res)) return
  try {
    const { room_id, room_topic, mention_ids, text } = req.body || {}
    const room = await resolveRoom(room_id, room_topic)
    if (!room) return res.status(404).json({ success: false, error: 'room not found' })

    const ids = Array.isArray(mention_ids) ? mention_ids : []
    if (ids.length === 0) {
      return res.status(400).json({ success: false, error: 'mention_ids is required' })
    }

    const members = await room.memberAll()
    const mentions = members.filter((m) => ids.includes(m.id))
    if (mentions.length === 0) {
      return res.status(404).json({ success: false, error: 'mention contacts not found in room' })
    }

    await room.say(text || '', ...mentions)
    return res.json({ success: true, message_id: '' })
  } catch (err) {
    return res.status(500).json({ success: false, error: String(err?.message || err) })
  }
})

app.listen(port, async () => {
  console.log(`[bridge] listening on :${port}`)
  // #region agent log
  fetch('http://127.0.0.1:7484/ingest/c35e479e-0584-4487-8e6b-3dd35b9a1ce2',{method:'POST',headers:{'Content-Type':'application/json','X-Debug-Session-Id':'6b26ca'},body:JSON.stringify({sessionId:'6b26ca',runId:'pre-fix',hypothesisId:'H3',location:'scripts/wechaty-bridge/server.js:156',message:'bridge listen started',data:{port,puppet,hasToken:Boolean(puppetToken)},timestamp:Date.now()})}).catch(()=>{});
  // #endregion
  try {
    await bot.start()
    console.log('[bridge] bot started')
  } catch (err) {
    console.error('[bridge] bot start failed:', err)
  }
})
