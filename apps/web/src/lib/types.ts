export interface Course {
	uuid: string;
	type: string;
	name: string;
	description: string;
	materials: Material[];
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
