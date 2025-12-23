import { ref, computed } from "vue";

export type Locale = "de" | "en";

const currentLocale = ref<Locale>(
  (localStorage.getItem("samla-locale") as Locale) || "de"
);

export const translations = {
  de: {
    // App
    appName: "Samla",

    // Navigation & Actions
    newSet: "Neues Set",
    edit: "Bearbeiten",
    delete: "Löschen",
    save: "Speichern",
    cancel: "Abbrechen",
    close: "Schließen",
    back: "Zurück",
    add: "Hinzufügen",
    remove: "Entfernen",
    confirm: "Bestätigen",
    search: "Suchen",

    // Menu
    settings: "Einstellungen",
    exportData: "Daten exportieren",
    importData: "Daten importieren",

    // Search
    searchPlaceholder: "Suchen... (@Box, @Produkt, @Hersteller, @Tag, @Ort)",
    sortBy: "Sortieren",
    sortName: "Name",
    sortBox: "Box",
    sortLocation: "Ort",
    sortAdded: "Hinzugefügt",

    // Set Form
    setName: "Set-Name",
    setNamePlaceholder: "z.B. Briefmarken Deutschland 2020",
    boxNumber: "Karton-Nr.",
    bagNumber: "Beutel-Nr.",
    auto: "Auto",

    // Location
    location: "Ort",
    room: "Raum",
    roomPlaceholder: "z.B. Wohnzimmer",
    shelf: "Regal",
    shelfPlaceholder: "z.B. Regal 3",
    compartment: "Fach",
    compartmentPlaceholder: "z.B. Fach A",

    // Products
    products: "Produkte",
    product: "Produkt",
    productName: "Produktname",
    productNamePlaceholder: "z.B. Briefmarke",
    manufacturer: "Hersteller",
    manufacturerPlaceholder: "Hersteller wählen",
    type: "Typ",
    typePlaceholder: "Typ wählen",
    quantity: "Menge",
    addProduct: "Produkt hinzufügen",
    noProducts: "Keine Produkte",

    // Tags
    tags: "Tags",
    tagsPlaceholder: "Tag eingeben und Enter drücken",

    // Images
    photo: "Foto",
    addPhoto: "Foto hinzufügen",
    changePhoto: "Foto ändern",
    removePhoto: "Foto entfernen",
    cropImage: "Bild zuschneiden",

    // Master Data
    masterData: "Stammdaten",
    manufacturers: "Hersteller",
    types: "Typen",
    newManufacturer: "Neuer Hersteller",
    newType: "Neuer Typ",
    newTag: "Neuer Tag",

    // Empty States
    noSets: "Keine Sets vorhanden",
    noSetsDescription: "Erstellen Sie Ihr erstes Set, um zu beginnen.",
    noResults: "Keine Ergebnisse",
    noResultsDescription: "Versuchen Sie andere Suchbegriffe.",

    // Confirmations
    confirmDelete: "Löschen bestätigen",
    confirmDeleteSet: "Möchten Sie dieses Set wirklich löschen?",
    confirmDeleteProduct: "Möchten Sie dieses Produkt wirklich löschen?",
    confirmImport: "Import bestätigen",
    confirmImportMessage:
      "Beim Import werden alle bestehenden Daten überschrieben. Fortfahren?",

    // Settings Panel
    settingsTitle: "Einstellungen",
    language: "Sprache",
    german: "Deutsch",
    english: "English",
    dataManagement: "Datenverwaltung",
    openDataFolder: "Datenordner öffnen",
    storagePaths: "Speicherpfade",
    baseDirectory: "Basisverzeichnis",
    database: "Datenbank",
    images: "Bilder",
    statistics: "Statistiken",
    statsSets: "Sets",
    statsProducts: "Produkte",
    statsManufacturers: "Hersteller",
    statsTypes: "Typen",
    statsTags: "Tags",
    statsImages: "Bilder",
    about: "Über",
    version: "Version",

    // Messages
    exportSuccess: "Export erfolgreich",
    exportError: "Export fehlgeschlagen",
    importSuccess: "Import erfolgreich",
    importError: "Import fehlgeschlagen",
    saveSuccess: "Gespeichert",
    saveError: "Speichern fehlgeschlagen",
    deleteSuccess: "Gelöscht",
    deleteError: "Löschen fehlgeschlagen",
  },

  en: {
    // App
    appName: "Samla",

    // Navigation & Actions
    newSet: "New Set",
    edit: "Edit",
    delete: "Delete",
    save: "Save",
    cancel: "Cancel",
    close: "Close",
    back: "Back",
    add: "Add",
    remove: "Remove",
    confirm: "Confirm",
    search: "Search",

    // Menu
    settings: "Settings",
    exportData: "Export Data",
    importData: "Import Data",

    // Search
    searchPlaceholder:
      "Search... (@Box, @Product, @Manufacturer, @Tag, @Location)",
    sortBy: "Sort by",
    sortName: "Name",
    sortBox: "Box",
    sortLocation: "Location",
    sortAdded: "Added",

    // Set Form
    setName: "Set Name",
    setNamePlaceholder: "e.g. Stamp Collection 2020",
    boxNumber: "Box No.",
    bagNumber: "Bag No.",
    auto: "Auto",

    // Location
    location: "Location",
    room: "Room",
    roomPlaceholder: "e.g. Living Room",
    shelf: "Shelf",
    shelfPlaceholder: "e.g. Shelf 3",
    compartment: "Compartment",
    compartmentPlaceholder: "e.g. Compartment A",

    // Products
    products: "Products",
    product: "Product",
    productName: "Product Name",
    productNamePlaceholder: "e.g. Stamp",
    manufacturer: "Manufacturer",
    manufacturerPlaceholder: "Select manufacturer",
    type: "Type",
    typePlaceholder: "Select type",
    quantity: "Quantity",
    addProduct: "Add Product",
    noProducts: "No products",

    // Tags
    tags: "Tags",
    tagsPlaceholder: "Enter tag and press Enter",

    // Images
    photo: "Photo",
    addPhoto: "Add Photo",
    changePhoto: "Change Photo",
    removePhoto: "Remove Photo",
    cropImage: "Crop Image",

    // Master Data
    masterData: "Master Data",
    manufacturers: "Manufacturers",
    types: "Types",
    newManufacturer: "New Manufacturer",
    newType: "New Type",
    newTag: "New Tag",

    // Empty States
    noSets: "No sets yet",
    noSetsDescription: "Create your first set to get started.",
    noResults: "No results",
    noResultsDescription: "Try different search terms.",

    // Confirmations
    confirmDelete: "Confirm Delete",
    confirmDeleteSet: "Are you sure you want to delete this set?",
    confirmDeleteProduct: "Are you sure you want to delete this product?",
    confirmImport: "Confirm Import",
    confirmImportMessage:
      "Importing will overwrite all existing data. Continue?",

    // Settings Panel
    settingsTitle: "Settings",
    language: "Language",
    german: "Deutsch",
    english: "English",
    dataManagement: "Data Management",
    openDataFolder: "Open Data Folder",
    storagePaths: "Storage Paths",
    baseDirectory: "Base Directory",
    database: "Database",
    images: "Images",
    statistics: "Statistics",
    statsSets: "Sets",
    statsProducts: "Products",
    statsManufacturers: "Manufacturers",
    statsTypes: "Types",
    statsTags: "Tags",
    statsImages: "Images",
    about: "About",
    version: "Version",

    // Messages
    exportSuccess: "Export successful",
    exportError: "Export failed",
    importSuccess: "Import successful",
    importError: "Import failed",
    saveSuccess: "Saved",
    saveError: "Save failed",
    deleteSuccess: "Deleted",
    deleteError: "Delete failed",
  },
};

// Search prefix mappings - both languages map to the same internal keys
export const searchPrefixes: Record<Locale, Record<string, string>> = {
  de: {
    "@box": "box",
    "@karton": "box",
    "@produkt": "product",
    "@hersteller": "manufacturer",
    "@tag": "tag",
    "@ort": "location",
    "@raum": "location",
  },
  en: {
    "@box": "box",
    "@product": "product",
    "@manufacturer": "manufacturer",
    "@tag": "tag",
    "@location": "location",
    "@room": "location",
  },
};

// Get all valid prefixes for current locale (for parsing)
export function getSearchPrefixes(): Record<string, string> {
  // Return prefixes for both languages so users can use either
  return { ...searchPrefixes.de, ...searchPrefixes.en };
}

export function useI18n() {
  const locale = computed(() => currentLocale.value);

  const t = (key: string): string => {
    const trans = translations[currentLocale.value] as Record<string, string>;
    return trans[key] || key;
  };

  const setLocale = (newLocale: Locale) => {
    currentLocale.value = newLocale;
    localStorage.setItem("samla-locale", newLocale);
  };

  return {
    locale,
    t,
    setLocale,
    getSearchPrefixes,
  };
}
