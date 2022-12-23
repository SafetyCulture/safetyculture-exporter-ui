export namespace api {
	
	export class ExportStatusResponseItem {
	    feed_name: string;
	    has_started: boolean;
	    debug_string: string;
	
	    static createFrom(source: any = {}) {
	        return new ExportStatusResponseItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.feed_name = source["feed_name"];
	        this.has_started = source["has_started"];
	        this.debug_string = source["debug_string"];
	    }
	}
	export class ExportStatusResponse {
	    export_started: boolean;
	    export_completed: boolean;
	    feeds: ExportStatusResponseItem[];
	
	    static createFrom(source: any = {}) {
	        return new ExportStatusResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.export_started = source["export_started"];
	        this.export_completed = source["export_completed"];
	        this.feeds = this.convertValues(source["feeds"], ExportStatusResponseItem);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
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
	
	
	export class TemplateResponseItem {
	    id: string;
	    name: string;
	    // Go type: time.Time
	    modified_at: any;
	
	    static createFrom(source: any = {}) {
	        return new TemplateResponseItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.modified_at = this.convertValues(source["modified_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
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

