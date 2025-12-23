<script lang="ts" setup>
import {
  computed,
  nextTick,
  onBeforeUnmount,
  onMounted,
  reactive,
  ref,
  watch,
} from "vue";
import Fuse from "fuse.js";
import { useI18n, getSearchPrefixes } from "./i18n";
import SearchBar from "./components/SearchBar.vue";
import SetCard from "./components/SetCard.vue";
import SetOverview from "./components/SetOverview.vue";
import EmptyState from "./components/EmptyState.vue";
import LocationPicker from "./components/LocationPicker.vue";
import ProductList from "./components/ProductList.vue";
import TagInput from "./components/TagInput.vue";
import ImageUpload from "./components/ImageUpload.vue";
import CropModal from "./components/CropModal.vue";
import ConfirmModal from "./components/ConfirmModal.vue";
import MasterDataPanel from "./components/MasterDataPanel.vue";
import SettingsPanel from "./components/SettingsPanel.vue";

import {
  AddProduct,
  AttachImageFromFile,
  AttachImageFromURL,
  AttachScannedImage,
  ChooseImageFile,
  CreateBagWithSet,
  CreateBox,
  CreateLocation,
  CreateManufacturer,
  CreateTag,
  CreateType,
  DeleteBox,
  DeleteManufacturer,
  DeleteProduct,
  DeleteLocation,
  DeleteSet,
  DeleteTag,
  DeleteType,
  ExportData,
  GetAppPaths,
  GetImageAsBase64,
  GetNextBagSerial,
  GetSet,
  GetStats,
  ImportData,
  ListBoxes,
  ListLocations,
  ListManufacturers,
  ListTags,
  ListTagsFull,
  ListTypes,
  OpenAppFolder,
  ReadFileAsBase64,
  RemoveImage,
  SaveCroppedImage,
  ScanImageToBase64,
  SearchSets,
  SetTags,
  UpdateBox,
  UpdateManufacturer,
  UpdateProduct,
  UpdateLocation,
  UpdateSet,
  UpdateTag,
  UpdateType,
} from "../wailsjs/go/main/App";

import samlaIcon from "./assets/images/samla-icon.svg";

// Types
type SearchResult = {
  setId: number;
  setName: string;
  manufacturerName: string;
  boxCode: string;
  boxName: string;
  bagSerial: string;
  locationName: string;
  tags: string[];
  thumbnailPath: string;
};

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

type ProductItem = {
  id: number;
  setId: number;
  name: string;
  kind: string;
};

type AppPaths = {
  baseDir: string;
  dataDir: string;
  imagesDir: string;
  dbPath: string;
};

type MasterDataItem = { id: number; name: string };

// State
const view = ref<"list" | "overview" | "detail">("list");
const searchQuery = ref("");
const sortBy = ref<"name" | "box" | "location" | "added">("name");
const allSets = ref<SearchResult[]>([]); // All sets from backend
const searchResults = ref<SearchResult[]>([]); // Filtered results
const searchLoading = ref(false);
const searchBarRef = ref<InstanceType<typeof SearchBar> | null>(null);

// Fuse.js instance for fuzzy search
const fuseInstance = ref<Fuse<SearchResult> | null>(null);
const fuseOptions = {
  keys: [
    { name: "setName", weight: 0.3 },
    { name: "manufacturerName", weight: 0.15 },
    { name: "boxCode", weight: 0.15 },
    { name: "boxName", weight: 0.1 },
    { name: "bagSerial", weight: 0.1 },
    { name: "locationName", weight: 0.1 },
    { name: "tags", weight: 0.1 },
  ],
  threshold: 0.4, // 0 = exact match, 1 = match anything
  ignoreLocation: true,
  includeScore: true,
};

const locations = ref<Location[]>([]);
const boxes = ref<Box[]>([]);
const manufacturersList = ref<MasterDataItem[]>([]);
const typesList = ref<MasterDataItem[]>([]);
const tagsList = ref<MasterDataItem[]>([]);
const manufacturers = ref<string[]>([]);
const types = ref<string[]>([]);
const tagSuggestions = ref<string[]>([]);
const appPaths = ref<AppPaths | null>(null);

// Confirm Modal
const confirmModal = ref({
  visible: false,
  message: "",
  danger: false,
  onConfirm: () => {},
});

// Master Data Panels
const masterDataPanel = ref<"manufacturers" | "types" | "tags" | null>(null);

// Settings & Menu
const menuOpen = ref(false);
const settingsOpen = ref(false);
const stats = ref<Record<string, number> | null>(null);

// i18n
const { t, locale } = useI18n();

const toast = ref({
  show: false,
  message: "",
  type: "success" as "success" | "error",
});

const form = reactive({
  id: null as number | null,
  name: "",
  manufacturer: "",
  type: "",
  locationId: null as number | null,
  boxId: null as number | null,
  bagSerial: "",
  tags: [] as string[],
  products: [] as ProductItem[],
  photoPath: "",
  photoSource: "",
});

const cropVisible = ref(false);
const cropTask = reactive({
  src: "",
  ext: "png",
  origin: "file" as "file" | "url" | "scan",
  skipHandler: null as (() => Promise<void> | void) | null,
});

