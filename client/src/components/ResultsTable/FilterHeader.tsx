import React from 'react'
interface FilterHeaderProps {
	filters: Partial<Record<string, string>>
	handleFilterChange: (field: string, value: string) => void
}

const FilterHeader: React.FC<FilterHeaderProps> = ({
	filters,
	handleFilterChange
}) => (
	<tr>
		<th>
			<input
				type='text'
				className='form-control form-control-sm'
				placeholder='Filter URL'
				value={filters.url || ''}
				onChange={e => handleFilterChange('url', e.target.value)}
			/>
		</th>
		<th>
			<input
				type='text'
				className='form-control form-control-sm'
				placeholder='Title'
				value={filters.title || ''}
				onChange={e => handleFilterChange('title', e.target.value)}
			/>
		</th>
		<th>
			<input
				type='text'
				className='form-control form-control-sm'
				placeholder='Version'
				value={filters.html_version || ''}
				onChange={e => handleFilterChange('html_version', e.target.value)}
			/>
		</th>
		<th>
			<input
				type='text'
				className='form-control form-control-sm'
				placeholder='H1'
				value={filters['headings_count.h1'] || ''}
				onChange={e => handleFilterChange('headings_count.h1', e.target.value)}
			/>
		</th>
		<th>
			<input
				type='text'
				className='form-control form-control-sm'
				placeholder='H2'
				value={filters['headings_count.h2'] || ''}
				onChange={e => handleFilterChange('headings_count.h2', e.target.value)}
			/>
		</th>
		<th>
			<input
				type='text'
				className='form-control form-control-sm'
				placeholder='H3'
				value={filters['headings_count.h3'] || ''}
				onChange={e => handleFilterChange('headings_count.h3', e.target.value)}
			/>
		</th>
		<th>
			<input
				type='text'
				className='form-control form-control-sm'
				placeholder='H4'
				value={filters['headings_count.h4'] || ''}
				onChange={e => handleFilterChange('headings_count.h4', e.target.value)}
			/>
		</th>
		<th>
			<input
				type='text'
				className='form-control form-control-sm'
				placeholder='H5'
				value={filters['headings_count.h5'] || ''}
				onChange={e => handleFilterChange('headings_count.h5', e.target.value)}
			/>
		</th>
		<th>
			<input
				type='text'
				className='form-control form-control-sm'
				placeholder='H6'
				value={filters['headings_count.h6'] || ''}
				onChange={e => handleFilterChange('headings_count.h6', e.target.value)}
			/>
		</th>
		<th>
			<input
				type='text'
				className='form-control form-control-sm'
				placeholder='Internal'
				value={filters.internal_links || ''}
				onChange={e => handleFilterChange('internal_links', e.target.value)}
			/>
		</th>
		<th>
			<input
				type='text'
				className='form-control form-control-sm'
				placeholder='External'
				value={filters.external_links || ''}
				onChange={e => handleFilterChange('external_links', e.target.value)}
			/>
		</th>
		<th>
			<input
				type='text'
				className='form-control form-control-sm'
				placeholder='Inaccessible'
				value={filters.inaccessible_links || ''}
				onChange={e => handleFilterChange('inaccessible_links', e.target.value)}
			/>
		</th>
		<th>
			<select
				className='form-select form-select-sm'
				value={filters.has_login_form || ''}
				onChange={e => handleFilterChange('has_login_form', e.target.value)}
			>
				<option value=''>All</option>
				<option value='true'>Yes</option>
				<option value='false'>No</option>
			</select>
		</th>
		<th>
			<input
				type='text'
				className='form-control form-control-sm'
				placeholder='Status'
				value={filters.status || ''}
				onChange={e => handleFilterChange('status', e.target.value)}
			/>
		</th>
		<th></th>
	</tr>
)

export default FilterHeader
