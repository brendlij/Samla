<script setup lang="ts">
defineProps<{
  setId: number;
  setName: string;
  boxCode: string;
  boxName: string;
  bagSerial: string;
  locationName: string;
  tags: string[];
  thumbnailPath: string;
  selected?: boolean;
}>();

const emit = defineEmits<{
  click: [];
}>();

function toFileUrl(path: string) {
  if (!path) return "";
  if (path.startsWith("file://")) return path;
  return "file:///" + path.replace(/\\/g, "/");
}
</script>

<template>
  <div class="set-card" :class="{ selected }" @click="emit('click')">
    <div class="card-thumb">
      <img v-if="thumbnailPath" :src="toFileUrl(thumbnailPath)" alt="" />
      <i v-else class="mdi mdi-package-variant"></i>
    </div>

    <div class="card-content">
      <h3 class="card-title">{{ setName }}</h3>

      <div class="card-location">
        <span class="location-item">
          <i class="mdi mdi-map-marker-outline"></i>
          {{ locationName || "–" }}
        </span>
        <span class="location-item">
          <i class="mdi mdi-archive-outline"></i>
          {{ boxCode }}
        </span>
        <span class="location-item">
          <i class="mdi mdi-tag-outline"></i>
          {{ bagSerial }}
        </span>
      </div>

      <div class="card-tags" v-if="tags?.length">
        <span class="tag" v-for="tag in tags.slice(0, 4)" :key="tag">{{
          tag
        }}</span>
        <span class="tag more" v-if="tags.length > 4"
          >+{{ tags.length - 4 }}</span
        >
      </div>
    </div>

    <i class="mdi mdi-chevron-right card-arrow"></i>
  </div>
</template>

<style scoped>
.set-card {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 14px 16px;
  background: white;
  border: 1px solid #eee;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.15s;
}

.set-card:hover {
  border-color: #ccc;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.set-card.selected {
  border-color: #111;
  background: #fafafa;
}

.card-thumb {
  width: 60px;
  height: 60px;
  border-radius: 10px;
  background: #f5f5f5;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  overflow: hidden;
}

.card-thumb img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.card-thumb i {
  font-size: 28px;
  color: #bbb;
}

.card-content {
  flex: 1;
  min-width: 0;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  margin: 0 0 6px 0;
  color: #111;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.card-location {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  font-size: 13px;
  color: #666;
  margin-bottom: 8px;
}

.location-item {
  display: flex;
  align-items: center;
  gap: 4px;
}

.location-item i {
  font-size: 16px;
  color: #999;
}

.card-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.tag {
  padding: 3px 8px;
  background: #f0f0f0;
  border-radius: 6px;
  font-size: 11px;
  color: #555;
}

.tag.more {
  background: #e0e0e0;
  color: #333;
}

.card-arrow {
  font-size: 24px;
  color: #ccc;
  flex-shrink: 0;
}

.set-card:hover .card-arrow {
  color: #999;
}
</style>
