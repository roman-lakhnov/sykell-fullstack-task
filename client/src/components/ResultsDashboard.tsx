import React, { useCallback, useEffect, useState } from 'react'
import { FaSort, FaSortDown, FaSortUp } from 'react-icons/fa'
import { toast } from 'react-toastify'
import {
	FilterHeader,
	Pagination,
	TableHeader,
	TableRows
} from './ResultsTable'
import type { Link, SortDirection, SortField } from './ResultsTable/types'

// Main Component
const ResultsDashboard: React.FC = () => {
	const [links, setLinks] = useState<Link[]>([])

	// Pagination state
	const [currentPage, setCurrentPage] = useState<number>(1)
	const [pageSize, setPageSize] = useState<number>(10)
	const [totalPages, setTotalPages] = useState<number>(1)

	// Sorting state
	const [sortField, setSortField] = useState<SortField>('')
	const [sortDirection, setSortDirection] = useState<SortDirection>('')

	// Filtering state
	const [filters, setFilters] = useState<Partial<Record<keyof Link, string>>>(
		{}
	)
	const [globalSearch, setGlobalSearch] = useState<string>('')

	const fetchLinks = useCallback(async () => {
		try {
			const response = await fetch(
				`http://localhost:8080/links?page=${currentPage}&amount=${pageSize}`
			)

			if (!response.ok) {
				throw new Error('Failed to fetch data')
			}

			const data = await response.json()
			setLinks(data.links || [])
			setTotalPages(data.pagination?.total_pages)
			// toast.success('Data Refreshed!')
		} catch (err) {
			console.error('Error fetching links:', err)
			toast.error('Failed to fetch links data. Please try again later.')
		}	
	}, [currentPage, pageSize])

	// Handle action (stop/analyze) for a link
	const handleAction = async (id: number, newStatus: string) => {
		try {
			const response = await fetch('http://localhost:8080/links', {
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					id: id,
					status: newStatus
				})
			})

			if (!response.ok) {
				const errorData = await response.json()
				throw new Error(errorData.error || 'Failed to update link status')
			}

			toast.success(
				`Link ${
					newStatus === 'stop' ? 'stopped' : 'queued for analysis'
				} successfully`
			)
			// Refresh data after action
			fetchLinks()
		} catch (error) {
			console.error('Error updating link status:', error)
			toast.error(
				`Failed to update link: ${
					error instanceof Error ? error.message : 'Unknown error'
				}`
			)
		}
	}

	useEffect(() => {
		fetchLinks()
	}, [currentPage, pageSize, fetchLinks])

	// Function to handle sorting
	const handleSort = (field: SortField) => {
		if (sortField === field) {
			// Toggle direction if already sorting by this field
			setSortDirection(
				sortDirection === 'asc' ? 'desc' : sortDirection === 'desc' ? '' : 'asc'
			)
			if (sortDirection === 'desc') {
				setSortField('')
			}
		} else {
			// Start with ascending sort for new field
			setSortField(field)
			setSortDirection('asc')
		}
	}

	// Update the handleFilterChange function to accept string instead of keyof Link
	const handleFilterChange = (field: string, value: string) => {
		setFilters(prev => ({
			...prev,
			[field]: value
		}))
		setCurrentPage(1) // Reset to first page when filters change
	}

	// Function to handle global search
	const handleGlobalSearch = (e: React.ChangeEvent<HTMLInputElement>) => {
		setGlobalSearch(e.target.value)
		setCurrentPage(1) // Reset to first page when search changes
	}

	// Apply sorting and filtering to links
	const filteredAndSortedLinks = React.useMemo(() => {
		// First apply filters
		let result = [...links]

		// Apply column filters
		Object.entries(filters).forEach(([field, value]) => {
			if (value && value.trim() !== '') {
				result = result.filter(link => {
					// Handle nested properties like headings_count.h1
					if (field.includes('.')) {
						const [parent, child] = field.split('.')
						const parentValue = link[parent as keyof Link]
						if (parentValue && typeof parentValue === 'object') {
							const childValue = (parentValue as Record<string, unknown>)[child]
							if (
								typeof childValue === 'number' ||
								typeof childValue === 'string'
							) {
								return String(childValue).includes(value)
							}
						}
						return false
					} else {
						// Handle regular properties
						const fieldValue = link[field as keyof Link]
						if (typeof fieldValue === 'string') {
							return fieldValue.toLowerCase().includes(value.toLowerCase())
						} else if (
							typeof fieldValue === 'number' ||
							typeof fieldValue === 'boolean'
						) {
							return String(fieldValue).includes(value)
						}
					}
					return true
				})
			}
		})

		// Apply global search across multiple fields
		if (globalSearch.trim() !== '') {
			result = result.filter(
				link =>
					link.url.toLowerCase().includes(globalSearch.toLowerCase()) ||
					link.title?.toLowerCase().includes(globalSearch.toLowerCase()) ||
					link.html_version
						?.toLowerCase()
						.includes(globalSearch.toLowerCase()) ||
					link.status?.toLowerCase().includes(globalSearch.toLowerCase()) ||
					link.headings_count?.h1?.toString().includes(globalSearch) ||
					link.headings_count?.h2?.toString().includes(globalSearch) ||
					link.headings_count?.h3?.toString().includes(globalSearch) ||
					link.headings_count?.h4?.toString().includes(globalSearch) ||
					link.headings_count?.h5?.toString().includes(globalSearch) ||
					link.headings_count?.h6?.toString().includes(globalSearch) ||
					link.inaccessible_links?.toString().includes(globalSearch) ||
					link.external_links?.toString().includes(globalSearch) ||
					link.internal_links?.toString().includes(globalSearch)
			)
		}

		// Then apply sorting
		if (sortField) {
			result.sort((a, b) => {
				// Handle nested properties like headings_count.h1
				if (sortField.includes('.')) {
					const [parent, child] = sortField.split('.')
					const parentA = a[parent as keyof Link] as
						| Record<string, unknown>
						| undefined
					const parentB = b[parent as keyof Link] as
						| Record<string, unknown>
						| undefined
					const valueA = parentA?.[child] ?? 0
					const valueB = parentB?.[child] ?? 0

					return sortDirection === 'asc'
						? Number(valueA) - Number(valueB)
						: Number(valueB) - Number(valueA)
				}

				const valueA = a[sortField as keyof Link]
				const valueB = b[sortField as keyof Link]

				if (typeof valueA === 'string' && typeof valueB === 'string') {
					return sortDirection === 'asc'
						? valueA.localeCompare(valueB)
						: valueB.localeCompare(valueA)
				} else if (
					(typeof valueA === 'number' && typeof valueB === 'number') ||
					(typeof valueA === 'boolean' && typeof valueB === 'boolean')
				) {
					return sortDirection === 'asc'
						? Number(valueA) - Number(valueB)
						: Number(valueB) - Number(valueA)
				}
				return 0
			})
		}

		return result
	}, [links, filters, globalSearch, sortField, sortDirection])

	const getSortIcon = (field: SortField) => {
		if (sortField !== field)
			return <FaSort className='ms-1 text-muted opacity-50' />
		if (sortDirection === 'asc')
			return <FaSortUp className='ms-1 text-primary' />
		if (sortDirection === 'desc')
			return <FaSortDown className='ms-1 text-primary' />
		return <FaSort className='ms-1 text-muted opacity-50' />
	}

	// Render function for the table content
	const renderTableContent = () => {
		return (
			<>
				<div className='table-responsive'>
					<table className='table table-striped table-hover'>
						<thead>
							<FilterHeader
								filters={filters}
								handleFilterChange={handleFilterChange}
							/>
							<TableHeader handleSort={handleSort} getSortIcon={getSortIcon} />
						</thead>
						<tbody>
							<TableRows
								links={filteredAndSortedLinks}
								handleAction={handleAction}
							/>
						</tbody>
					</table>
				</div>

				<Pagination
					currentPage={currentPage}
					totalPages={totalPages}
					pageSize={pageSize}
					setCurrentPage={setCurrentPage}
					setPageSize={setPageSize}
					fetchLinks={fetchLinks}
				/>
			</>
		)
	}

	return (
		<div className='container mt-5 mb-5'>
			<div className='row'>
				<div className='col-12'>
					<div className='card'>
						<div className='card-header bg-secondary text-white'>
							<h3 className='mb-0'>Results Dashboard</h3>
						</div>
						<div className='card-body d-flex flex-column gap-3'>
							<div className='mb-3'>
								<input
									type='text'
									className='form-control'
									placeholder='Search across all fields...'
									value={globalSearch}
									onChange={handleGlobalSearch}
								/>
							</div>
							{renderTableContent()}
						</div>
					</div>
				</div>
			</div>
		</div>
	)
}

export default ResultsDashboard
