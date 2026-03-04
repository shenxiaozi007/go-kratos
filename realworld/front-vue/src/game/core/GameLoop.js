const KEY_MAP = {
  ArrowUp: 'up',
  ArrowDown: 'down',
  ArrowLeft: 'left',
  ArrowRight: 'right',
  w: 'up',
  s: 'down',
  a: 'left',
  d: 'right',
  ' ': 'fire'
}

// 无后端配置时的默认值（已调低）
const DEFAULT_PLAYER_SPEED = 85
const DEFAULT_PLAYER_BULLET = 220
const DEFAULT_PLAYER_COOLDOWN = 0.5
const DEFAULT_ENEMY_SPEED = 55
const DEFAULT_ENEMY_BULLET = 180
const DEFAULT_ENEMY_SHOOT = 1.2
const TANK_RADIUS = 14
const BULLET_RADIUS = 3
const MOVE_SCALE = 70
const BULLET_SCALE = 100

// 与后端 game.v1.TileType 一致
const TILE_EMPTY = 0
const TILE_WALL_BREAKABLE = 1
const TILE_WALL_UNBREAKABLE = 2
const TILE_WATER = 3
const TILE_GRASS = 4

function parseTankConfig(tankConfig) {
  if (!tankConfig || !tankConfig.player) {
    return {
      playerSpeed: DEFAULT_PLAYER_SPEED,
      playerBullet: DEFAULT_PLAYER_BULLET,
      playerCooldown: DEFAULT_PLAYER_COOLDOWN,
      enemySpeed: DEFAULT_ENEMY_SPEED,
      enemyBullet: DEFAULT_ENEMY_BULLET,
      enemyShoot: DEFAULT_ENEMY_SHOOT
    }
  }
  const p = tankConfig.player
  const e = tankConfig.enemy || {}
  return {
    playerSpeed: (p.move_speed != null ? p.move_speed : 1.2) * MOVE_SCALE,
    playerBullet: (p.bullet_speed != null ? p.bullet_speed : 2.5) * BULLET_SCALE,
    playerCooldown: p.fire_cooldown != null ? p.fire_cooldown : DEFAULT_PLAYER_COOLDOWN,
    enemySpeed: (e.move_speed != null ? e.move_speed : 1.0) * MOVE_SCALE,
    enemyBullet: (e.bullet_speed != null ? e.bullet_speed : 2.0) * BULLET_SCALE,
    enemyShoot: e.fire_cooldown != null ? e.fire_cooldown : DEFAULT_ENEMY_SHOOT
  }
}

function normalizeMap(map) {
  if (!map) return null
  const tiles = map.tiles || map.Tiles || []
  const mapWidth = Number(map.width ?? map.Width ?? 20)
  const mapHeight = Number(map.height ?? map.Height ?? 15)
  if (!Array.isArray(tiles) || mapWidth <= 0 || mapHeight <= 0) return null
  return { width: mapWidth, height: mapHeight, tiles }
}

function buildMapGrid(map) {
  const normalized = normalizeMap(map)
  if (!normalized) return null
  const { width: mapWidth, height: mapHeight, tiles } = normalized
  const tileByKey = new Map()
  const breakableKeys = new Set()
  for (const t of tiles) {
    const x = Number(t.x ?? t.X ?? 0)
    const y = Number(t.y ?? t.Y ?? 0)
    const type = t.type != null ? Number(t.type) : (t.Type != null ? Number(t.Type) : TILE_EMPTY)
    const key = `${x},${y}`
    tileByKey.set(key, type)
    if (type === TILE_WALL_BREAKABLE) breakableKeys.add(key)
  }
  return {
    mapWidth,
    mapHeight,
    tileByKey,
    breakableKeys
  }
}

/**
 * @param {HTMLCanvasElement} canvas
 * @param {Object} options
 * @param {Object} [options.tankConfig] - 后端返回的 config
 * @param {Object} [options.map] - 后端返回的 map（width, height, tiles），用于地形渲染与碰撞
 * @param {Function} [options.onStats] - 每帧回调
 * @param {Function} [options.onGameEnd] - 对局结束回调
 * @param {number} [options.autoEndSeconds] - 可选，多少秒后自动结束
 */
