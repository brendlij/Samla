<script setup lang="ts">
import { ref, computed, watch } from "vue";

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
  if (
    locationForm.value.id &&
    confirm(
      "Standort wirklich löschen? Alle Boxen und Sets werden auch gelöscht."
    )
  ) {
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
  if (
    boxForm.value.id &&
    confirm("Box wirklich löschen? Alle Sets werden auch gelöscht.")
  ) {
    emit("delete-box", boxForm.value.id);
    showBoxForm.value = false;
  }
}

function getLocationLabel(loc: Location) {
  const parts = [loc.friendlyName];
  if (loc.room) parts.push(loc.room);
  if (loc.shelf) parts.push(`Regal ${loc.shelf}`);
  if (loc.compartment) parts.push(`Fach ${loc.compartment}`);
  return parts.join(" · ");
}
</script>

<template>
  <div class="location-picker">
    <!-- Standort Auswahl -->
    <div class="field">
      <label>Standort</label>
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
          <option :value="null">Standort wählen...</option>
          <option v-for="loc in locations" :key="loc.id" :value="loc.id">
            {{ getLocationLabel(loc) }}
          </option>
        </select>
        <button
          v-if="locationId"
          class="btn-icon"
          @click="editLocation"
          title="Bearbeiten"
        >
          <i class="mdi mdi-pencil"></i>
        </button>
        <button class="btn-icon" @click="newLocation" title="Neuer Standort">
          <i class="mdi mdi-plus"></i>
        </button>
      </div>
    </div>

    <!-- Box Auswahl -->
    <div class="field">
      <label>Box</label>
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
          <option :value="null">Box wählen...</option>
          <option v-for="box in filteredBoxes" :key="box.id" :value="box.id">
            {{ box.code }}{{ box.name ? ` – ${box.name}` : "" }}
          </option>
        </select>
        <button
          v-if="boxId"
          class="btn-icon"
          @click="editBox"
          title="Bearbeiten"
        >
          <i class="mdi mdi-pencil"></i>
        </button>
        <button
          class="btn-icon"
          @click="newBox"
          :disabled="!locationId"
          title="Neue Box"
        >
          <i class="mdi mdi-plus"></i>
        </button>
      </div>
    </div>

    <!-- Beutel Nummer -->
    <div class="field">
      <label>Beutel-Nr.</label>
      <input
        type="text"
        :value="bagSerial"
        @input="
          emit('update:bagSerial', ($event.target as HTMLInputElement).value)
        "
        placeholder="z.B. 01, A1..."
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
          {{ locationForm.id ? "Standort bearbeiten" : "Neuer Standort" }}
        </h3>

        <div class="field">
          <label>Name *</label>
          <input
            v-model="locationForm.name"
            type="text"
            placeholder="z.B. Wohnzimmer"
            class="input"
          />
        </div>

        <div class="field-row">
          <div class="field">
            <label>Raum</label>
            <input
              v-model="locationForm.room"
              type="text"
              placeholder="z.B. Büro"
              class="input"
            />
          </div>
          <div class="field">
            <label>Regal</label>
            <input
              v-model="locationForm.shelf"
              type="text"
              placeholder="z.B. A"
              class="input"
            />
          </div>
          <div class="field">
            <label>Fach</label>
            <input
              v-model="locationForm.compartment"
              type="text"
              placeholder="z.B. 1"
              class="input"
            />
          </div>
        </div>

        <div class="field">
          <label>Notiz</label>
          <input
            v-model="locationForm.note"
            type="text"
            placeholder="Zusätzliche Info..."
            class="input"
          />
        </div>

        <div class="modal-actions">
          <button
            v-if="locationForm.id"
            class="btn btn-danger"
            @click="deleteLocation"
          >
            <i class="mdi mdi-delete"></i> Löschen
          </button>
          <div class="spacer"></div>
          <button class="btn btn-secondary" @click="showLocationForm = false">
            Abbrechen
          </button>
          <button
            class="btn btn-primary"
            @click="saveLocation"
            :disabled="!locationForm.name.trim()"
          >
            Speichern
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
        <h3>{{ boxForm.id ? "Box bearbeiten" : "Neue Box" }}</h3>

        <div class="field-row">
          <div class="field">
            <label>Code *</label>
            <input
              v-model="boxForm.code"
              type="text"
              placeholder="z.B. A, B, 1..."
              class="input"
            />
          </div>
          <div class="field">
            <label>Name</label>
            <input
              v-model="boxForm.name"
              type="text"
              placeholder="Optional..."
              class="input"
            />
          </div>
        </div>

        <div class="modal-actions">
          <button v-if="boxForm.id" class="btn btn-danger" @click="deleteBox">
            <i class="mdi mdi-delete"></i> Löschen
          </button>
          <div class="spacer"></div>
          <button class="btn btn-secondary" @click="showBoxForm = false">
            Abbrechen
          </button>
          <button
            class="btn btn-primary"
            @click="saveBox"
            :disabled="!boxForm.code.trim()"
          >
            Speichern
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
