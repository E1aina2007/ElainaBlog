<script setup lang="ts">
import { computed } from 'vue'
import { useSiteStore } from '@/stores/site'

const siteStore = useSiteStore()
const currentYear = new Date().getFullYear()
const siteName = computed(() => siteStore.get('site_name'))
const siteTitle = computed(() => siteStore.get('site_title'))
const icpBeian = computed(() => siteStore.get('icp_beian'))
const govPoliceRecord = computed(() => siteStore.get('gov_police_record'))
</script>

<template>
  <footer class="footer">
    <div class="footer-container">
      <!-- Logo -->
      <div class="footer-logo">
        <img class="logo-icon" src="/favicon.ico" alt="logo" />
        <span class="logo-text">{{ siteName }}</span>
      </div>

      <!-- 版权信息 -->
      <div class="copyright">
        <p>© {{ currentYear }} {{ siteTitle }}. All rights reserved.</p>
      </div>

      <!-- 备案信息 -->
      <div class="beian">
        <a v-if="icpBeian && icpBeian !== '京ICP备xxxxxxxx号'" :href="'https://beian.miit.gov.cn'" target="_blank">{{ icpBeian }}</a>
        <span v-else>{{ icpBeian }}</span>
        <span v-if="icpBeian && govPoliceRecord" class="beian-divider"> | </span>
        <a v-if="govPoliceRecord" :href="'https://www.beian.gov.cn'" target="_blank">{{ govPoliceRecord }}</a>
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
  width: 24px;
  height: 24px;
  object-fit: contain;
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

.beian-divider {
  color: var(--text-muted);
  margin: 0 4px;
}
</style>