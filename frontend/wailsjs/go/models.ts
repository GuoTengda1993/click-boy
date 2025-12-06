export namespace click {
	
	export class Params {
	    times: number;
	    interval: number;
	    duration: number;
	    package: string;
	    activity: string;
	
	    static createFrom(source: any = {}) {
	        return new Params(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.times = source["times"];
	        this.interval = source["interval"];
	        this.duration = source["duration"];
	        this.package = source["package"];
	        this.activity = source["activity"];
	    }
	}
	export class Point {
	    x: number;
	    y: number;
	    event: string;
	    edit: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Point(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.x = source["x"];
	        this.y = source["y"];
	        this.event = source["event"];
	        this.edit = source["edit"];
	    }
	}

}

