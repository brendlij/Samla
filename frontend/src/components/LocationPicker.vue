<script setup lang="ts">
import { ref, computed, watch } from "vue";
import { useI18n } from "../i18n";

type Location = {
  id: number;
  friendlyName: string;
  room: string;
  shelf: string;
  compartment: string;
  note: string;
};

type Box = {
  id: number;
  locationId: number;
  code: string;
  name: string;
};

const props = defineProps<{
  locations: Location[];
  boxes: Box[];
  locationId: number | null;
  boxId: number | null;
  bagSerial: string;
}>();

const emit = defineEmits<{
  "update:locationId": [value: number | null];
  "update:boxId": [value: number | null];
  "update:bagSerial": [value: string];
  "save-location": [
    data: {
      id: number | null;
      name: string;
      room: string;
      shelf: string;
      compartment: string;
      note: string;
    }
  ];
  "delete-location": [id: number];
  "save-box": [
    data: { id: number | null; locationId: number; code: string; name: string }
  ];
  "delete-box": [id: number];
}>();

const { t, locale } = useI18n();

const filteredBoxes = computed(() =>
  props.locationId
    ? props.boxes.filter((b) => b.locationId === props.locationId)
    : []
);

const showLocationForm = ref(false);
const showBoxForm = ref(false);

const locationForm = ref({
  id: null as number | null,
  name: "",
  room: "",
  shelf: "",
  compartment: "",
  note: "",
});

const boxForm = ref({
  id: null as number | null,
  code: "",
  name: "",
});

watch(
  () => props.locationId,
  (id) => {
    if (id) {
      const loc = props.locations.find((l) => l.id === id);
      if (loc) {
        locationForm.value = {
          id: loc.id,
          name: loc.friendlyName,
          room: loc.room || "",
          shelf: loc.shelf || "",
          compartment: loc.compartment || "",
          note: loc.note || "",
        };
      }
    }
  }
);

watch(
  () => props.boxId,
  (id) => {
    if (id) {
      const box = props.boxes.find((b) => b.id === id);
      if (box) {
        boxForm.value = {
          id: box.id,
          code: box.code,
          name: box.name || "",
        };
      }
    }
  }
);

function newLocation() {
  locationForm.value = {
    id: null,
    name: "",
    room: "",
    shelf: "",
    compartment: "",
    note: "",
  };
  showLocationForm.value = true;
}

function editLocation() {
  showLocationForm.value = true;
}

function saveLocation() {
  emit("save-location", { ...locationForm.value });
  showLocationForm.value = false;
}

function deleteLocation() {
  const msg =
    locale.value === "de"
      ? "Standort wirklich löschen? Alle Boxen und Sets werden auch gelöscht."
      : "Really delete this location? All boxes and sets will also be deleted.";
  if (locationForm.value.id && confirm(msg)) {
    emit("delete-location", locationForm.value.id);
    showLocationForm.value = false;
  }
}

function newBox() {
  boxForm.value = { id: null, code: "", name: "" };
  showBoxForm.value = true;
}

function editBox() {
  showBoxForm.value = true;
}

function saveBox() {
  if (props.locationId) {
    emit("save-box", { ...boxForm.value, locationId: props.locationId });
    showBoxForm.value = false;
  }
}

function deleteBox() {
  const msg =
    locale.value === "de"
      ? "Box wirklich löschen? Alle Sets werden auch gelöscht."
      : "Really delete this box? All sets will also be deleted.";
  if (boxForm.value.id && confirm(msg)) {
    emit("delete-box", boxForm.value.id);
    showBoxForm.value = false;
  }
}

function getLocationLabel(loc: Location) {
  const parts = [loc.friendlyName];
  if (loc.room) parts.push(loc.room);
  const shelfLabel = locale.value === "de" ? "Regal" : "Shelf";
  const compLabel = locale.value === "de" ? "Fach" : "Comp.";
  if (loc.shelf) parts.push(`${shelfLabel} ${loc.shelf}`);
  if (loc.compartment) parts.push(`${compLabel} ${loc.compartment}`);
  return parts.join(" · ");
}
</script>

