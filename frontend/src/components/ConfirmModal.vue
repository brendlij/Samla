<script lang="ts" setup>
import { computed } from "vue";
import { useI18n } from "../i18n";

const props = defineProps<{
  visible: boolean;
  title?: string;
  message: string;
  confirmText?: string;
  cancelText?: string;
  danger?: boolean;
  icon?: string;
}>();

const emit = defineEmits<{
  confirm: [];
  cancel: [];
}>();

const { t, locale } = useI18n();

const iconClass = computed(() => {
  if (props.icon) return props.icon;
  return props.danger ? "mdi-alert-circle-outline" : "mdi-help-circle-outline";
});
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="visible" class="modal-backdrop" @click.self="emit('cancel')">
        <div class="modal">
          <div class="modal-icon" :class="{ danger }">
            <i :class="['mdi', iconClass]"></i>
          </div>
          <h3 class="modal-title">
            {{ title || t("confirmDelete") }}
          </h3>
          <p class="modal-message">{{ message }}</p>
          <div class="modal-actions">
            <button class="btn btn-secondary" @click="emit('cancel')">
              {{ cancelText || t("cancel") }}
            </button>
            <button
              class="btn"
              :class="danger ? 'btn-danger' : 'btn-primary'"
              @click="emit('confirm')"
            >
              {{ confirmText || (danger ? t("delete") : t("confirm")) }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  backdrop-filter: blur(2px);
}

.modal {
  background: white;
  border-radius: 16px;
  padding: 28px;
  width: 100%;
  max-width: 380px;
  text-align: center;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.2);
}

.modal-icon {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: #f3f4f6;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 16px;
}

.modal-icon i {
  font-size: 28px;
  color: #374151;
}

.modal-icon.danger {
  background: #fef2f2;
}

.modal-icon.danger i {
  color: #dc2626;
}

.modal-title {
  font-size: 18px;
  font-weight: 600;
  margin: 0 0 8px;
  color: #111;
}

.modal-message {
  font-size: 14px;
  color: #6b7280;
  margin: 0 0 24px;
  line-height: 1.5;
}

.modal-actions {
  display: flex;
  gap: 10px;
}

.btn {
  flex: 1;
  padding: 12px 16px;
  font-size: 14px;
  font-weight: 500;
  border-radius: 10px;
  border: none;
  cursor: pointer;
  transition: all 0.15s;
}

.btn-secondary {
  background: #f3f4f6;
  color: #374151;
}

.btn-secondary:hover {
  background: #e5e7eb;
}

.btn-primary {
  background: #111;
  color: white;
}

.btn-primary:hover {
  background: #333;
}

.btn-danger {
  background: #dc2626;
  color: white;
}

.btn-danger:hover {
  background: #b91c1c;
}

/* Transitions */
.modal-enter-active,
.modal-leave-active {
  transition: all 0.2s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .modal,
.modal-leave-to .modal {
  transform: scale(0.95);
}
</style>
