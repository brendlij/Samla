<script lang="ts" setup>
import { nextTick, onBeforeUnmount, ref, watch } from 'vue'
import Cropper from 'cropperjs'
import 'cropperjs/dist/cropper.css'

const props = defineProps<{
  visible: boolean
  imageSrc: string
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'confirm', dataUrl: string): void
  (e: 'skip'): void
}>()

const imgRef = ref<HTMLImageElement | null>(null)
let cropper: Cropper | null = null

function destroyCropper() {
  if (cropper) {
    cropper.destroy()
    cropper = null
  }
}

function initCropper() {
  if (!imgRef.value) return
  destroyCropper()
  cropper = new Cropper(imgRef.value, {
    viewMode: 1,
    autoCropArea: 0.92,
    responsive: true,
    background: false,
  })
}

watch(
  () => props.visible,
  (visible) => {
    if (visible) {
      nextTick(initCropper)
    } else {
      destroyCropper()
    }
  },
)

watch(
  () => props.imageSrc,
  (src) => {
    if (cropper && src) {
      cropper.replace(src)
    }
  },
)

function close() {
  emit('close')
}

function confirmCrop() {
  if (!cropper) return
  const canvas = cropper.getCroppedCanvas({ maxWidth: 4000, maxHeight: 4000 })
  const dataUrl = canvas.toDataURL('image/png')
  emit('confirm', dataUrl)
}

function skipCrop() {
  emit('skip')
}

onBeforeUnmount(destroyCropper)
</script>

<template>
  <div v-if="visible" class="cropper-overlay">
    <div class="cropper-card">
      <div class="header">
        <div>
          <div class="title">Bild zuschneiden</div>
          <div class="hint">Ziehe den Rahmen oder zoome für einen passenden Ausschnitt.</div>
        </div>
        <button class="ghost" @click="close">Schließen</button>
      </div>
      <div class="body">
        <img v-if="imageSrc" :src="imageSrc" ref="imgRef" class="image" />
        <div v-else class="placeholder">Kein Bild geladen.</div>
      </div>
      <div class="actions">
        <button class="ghost" @click="skipCrop">Ohne Zuschnitt</button>
        <button class="primary" @click="confirmCrop">Crop &amp; speichern</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.cropper-overlay {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.55);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 20;
}

.cropper-card {
  background: #fff;
  border-radius: 14px;
  width: min(900px, 96vw);
  padding: 14px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.25);
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title {
  font-weight: 700;
  font-size: 18px;
}

.hint {
  color: var(--muted);
  font-size: 13px;
}

.body {
  background: #0f172a;
  border-radius: 10px;
  padding: 10px;
  min-height: 320px;
}

.image {
  max-width: 100%;
  max-height: 480px;
  display: block;
  margin: 0 auto;
}

.placeholder {
  color: #cbd5e1;
  text-align: center;
  padding: 40px 0;
}

.actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

button {
  border: 1px solid var(--border);
  background: #fff;
  color: var(--ink);
  border-radius: 8px;
  padding: 8px 12px;
  cursor: pointer;
  font-weight: 600;
}

button.primary {
  background: var(--accent);
  color: #fff;
  border-color: var(--accent);
}

button.ghost {
  background: transparent;
}
</style>
