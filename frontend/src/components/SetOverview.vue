<script setup lang="ts">
import { computed } from "vue";

type Product = {
  id: number;
  setId: number;
  name: string;
  kind: string;
};

type BagInfo = {
  id: number;
  serialNo: string;
  boxId: number;
  boxCode: string;
  boxName: string;
  locationId: number;
  locationName: string;
  locationNote: string;
  locationRoom: string;
  locationShelf: string;
  locationCompartment: string;
};

const props = defineProps<{
  id: number;
  name: string;
  manufacturerName: string;
  typeName: string;
  photoPath: string;
  bag: BagInfo | null;
  tags: string[];
  products: Product[];
  basePath: string;
}>();

const emit = defineEmits<{
  edit: [];
  back: [];
  delete: [];
}>();

function toFileUrl(path: string) {
  if (!path) return "";
  if (path.startsWith("file://")) return path;
  const isAbs = /^[a-zA-Z]:[\\/]|^\\/.test(path);
  const combined = isAbs
    ? path
    : props.basePath
    ? props.basePath + "/" + path
    : path;
  return "file:///" + combined.replace(/\\/g, "/");
}

const imageUrl = computed(() =>
  props.photoPath ? toFileUrl(props.photoPath) : ""
);

const locationDisplay = computed(() => {
  if (!props.bag) return "-";
  const parts = [props.bag.locationName];
  if (props.bag.locationRoom) parts.push(props.bag.locationRoom);
  if (props.bag.locationShelf) parts.push(`Regal ${props.bag.locationShelf}`);
  if (props.bag.locationCompartment)
    parts.push(`Fach ${props.bag.locationCompartment}`);
  return parts.join(" • ");
});

const stempelCount = computed(
  () => props.products.filter((p) => p.kind === "stempel").length
);
const stanzeCount = computed(
  () => props.products.filter((p) => p.kind === "stanze").length
);
const otherCount = computed(() => props.products.filter((p) => !p.kind).length);
</script>

<template>
  <div class="overview">
    <!-- Header -->
    <div class="overview-header">
      <button class="btn-back" @click="emit('back')">
        <i class="mdi mdi-arrow-left"></i>
      </button>
      <div class="header-actions">
        <button class="btn-action" @click="emit('edit')">
          <i class="mdi mdi-pencil-outline"></i>
          Bearbeiten
        </button>
        <button class="btn-action danger" @click="emit('delete')">
          <i class="mdi mdi-delete-outline"></i>
        </button>
      </div>
    </div>

    <!-- Main Content -->
    <div class="overview-content">
      <!-- Hero Section with Image -->
      <div class="hero">
        <div v-if="imageUrl" class="hero-image">
          <img :src="imageUrl" :alt="name" />
        </div>
        <div v-else class="hero-placeholder">
          <i class="mdi mdi-image-outline"></i>
        </div>
        <div class="hero-info">
          <h1 class="set-name">{{ name }}</h1>
          <div class="set-meta">
            <span v-if="manufacturerName" class="meta-item">
              <i class="mdi mdi-factory"></i>
              {{ manufacturerName }}
            </span>
            <span v-if="typeName" class="meta-item">
              <i class="mdi mdi-shape-outline"></i>
              {{ typeName }}
            </span>
          </div>
        </div>
      </div>

      <!-- Info Cards -->
      <div class="info-grid">
        <!-- Location Card -->
        <div class="info-card">
          <div class="card-icon">
            <i class="mdi mdi-map-marker-outline"></i>
          </div>
          <div class="card-content">
            <span class="card-label">Lagerort</span>
            <span class="card-value">{{ locationDisplay }}</span>
          </div>
        </div>

        <!-- Box Card -->
        <div class="info-card">
          <div class="card-icon">
            <i class="mdi mdi-package-variant-closed"></i>
          </div>
          <div class="card-content">
            <span class="card-label">Box</span>
            <span class="card-value"
              >{{ bag?.boxCode || "-" }}
              {{ bag?.boxName ? `(${bag.boxName})` : "" }}</span
            >
          </div>
        </div>

        <!-- Bag Card -->
        <div class="info-card">
          <div class="card-icon">
            <i class="mdi mdi-shopping-outline"></i>
          </div>
          <div class="card-content">
            <span class="card-label">Beutel-Nr.</span>
            <span class="card-value">{{ bag?.serialNo || "-" }}</span>
          </div>
        </div>

        <!-- Products Summary Card -->
        <div class="info-card">
          <div class="card-icon">
            <i class="mdi mdi-view-grid-outline"></i>
          </div>
          <div class="card-content">
            <span class="card-label">Produkte</span>
            <span class="card-value">
              {{ products.length }} gesamt
              <span v-if="stempelCount" class="product-badge"
                >{{ stempelCount }} Stempel</span
              >
              <span v-if="stanzeCount" class="product-badge"
                >{{ stanzeCount }} Stanzen</span
              >
            </span>
          </div>
        </div>
      </div>

      <!-- Tags -->
      <div v-if="tags.length" class="section">
        <h2><i class="mdi mdi-tag-multiple-outline"></i> Tags</h2>
        <div class="tags-list">
          <span v-for="tag in tags" :key="tag" class="tag">{{ tag }}</span>
        </div>
      </div>

      <!-- Products -->
      <div v-if="products.length" class="section">
        <h2>
          <i class="mdi mdi-shape-outline"></i> Produkte ({{ products.length }})
        </h2>
        <div class="products-grid">
          <div
            v-for="product in products"
            :key="product.id"
            class="product-item"
          >
            <i
              :class="[
                'mdi',
                product.kind === 'stempel'
                  ? 'mdi-stamper'
                  : product.kind === 'stanze'
                  ? 'mdi-content-cut'
                  : 'mdi-help-circle-outline',
              ]"
            ></i>
            <span class="product-name">{{ product.name }}</span>
            <span v-if="product.kind" class="product-kind">{{
              product.kind
            }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.overview {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: #f5f5f5;
  overflow: hidden;
}

.overview-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 20px;
  background: white;
  border-bottom: 1px solid #eee;
}

