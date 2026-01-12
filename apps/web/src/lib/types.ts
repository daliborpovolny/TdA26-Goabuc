export interface Course {
	uuid: string;
	type: string;
	name: string;
	description: string;
	materials: Material[];
	quizzes: Quiz[];
}

export interface BaseMaterial {
	uuid: string;
	name: string;
	description: string;
	type: 'file' | 'url';
}

export interface FileMaterial extends BaseMaterial {
	type: 'file';
	fileUrl: string;
	mimeType: string;
	sizeBytes: number;
}

export interface UrlMaterial extends BaseMaterial {
	type: 'url';
	url: string;
	faviconUrl: string;
}

export type Material = FileMaterial | UrlMaterial;

export interface BaseQuestion {
	uuid: string;
	type: 'singleChoice' | 'multipleChoice';
	question: string;
	options: string[];
}

export interface SingleChoiceQuestion extends BaseQuestion {
	type: 'singleChoice';
	correctIndex: number;
}

export interface MultipleChoiceQuestion extends BaseQuestion {
	type: 'multipleChoice';
	correctIndices: number[];
}

export type Question = SingleChoiceQuestion | MultipleChoiceQuestion;

export interface Quiz {
	uuid: string;
	title: string;
	attemptsCount: number;
	questions: Question[];
}

export interface Answer {
	uuid: string;
	selectedIndex?: number;
	selectedIndices?: number[];
	comment?: string;
}

export interface QuizSubmit {
	answers: Answer[];
}

export interface QuizMarked {
	quizUuid: string;
	score: number;
	maxScore: number;
	correctPerQuestion: boolean[];
	submittedAt: string;
}