// Image data (base64)
const currentImageData = ref("");

// Computed
const isEditing = computed(() => form.id !== null);

const currentBagInfo = computed(() => {
  if (!form.boxId) return null;
  const box = boxes.value.find((b) => b.id === form.boxId);
  const location = form.locationId
    ? locations.value.find((l) => l.id === form.locationId)
    : null;
  return {
    id: 0,
    serialNo: form.bagSerial,
    boxId: form.boxId,
    boxCode: box?.code || "",
    boxName: box?.name || "",
    locationId: form.locationId || 0,
    locationName: location?.friendlyName || "",
    locationNote: location?.note || "",
    locationRoom: location?.room || "",
    locationShelf: location?.shelf || "",
    locationCompartment: location?.compartment || "",
  };
});

// Watch
let searchTimer: number | undefined;
watch([searchQuery, sortBy], () => {
  if (searchTimer) clearTimeout(searchTimer);
  searchTimer = window.setTimeout(runSearch, 200);
});

// Watch photoPath changes to load image as base64
watch(
  () => form.photoPath,
  async (newPath) => {
    if (newPath) {
      try {
        currentImageData.value = await GetImageAsBase64(newPath);
      } catch {
        currentImageData.value = "";
      }
    } else {
      currentImageData.value = "";
    }
  },
  { immediate: true }
);

// Functions
function showToast(message: string, type: "success" | "error" = "success") {
  toast.value = { show: true, message, type };
  setTimeout(() => {
    toast.value.show = false;
  }, 3000);
}

function resetForm() {
  form.id = null;
  form.name = "";
  form.manufacturer = "";
  form.type = "";
  form.locationId = null;
  form.boxId = null;
  form.bagSerial = "";
  form.tags = [];
  form.products = [];
  form.photoPath = "";
  form.photoSource = "";
}

function startNewSet() {
  resetForm();
  view.value = "detail";
}

function backToList() {
  view.value = "list";
  resetForm();
}

async function loadInitial() {
  try {
    appPaths.value = await GetAppPaths();
    await Promise.all([
      refreshLocations(),
      refreshBoxes(),
      refreshManufacturers(),
      refreshTypes(),
      refreshTags(),
    ]);
    await loadAllSets();
    await runSearch();
  } catch (err: any) {
    showToast(err?.message ?? String(err), "error");
  }
}

async function refreshLocations() {
  locations.value = (await ListLocations()) || [];
}

async function refreshBoxes() {
  boxes.value = (await ListBoxes(0)) || [];
}

async function refreshManufacturers() {
  const list = await ListManufacturers();
  manufacturersList.value = (list || []).map((m: any) => ({
    id: m.id,
    name: m.name,
  }));
  manufacturers.value = manufacturersList.value.map((m) => m.name);
}

async function refreshTypes() {
  try {
    const list = await ListTypes();
    typesList.value = (list || []).map((t: any) => ({
      id: t.id,
      name: t.name,
    }));
    types.value = typesList.value.map((t) => t.name);
  } catch {
    typesList.value = [];
    types.value = [];
  }
}

async function refreshTags() {
  tagSuggestions.value = (await ListTags()) || [];
  try {
    const list = await ListTagsFull();
    tagsList.value = (list || []).map((t: any) => ({ id: t.id, name: t.name }));
  } catch {
    tagsList.value = [];
  }
}

// Check if query has special @ prefix
function hasSpecialFilter(query: string): boolean {
  const prefixes = [
    "@box",
    "@produkt",
    "@product",
    "@hersteller",
    "@tag",
    "@ort",
    "@standort",
  ];
  const lowerQuery = query.toLowerCase().trim();
  return prefixes.some((p) => lowerQuery.startsWith(p + " "));
}

// Sort results based on sortBy
function sortResults(results: SearchResult[]): SearchResult[] {
  const sorted = [...results];
  switch (sortBy.value) {
    case "box":
      sorted.sort((a, b) => {
        const cmp = a.boxCode.localeCompare(b.boxCode);
        return cmp !== 0 ? cmp : a.bagSerial.localeCompare(b.bagSerial);
      });
      break;
    case "location":
      sorted.sort((a, b) => a.locationName.localeCompare(b.locationName));
      break;
    case "added":
      sorted.sort((a, b) => b.setId - a.setId); // Higher ID = newer
      break;
    default:
      sorted.sort((a, b) => a.setName.localeCompare(b.setName));
  }
  return sorted;
}

// Load all sets from backend
async function loadAllSets() {
  try {
    allSets.value = (await SearchSets("", sortBy.value)) || [];
    // Initialize Fuse.js with all sets
    fuseInstance.value = new Fuse(allSets.value, fuseOptions);
  } catch (err: any) {
    showToast(err?.message ?? String(err), "error");
  }
}

