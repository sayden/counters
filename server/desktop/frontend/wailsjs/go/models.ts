export namespace server {
	
	export class CounterImage {
	    filename: string;
	    counter: string;
	    id: string;
	    pretty_name: string;
	
	    static createFrom(source: any = {}) {
	        return new CounterImage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.filename = source["filename"];
	        this.counter = source["counter"];
	        this.id = source["id"];
	        this.pretty_name = source["pretty_name"];
	    }
	}

}

