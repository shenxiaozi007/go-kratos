<template>
  <div class="app-root">
    <div class="app-background" aria-hidden="true"></div>

    <a href="#main-content" class="skip-link">跳到主内容</a>

    <main id="main-content" class="app-shell">
      <header class="shell-header">
        <div class="brand">
          <div class="brand-mark" aria-hidden="true">
            <span class="brand-badge"></span>
          </div>
          <div class="brand-text">
            <p class="brand-title">坦克大战</p>
            <p class="brand-subtitle">前端单人模式原型（Vue 3）</p>
          </div>
        </div>

        <nav class="nav" aria-label="主导航">
          <router-link class="nav-link" to="/">首页</router-link>
          <router-link class="nav-link" to="/game">游戏</router-link>
          <router-link class="nav-link" to="/history">历史战绩</router-link>
          <router-link class="nav-link" to="/leaderboard">排行榜</router-link>
        </nav>
      </header>

      <section class="app-content" aria-label="页面内容">
        <router-view />
      </section>
    </main>
  </div>
</template>

<script>
export default {
  name: 'App'
}
</script>

<style>
.app-root {
  position: relative;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: radial-gradient(circle at top, #0f172a, #020617);
  color: #e5e7eb;
}

.app-background {
  position: fixed;
  inset: 0;
  pointer-events: none;
  background:
    radial-gradient(circle at 20% 0%, rgba(56, 189, 248, 0.18), transparent 55%),
    radial-gradient(circle at 80% 0%, rgba(129, 140, 248, 0.18), transparent 55%),
    radial-gradient(circle at 0% 100%, rgba(34, 197, 94, 0.12), transparent 55%);
  filter: blur(40px);
  opacity: 0.85;
  z-index: 0;
}

.app-shell {
  position: relative;
  z-index: 1;
  max-width: 1120px;
  margin: 0 auto;
  padding: 24px 16px 32px;
}

.shell-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px 12px;
  border-radius: 24px 24px 16px 16px;
  border: 1px solid rgba(148, 163, 184, 0.4);
  background: radial-gradient(circle at 0 0, rgba(248, 250, 252, 0.16), rgba(15, 23, 42, 0.95));
  box-shadow:
    0 22px 60px rgba(15, 23, 42, 0.85),
    0 0 0 1px rgba(15, 23, 42, 0.7);
}

.brand {
  display: flex;
  align-items: center;
  gap: 10px;
}

.brand-mark {
  width: 34px;
  height: 34px;
  border-radius: 14px;
  background: radial-gradient(circle at 30% 0, #4f46e5, #0f172a);
  box-shadow:
    0 14px 30px rgba(79, 70, 229, 0.8),
    0 0 0 1px rgba(129, 140, 248, 0.65);
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
}

.brand-badge {
  width: 18px;
  height: 18px;
  border-radius: 999px;
  background: conic-gradient(from 140deg, #22c55e, #38bdf8, #a855f7, #22c55e);
  box-shadow:
    0 8px 20px rgba(15, 23, 42, 0.9),
    0 0 0 1px rgba(15, 23, 42, 0.9);
}

.brand-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.brand-title {
  font-size: 18px;
  font-weight: 700;
  letter-spacing: 0.03em;
}

.brand-subtitle {
  font-size: 13px;
  color: #b8c5d6;
}

.nav {
  display: flex;
  gap: 12px;
}

.nav-link {
  font-size: 14px;
  color: #e5e7eb;
  text-decoration: none;
  padding: 6px 10px;
  border-radius: 999px;
  transition: all 0.15s ease;
}

.nav-link.router-link-exact-active {
  background: rgba(15, 23, 42, 0.95);
  box-shadow:
    0 1px 2px rgba(15, 23, 42, 0.5),
    0 4px 12px rgba(15, 23, 42, 0.6);
}

.nav-link:hover {
  background: rgba(15, 23, 42, 0.9);
}

.nav-link:focus-visible {
  outline: 2px solid #38bdf8;
  outline-offset: 2px;
}

.skip-link {
  position: absolute;
  top: -100px;
  left: 16px;
  padding: 8px 12px;
  background: #0f172a;
  color: #e5e7eb;
  border: 2px solid #38bdf8;
  border-radius: 8px;
  font-size: 14px;
  text-decoration: none;
  z-index: 100;
  transition: top 0.2s ease;
}
.skip-link:focus {
  top: 16px;
  outline: none;
}

.app-content {
  margin-top: 16px;
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  margin-bottom: 6px;
}

.page-subtitle {
  font-size: 14px;
  color: #b8c5d6;
}
</style>

