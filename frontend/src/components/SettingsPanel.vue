<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useI18n, type Locale } from "../i18n";

type AppPaths = {
  baseDir: string;
  dataDir: string;
  imagesDir: string;
  dbPath: string;
};

type Stats = Record<string, number>;

const props = defineProps<{
  visible: boolean;
  paths: AppPaths | null;
  stats: Stats | null;
}>();

const emit = defineEmits<{
  close: [];
  "open-folder": [];
  export: [];
  import: [];
}>();

const { t, locale, setLocale } = useI18n();
const appVersion = "1.0.0";

function handleLanguageChange(e: Event) {
  const value = (e.target as HTMLSelectElement).value as Locale;
  setLocale(value);
}
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="visible" class="modal-overlay" @click.self="emit('close')">
        <div class="settings-panel">
          <div class="panel-header">
            <h2><i class="mdi mdi-cog"></i> {{ t("settingsTitle") }}</h2>
            <button class="btn-close" @click="emit('close')">
              <i class="mdi mdi-close"></i>
            </button>
          </div>

          <div class="panel-content">
            <!-- Language -->
            <section class="settings-section">
              <h3><i class="mdi mdi-translate"></i> {{ t("language") }}</h3>
              <div class="language-selector">
                <select
                  :value="locale"
                  @change="handleLanguageChange"
                  class="lang-select"
                >
                  <option value="de">🇩🇪 {{ t("german") }}</option>
                  <option value="en">🇬🇧 {{ t("english") }}</option>
                </select>
              </div>
            </section>

            <!-- Data Management -->
            <section class="settings-section">
              <h3>
                <i class="mdi mdi-database"></i> {{ t("dataManagement") }}
              </h3>

              <div class="action-buttons">
                <button class="action-btn" @click="emit('export')">
                  <i class="mdi mdi-export"></i>
                  <span>{{ t("exportData") }}</span>
                  <small>{{
                    locale === "de"
                      ? "Datenbank & Bilder als ZIP speichern"
                      : "Save database & images as ZIP"
                  }}</small>
                </button>

                <button class="action-btn" @click="emit('import')">
                  <i class="mdi mdi-import"></i>
                  <span>{{ t("importData") }}</span>
                  <small>{{
                    locale === "de"
                      ? "Aus ZIP-Backup wiederherstellen"
                      : "Restore from ZIP backup"
                  }}</small>
                </button>

                <button class="action-btn" @click="emit('open-folder')">
                  <i class="mdi mdi-folder-open"></i>
                  <span>{{ t("openDataFolder") }}</span>
                  <small>{{
                    locale === "de"
                      ? "Dateien im Explorer anzeigen"
                      : "View files in explorer"
                  }}</small>
                </button>
              </div>
            </section>

            <!-- Storage Paths -->
            <section class="settings-section">
              <h3>
                <i class="mdi mdi-folder-cog"></i> {{ t("storagePaths") }}
              </h3>

              <div class="path-info" v-if="paths">
                <div class="path-item">
                  <label>{{ t("baseDirectory") }}</label>
                  <code>{{ paths.baseDir }}</code>
                </div>
                <div class="path-item">
                  <label>{{ t("database") }}</label>
                  <code>{{ paths.dbPath }}</code>
                </div>
                <div class="path-item">
                  <label>{{ t("images") }}</label>
                  <code>{{ paths.imagesDir }}</code>
                </div>
              </div>
            </section>

            <!-- Statistics -->
            <section class="settings-section" v-if="stats">
              <h3><i class="mdi mdi-chart-box"></i> {{ t("statistics") }}</h3>

              <div class="stats-grid">
                <div class="stat-item">
                  <span class="stat-value">{{ stats.sets || 0 }}</span>
                  <span class="stat-label">{{ t("statsSets") }}</span>
                </div>
                <div class="stat-item">
                  <span class="stat-value">{{ stats.products || 0 }}</span>
                  <span class="stat-label">{{ t("statsProducts") }}</span>
                </div>
                <div class="stat-item">
                  <span class="stat-value">{{ stats.boxes || 0 }}</span>
                  <span class="stat-label">{{
                    locale === "de" ? "Kartons" : "Boxes"
                  }}</span>
                </div>
                <div class="stat-item">
                  <span class="stat-value">{{ stats.locations || 0 }}</span>
                  <span class="stat-label">{{
                    locale === "de" ? "Orte" : "Locations"
                  }}</span>
                </div>
                <div class="stat-item">
                  <span class="stat-value">{{ stats.tags || 0 }}</span>
                  <span class="stat-label">{{ t("statsTags") }}</span>
                </div>
                <div class="stat-item">
                  <span class="stat-value">{{ stats.images || 0 }}</span>
                  <span class="stat-label">{{ t("statsImages") }}</span>
                </div>
              </div>
            </section>

            <!-- About -->
            <section class="settings-section">
              <h3><i class="mdi mdi-information"></i> {{ t("about") }}</h3>
              <div class="about-info">
                <p><strong>Samla</strong> v{{ appVersion }}</p>
                <p class="about-desc">
                  {{
                    locale === "de"
                      ? "Ein universeller Sammlungsorganizer für die Verwaltung Ihrer Artikel, Lagerorte und Bestände."
                      : "A universal collection organizer for managing your items, storage locations, and inventory."
                  }}
                </p>
                <p class="about-tech">Built with Wails, Vue 3, and SQLite</p>
              </div>
            </section>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.settings-panel {
  background: white;
  border-radius: 16px;
  width: 90%;
  max-width: 520px;
  max-height: 85vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid #eee;
}