async function runSearch() {
  searchLoading.value = true;
  try {
    const query = searchQuery.value.trim();

    // If query has special @ prefix, use backend filtering
    if (hasSpecialFilter(query)) {
      searchResults.value = (await SearchSets(query, sortBy.value)) || [];
    } else if (query === "") {
      // No query - show all sets sorted
      searchResults.value = sortResults(allSets.value);
    } else {
      // Use Fuse.js for fuzzy search
      if (fuseInstance.value) {
        const results = fuseInstance.value.search(query);
        searchResults.value = sortResults(results.map((r) => r.item));
      } else {
        searchResults.value = [];
      }
    }
  } catch (err: any) {
    showToast(err?.message ?? String(err), "error");
  } finally {
    searchLoading.value = false;
  }
}

async function selectResult(result: SearchResult) {
  try {
    const details = await GetSet(result.setId);
    form.id = details.id;
    form.name = details.name;
    form.manufacturer = details.manufacturerName || "";
    form.type = details.typeName || "";
    form.locationId = details.bag?.locationId ?? null;
    form.boxId = details.bag?.boxId ?? null;
    form.bagSerial = details.bag?.serialNo ?? "";
    form.tags = details.tags || [];
    form.products = details.products || [];
    form.photoPath = details.photoPath || "";
    form.photoSource = details.photoSource || "";
    view.value = "overview";
  } catch (err: any) {
    showToast(err?.message ?? String(err), "error");
  }
}

function openEditFromOverview() {
  view.value = "detail";
}

function backFromOverview() {
  view.value = "list";
  resetForm();
}

function requestDeleteFromOverview() {
  if (!form.id) return;
  confirmModal.value = {
    visible: true,
    message: `Set "${form.name}" wirklich löschen?`,
    danger: true,
    onConfirm: async () => {
      confirmModal.value.visible = false;
      await executeDeleteSet();
    },
  };
}

async function saveSet() {
  if (!form.name.trim()) {
    showToast("Bitte einen Namen eingeben", "error");
    return;
  }
  if (!form.locationId || !form.boxId) {
    showToast("Bitte Standort und Box wählen", "error");
    return;
  }
  if (!form.bagSerial.trim()) {
    showToast("Bitte Beutel-Nr. eingeben", "error");
    return;
  }

  try {
    if (form.id) {
      await UpdateSet(
        form.id,
        form.name.trim(),
        form.manufacturer.trim(),
        form.type.trim(),
        form.boxId,
        form.bagSerial.trim()
      );
    } else {
      const createdId = await CreateBagWithSet(
        form.boxId,
        form.bagSerial.trim(),
        form.name.trim(),
        form.manufacturer.trim(),
        form.type.trim()
      );
      form.id = createdId;
    }

    if (form.id) {
      await SetTags(form.id, form.tags);
      for (const prod of form.products) {
        if (prod.id <= 0) {
          const newId = await AddProduct(form.id, prod.name, prod.kind);
          prod.id = newId;
          prod.setId = form.id;
        } else {
          await UpdateProduct(prod.id, prod.name, prod.kind);
        }
      }
    }

    await refreshTags();
    await refreshTypes();
    await loadAllSets(); // Reload all sets to update Fuse index
    await runSearch();
    showToast("Gespeichert!");
  } catch (err: any) {
    showToast(err?.message ?? String(err), "error");
  }
}

function requestDeleteSet() {
  if (!form.id) return;
  confirmModal.value = {
    visible: true,
    message: `Set "${form.name}" wirklich löschen?`,
    danger: true,
    onConfirm: async () => {
      confirmModal.value.visible = false;
      await executeDeleteSet();
    },
  };
}

async function executeDeleteSet() {
  if (!form.id) return;
  try {
    await DeleteSet(form.id);
    showToast("Gelöscht");
    backToList();
    await loadAllSets(); // Reload all sets to update Fuse index
    await runSearch();
  } catch (err: any) {
    showToast(err?.message ?? String(err), "error");
  }
}

// Master Data handlers
async function handleCreateManufacturer(name: string) {
  try {
    await CreateManufacturer(name);
    await refreshManufacturers();
    showToast("Hersteller erstellt");
  } catch (err: any) {
    showToast(err?.message ?? String(err), "error");
  }
}

async function handleUpdateManufacturer(id: number, name: string) {
  try {
    await UpdateManufacturer(id, name);
    await refreshManufacturers();
    showToast("Hersteller aktualisiert");
  } catch (err: any) {
    showToast(err?.message ?? String(err), "error");
  }
}

async function handleDeleteManufacturer(id: number) {
  try {
    await DeleteManufacturer(id);
    await refreshManufacturers();
    showToast("Hersteller gelöscht");
  } catch (err: any) {
    showToast(err?.message ?? String(err), "error");
  }
}

async function handleCreateType(name: string) {
  try {
    await CreateType(name);
    await refreshTypes();
    showToast("Typ erstellt");
  } catch (err: any) {
    showToast(err?.message ?? String(err), "error");
  }
}

async function handleUpdateType(id: number, name: string) {
  try {
    await UpdateType(id, name);
    await refreshTypes();
    showToast("Typ aktualisiert");
  } catch (err: any) {
    showToast(err?.message ?? String(err), "error");
  }
}

async function handleDeleteType(id: number) {
  try {
    await DeleteType(id);
    await refreshTypes();
    showToast("Typ gelöscht");
  } catch (err: any) {
    showToast(err?.message ?? String(err), "error");
  }
}

