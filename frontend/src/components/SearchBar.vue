<script setup lang="ts">
import { ref, computed } from "vue";
import { useI18n } from "../i18n";

const props = defineProps<{
  modelValue: string;
  sortBy: "name" | "box" | "location" | "added";
  loading?: boolean;
  resultCount?: number;
}>();

const emit = defineEmits<{
  "update:modelValue": [value: string];
  "update:sortBy": [value: "name" | "box" | "location" | "added"];
  "new-set": [];
}>();

const { t, locale } = useI18n();
const inputRef = ref<HTMLInputElement | null>(null);

const sortOptions = computed(
  () =>
    [
      {
        value: "name",
        label: t("sortName"),
        icon: "mdi-sort-alphabetical-ascending",
      },
      { value: "box", label: "Box", icon: "mdi-package-variant" },
      { value: "location", label: t("sortLocation"), icon: "mdi-map-marker" },
      {
        value: "added",
        label: locale.value === "de" ? "Neueste" : "Newest",
        icon: "mdi-clock-outline",
      },
    ] as const
);

function focus() {
  inputRef.value?.focus();
}

defineExpose({ focus });
</script>

<template>
  <div class="search-bar">
    <div class="search-wrapper">
      <i class="mdi mdi-magnify search-icon"></i>
      <input
        ref="inputRef"
        type="text"
        :value="modelValue"
        @input="
          emit('update:modelValue', ($event.target as HTMLInputElement).value)
        "
        :placeholder="t('searchPlaceholder')"
        class="search-input"
      />
      <span v-if="loading" class="search-spinner">
        <i class="mdi mdi-loading mdi-spin"></i>
      </span>
      <span v-else-if="resultCount !== undefined" class="result-count">
        {{ resultCount }}
      </span>
    </div>

    <div class="sort-wrapper">
      <i class="mdi mdi-sort sort-icon"></i>
      <select
        :value="sortBy"
        @change="
          emit(
            'update:sortBy',
            ($event.target as HTMLSelectElement).value as any
          )
        "
        class="sort-select"
      >
        <option v-for="opt in sortOptions" :key="opt.value" :value="opt.value">
          {{ opt.label }}
        </option>
      </select>
    </div>

    <button class="btn-new" @click="emit('new-set')">
      <i class="mdi mdi-plus"></i>
      <span>{{ locale === "de" ? "Neu" : "New" }}</span>
    </button>
  </div>
</template>

<style scoped>
.search-bar {
  display: flex;
  gap: 12px;
  padding: 16px 20px;
  background: white;
  border-bottom: 1px solid #e5e5e5;
  position: sticky;
  top: 0;
  z-index: 100;
}

.search-wrapper {
  flex: 1;
  position: relative;
  display: flex;
  align-items: center;
}

.search-icon {
  position: absolute;
  left: 14px;
  font-size: 22px;
  color: #888;
  pointer-events: none;
}

.search-input {
  width: 100%;
  padding: 14px 50px 14px 48px;
  font-size: 16px;
  border: 2px solid #e5e5e5;
  border-radius: 12px;
  background: #f8f8f8;
  transition: all 0.2s;
}

.search-input:focus {
  outline: none;
  border-color: #111;
  background: white;
}

.search-input::placeholder {
  color: #999;
}

.search-spinner,
.result-count {
  position: absolute;
  right: 14px;
  font-size: 14px;
  color: #666;
}

.search-spinner i {
  font-size: 20px;
}

.sort-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.sort-icon {
  position: absolute;
  left: 12px;
  font-size: 18px;
  color: #666;
  pointer-events: none;
}

.sort-select {
  padding: 14px 14px 14px 38px;
  font-size: 14px;
  border: 2px solid #e5e5e5;
  border-radius: 12px;
  background: #f8f8f8;
  cursor: pointer;
  min-width: 120px;
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 24 24'%3E%3Cpath fill='%23666' d='M7 10l5 5 5-5z'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 10px center;
}

.sort-select:focus {
  outline: none;
  border-color: #111;
  background-color: white;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 24 24'%3E%3Cpath fill='%23666' d='M7 10l5 5 5-5z'/%3E%3C/svg%3E");
}

.btn-new {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 14px 20px;
  background: #111;
  color: white;
  border: none;
  border-radius: 12px;
  font-size: 15px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-new:hover {
  background: #333;
}

.btn-new i {
  font-size: 20px;
}

@media (max-width: 600px) {
  .btn-new span {
    display: none;
  }

  .btn-new {
    padding: 14px;
  }

  .sort-wrapper {
    display: none;
  }
}
</style>
