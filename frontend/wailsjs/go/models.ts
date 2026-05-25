export namespace main {
	
	export class Task {
	    id: number;
	    title: string;
	    priority: string;
	    done: boolean;
	    createdAt: string;
	
	    static createFrom(source: any = {}) {
	        return new Task(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.priority = source["priority"];
	        this.done = source["done"];
	        this.createdAt = source["createdAt"];
	    }
	}
	export class InteractionState {
	    callCount: number;
	    lastMessage: string;
	    tasks: Task[];
	
	    static createFrom(source: any = {}) {
	        return new InteractionState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.callCount = source["callCount"];
	        this.lastMessage = source["lastMessage"];
	        this.tasks = this.convertValues(source["tasks"], Task);
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
	export class ProfileInput {
	    name: string;
	    role: string;
	    years: number;
	
	    static createFrom(source: any = {}) {
	        return new ProfileInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.role = source["role"];
	        this.years = source["years"];
	    }
	}
	export class ProfileSummary {
	    title: string;
	    message: string;
	    score: number;
	    tags: string[];
	
	    static createFrom(source: any = {}) {
	        return new ProfileSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.message = source["message"];
	        this.score = source["score"];
	        this.tags = source["tags"];
	    }
	}
	
	export class TaskInput {
	    title: string;
	    priority: string;
	
	    static createFrom(source: any = {}) {
	        return new TaskInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.priority = source["priority"];
	    }
	}
	export class TrayStatus {
	    mode: string;
	    supported: boolean;
	    notes: string[];
	    menuItems: string[];
	
	    static createFrom(source: any = {}) {
	        return new TrayStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.mode = source["mode"];
	        this.supported = source["supported"];
	        this.notes = source["notes"];
	        this.menuItems = source["menuItems"];
	    }
	}
	export class WindowCommandResult {
	    action: string;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new WindowCommandResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.action = source["action"];
	        this.message = source["message"];
	    }
	}
	export class WindowSizeInput {
	    width: number;
	    height: number;
	
	    static createFrom(source: any = {}) {
	        return new WindowSizeInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.width = source["width"];
	        this.height = source["height"];
	    }
	}

}