async function handleCreateTag(name: string) {
  try {
    await CreateTag(name);
    await refreshTags();
    showToast("Tag erstellt");
  } catch (err: any) {
    showToast(err?.message ?? String(err), "error");
  }
}

async function handleUpdateTag(id: number, name: string) {
  try {
    await UpdateTag(id, name);
    await refreshTags();
    showToast("Tag aktualisiert");
  } catch (err: any) {
    showToast(err?.message ?? String(err), "error");
  }
}

async function handleDeleteTag(id: number) {
  try {
    await DeleteTag(id);
    await refreshTags();
    showToast("Tag gelöscht");
  } catch (err: any) {
    showToast(err?.message ?? String(err), "error");
  }
}

// Location handlers
async function handleSaveLocation(data: {
  id: number | null;
  name: string;
  room: string;
  shelf: string;
  compartment: string;
  note: string;
}) {
  try {
    if (data.id) {
      await UpdateLocation(
        data.id,
        data.name,
        data.room,
        data.shelf,
        data.compartment,
        data.note
      );
    } else {
      const newId = await CreateLocation(
        data.name,
        data.room,
        data.shelf,
        data.compartment,
        data.note
      );
      form.locationId = newId;
    }
    await refreshLocations();
    showToast("Standort gespeichert");
  } catch (err: any) {
    showToast(err?.message ?? String(err), "error");
  }
}

async function handleDeleteLocation(id: number) {
  try {
    await DeleteLocation(id);
    form.locationId = null;
    form.boxId = null;
    await refreshLocations();
    await refreshBoxes();
    await runSearch();
    showToast("Standort gelöscht");
  } catch (err: any) {
    showToast(err?.message ?? String(err), "error");
  }
}

async function handleSaveBox(data: {
  id: number | null;
  locationId: number;
  code: string;
  name: string;
}) {
  try {
    if (data.id) {
      await UpdateBox(data.id, data.locationId, data.code, data.name);
    } else {
      const newId = await CreateBox(data.locationId, data.code, data.name);
      form.boxId = newId;
      // Auto-set next bag serial for new box
      await updateBagSerialForBox(newId);
    }
    await refreshBoxes();
    showToast("Box gespeichert");
  } catch (err: any) {
    showToast(err?.message ?? String(err), "error");
  }
}

async function handleDeleteBox(id: number) {
  try {
    await DeleteBox(id);
    form.boxId = null;
    form.bagSerial = "";
    await refreshBoxes();
    await runSearch();
    showToast("Box gelöscht");
  } catch (err: any) {
    showToast(err?.message ?? String(err), "error");
  }
}

// Auto-increment bag serial when box changes
async function updateBagSerialForBox(boxId: number | null) {
  if (!boxId || form.id) return; // Only for new sets
  try {
    const nextSerial = await GetNextBagSerial(boxId);
    form.bagSerial = nextSerial;
  } catch (err) {
    console.error("Failed to get next bag serial:", err);
  }
}

async function handleBoxIdChange(boxId: number | null) {
  form.boxId = boxId;
  await updateBagSerialForBox(boxId);
}

// Product handlers
async function handleAddProduct(product: { name: string; kind: string }) {
  if (form.id) {
    try {
      const newId = await AddProduct(form.id, product.name, product.kind);
      form.products.push({
        id: newId,
        setId: form.id,
        name: product.name,
        kind: product.kind,
      });
    } catch (err: any) {
      showToast(err?.message ?? String(err), "error");
    }
  } else {
    const tempId = -Date.now();
    form.products.push({
      id: tempId,
      setId: -1,
      name: product.name,
      kind: product.kind,
    });
  }
}

async function handleUpdateProduct(product: ProductItem) {
  if (product.id > 0) {
    try {
      await UpdateProduct(product.id, product.name, product.kind);
    } catch (err: any) {
      showToast(err?.message ?? String(err), "error");
    }
  }
}

async function handleDeleteProduct(id: number) {
  if (id > 0) {
    try {
      await DeleteProduct(id);
    } catch (err: any) {
      showToast(err?.message ?? String(err), "error");
      return;
    }
  }
  form.products = form.products.filter((p) => p.id !== id);
}

// Image handlers
function toFileUrl(path: string) {
  if (!path) return "";
  if (path.startsWith("file://")) return path;
  const base = appPaths.value?.baseDir ?? "";
  const isAbs = /^[a-zA-Z]:[\\/]|^\\/.test(path);
  const combined = isAbs ? path : base ? base + "/" + path : path;
  return "file:///" + combined.replace(/\\/g, "/");
}

function extFromPath(path: string) {
  const idx = path.lastIndexOf(".");
  return idx === -1 ? ".png" : path.slice(idx);
}