export function createGameLoop(canvas, { tankConfig, map: mapData, onStats, onGameEnd, autoEndSeconds } = {}) {
  const ctx = canvas.getContext('2d')
  const width = canvas.width
  const height = canvas.height
  const cfg = parseTankConfig(tankConfig)

  const keys = new Set()

  const tank = {
    x: width / 2,
    y: height / 2,
    angle: 0,
    speed: cfg.playerSpeed,
    hp: 1
  }

  const bullets = []
  const enemies = []
  let enemyKills = 0
  let fireCooldown = 0
  let lastTime = 0
  let frameId = null
  let fps = 0
  let startTime = null
  let ended = false
  let enemyAccum = 0

  const grid = buildMapGrid(mapData)
  let cellW = width / 20
  let cellH = height / 15
  const destroyedBreakable = new Set()
  if (grid) {
    cellW = width / grid.mapWidth
    cellH = height / grid.mapHeight
  }

  function pixelToCell(px, py) {
    const cx = Math.floor(px / cellW)
    const cy = Math.floor(py / cellH)
    return { cx, cy }
  }

  function isBlocking(cx, cy) {
    if (!grid) return false
    const key = `${cx},${cy}`
    const type = grid.tileByKey.get(key)
    if (type === TILE_WALL_UNBREAKABLE) return true
    if (type === TILE_WALL_BREAKABLE && !destroyedBreakable.has(key)) return true
    return false
  }

  function tankHitsWall(tankX, tankY) {
    if (!grid) return false
    const r = TANK_RADIUS
    const cells = new Set()
    for (let dx = -1; dx <= 1; dx++) {
      for (let dy = -1; dy <= 1; dy++) {
        const { cx, cy } = pixelToCell(tankX + dx * r * 0.6, tankY + dy * r * 0.6)
        cells.add(`${cx},${cy}`)
      }
    }
    for (const key of cells) {
      const [cx, cy] = key.split(',').map(Number)
      if (isBlocking(cx, cy)) return true
    }
    return false
  }

  function findSpawnPosition(centerX, centerY) {
    if (!grid) return { x: centerX, y: centerY }
    const { cx, cy } = pixelToCell(centerX, centerY)
    if (!isBlocking(cx, cy)) return { x: centerX, y: centerY }
    for (let r = 1; r <= 5; r++) {
      for (let dy = -r; dy <= r; dy++) {
        for (let dx = -r; dx <= r; dx++) {
          if (Math.abs(dx) !== r && Math.abs(dy) !== r) continue
          const nc = cx + dx
          const mc = cy + dy
          if (nc < 0 || nc >= grid.mapWidth || mc < 0 || mc >= grid.mapHeight) continue
          if (!isBlocking(nc, mc)) {
            return { x: (nc + 0.5) * cellW, y: (mc + 0.5) * cellH }
          }
        }
      }
    }
    return { x: centerX, y: centerY }
  }

  function spawnEnemy() {
    const side = Math.floor(Math.random() * 4)
    let x, y
    if (side === 0) x = 40 + Math.random() * (width - 80)
    else if (side === 1) x = width - 40
    else if (side === 2) x = 40 + Math.random() * (width - 80)
    else x = 40
    if (side === 0) y = 40
    else if (side === 1) y = 40 + Math.random() * (height - 80)
    else if (side === 2) y = height - 40
    else y = 40 + Math.random() * (height - 80)
    if (grid && tankHitsWall(x, y)) {
      const { cx, cy } = pixelToCell(x, y)
      const nx = (cx + 0.5) * cellW
      const ny = (cy + 0.5) * cellH
      if (!tankHitsWall(nx, ny)) {
        x = nx
        y = ny
      }
    }
    enemies.push({
      x,
      y,
      angle: Math.random() * Math.PI * 2,
      speed: cfg.enemySpeed,
      hp: 1,
      nextShoot: Math.random() * cfg.enemyShoot
    })
  }

  function spawnBullet(ox, oy, angle, owner) {
    const dist = 18
    const bulletSpeed = owner === 'player' ? cfg.playerBullet : cfg.enemyBullet
    bullets.push({
      x: ox + Math.cos(angle) * dist,
      y: oy + Math.sin(angle) * dist,
      vx: Math.cos(angle) * bulletSpeed,
      vy: Math.sin(angle) * bulletSpeed,
      owner
    })
  }

  function triggerGameEnd(isWin = true) {
    if (ended || typeof onGameEnd !== 'function') return
    ended = true
    const durationSeconds = startTime ? Math.max(0, Math.floor((Date.now() - startTime) / 1000)) : 0
    onGameEnd({
      is_win: isWin,
      player1_kills: 0,
      player2_kills: 0,
      enemy_kills: enemyKills,
      duration_seconds: durationSeconds
    })
  }

  function handleKeyDown(e) {
    if (e.key === 'e' || e.key === 'E') {
      e.preventDefault()
      triggerGameEnd(true)
      return
    }
    const key = KEY_MAP[e.key]
    if (key) {
      e.preventDefault()
      keys.add(key)
    }
  }

  function handleKeyUp(e) {
    const key = KEY_MAP[e.key]
    if (key) {
      e.preventDefault()
      keys.delete(key)
    }
  }

  function hitTest(x1, y1, r1, x2, y2, r2) {
    const d = (x1 - x2) ** 2 + (y1 - y2) ** 2
    return d <= (r1 + r2) ** 2
  }

  function update(dt) {
    if (ended) return

    if (fireCooldown > 0) fireCooldown -= dt
    if (keys.has('fire') && fireCooldown <= 0) {
      spawnBullet(tank.x, tank.y, tank.angle, 'player')
      fireCooldown = cfg.playerCooldown
    }

    const distance = tank.speed * dt
    if (keys.has('left')) tank.angle -= 2 * dt
    if (keys.has('right')) tank.angle += 2 * dt
    const dx = Math.cos(tank.angle) * distance
    const dy = Math.sin(tank.angle) * distance
    if (keys.has('up')) {
      const nx = tank.x + dx
      const ny = tank.y + dy
      if (!tankHitsWall(nx, ny)) {
        tank.x = nx
        tank.y = ny
      }
    }
    if (keys.has('down')) {
      const nx = tank.x - dx
      const ny = tank.y - dy
      if (!tankHitsWall(nx, ny)) {
        tank.x = nx
        tank.y = ny
      }
    }
    const padding = 20
    tank.x = Math.max(padding, Math.min(width - padding, tank.x))
    tank.y = Math.max(padding, Math.min(height - padding, tank.y))

    for (let i = bullets.length - 1; i >= 0; i--) {
      const b = bullets[i]
      b.x += b.vx * dt
      b.y += b.vy * dt
      if (b.x < -10 || b.x > width + 10 || b.y < -10 || b.y > height + 10) {
        bullets.splice(i, 1)
        continue
      }
      if (grid) {
        const { cx, cy } = pixelToCell(b.x, b.y)
        if (isBlocking(cx, cy)) {
          const key = `${cx},${cy}`
          if (grid.breakableKeys.has(key)) destroyedBreakable.add(key)
          bullets.splice(i, 1)
          continue
        }
      }
      if (b.owner === 'player') {
        for (let j = enemies.length - 1; j >= 0; j--) {
          if (hitTest(b.x, b.y, BULLET_RADIUS, enemies[j].x, enemies[j].y, TANK_RADIUS)) {
            enemies.splice(j, 1)
            bullets.splice(i, 1)
            enemyKills++
            break
          }
        }
      } else {
        if (hitTest(b.x, b.y, BULLET_RADIUS, tank.x, tank.y, TANK_RADIUS)) {
          bullets.splice(i, 1)
          tank.hp--
          if (tank.hp <= 0) triggerGameEnd(false)
        }
      }
    }

    enemyAccum += dt
    if (enemyAccum >= 4) {
      enemyAccum = 0
      if (enemies.length < 5) spawnEnemy()
    }

    for (const e of enemies) {
      const toPlayerX = tank.x - e.x
      const toPlayerY = tank.y - e.y
      const dist = Math.hypot(toPlayerX, toPlayerY) || 1
      const wantAngle = Math.atan2(toPlayerY, toPlayerX)
      let da = wantAngle - e.angle
      while (da > Math.PI) da -= 2 * Math.PI
      while (da < -Math.PI) da += 2 * Math.PI
      e.angle += Math.sign(da) * 1.5 * dt
      const ex = e.x + Math.cos(e.angle) * e.speed * dt
      const ey = e.y + Math.sin(e.angle) * e.speed * dt
      const nx = Math.max(20, Math.min(width - 20, ex))
      const ny = Math.max(20, Math.min(height - 20, ey))
      if (!tankHitsWall(nx, ny)) {
        e.x = nx
        e.y = ny
      }
      e.nextShoot -= dt
      if (e.nextShoot <= 0 && dist < 300) {
        spawnBullet(e.x, e.y, e.angle, 'enemy')
        e.nextShoot = cfg.enemyShoot
      }
    }
  }

  function renderMap() {
    if (!grid) return
    const { mapWidth, mapHeight, tileByKey } = grid
    for (let cy = 0; cy < mapHeight; cy++) {
      for (let cx = 0; cx < mapWidth; cx++) {
        const key = `${cx},${cy}`
        const type = tileByKey.get(key)
        if (type === undefined) continue
        const px = cx * cellW
        const py = cy * cellH
        if (type === TILE_WATER) {
          ctx.fillStyle = '#1e3a5f'
          ctx.fillRect(px, py, cellW + 1, cellH + 1)
          ctx.fillStyle = '#3b82f6'
          ctx.globalAlpha = 0.6
          ctx.fillRect(px, py, cellW + 1, cellH + 1)
          ctx.globalAlpha = 1
        } else if (type === TILE_GRASS) {
          ctx.fillStyle = '#14532d'
          ctx.fillRect(px, py, cellW + 1, cellH + 1)
        } else if (type === TILE_WALL_BREAKABLE && !destroyedBreakable.has(key)) {
          ctx.fillStyle = '#78350f'
          ctx.fillRect(px, py, cellW + 1, cellH + 1)
          ctx.strokeStyle = '#92400e'
          ctx.lineWidth = 1
          ctx.strokeRect(px, py, cellW, cellH)
        } else if (type === TILE_WALL_UNBREAKABLE) {
          ctx.fillStyle = '#374151'
          ctx.fillRect(px, py, cellW + 1, cellH + 1)
          ctx.strokeStyle = '#1f2937'
          ctx.lineWidth = 1
          ctx.strokeRect(px, py, cellW, cellH)
        }
      }
    }
  }

  function render() {
    ctx.clearRect(0, 0, width, height)
    ctx.fillStyle = '#020617'
    ctx.fillRect(0, 0, width, height)
    renderMap()
    if (!grid) {
      ctx.strokeStyle = '#1e293b'
      ctx.lineWidth = 1
      const gridSize = 32
      for (let x = 0; x <= width; x += gridSize) {
        ctx.beginPath()
        ctx.moveTo(x, 0)
        ctx.lineTo(x, height)
        ctx.stroke()
      }
      for (let y = 0; y <= height; y += gridSize) {
        ctx.beginPath()
        ctx.moveTo(0, y)
        ctx.lineTo(width, y)
        ctx.stroke()
      }
    }

    for (const b of bullets) {
      ctx.fillStyle = b.owner === 'player' ? '#38bdf8' : '#f97316'
      ctx.beginPath()
      ctx.arc(b.x, b.y, BULLET_RADIUS, 0, Math.PI * 2)
      ctx.fill()
    }

    for (const e of enemies) {
      ctx.save()
      ctx.translate(e.x, e.y)
      ctx.rotate(e.angle)
      ctx.fillStyle = '#dc2626'
      ctx.fillRect(-14, -10, 28, 20)
      ctx.fillStyle = '#7f1d1d'
      ctx.fillRect(-12, -8, 24, 16)
      ctx.fillStyle = '#f97316'
      ctx.fillRect(0, -3, 20, 6)
      ctx.restore()
    }

    ctx.save()
    ctx.translate(tank.x, tank.y)
    ctx.rotate(tank.angle)
    ctx.fillStyle = '#22c55e'
    ctx.fillRect(-14, -10, 28, 20)
    ctx.fillStyle = '#0f172a'
    ctx.fillRect(-12, -8, 24, 16)
    ctx.fillStyle = '#38bdf8'
    ctx.fillRect(0, -3, 20, 6)
    ctx.restore()
  }

  function loop(timestamp) {
    if (!lastTime) {
      lastTime = timestamp
    }
    const dt = (timestamp - lastTime) / 1000
    lastTime = timestamp

    if (autoEndSeconds != null && startTime != null && !ended) {
      const elapsed = (Date.now() - startTime) / 1000
      if (elapsed >= autoEndSeconds) {
        triggerGameEnd(true)
      }
    }

    const alpha = 0.1
    const currentFps = 1 / Math.max(dt, 1e-6)
    fps = fps * (1 - alpha) + currentFps * alpha

    update(dt)
    render()

    if (typeof onStats === 'function') {
      onStats({
        fps: Math.round(fps),
        tank: { x: tank.x, y: tank.y, angle: tank.angle }
      })
    }

    frameId = window.requestAnimationFrame(loop)
  }

  function start() {
    if (frameId != null) return
    startTime = Date.now()
    lastTime = 0
    ended = false
    enemyKills = 0
    bullets.length = 0
    enemies.length = 0
    enemyAccum = 0
    destroyedBreakable.clear()
    const spawn = findSpawnPosition(width / 2, height / 2)
    tank.x = spawn.x
    tank.y = spawn.y
    tank.angle = 0
    tank.hp = 1
    for (let i = 0; i < 2; i++) spawnEnemy()
    window.addEventListener('keydown', handleKeyDown)
    window.addEventListener('keyup', handleKeyUp)
    frameId = window.requestAnimationFrame(loop)
  }

  function stop() {
    if (frameId != null) {
      window.cancelAnimationFrame(frameId)
      frameId = null
    }
    window.removeEventListener('keydown', handleKeyDown)
    window.removeEventListener('keyup', handleKeyUp)
  }

  return { start, stop }
}

