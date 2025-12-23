<script setup lang="ts">
import { ref } from "vue";

const props = defineProps<{
  photoPath: string;
  setId: number | null;
  basePath: string;
}>();

const emit = defineEmits<{
  "choose-file": [];
  "from-url": [url: string];
  crop: [src: string, ext: string];
}>();

const urlInput = ref("");
const showUrlInput = ref(false);

function toFileUrl(path: string) {
  if (!path) return "";
  if (path.startsWith("file://")) return path;
  const base = props.basePath;
  const isAbs = /^[a-zA-Z]:[\\/]|^\\/.test(path);
  const combined = isAbs ? path : base ? base + "/" + path : path;
  return "file:///" + combined.replace(/\\/g, "/");
}

function submitUrl() {
  if (urlInput.value.trim()) {
    emit("from-url", urlInput.value.trim());
    urlInput.value = "";
    showUrlInput.value = false;
  }
}
</script>

<template>
  <div class="image-upload">
    <!-- Preview -->
    <div class="preview" v-if="photoPath">
      <img :src="toFileUrl(photoPath)" alt="Set Bild" />
    </div>
    <div class="preview empty" v-else>
      <i class="mdi mdi-image-outline"></i>
      <span v-if="!setId">Speichere das Set zuerst</span>
      <span v-else>Kein Bild</span>
    </div>

    <!-- Actions -->
    <div class="actions">
      <button class="btn" @click="emit('choose-file')" :disabled="!setId">
        <i class="mdi mdi-folder-open"></i>
        Datei wählen
      </button>
      <button
        class="btn"
        @click="showUrlInput = !showUrlInput"
        :disabled="!setId"
      >
        <i class="mdi mdi-link"></i>
        Von URL
      </button>
    </div>

    <!-- URL Input -->
    <div class="url-input" v-if="showUrlInput">
      <input
        v-model="urlInput"
        type="text"
        placeholder="Bild-URL eingeben..."
        class="input"
        @keyup.enter="submitUrl"
      />
      <button
        class="btn-submit"
        @click="submitUrl"
        :disabled="!urlInput.trim()"
      >
        <i class="mdi mdi-check"></i>
      </button>
    </div>
  </div>
</template>

<style scoped>
.image-upload {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.preview {
  width: 100%;
  aspect-ratio: 4/3;
  border-radius: 12px;
  overflow: hidden;
  background: #f5f5f5;
  border: 1px solid #eee;
}

.preview img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.preview.empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  color: #999;
}

.preview.empty i {
  font-size: 48px;
  color: #ccc;
}

.preview.empty span {
  font-size: 14px;
}

.actions {
  display: flex;
  gap: 8px;
}

.btn {
  flex: 1;
  padding: 12px 16px;
  font-size: 14px;
  font-weight: 500;
  border: 1px solid #ddd;
  border-radius: 8px;
  background: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  transition: all 0.15s;
}

.btn:hover:not(:disabled) {
  border-color: #111;
  background: #f8f8f8;
}

.btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.btn i {
  font-size: 18px;
}

.url-input {
  display: flex;
  gap: 8px;
}

.input {
  flex: 1;
  padding: 10px 14px;
  font-size: 14px;
  border: 1px solid #ddd;
  border-radius: 8px;
}

.input:focus {
  outline: none;
  border-color: #111;
}

.btn-submit {
  width: 42px;
  height: 42px;
  border: none;
  border-radius: 8px;
  background: #111;
  color: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-submit:hover:not(:disabled) {
  background: #333;
}

.btn-submit:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

.btn-submit i {
  font-size: 20px;
}
</style>