async function handleChooseFile() {
  if (!form.id) {
    const saved = await saveSetFirst();
    if (!saved) return;
  }

  try {
    const filePath = await ChooseImageFile();
    if (!filePath) return;

    // Read file as base64 for crop preview
    const base64Data = await ReadFileAsBase64(filePath);

    cropTask.src = base64Data;
    cropTask.ext = extFromPath(filePath);
    cropTask.origin = "file";
    cropTask.skipHandler = async () => {
      if (!form.id) return;
      const rel = await AttachImageFromFile(form.id, filePath);
      form.photoPath = rel;
      showToast(t("saveSuccess"));
    };
    cropVisible.value = true;
  } catch (err: any) {
    showToast(err?.message ?? String(err), "error");
  }
}

async function handleScan() {
  if (!form.id) {
    const saved = await saveSetFirst();
    if (!saved) return;
  }

  try {
    showToast(locale.value === "de" ? "Scanne..." : "Scanning...");
    const result = await ScanImageToBase64();
    if (!result || !result.base64Data) return;

    cropTask.src = result.base64Data;
    cropTask.ext = ".png";
    cropTask.origin = "scan";
    cropTask.skipHandler = async () => {
      if (!form.id) return;
      // Attach the scanned image to the set
      await AttachScannedImage(form.id, result.relPath);
      form.photoPath = result.relPath;
      showToast(t("saveSuccess"));
    };
    cropVisible.value = true;
  } catch (err: any) {
    const msg = err?.message ?? String(err);
    if (msg.includes("no scanner") || msg.includes("No scanner")) {
      showToast(
        locale.value === "de" ? "Kein Scanner gefunden" : "No scanner found",
        "error"
      );
    } else {
      showToast(msg, "error");
    }
  }
}

