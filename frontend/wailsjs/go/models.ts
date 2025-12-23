export namespace main {
	
	export class AppPaths {
	    baseDir: string;
	    dataDir: string;
	    imagesDir: string;
	    dbPath: string;
	
	    static createFrom(source: any = {}) {
	        return new AppPaths(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.baseDir = source["baseDir"];
	        this.dataDir = source["dataDir"];
	        this.imagesDir = source["imagesDir"];
	        this.dbPath = source["dbPath"];
	    }
	}
	export class BagInfo {
	    id: number;
	    serialNo: string;
	    boxId: number;
	    boxCode: string;
	    boxName: string;
	    locationId: number;
	    locationName: string;
	    locationRoom: string;
	    locationShelf: string;
	    locationCompartment: string;
	    locationNote: string;
	
	    static createFrom(source: any = {}) {
	        return new BagInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.serialNo = source["serialNo"];
	        this.boxId = source["boxId"];
	        this.boxCode = source["boxCode"];
	        this.boxName = source["boxName"];
	        this.locationId = source["locationId"];
	        this.locationName = source["locationName"];
	        this.locationRoom = source["locationRoom"];
	        this.locationShelf = source["locationShelf"];
	        this.locationCompartment = source["locationCompartment"];
	        this.locationNote = source["locationNote"];
	    }
	}
	export class Box {
	    id: number;
	    locationId: number;
	    code: string;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new Box(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.locationId = source["locationId"];
	        this.code = source["code"];
	        this.name = source["name"];
	    }
	}
	export class Manufacturer {
	    id: number;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new Manufacturer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	    }
	}
	export class Product {
	    id: number;
	    setId: number;
	    name: string;
	    kind: string;
	
	    static createFrom(source: any = {}) {
	        return new Product(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.setId = source["setId"];
	        this.name = source["name"];
	        this.kind = source["kind"];
	    }
	}
	export class ScanResult {
	    base64Data: string;
	    relPath: string;
	
	    static createFrom(source: any = {}) {
	        return new ScanResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.base64Data = source["base64Data"];
	        this.relPath = source["relPath"];
	    }
	}
	export class SetDetails {
	    id: number;
	    name: string;
	    manufacturerId?: number;
	    manufacturerName: string;
	    typeId?: number;
	    typeName: string;
	    bag: BagInfo;
	    photoPath: string;
	    photoSource: string;
	    tags: string[];
	    products: Product[];
	
	    static createFrom(source: any = {}) {
	        return new SetDetails(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.manufacturerId = source["manufacturerId"];
	        this.manufacturerName = source["manufacturerName"];
	        this.typeId = source["typeId"];
	        this.typeName = source["typeName"];
	        this.bag = this.convertValues(source["bag"], BagInfo);
	        this.photoPath = source["photoPath"];
	        this.photoSource = source["photoSource"];
	        this.tags = source["tags"];
	        this.products = this.convertValues(source["products"], Product);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SetSearchResult {
	    setId: number;
	    setName: string;
	    manufacturerName: string;
	    boxCode: string;
	    boxName: string;
	    bagSerial: string;
	    locationName: string;
	    tags: string[];
	    thumbnailPath: string;
	
	    static createFrom(source: any = {}) {
	        return new SetSearchResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.setId = source["setId"];
	        this.setName = source["setName"];
	        this.manufacturerName = source["manufacturerName"];
	        this.boxCode = source["boxCode"];
	        this.boxName = source["boxName"];
	        this.bagSerial = source["bagSerial"];
	        this.locationName = source["locationName"];
	        this.tags = source["tags"];
	        this.thumbnailPath = source["thumbnailPath"];
	    }
	}
	export class StorageLocation {
	    id: number;
	    friendlyName: string;
	    room: string;
	    shelf: string;
	    compartment: string;
	    note: string;
	
	    static createFrom(source: any = {}) {
	        return new StorageLocation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.friendlyName = source["friendlyName"];
	        this.room = source["room"];
	        this.shelf = source["shelf"];
	        this.compartment = source["compartment"];
	        this.note = source["note"];
	    }
	}
	export class Tag {
	    id: number;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new Tag(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	    }
	}
	export class Type {
	    id: number;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new Type(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	    }
	}

}

