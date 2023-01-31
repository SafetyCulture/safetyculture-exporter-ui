export namespace api {
	
	
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

export namespace main {
	
	export class VersionResponse {
	    current: string;
	    latest: string;
	    download_url: string;
	    should_update: boolean;
	
	    static createFrom(source: any = {}) {
	        return new VersionResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.current = source["current"];
	        this.latest = source["latest"];
	        this.download_url = source["download_url"];
	        this.should_update = source["should_update"];
	    }
	}

}

