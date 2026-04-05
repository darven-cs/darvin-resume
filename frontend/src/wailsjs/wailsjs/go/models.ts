export namespace ai {
	
	export class ChatMessage {
	    id: string;
	    resumeId: string;
	    role: string;
	    content: string;
	    quotedText?: string;
	    createdAt: number;
	
	    static createFrom(source: any = {}) {
	        return new ChatMessage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.resumeId = source["resumeId"];
	        this.role = source["role"];
	        this.content = source["content"];
	        this.quotedText = source["quotedText"];
	        this.createdAt = source["createdAt"];
	    }
	}

}

export namespace model {
	
	export class BasicInfo {
	    name: string;
	    phone: string;
	    email: string;
	    avatar: string;
	    website: string;
	    github: string;
	    address: string;
	    summary: string;
	
	    static createFrom(source: any = {}) {
	        return new BasicInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.phone = source["phone"];
	        this.email = source["email"];
	        this.avatar = source["avatar"];
	        this.website = source["website"];
	        this.github = source["github"];
	        this.address = source["address"];
	        this.summary = source["summary"];
	    }
	}
	export class Module {
	    type: string;
	    title: string;
	    order: number;
	    items: number[];
	    visible: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Module(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.title = source["title"];
	        this.order = source["order"];
	        this.items = source["items"];
	        this.visible = source["visible"];
	    }
	}
	export class Resume {
	    id: string;
	    title: string;
	    basicInfo: BasicInfo;
	    modules: Module[];
	    templateId: string;
	    customCss: string;
	    markdownContent: string;
	    jobTarget: string;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Resume(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.basicInfo = this.convertValues(source["basicInfo"], BasicInfo);
	        this.modules = this.convertValues(source["modules"], Module);
	        this.templateId = source["templateId"];
	        this.customCss = source["customCss"];
	        this.markdownContent = source["markdownContent"];
	        this.jobTarget = source["jobTarget"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
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
	export class ResumeListItem {
	    id: string;
	    title: string;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new ResumeListItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
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

