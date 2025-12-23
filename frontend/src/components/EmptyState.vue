<script setup lang="ts">
defineProps<{
  icon?: string;
  title?: string;
}>();
</script>

<template>
  <div class="empty-state">
    <i :class="['mdi', icon || 'mdi-inbox-outline']"></i>
    <h3>{{ title || "Nichts gefunden" }}</h3>
    <p><slot>Versuche eine andere Suche oder erstelle ein neues Set.</slot></p>
  </div>
</template>

<style scoped>
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  text-align: center;
  color: #888;
}

.empty-state i {
  font-size: 56px;
  margin-bottom: 16px;
  color: #ccc;
}

.empty-state h3 {
  font-size: 18px;
  font-weight: 500;
  margin: 0 0 8px 0;
  color: #666;
}

.empty-state p {
  font-size: 14px;
  margin: 0;
  max-width: 300px;
}
</style>
