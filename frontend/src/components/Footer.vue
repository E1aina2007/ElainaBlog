<script setup lang="ts">
const currentYear = new Date().getFullYear()
</script>

<template>
  <footer class="footer">
    <div class="footer-container">
      <!-- Logo -->
      <div class="footer-logo">
        <span class="logo-icon">🌿</span>
        <span class="logo-text">ElainaBlog</span>
      </div>

      <!-- 版权信息 -->
      <div class="copyright">
        <p>© {{ currentYear }} ElainaBlog. All rights reserved.</p>
      </div>

      <!-- 备案信息 -->
      <div class="beian">
        <a href="#" target="_blank">京ICP备xxxxxxxx号</a>
      </div>
    </div>
  </footer>
</template> 

<style scoped>
.footer {
  padding: 32px 24px;
  text-align: center;
  background: var(--bg-secondary);
}

.footer-container {
  max-width: 1200px;
  margin: 0 auto;
}

/* Logo */
.footer-logo {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin-bottom: 16px;
}

.logo-icon {
  font-size: 24px;
}

.logo-text {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
}

/* 版权 */
.copyright {
  font-size: 13px;
  color: var(--text-muted);
  margin-bottom: 8px;
}

/* 备案 */
.beian {
  font-size: 13px;
}

.beian a {
  color: var(--text-muted);
  text-decoration: none;
  transition: color var(--transition-fast);
}

.beian a:hover {
  color: var(--primary);
}
</style>