async function handleFromUrl(inputUrl: string) {
  if (!form.id) {
    const saved = await saveSetFirst();
    if (!saved) return;
  }

  try {
    // First download the image
    const rel = await AttachImageFromURL(form.id!, inputUrl);
    form.photoPath = rel;

    // Read as base64 for crop dialog
    const fullPath = (appPaths.value?.baseDir || "") + "/" + rel;
    const base64Data = await ReadFileAsBase64(fullPath.replace(/\//g, "\\"));

    // Show crop dialog
    cropTask.src = base64Data;
    cropTask.ext = extFromPath(rel);
    cropTask.origin = "url";
    cropTask.skipHandler = null; // Image already saved
    cropVisible.value = true;

    showToast(locale.value === "de" ? "Bild geladen" : "Image loaded");
  } catch (err: any) {
    showToast(err?.message ?? String(err), "error");
  }
}

async function saveSetFirst(): Promise<boolean> {
  if (
    !form.name.trim() ||
    !form.locationId ||
    !form.boxId ||
    !form.bagSerial.trim()
  ) {
    showToast("Bitte zuerst alle Pflichtfelder ausfüllen", "error");
    return false;
  }
  await saveSet();
  return form.id !== null;
}

async function handleRemoveImage() {
  if (!form.id || !form.photoPath) return;

  try {
    await RemoveImage(form.id);
    form.photoPath = "";
    showToast(locale.value === "de" ? "Bild entfernt" : "Image removed");
  } catch (err: any) {
    showToast(err?.message ?? String(err), "error");
  }
}

async function handleCropConfirm(dataUrl: string) {
  if (!form.id) return;
  try {
    const rel = await SaveCroppedImage(
      form.id,
      dataUrl,
      cropTask.ext.replace(".", "")
    );
    form.photoPath = rel;
    showToast("Bild gespeichert");
  } catch (err: any) {
    showToast(err?.message ?? String(err), "error");
  } finally {
    cropVisible.value = false;
  }
}

async function handleCropSkip() {
  cropVisible.value = false;
  if (cropTask.skipHandler) {
    await cropTask.skipHandler();
  }
}

// Settings & Menu handlers
function toggleMenu() {
  menuOpen.value = !menuOpen.value;
}

function closeMenu() {
  menuOpen.value = false;
}

async function openSettings() {
  closeMenu();
  try {
    stats.value = await GetStats();
  } catch {
    stats.value = null;
  }
  settingsOpen.value = true;
}

async function handleExport() {
  try {
    const path = await ExportData();
    if (path) {
      showToast(t("exportSuccess"));
    }
  } catch (err: any) {
    showToast(err?.message ?? String(err), "error");
  }
}

async function handleImport() {
  confirmModal.value = {
    visible: true,
    message: t("confirmImportMessage"),
    danger: true,
    onConfirm: async () => {
      confirmModal.value.visible = false;
      try {
        const path = await ImportData();
        if (path) {
          showToast(t("importSuccess"));
          await loadInitial();
          await loadAllSets();
          await runSearch();
        }
      } catch (err: any) {
        showToast(err?.message ?? String(err), "error");
      }
    },
  };
}

// Keyboard shortcuts
function handleKeydown(e: KeyboardEvent) {
  if (e.ctrlKey && e.key.toLowerCase() === "f") {
    e.preventDefault();
    if (view.value === "list") {
      searchBarRef.value?.focus();
    }
  }
  if (e.ctrlKey && e.key.toLowerCase() === "n") {
    e.preventDefault();
    startNewSet();
  }
  if (e.key === "Escape") {
    if (view.value === "detail") {
      backToList();
    } else if (view.value === "overview") {
      backFromOverview();
    }
  }
}

onMounted(() => {
  loadInitial();
  document.addEventListener("keydown", handleKeydown);
  nextTick(() => searchBarRef.value?.focus());
});

onBeforeUnmount(() => {
  document.removeEventListener("keydown", handleKeydown);
});
</script>

<template>
  <div class="app" @click="closeMenu">
    <!-- Header -->
    <header class="header">
      <div class="header-brand">
        <img :src="samlaIcon" alt="" class="logo" />
        <span class="app-name">Samla</span>
      </div>
      <div class="header-actions">
        <button
          class="header-btn"
          @click="masterDataPanel = 'manufacturers'"
          title="Manufacturers"
        >
          <i class="mdi mdi-factory"></i>
        </button>
        <button
          class="header-btn"
          @click="masterDataPanel = 'types'"
          title="Types"
        >
          <i class="mdi mdi-shape-outline"></i>
        </button>
        <button
          class="header-btn"
          @click="masterDataPanel = 'tags'"
          title="Tags"
        >
          <i class="mdi mdi-tag-outline"></i>
        </button>
        <div class="header-divider"></div>

        <!-- Hamburger Menu -->
        <div class="menu-container">
          <button
            class="header-btn menu-btn"
            @click.stop="toggleMenu"
            title="Menu"
          >
            <i class="mdi mdi-menu"></i>
          </button>

          <Transition name="menu">
            <div v-if="menuOpen" class="dropdown-menu" @click.stop>
              <button class="menu-item" @click="openSettings">
                <i class="mdi mdi-cog"></i>
                <span>{{ t("settings") }}</span>
              </button>
              <button
                class="menu-item"
                @click="
                  OpenAppFolder();
                  closeMenu();
                "
              >
                <i class="mdi mdi-folder-open"></i>
                <span>{{ t("openDataFolder") }}</span>
              </button>
              <div class="menu-divider"></div>
              <button
                class="menu-item"
                @click="
                  handleExport();
                  closeMenu();
                "
              >
                <i class="mdi mdi-export"></i>
                <span>{{ t("exportData") }}</span>
              </button>
              <button
                class="menu-item"
                @click="
                  handleImport();
                  closeMenu();
                "
              >
                <i class="mdi mdi-import"></i>
                <span>{{ t("importData") }}</span>
              </button>
            </div>
          </Transition>
        </div>
      </div>
    </header>

    <!-- Toast -->
    <Transition name="toast">
      <div v-if="toast.show" :class="['toast', toast.type]">
        <i
          :class="[
            'mdi',
            toast.type === 'success' ? 'mdi-check-circle' : 'mdi-alert-circle',
          ]"
        ></i>
        {{ toast.message }}
      </div>
    </Transition>

    <!-- List View -->
    <div v-if="view === 'list'" class="view-list">
      <SearchBar
        ref="searchBarRef"
        v-model="searchQuery"
        v-model:sort-by="sortBy"
        :loading="searchLoading"
        :result-count="searchResults.length"
        @new-set="startNewSet"
      />

      <div class="results">
        <SetCard
          v-for="item in searchResults"
          :key="item.setId"
          v-bind="item"
          @click="selectResult(item)"
        />

        <EmptyState v-if="!searchResults.length && !searchLoading" />
      </div>
    </div>

    <!-- Overview View -->
    <SetOverview
      v-else-if="view === 'overview'"
      :id="form.id!"
      :name="form.name"
      :manufacturer-name="form.manufacturer"
      :type-name="form.type"
      :photo-path="form.photoPath"
      :bag="currentBagInfo"
      :tags="form.tags"
      :products="form.products"
      @edit="openEditFromOverview"
      @back="backFromOverview"
      @delete="requestDeleteFromOverview"
    />

    <!-- Detail View -->
    <div v-else-if="view === 'detail'" class="view-detail">
      <div class="detail-header">
        <button class="btn-back" @click="backToList">
          <i class="mdi mdi-arrow-left"></i>
          {{ t("back") }}
        </button>
        <h1>
          {{
            isEditing
              ? locale === "de"
                ? "Set bearbeiten"
                : "Edit Set"
              : t("newSet")
          }}
        </h1>
        <div class="detail-actions">
          <button v-if="isEditing" class="btn-delete" @click="requestDeleteSet">
            <i class="mdi mdi-delete-outline"></i>
          </button>
          <button class="btn-save" @click="saveSet">
            <i class="mdi mdi-check"></i>
            {{ t("save") }}
          </button>
        </div>
      </div>

      <div class="detail-content">
        <div class="detail-main">
          <!-- Name & Info -->
          <section class="section">
            <h2>
              <i class="mdi mdi-information-outline"></i>
              {{ locale === "de" ? "Grunddaten" : "Basic Info" }}
            </h2>
            <div class="field">
              <label>{{ t("setName") }} *</label>
              <input
                v-model="form.name"
                type="text"
                :placeholder="t('setNamePlaceholder')"
                class="input large"
              />
            </div>
            <div class="field-row">
              <div class="field">
                <label>{{ t("manufacturer") }}</label>
                <div class="select-with-action">
                  <select v-model="form.manufacturer" class="input">
                    <option value="">
                      — {{ locale === "de" ? "Auswählen" : "Select" }} —
                    </option>
                    <option v-for="m in manufacturers" :key="m" :value="m">
                      {{ m }}
                    </option>
                  </select>
                  <button
                    class="btn-inline"
                    @click="masterDataPanel = 'manufacturers'"
                    :title="
                      locale === 'de'
                        ? 'Hersteller verwalten'
                        : 'Manage manufacturers'
                    "
                    type="button"
                  >
                    <i class="mdi mdi-cog-outline"></i>
                  </button>
                </div>
              </div>
              <div class="field">
                <label>{{ t("type") }}</label>
                <div class="select-with-action">
                  <select v-model="form.type" class="input">
                    <option value="">
                      — {{ locale === "de" ? "Auswählen" : "Select" }} —
                    </option>
                    <option v-for="tp in types" :key="tp" :value="tp">
                      {{ tp }}
                    </option>
                  </select>
                  <button
                    class="btn-inline"
                    @click="masterDataPanel = 'types'"
                    :title="
                      locale === 'de' ? 'Typen verwalten' : 'Manage types'
                    "
                    type="button"
                  >
                    <i class="mdi mdi-cog-outline"></i>
                  </button>
                </div>
              </div>
            </div>
          </section>

          <!-- Location -->
          <section class="section">
            <h2>
              <i class="mdi mdi-map-marker-outline"></i> {{ t("location") }}
            </h2>
            <LocationPicker
              :locations="locations"
              :boxes="boxes"
              :location-id="form.locationId"
              :box-id="form.boxId"
              :bag-serial="form.bagSerial"
              @update:location-id="form.locationId = $event"
              @update:box-id="handleBoxIdChange"
              @update:bag-serial="form.bagSerial = $event"
              @save-location="handleSaveLocation"
              @delete-location="handleDeleteLocation"
              @save-box="handleSaveBox"
              @delete-box="handleDeleteBox"
            />
          </section>

          <!-- Tags -->
          <section class="section">
            <h2><i class="mdi mdi-tag-outline"></i> {{ t("tags") }}</h2>
            <TagInput
              :tags="form.tags"
              :suggestions="tagSuggestions"
              @update:tags="form.tags = $event"
            />
          </section>

          <!-- Products -->
          <section class="section">
            <h2><i class="mdi mdi-shape-outline"></i> {{ t("products") }}</h2>
            <ProductList
              :products="form.products"
              :set-id="form.id"
              @add="handleAddProduct"
              @update="handleUpdateProduct"
              @delete="handleDeleteProduct"
            />
          </section>
        </div>

        <div class="detail-side">
          <!-- Image -->
          <section class="section">
            <h2><i class="mdi mdi-image-outline"></i> {{ t("photo") }}</h2>
            <ImageUpload
              :photo-path="form.photoPath"
              :image-data="currentImageData"
              :set-id="form.id"
              @choose-file="handleChooseFile"
              @from-url="handleFromUrl"
              @scan="handleScan"
              @remove="handleRemoveImage"
            />
          </section>
        </div>
      </div>
    </div>

    <!-- Crop Modal -->
    <CropModal
      :visible="cropVisible"
      :image-src="cropTask.src"
      @close="cropVisible = false"
      @confirm="handleCropConfirm"
      @skip="handleCropSkip"
    />

    <!-- Confirm Modal -->
    <ConfirmModal
      :visible="confirmModal.visible"
      :message="confirmModal.message"
      :danger="confirmModal.danger"
      @confirm="confirmModal.onConfirm"
      @cancel="confirmModal.visible = false"
    />

    <!-- Master Data Panels -->
    <MasterDataPanel
      :visible="masterDataPanel === 'manufacturers'"
      :title="t('manufacturers')"
      icon="mdi-factory"
      :items="manufacturersList"
      @close="masterDataPanel = null"
      @create="handleCreateManufacturer"
      @update="handleUpdateManufacturer"
      @delete="handleDeleteManufacturer"
    />

    <MasterDataPanel
      :visible="masterDataPanel === 'types'"
      :title="t('types')"
      icon="mdi-shape-outline"
      :items="typesList"
      @close="masterDataPanel = null"
      @create="handleCreateType"
      @update="handleUpdateType"
      @delete="handleDeleteType"
    />

    <MasterDataPanel
      :visible="masterDataPanel === 'tags'"
      :title="t('tags')"
      icon="mdi-tag-outline"
      :items="tagsList"
      @close="masterDataPanel = null"
      @create="handleCreateTag"
      @update="handleUpdateTag"
      @delete="handleDeleteTag"
    />

    <!-- Settings Panel -->
    <SettingsPanel
      :visible="settingsOpen"
      :paths="appPaths"
      :stats="stats"
      @close="settingsOpen = false"
      @open-folder="OpenAppFolder"
      @export="handleExport"
      @import="handleImport"
    />
  </div>
</template>

<style>
@import "@mdi/font/css/materialdesignicons.css";

* {
  box-sizing: border-box;
}

body {
  margin: 0;
  font-family: "Inter", -apple-system, BlinkMacSystemFont, "Segoe UI",
    sans-serif;
  background: #f5f5f5;
  color: #111;
  -webkit-font-smoothing: antialiased;
}

.app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

/* Header */
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 20px;
  background: #111;
  color: white;
}

