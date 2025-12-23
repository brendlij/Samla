<script lang="ts" setup>
import { ref, computed, onMounted } from "vue";
import ConfirmModal from "./ConfirmModal.vue";

type Item = { id: number; name: string };

const props = defineProps<{
  visible: boolean;
  title: string;
  icon: string;
  items: Item[];
}>();

const emit = defineEmits<{
  close: [];
  create: [name: string];
  update: [id: number, name: string];
  delete: [id: number];
}>();

const search = ref("");
const editingId = ref<number | null>(null);
const editingName = ref("");
const newName = ref("");
const deleteModal = ref({ visible: false, id: 0, name: "" });

const filteredItems = computed(() => {
  const q = search.value.toLowerCase().trim();
  if (!q) return props.items;
  return props.items.filter((i) => i.name.toLowerCase().includes(q));
});

function startEdit(item: Item) {
  editingId.value = item.id;
  editingName.value = item.name;
}

function cancelEdit() {
  editingId.value = null;
  editingName.value = "";
}

function saveEdit() {
  if (!editingId.value || !editingName.value.trim()) return;
  emit("update", editingId.value, editingName.value.trim());
  cancelEdit();
}

function handleCreate() {
  if (!newName.value.trim()) return;
  emit("create", newName.value.trim());
  newName.value = "";
}

function confirmDelete(item: Item) {
  deleteModal.value = { visible: true, id: item.id, name: item.name };
}

function executeDelete() {
  emit("delete", deleteModal.value.id);
  deleteModal.value.visible = false;
}

function handleKeydown(e: KeyboardEvent) {
  if (e.key === "Escape") {
    if (editingId.value) {
      cancelEdit();
    } else {
      emit("close");
    }
  }
}
</script>

<template>
  <Teleport to="body">
    <Transition name="slide">
      <div
        v-if="visible"
        class="panel-backdrop"
        @click.self="emit('close')"
        @keydown="handleKeydown"
      >
        <div class="panel">
          <header class="panel-header">
            <div class="panel-title">
              <i :class="['mdi', icon]"></i>
              {{ title }}
            </div>
            <button class="btn-close" @click="emit('close')">
              <i class="mdi mdi-close"></i>
            </button>
          </header>

          <div class="panel-body">
            <!-- Add new -->
            <div class="add-form">
              <input
                v-model="newName"
                type="text"
                class="input"
                :placeholder="`Neue/n ${title} hinzufügen...`"
                @keydown.enter="handleCreate"
              />
              <button
                class="btn-add"
                @click="handleCreate"
                :disabled="!newName.trim()"
              >
                <i class="mdi mdi-plus"></i>
              </button>
            </div>

            <!-- Search -->
            <div class="search-box" v-if="items.length > 5">
              <i class="mdi mdi-magnify"></i>
              <input v-model="search" type="text" placeholder="Suchen..." />
            </div>

            <!-- List -->
            <div class="item-list">
              <div v-if="!filteredItems.length" class="empty">
                <i class="mdi mdi-folder-open-outline"></i>
                <p>Keine Einträge</p>
              </div>

              <div v-for="item in filteredItems" :key="item.id" class="item">
                <template v-if="editingId === item.id">
                  <input
                    v-model="editingName"
                    type="text"
                    class="input edit-input"
                    @keydown.enter="saveEdit"
                    @keydown.escape="cancelEdit"
                    autofocus
                  />
                  <button
                    class="btn-icon success"
                    @click="saveEdit"
                    title="Speichern"
                  >
                    <i class="mdi mdi-check"></i>
                  </button>
                  <button
                    class="btn-icon"
                    @click="cancelEdit"
                    title="Abbrechen"
                  >
                    <i class="mdi mdi-close"></i>
                  </button>
                </template>
                <template v-else>
                  <span class="item-name">{{ item.name }}</span>
                  <button
                    class="btn-icon"
                    @click="startEdit(item)"
                    title="Bearbeiten"
                  >
                    <i class="mdi mdi-pencil-outline"></i>
                  </button>
                  <button
                    class="btn-icon danger"
                    @click="confirmDelete(item)"
                    title="Löschen"
                  >
                    <i class="mdi mdi-delete-outline"></i>
                  </button>
                </template>
              </div>
            </div>

            <div class="item-count">{{ items.length }} Einträge</div>
          </div>
        </div>

        <ConfirmModal
          :visible="deleteModal.visible"
          :message="`'${deleteModal.name}' wirklich löschen?`"
          :danger="true"
          @confirm="executeDelete"
          @cancel="deleteModal.visible = false"
        />
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.panel-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.3);
  display: flex;
  justify-content: flex-end;
  z-index: 900;
}

