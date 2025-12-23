<script setup lang="ts">
import { ref, computed } from "vue";

type ProductItem = {
  id: number;
  setId: number;
  name: string;
  kind: string;
};

const props = defineProps<{
  products: ProductItem[];
  setId: number | null;
}>();

const emit = defineEmits<{
  add: [product: { name: string; kind: string }];
  update: [product: ProductItem];
  delete: [id: number];
}>();

const newName = ref("");
const newKind = ref("");
const editingId = ref<number | null>(null);

function addProduct() {
  if (!newName.value.trim()) return;
  emit("add", { name: newName.value.trim(), kind: newKind.value });
  newName.value = "";
  newKind.value = "";
}

function startEdit(product: ProductItem) {
  editingId.value = product.id;
}

function saveEdit(product: ProductItem) {
  emit("update", product);
  editingId.value = null;
}

function deleteProduct(id: number) {
  if (confirm("Produkt entfernen?")) {
    emit("delete", id);
  }
}

function getKindIcon(kind: string) {
  switch (kind) {
    case "stempel":
      return "mdi-stamp";
    case "stanze":
      return "mdi-content-cut";
    default:
      return "mdi-shape";
  }
}
</script>

<template>
  <div class="product-list">
    <!-- Add Form -->
    <div class="add-form">
      <input
        v-model="newName"
        type="text"
        placeholder="Neues Produkt..."
        class="input"
        @keyup.enter="addProduct"
      />
      <select v-model="newKind" class="select-kind">
        <option value="">Typ</option>
        <option value="stempel">Stempel</option>
        <option value="stanze">Stanze</option>
      </select>
      <button class="btn-add" @click="addProduct" :disabled="!newName.trim()">
        <i class="mdi mdi-plus"></i>
      </button>
    </div>

    <!-- Products -->
    <div class="products" v-if="products.length">
      <div v-for="product in products" :key="product.id" class="product-item">
        <i :class="['mdi', getKindIcon(product.kind), 'product-icon']"></i>

        <template v-if="editingId === product.id">
          <input
            v-model="product.name"
            class="edit-input"
            @blur="saveEdit(product)"
            @keyup.enter="saveEdit(product)"
          />
          <select
            v-model="product.kind"
            class="edit-select"
            @change="saveEdit(product)"
          >
            <option value="">–</option>
            <option value="stempel">Stempel</option>
            <option value="stanze">Stanze</option>
          </select>
        </template>

        <template v-else>
          <span class="product-name" @click="startEdit(product)">{{
            product.name
          }}</span>
          <span class="product-kind">{{ product.kind || "–" }}</span>
        </template>

        <button class="btn-delete" @click="deleteProduct(product.id)">
          <i class="mdi mdi-close"></i>
        </button>
      </div>
    </div>

    <div v-else class="empty">
      <i class="mdi mdi-shape-outline"></i>
      <span>Noch keine Produkte</span>
    </div>
  </div>
</template>

<style scoped>
.product-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.add-form {
  display: flex;
  gap: 8px;
}

.input {
  flex: 1;
  padding: 12px 14px;
  font-size: 15px;
  border: 1px solid #ddd;
  border-radius: 8px;
}

.input:focus {
  outline: none;
  border-color: #111;
}

.select-kind {
  width: 110px;
  padding: 12px 10px;
  font-size: 14px;
  border: 1px solid #ddd;
  border-radius: 8px;
  background: white;
}

.btn-add {
  width: 46px;
  height: 46px;
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
  opacity: 0.3;
  cursor: not-allowed;
}

.btn-add i {
  font-size: 22px;
}

.products {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.product-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  background: #f8f8f8;
  border-radius: 8px;
}

.product-icon {
  font-size: 20px;
  color: #666;
}

.product-name {
  flex: 1;
  font-size: 15px;
  cursor: pointer;
}

.product-name:hover {
  color: #111;
}

.product-kind {
  font-size: 13px;
  color: #888;
  min-width: 60px;
}

.edit-input {
  flex: 1;
  padding: 6px 10px;
  font-size: 15px;
  border: 1px solid #111;
  border-radius: 6px;
}

.edit-select {
  padding: 6px 8px;
  font-size: 13px;
  border: 1px solid #ddd;
  border-radius: 6px;
}

.btn-delete {
  width: 30px;
  height: 30px;
  border: none;
  border-radius: 6px;
  background: transparent;
  color: #999;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-delete:hover {
  background: #fee;
  color: #c00;
}

.empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 32px;
  color: #999;
}

.empty i {
  font-size: 32px;
}

.empty span {
  font-size: 14px;
}
</style>
