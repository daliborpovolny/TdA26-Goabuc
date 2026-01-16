/// Types for interaction with backend

// Course

export interface Course {
	uuid: string;
	type: string;
	name: string;
	description: string;
	materials: Material[];
	quizzes: Quiz[];
}

// Materials

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

// Quizzes

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

// Feed

export interface FeedPost {
	uuid: string;
	type: 'manual' | 'system';
	message: string;
	edited: boolean;
	createdAt: string;
	updatedAt: string;
}

// CREATE TABLE IF NOT EXISTS answer (
//     quiz_uuid TEXT NOT NULL,
//     comment TEXT,

//     score INTEGER NOT NULL,
//     max_score INTEGER NOT NULL,

//     user_id INTEGER,
//     attempt_number INTEGER NOT NULL,

//     submitted_at INTEGER NOT NULL,

//     FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
//     FOREIGN KEY (quiz_uuid) REFERENCES quizz(uuid) ON DELETE CASCADE
// );


export interface QuizOutcome {
	quiz_uuid: string;
	comment: string;

	score: number;
	max_score: number

	user_id: number;
	attempt_number: number;

	submitted_at: number;
}
