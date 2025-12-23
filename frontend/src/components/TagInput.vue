<script setup lang="ts">
import { ref, computed } from "vue";

const props = defineProps<{
  tags: string[];
  suggestions: string[];
}>();

const emit = defineEmits<{
  "update:tags": [tags: string[]];
}>();

const inputValue = ref("");

function addTag() {
  const tag = inputValue.value.trim().toLowerCase();
  if (!tag) return;
  if (!props.tags.includes(tag)) {
    emit("update:tags", [...props.tags, tag]);
  }
  inputValue.value = "";
}

function removeTag(tag: string) {
  emit(
    "update:tags",
    props.tags.filter((t) => t !== tag)
  );
}

function addSuggestion(tag: string) {
  if (!props.tags.includes(tag)) {
    emit("update:tags", [...props.tags, tag]);
  }
}

const filteredSuggestions = computed(() =>
  props.suggestions.filter((s) => !props.tags.includes(s)).slice(0, 8)
);
</script>

<template>
  <div class="tag-input">
    <!-- Current Tags -->
    <div class="tags" v-if="tags.length">
      <span v-for="tag in tags" :key="tag" class="tag" @click="removeTag(tag)">
        {{ tag }}
        <i class="mdi mdi-close"></i>
      </span>
    </div>

    <!-- Input -->
    <div class="input-row">
      <input
        v-model="inputValue"
        type="text"
        placeholder="Tag hinzufügen..."
        class="input"
        @keyup.enter="addTag"
      />
      <button class="btn-add" @click="addTag" :disabled="!inputValue.trim()">
        <i class="mdi mdi-plus"></i>
      </button>
    </div>

    <!-- Suggestions -->
    <div class="suggestions" v-if="filteredSuggestions.length">
      <span class="label">Vorschläge:</span>
      <span
        v-for="tag in filteredSuggestions"
        :key="tag"
        class="suggestion"
        @click="addSuggestion(tag)"
      >
        {{ tag }}
      </span>
    </div>
  </div>
</template>

<style scoped>
.tag-input {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.tag {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 10px;
  background: #111;
  color: white;
  border-radius: 6px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.15s;
}

.tag:hover {
  background: #c00;
}

.tag i {
  font-size: 14px;
}

.input-row {
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

.btn-add {
  width: 42px;
  height: 42px;
  border: 1px solid #ddd;
  border-radius: 8px;
  background: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-add:hover:not(:disabled) {
  border-color: #111;
  background: #f5f5f5;
}

.btn-add:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

.btn-add i {
  font-size: 20px;
  color: #333;
}

.suggestions {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  align-items: center;
}

.label {
  font-size: 12px;
  color: #888;
}

.suggestion {
  padding: 4px 10px;
  background: #f0f0f0;
  border-radius: 6px;
  font-size: 12px;
  color: #666;
  cursor: pointer;
  transition: all 0.15s;
}

.suggestion:hover {
  background: #e0e0e0;
  color: #111;
}
</style>