.btn-back {
  width: 40px;
  height: 40px;
  border: 1px solid #ddd;
  border-radius: 10px;
  background: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-back:hover {
  border-color: #111;
  background: #f8f8f8;
}

.btn-back i {
  font-size: 22px;
  color: #333;
}

.header-actions {
  display: flex;
  gap: 8px;
}

.btn-action {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 16px;
  font-size: 14px;
  font-weight: 500;
  border: 1px solid #ddd;
  border-radius: 10px;
  background: white;
  cursor: pointer;
}

.btn-action:hover {
  border-color: #111;
  background: #f8f8f8;
}

.btn-action i {
  font-size: 18px;
}

.btn-action.danger {
  border-color: #fcc;
  color: #c00;
}

.btn-action.danger:hover {
  background: #fee;
}

.overview-content {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
}

/* Hero */
.hero {
  display: flex;
  gap: 24px;
  background: white;
  border-radius: 16px;
  padding: 24px;
  margin-bottom: 20px;
  border: 1px solid #eee;
}

.hero-image {
  width: 200px;
  height: 200px;
  border-radius: 12px;
  overflow: hidden;
  flex-shrink: 0;
  background: #f5f5f5;
}

.hero-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.hero-placeholder {
  width: 200px;
  height: 200px;
  border-radius: 12px;
  background: #f5f5f5;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.hero-placeholder i {
  font-size: 64px;
  color: #ccc;
}

.hero-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.set-name {
  font-size: 28px;
  font-weight: 700;
  margin: 0 0 12px;
  color: #111;
}

.set-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 15px;
  color: #666;
}

.meta-item i {
  font-size: 18px;
  color: #999;
}

/* Info Grid */
.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 12px;
  margin-bottom: 20px;
}

.info-card {
  display: flex;
  align-items: center;
  gap: 14px;
  background: white;
  border-radius: 12px;
  padding: 16px;
  border: 1px solid #eee;
}

.card-icon {
  width: 44px;
  height: 44px;
  border-radius: 10px;
  background: #f5f5f5;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.card-icon i {
  font-size: 22px;
  color: #666;
}

.card-content {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.card-label {
  font-size: 12px;
  color: #999;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.card-value {
  font-size: 14px;
  font-weight: 500;
  color: #333;
}

.product-badge {
  display: inline-block;
  font-size: 11px;
  padding: 2px 6px;
  background: #f0f0f0;
  border-radius: 4px;
  margin-left: 6px;
  font-weight: 400;
}

/* Sections */
.section {
  background: white;
  border-radius: 12px;
  padding: 20px;
  border: 1px solid #eee;
  margin-bottom: 16px;
}

.section h2 {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 15px;
  font-weight: 600;
  margin: 0 0 16px;
  color: #333;
}

.section h2 i {
  font-size: 20px;
  color: #666;
}

/* Tags */
.tags-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.tag {
  padding: 6px 12px;
  background: #f0f0f0;
  border-radius: 20px;
  font-size: 13px;
  font-weight: 500;
  color: #555;
}

/* Products */
.products-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 10px;
}

.product-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 14px;
  background: #fafafa;
  border-radius: 10px;
  border: 1px solid #eee;
}

.product-item i {
  font-size: 20px;
  color: #888;
}

.product-name {
  flex: 1;
  font-size: 14px;
  font-weight: 500;
}

.product-kind {
  font-size: 11px;
  padding: 3px 8px;
  background: #e8e8e8;
  border-radius: 4px;
  color: #666;
  text-transform: capitalize;
}

/* Responsive */
@media (max-width: 600px) {
  .hero {
    flex-direction: column;
    align-items: center;
    text-align: center;
  }

  .hero-image,
  .hero-placeholder {
    width: 160px;
    height: 160px;
  }

  .set-meta {
    justify-content: center;
  }

  .btn-action span {
    display: none;
  }
}
</style>
