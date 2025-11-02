<script setup>
import { ShieldAlert, Mail } from 'lucide-vue-next'
</script>

<template>
  <main class="install-gate">
    <div class="card">
      <ShieldAlert class="icon" :size="64" :stroke-width="1.5" />
      <h1>Setup Required</h1>

      <p>This open-source app requires installation & configuration.</p>
      <p class="secondary">An administrator must perform the install steps.</p>

      <a class="cta">
        <Mail :size="20" :stroke-width="1.5" />
        <span>Contact Support</span>
      </a>
    </div>
  </main>
</template>

<style scoped>
.install-gate {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  background: #0b1020;
  color: #e6e9f2;
}

.card {
  text-align: center;
  max-width: 700px;
  padding: 2rem 2.5rem;
  border-radius: 1.25rem;
  border: 1px solid rgba(255, 255, 255, 0.12);
  background: rgba(255, 255, 255, 0.04);
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.28);
}

.icon {
  display: block;
  margin: 0 auto 1rem;
}

h1 {
  margin: 0 0 0.75rem;
  font-size: 2rem;
  line-height: 1.2;
}

p {
  margin: 0.25rem 0;
}
.secondary {
  opacity: 0.85;
}

.cta {
  margin-top: 1rem;
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.6rem 1rem;
  border-radius: 0.75rem;
  border: 1px solid rgba(255, 255, 255, 0.18);
  text-decoration: none;
  color: inherit;
}

.cta:hover {
  background: rgba(255, 255, 255, 0.06);
}

.hint {
  margin-top: 0.75rem;
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  opacity: 0.8;
}
</style>
