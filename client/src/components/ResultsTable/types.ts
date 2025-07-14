export interface LinkIssue {
	url: string
	status_code: number
}

export interface Link {
	id: number
	url: string
	post_time?: string
	status: string
	check_time?: string
	title: string
	html_version: string
	headings_count: {
		h1: number
		h2: number
		h3: number
		h4: number
		h5: number
		h6: number
	}
	internal_links: number
	external_links: number
	inaccessible_links: number
	inaccessible_details: LinkIssue[] 
	has_login_form: boolean
}

export type SortDirection = 'asc' | 'desc' | ''

export type SortField =
	| keyof Link
	| 'headings_count.h1'
	| 'headings_count.h2'
	| 'headings_count.h3'
	| 'headings_count.h4'
	| 'headings_count.h5'
	| 'headings_count.h6'
	| ''
