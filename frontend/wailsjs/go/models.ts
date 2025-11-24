export namespace main {
	
	export class GridItem {
	    id: number;
	    status: number;
	    is_marked: boolean;
	
	    static createFrom(source: any = {}) {
	        return new GridItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.status = source["status"];
	        this.is_marked = source["is_marked"];
	    }
	}
	export class QuestionView {
	    id: number;
	    type: string;
	    content: string;
	    options: string[];
	    user_answer: string;
	    status: number;
	    is_marked: boolean;
	    ai_explanation: string;
	    explanation?: string;
	    correct_answer?: string;
	
	    static createFrom(source: any = {}) {
	        return new QuestionView(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.type = source["type"];
	        this.content = source["content"];
	        this.options = source["options"];
	        this.user_answer = source["user_answer"];
	        this.status = source["status"];
	        this.is_marked = source["is_marked"];
	        this.ai_explanation = source["ai_explanation"];
	        this.explanation = source["explanation"];
	        this.correct_answer = source["correct_answer"];
	    }
	}
	export class Stats {
	    total: number;
	    done: number;
	    correct: number;
	    accuracy: string;
	
	    static createFrom(source: any = {}) {
	        return new Stats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total = source["total"];
	        this.done = source["done"];
	        this.correct = source["correct"];
	        this.accuracy = source["accuracy"];
	    }
	}
	export class SubmitResult {
	    correct: boolean;
	    explanation: string;
	    correct_answer: string;
	    ai_explanation: string;
	
	    static createFrom(source: any = {}) {
	        return new SubmitResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.correct = source["correct"];
	        this.explanation = source["explanation"];
	        this.correct_answer = source["correct_answer"];
	        this.ai_explanation = source["ai_explanation"];
	    }
	}

}