<template>
  <div class="location-picker">
    <!-- Standort Auswahl -->
    <div class="field">
      <label>{{ t("location") }}</label>
      <div class="select-row">
        <select
          :value="locationId"
          @change="
            emit(
              'update:locationId',
              Number(($event.target as HTMLSelectElement).value) || null
            )
          "
          class="select"
        >
          <option :value="null">
            {{ locale === "de" ? "Standort wählen..." : "Select location..." }}
          </option>
          <option v-for="loc in locations" :key="loc.id" :value="loc.id">
            {{ getLocationLabel(loc) }}
          </option>
        </select>
        <button
          v-if="locationId"
          class="btn-icon"
          @click="editLocation"
          :title="t('edit')"
        >
          <i class="mdi mdi-pencil"></i>
        </button>
        <button
          class="btn-icon"
          @click="newLocation"
          :title="locale === 'de' ? 'Neuer Standort' : 'New location'"
        >
          <i class="mdi mdi-plus"></i>
        </button>
      </div>
    </div>

    <!-- Box Auswahl -->
    <div class="field">
      <label>{{ t("boxNumber") }}</label>
      <div class="select-row">
        <select
          :value="boxId"
          @change="
            emit(
              'update:boxId',
              Number(($event.target as HTMLSelectElement).value) || null
            )
          "
          :disabled="!locationId"
          class="select"
        >
          <option :value="null">
            {{ locale === "de" ? "Box wählen..." : "Select box..." }}
          </option>
          <option v-for="box in filteredBoxes" :key="box.id" :value="box.id">
            {{ box.code }}{{ box.name ? ` – ${box.name}` : "" }}
          </option>
        </select>
        <button
          v-if="boxId"
          class="btn-icon"
          @click="editBox"
          :title="t('edit')"
        >
          <i class="mdi mdi-pencil"></i>
        </button>
        <button
          class="btn-icon"
          @click="newBox"
          :disabled="!locationId"
          :title="locale === 'de' ? 'Neue Box' : 'New box'"
        >
          <i class="mdi mdi-plus"></i>
        </button>
      </div>
    </div>

    <!-- Beutel Nummer -->
    <div class="field">
      <label>{{ t("bagNumber") }}</label>
      <input
        type="text"
        :value="bagSerial"
        @input="
          emit('update:bagSerial', ($event.target as HTMLInputElement).value)
        "
        :placeholder="locale === 'de' ? 'z.B. 01, A1...' : 'e.g. 01, A1...'"
        class="input"
      />
    </div>

    <!-- Standort Form Modal -->
    <div
      v-if="showLocationForm"
      class="modal-overlay"
      @click.self="showLocationForm = false"
    >
      <div class="modal">
        <h3>
          {{
            locationForm.id
              ? locale === "de"
                ? "Standort bearbeiten"
                : "Edit Location"
              : locale === "de"
              ? "Neuer Standort"
              : "New Location"
          }}
        </h3>

        <div class="field">
          <label>{{ locale === "de" ? "Name" : "Name" }} *</label>
          <input
            v-model="locationForm.name"
            type="text"
            :placeholder="
              locale === 'de' ? 'z.B. Wohnzimmer' : 'e.g. Living Room'
            "
            class="input"
          />
        </div>

        <div class="field-row">
          <div class="field">
            <label>{{ t("room") }}</label>
            <input
              v-model="locationForm.room"
              type="text"
              :placeholder="t('roomPlaceholder')"
              class="input"
            />
          </div>
          <div class="field">
            <label>{{ t("shelf") }}</label>
            <input
              v-model="locationForm.shelf"
              type="text"
              :placeholder="locale === 'de' ? 'z.B. A' : 'e.g. A'"
              class="input"
            />
          </div>
          <div class="field">
            <label>{{ t("compartment") }}</label>
            <input
              v-model="locationForm.compartment"
              type="text"
              :placeholder="locale === 'de' ? 'z.B. 1' : 'e.g. 1'"
              class="input"
            />
          </div>
        </div>

        <div class="field">
          <label>{{ locale === "de" ? "Notiz" : "Note" }}</label>
          <input
            v-model="locationForm.note"
            type="text"
            :placeholder="
              locale === 'de' ? 'Zusätzliche Info...' : 'Additional info...'
            "
            class="input"
          />
        </div>

        <div class="modal-actions">
          <button
            v-if="locationForm.id"
            class="btn btn-danger"
            @click="deleteLocation"
          >
            <i class="mdi mdi-delete"></i> {{ t("delete") }}
          </button>
          <div class="spacer"></div>
          <button class="btn btn-secondary" @click="showLocationForm = false">
            {{ t("cancel") }}
          </button>
          <button
            class="btn btn-primary"
            @click="saveLocation"
            :disabled="!locationForm.name.trim()"
          >
            {{ t("save") }}
          </button>
        </div>
      </div>
    </div>

    <!-- Box Form Modal -->
    <div
      v-if="showBoxForm"
      class="modal-overlay"
      @click.self="showBoxForm = false"
    >
      <div class="modal">
        <h3>
          {{
            boxForm.id
              ? locale === "de"
                ? "Box bearbeiten"
                : "Edit Box"
              : locale === "de"
              ? "Neue Box"
              : "New Box"
          }}
        </h3>

        <div class="field-row">
          <div class="field">
            <label>Code *</label>
            <input
              v-model="boxForm.code"
              type="text"
              :placeholder="
                locale === 'de' ? 'z.B. A, B, 1...' : 'e.g. A, B, 1...'
              "
              class="input"
            />
          </div>
          <div class="field">
            <label>{{ locale === "de" ? "Name" : "Name" }}</label>
            <input
              v-model="boxForm.name"
              type="text"
              :placeholder="locale === 'de' ? 'Optional...' : 'Optional...'"
              class="input"
            />
          </div>
        </div>

        <div class="modal-actions">
          <button v-if="boxForm.id" class="btn btn-danger" @click="deleteBox">
            <i class="mdi mdi-delete"></i> {{ t("delete") }}
          </button>
          <div class="spacer"></div>
          <button class="btn btn-secondary" @click="showBoxForm = false">
            {{ t("cancel") }}
          </button>
          <button
            class="btn btn-primary"
            @click="saveBox"
            :disabled="!boxForm.code.trim()"
          >
            {{ t("save") }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.location-picker {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.field label {
  font-size: 13px;
  font-weight: 500;
  color: #666;
}

.field-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
}

.select-row {
  display: flex;
  gap: 8px;
}

.select,
.input {
  flex: 1;
  padding: 12px 14px;
  font-size: 15px;
  border: 1px solid #ddd;
  border-radius: 8px;
  background: white;
}

.select:focus,
.input:focus {
  outline: none;
  border-color: #111;
}

.select:disabled {
  background: #f5f5f5;
  color: #999;
}

.btn-icon {
  width: 44px;
  height: 44px;
  border: 1px solid #ddd;
  border-radius: 8px;
  background: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s;
}

.btn-icon:hover:not(:disabled) {
  border-color: #111;
  background: #f5f5f5;
}

.btn-icon:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.btn-icon i {
  font-size: 20px;
  color: #333;
}

/* Modal */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.modal {
  background: white;
  border-radius: 16px;
  padding: 24px;
  width: 100%;
  max-width: 480px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.2);
}

.modal h3 {
  font-size: 18px;
  font-weight: 600;
  margin: 0 0 20px 0;
}

.modal .field {
  margin-bottom: 16px;
}

.modal-actions {
  display: flex;
  gap: 10px;
  margin-top: 24px;
  padding-top: 16px;
  border-top: 1px solid #eee;
}

.spacer {
  flex: 1;
}

.btn {
  padding: 10px 18px;
  font-size: 14px;
  font-weight: 500;
  border-radius: 8px;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 6px;
  transition: all 0.15s;
}

.btn-primary {
  background: #111;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: #333;
}

.btn-primary:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.btn-secondary {
  background: #f0f0f0;
  color: #333;
}

.btn-secondary:hover {
  background: #e5e5e5;
}

.btn-danger {
  background: #fee;
  color: #c00;
}

.btn-danger:hover {
  background: #fdd;
}
</style>