.header-brand {
  display: flex;
  align-items: center;
  gap: 10px;
}

.logo {
  width: 32px;
  height: 32px;
  background: white;
  border-radius: 8px;
  padding: 4px;
}

.app-name {
  font-size: 18px;
  font-weight: 600;
  letter-spacing: 0.5px;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 6px;
}

.header-divider {
  width: 1px;
  height: 24px;
  background: rgba(255, 255, 255, 0.2);
  margin: 0 6px;
}

.header-btn {
  width: 36px;
  height: 36px;
  border: none;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.1);
  color: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.header-btn:hover {
  background: rgba(255, 255, 255, 0.2);
}

.header-btn i {
  font-size: 20px;
}

/* Menu Container */
.menu-container {
  position: relative;
}

.dropdown-menu {
  position: absolute;
  top: 100%;
  right: 0;
  margin-top: 8px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
  min-width: 200px;
  padding: 8px;
  z-index: 1000;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
  padding: 12px 14px;
  border: none;
  background: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  color: #333;
  text-align: left;
}

.menu-item:hover {
  background: #f5f5f5;
}

.menu-item i {
  font-size: 20px;
  color: #666;
}

.menu-divider {
  height: 1px;
  background: #eee;
  margin: 6px 0;
}

.menu-enter-active,
.menu-leave-active {
  transition: all 0.15s ease;
}

