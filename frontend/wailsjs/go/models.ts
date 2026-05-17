export namespace models {
	
	export class AuditEntry {
	    id: string;
	    userId: string;
	    action: string;
	    entityType: string;
	    entityId: string;
	    details: string;
	    // Go type: time
	    createdAt: any;
	
	    static createFrom(source: any = {}) {
	        return new AuditEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.userId = source["userId"];
	        this.action = source["action"];
	        this.entityType = source["entityType"];
	        this.entityId = source["entityId"];
	        this.details = source["details"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
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
	export class Client {
	    id: string;
	    name: string;
	    phone: string;
	    debtLimit: number;
	
	    static createFrom(source: any = {}) {
	        return new Client(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.phone = source["phone"];
	        this.debtLimit = source["debtLimit"];
	    }
	}
	export class CloseShiftRequest {
	    vehicleId: string;
	    startOdometer: number;
	    endOdometer: number;
	    fuelStart: number;
	    fuelEnd: number;
	    fuelRefilled: number;
	    comment: string;
	
	    static createFrom(source: any = {}) {
	        return new CloseShiftRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.vehicleId = source["vehicleId"];
	        this.startOdometer = source["startOdometer"];
	        this.endOdometer = source["endOdometer"];
	        this.fuelStart = source["fuelStart"];
	        this.fuelEnd = source["fuelEnd"];
	        this.fuelRefilled = source["fuelRefilled"];
	        this.comment = source["comment"];
	    }
	}
	export class Waybill {
	    id: string;
	    driverId: string;
	    vehicleId: string;
	    date: string;
	    startOdometer: number;
	    endOdometer: number;
	    fuelStart: number;
	    fuelEnd: number;
	    fuelRefilled: number;
	    actualConsumption: number;
	    normConsumption: number;
	    overPercent: number;
	    comment: string;
	    closed: boolean;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    closedAt?: any;
	
	    static createFrom(source: any = {}) {
	        return new Waybill(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.driverId = source["driverId"];
	        this.vehicleId = source["vehicleId"];
	        this.date = source["date"];
	        this.startOdometer = source["startOdometer"];
	        this.endOdometer = source["endOdometer"];
	        this.fuelStart = source["fuelStart"];
	        this.fuelEnd = source["fuelEnd"];
	        this.fuelRefilled = source["fuelRefilled"];
	        this.actualConsumption = source["actualConsumption"];
	        this.normConsumption = source["normConsumption"];
	        this.overPercent = source["overPercent"];
	        this.comment = source["comment"];
	        this.closed = source["closed"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.closedAt = this.convertValues(source["closedAt"], null);
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
	export class CloseShiftResult {
	    waybill: Waybill;
	    blocked: boolean;
	    requireComment: boolean;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new CloseShiftResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.waybill = this.convertValues(source["waybill"], Waybill);
	        this.blocked = source["blocked"];
	        this.requireComment = source["requireComment"];
	        this.message = source["message"];
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
	export class CreateOrderRequest {
	    clientId: string;
	    vehicleId: string;
	    driverId: string;
	    fromAddr: string;
	    toAddr: string;
	    distanceKm: number;
	    idleHours: number;
	    urgent: boolean;
	    tariffId: string;
	    scheduledAt: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateOrderRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.clientId = source["clientId"];
	        this.vehicleId = source["vehicleId"];
	        this.driverId = source["driverId"];
	        this.fromAddr = source["fromAddr"];
	        this.toAddr = source["toAddr"];
	        this.distanceKm = source["distanceKm"];
	        this.idleHours = source["idleHours"];
	        this.urgent = source["urgent"];
	        this.tariffId = source["tariffId"];
	        this.scheduledAt = source["scheduledAt"];
	    }
	}
	export class DashboardStats {
	    ordersToday: number;
	    activeVehicles: number;
	    openWaybills: number;
	    revenueMonth: number;
	    fuelOverruns: number;
	
	    static createFrom(source: any = {}) {
	        return new DashboardStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ordersToday = source["ordersToday"];
	        this.activeVehicles = source["activeVehicles"];
	        this.openWaybills = source["openWaybills"];
	        this.revenueMonth = source["revenueMonth"];
	        this.fuelOverruns = source["fuelOverruns"];
	    }
	}
	export class LoginRequest {
	    login: string;
	    password: string;
	
	    static createFrom(source: any = {}) {
	        return new LoginRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.login = source["login"];
	        this.password = source["password"];
	    }
	}
	export class User {
	    id: string;
	    login: string;
	    role: string;
	    firstName: string;
	    lastName: string;
	    phone: string;
	    active: boolean;
	
	    static createFrom(source: any = {}) {
	        return new User(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.login = source["login"];
	        this.role = source["role"];
	        this.firstName = source["firstName"];
	        this.lastName = source["lastName"];
	        this.phone = source["phone"];
	        this.active = source["active"];
	    }
	}
	export class LoginResponse {
	    token: string;
	    user: User;
	
	    static createFrom(source: any = {}) {
	        return new LoginResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.token = source["token"];
	        this.user = this.convertValues(source["user"], User);
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
	export class Order {
	    id: string;
	    clientId: string;
	    vehicleId: string;
	    driverId: string;
	    fromAddr: string;
	    toAddr: string;
	    distanceKm: number;
	    idleHours: number;
	    urgent: boolean;
	    tariffId: string;
	    price: number;
	    status: string;
	    // Go type: time
	    scheduledAt: any;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    startedAt?: any;
	    // Go type: time
	    completedAt?: any;
	
	    static createFrom(source: any = {}) {
	        return new Order(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.clientId = source["clientId"];
	        this.vehicleId = source["vehicleId"];
	        this.driverId = source["driverId"];
	        this.fromAddr = source["fromAddr"];
	        this.toAddr = source["toAddr"];
	        this.distanceKm = source["distanceKm"];
	        this.idleHours = source["idleHours"];
	        this.urgent = source["urgent"];
	        this.tariffId = source["tariffId"];
	        this.price = source["price"];
	        this.status = source["status"];
	        this.scheduledAt = this.convertValues(source["scheduledAt"], null);
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.startedAt = this.convertValues(source["startedAt"], null);
	        this.completedAt = this.convertValues(source["completedAt"], null);
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
	export class Tariff {
	    id: string;
	    name: string;
	    baseFee: number;
	    pricePerKm: number;
	    pricePerIdleHr: number;
	    urgencyCoeff: number;
	    active: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Tariff(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.baseFee = source["baseFee"];
	        this.pricePerKm = source["pricePerKm"];
	        this.pricePerIdleHr = source["pricePerIdleHr"];
	        this.urgencyCoeff = source["urgencyCoeff"];
	        this.active = source["active"];
	    }
	}
	
	export class Vehicle {
	    id: string;
	    plate: string;
	    model: string;
	    fuelNorm: number;
	    active: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Vehicle(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.plate = source["plate"];
	        this.model = source["model"];
	        this.fuelNorm = source["fuelNorm"];
	        this.active = source["active"];
	    }
	}
	export class VehicleScheduleItem {
	    vehicleId: string;
	    plate: string;
	    orders: Order[];
	
	    static createFrom(source: any = {}) {
	        return new VehicleScheduleItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.vehicleId = source["vehicleId"];
	        this.plate = source["plate"];
	        this.orders = this.convertValues(source["orders"], Order);
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

}