.panel-header h2 {
  display: flex;
  align-items: center;
  gap: 10px;
  margin: 0;
  font-size: 20px;
  font-weight: 600;
}

.panel-header h2 i {
  font-size: 24px;
  color: #666;
}

.btn-close {
  width: 36px;
  height: 36px;
  border: none;
  background: #f5f5f5;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-close:hover {
  background: #eee;
}

.btn-close i {
  font-size: 20px;
  color: #666;
}

.panel-content {
  flex: 1;
  overflow-y: auto;
  padding: 20px 24px;
}

.settings-section {
  margin-bottom: 28px;
}

.settings-section:last-child {
  margin-bottom: 0;
}

.settings-section h3 {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 0 0 16px;
  font-size: 14px;
  font-weight: 600;
  color: #666;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.settings-section h3 i {
  font-size: 18px;
}

/* Action Buttons */
.action-buttons {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 14px 16px;
  background: #f8f8f8;
  border: 1px solid #e5e5e5;
  border-radius: 12px;
  cursor: pointer;
  text-align: left;
  transition: all 0.2s;
}

.action-btn:hover {
  background: #f0f0f0;
  border-color: #ddd;
}

.action-btn i {
  font-size: 24px;
  color: #555;
}

.action-btn span {
  font-size: 15px;
  font-weight: 500;
  color: #333;
}

.action-btn small {
  display: block;
  font-size: 12px;
  color: #888;
  margin-top: 2px;
}

/* Path Info */
.path-info {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.path-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.path-item label {
  font-size: 12px;
  color: #888;
  font-weight: 500;
}

.path-item code {
  font-size: 12px;
  padding: 8px 10px;
  background: #f5f5f5;
  border-radius: 6px;
  color: #555;
  word-break: break-all;
  font-family: "Consolas", "Monaco", monospace;
}

/* Stats Grid */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 14px;
  background: #f8f8f8;
  border-radius: 10px;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: #333;
}

.stat-label {
  font-size: 11px;
  color: #888;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-top: 2px;
}

/* About */
.about-info {
  text-align: center;
  padding: 16px;
  background: #f8f8f8;
  border-radius: 12px;
}

.about-info p {
  margin: 0;
}

.about-info strong {
  font-size: 18px;
}

.about-desc {
  margin-top: 8px !important;
  font-size: 13px;
  color: #666;
}

.about-tech {
  margin-top: 12px !important;
  font-size: 11px;
  color: #999;
}

/* Language Selector */
.language-selector {
  display: flex;
  justify-content: center;
}

.lang-select {
  padding: 10px 16px;
  font-size: 14px;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  background: #f8f8f8;
  cursor: pointer;
  min-width: 160px;
}

.lang-select:hover {
  border-color: #ccc;
}

.lang-select:focus {
  outline: none;
  border-color: #999;
}

/* Transitions */
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.2s ease;
}

.modal-enter-active .settings-panel,
.modal-leave-active .settings-panel {
  transition: transform 0.2s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .settings-panel,
.modal-leave-to .settings-panel {
  transform: scale(0.95);
}
</style>