.panel {
  width: 100%;
  max-width: 400px;
  height: 100%;
  background: white;
  display: flex;
  flex-direction: column;
  box-shadow: -4px 0 20px rgba(0, 0, 0, 0.1);
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid #eee;
  background: #fafafa;
}

.panel-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 16px;
  font-weight: 600;
}

.panel-title i {
  font-size: 22px;
  color: #666;
}

.btn-close {
  width: 36px;
  height: 36px;
  border: none;
  border-radius: 8px;
  background: transparent;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-close:hover {
  background: #eee;
}

.btn-close i {
  font-size: 22px;
  color: #666;
}

.panel-body {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 16px;
  overflow: hidden;
}

.add-form {
  display: flex;
  gap: 8px;
  margin-bottom: 12px;
}

.add-form .input {
  flex: 1;
}

.btn-add {
  width: 44px;
  height: 44px;
  border: none;
  border-radius: 8px;
  background: #111;
  color: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-add:hover:not(:disabled) {
  background: #333;
}

.btn-add:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.btn-add i {
  font-size: 22px;
}

.search-box {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  background: #f5f5f5;
  border-radius: 8px;
  margin-bottom: 12px;
}

.search-box i {
  color: #999;
  font-size: 18px;
}

.search-box input {
  flex: 1;
  border: none;
  background: transparent;
  font-size: 14px;
  outline: none;
}

.item-list {
  flex: 1;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.empty {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #999;
}

.empty i {
  font-size: 48px;
  margin-bottom: 8px;
}

.empty p {
  font-size: 14px;
  margin: 0;
}

.item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  background: #f9f9f9;
  border-radius: 8px;
  transition: background 0.15s;
}

.item:hover {
  background: #f0f0f0;
}

.item-name {
  flex: 1;
  font-size: 14px;
  font-weight: 500;
}

.edit-input {
  flex: 1;
  padding: 8px 10px !important;
  font-size: 14px !important;
}

.btn-icon {
  width: 32px;
  height: 32px;
  border: none;
  border-radius: 6px;
  background: transparent;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #666;
  transition: all 0.15s;
}

.btn-icon:hover {
  background: #e5e5e5;
}

.btn-icon.success {
  color: #059669;
}

.btn-icon.success:hover {
  background: #d1fae5;
}

.btn-icon.danger:hover {
  background: #fee2e2;
  color: #dc2626;
}

.btn-icon i {
  font-size: 18px;
}

.item-count {
  padding-top: 12px;
  text-align: center;
  font-size: 12px;
  color: #999;
  border-top: 1px solid #eee;
  margin-top: 12px;
}

.input {
  width: 100%;
  padding: 12px 14px;
  font-size: 15px;
  border: 1px solid #ddd;
  border-radius: 8px;
  background: white;
}

.input:focus {
  outline: none;
  border-color: #111;
}

/* Transitions */
.slide-enter-active,
.slide-leave-active {
  transition: all 0.25s ease;
}

.slide-enter-from .panel,
.slide-leave-to .panel {
  transform: translateX(100%);
}

.slide-enter-from,
.slide-leave-to {
  background: rgba(0, 0, 0, 0);
}
</style>