.menu-enter-from,
.menu-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}

/* Toast */
.toast {
  position: fixed;
  top: 70px;
  left: 50%;
  transform: translateX(-50%);
  padding: 12px 20px;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 8px;
  z-index: 1000;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
}

.toast.success {
  background: #111;
  color: white;
}

.toast.error {
  background: #dc2626;
  color: white;
}

.toast i {
  font-size: 18px;
}

.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translate(-50%, -20px);
}

/* List View */
.view-list {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: white;
}

.results {
  flex: 1;
  padding: 16px 20px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  overflow-y: auto;
}

/* Detail View */
.view-detail {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: #f5f5f5;
}

.detail-header {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px 20px;
  background: white;
  border-bottom: 1px solid #e5e5e5;
}

.detail-header h1 {
  flex: 1;
  font-size: 18px;
  font-weight: 600;
  margin: 0;
}

.btn-back {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  font-size: 14px;
  font-weight: 500;
  border: 1px solid #ddd;
  border-radius: 8px;
  background: white;
  cursor: pointer;
}

.btn-back:hover {
  border-color: #111;
  background: #f8f8f8;
}

.btn-back i {
  font-size: 18px;
}

.detail-actions {
  display: flex;
  gap: 8px;
}

.btn-delete {
  width: 40px;
  height: 40px;
  border: 1px solid #fcc;
  border-radius: 8px;
  background: #fff;
  color: #c00;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-delete:hover {
  background: #fee;
}

.btn-delete i {
  font-size: 20px;
}

.btn-save {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 20px;
  font-size: 14px;
  font-weight: 500;
  border: none;
  border-radius: 8px;
  background: #111;
  color: white;
  cursor: pointer;
}

.btn-save:hover {
  background: #333;
}

.btn-save i {
  font-size: 18px;
}

.detail-content {
  flex: 1;
  display: grid;
  grid-template-columns: 1fr 320px;
  gap: 20px;
  padding: 20px;
  overflow-y: auto;
}

.detail-main {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.detail-side {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.section {
  background: white;
  border-radius: 12px;
  padding: 20px;
  border: 1px solid #eee;
}

.section h2 {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 15px;
  font-weight: 600;
  margin: 0 0 16px 0;
  color: #333;
}

.section h2 i {
  font-size: 20px;
  color: #666;
}

.field {
  margin-bottom: 14px;
}

.field:last-child {
  margin-bottom: 0;
}

.field label {
  display: block;
  font-size: 13px;
  font-weight: 500;
  color: #666;
  margin-bottom: 6px;
}

.field-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
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

.input.large {
  padding: 14px 16px;
  font-size: 16px;
}

select.input {
  cursor: pointer;
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 24 24'%3E%3Cpath fill='%23666' d='M7 10l5 5 5-5z'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 12px center;
  padding-right: 36px;
}

.select-with-action {
  display: flex;
  gap: 8px;
}

.select-with-action select {
  flex: 1;
}

.btn-inline {
  width: 44px;
  height: 44px;
  border: 1px solid #ddd;
  border-radius: 8px;
  background: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #666;
  flex-shrink: 0;
}

.btn-inline:hover {
  border-color: #111;
  background: #f8f8f8;
}

.btn-inline i {
  font-size: 18px;
}

/* Responsive */
@media (max-width: 900px) {
  .detail-content {
    grid-template-columns: 1fr;
  }

  .detail-side {
    order: -1;
  }
}

@media (max-width: 600px) {
  .field-row {
    grid-template-columns: 1fr;
  }

  .detail-header h1 {
    display: none;
  }
}
</style>
