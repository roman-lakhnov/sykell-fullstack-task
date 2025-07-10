import React, { type JSX } from 'react'
import type { SortField } from './types'

interface TableHeaderProps {
	handleSort: (field: SortField) => void
	getSortIcon: (field: SortField) => JSX.Element
}

const TableHeader: React.FC<TableHeaderProps> = ({
	handleSort,
	getSortIcon
}) => (
	<tr>
		<th className='cursor-pointer' onClick={() => handleSort('url')}>
			<p className='m-0'>URL</p> {getSortIcon('url')}
		</th>
		<th className='cursor-pointer' onClick={() => handleSort('title')}>
			<p className='m-0'>Title</p> {getSortIcon('title')}
		</th>
		<th className='cursor-pointer' onClick={() => handleSort('html_version')}>
			<p className='m-0'>HTML Version</p> {getSortIcon('html_version')}
		</th>
		<th
			className='cursor-pointer'
			onClick={() => handleSort('headings_count.h1')}
		>
			<p className='m-0'>H1</p> {getSortIcon('headings_count.h1')}
		</th>
		<th
			className='cursor-pointer'
			onClick={() => handleSort('headings_count.h2')}
		>
			<p className='m-0'>H2</p>
			{getSortIcon('headings_count.h2')}
		</th>
		<th
			className='cursor-pointer'
			onClick={() => handleSort('headings_count.h3')}
		>
			<p className='m-0'>H3</p>
			{getSortIcon('headings_count.h3')}
		</th>
		<th
			className='cursor-pointer'
			onClick={() => handleSort('headings_count.h4')}
		>
			<p className='m-0'>H4</p>
			{getSortIcon('headings_count.h4')}
		</th>
		<th
			className='cursor-pointer'
			onClick={() => handleSort('headings_count.h5')}
		>
			<p className='m-0'>H5</p>
			{getSortIcon('headings_count.h5')}
		</th>
		<th
			className='cursor-pointer'
			onClick={() => handleSort('headings_count.h6')}
		>
			<p className='m-0'>H6</p>
			{getSortIcon('headings_count.h6')}
		</th>
		<th className='cursor-pointer' onClick={() => handleSort('internal_links')}>
			<p className='m-0'>Internal Links</p>
			{getSortIcon('internal_links')}
		</th>
		<th className='cursor-pointer' onClick={() => handleSort('external_links')}>
			<p className='m-0'>External Links</p>
			{getSortIcon('external_links')}
		</th>
		<th
			className='cursor-pointer'
			onClick={() => handleSort('inaccessible_links')}
		>
			<p className='m-0'>Inaccessible Links</p>
			{getSortIcon('inaccessible_links')}
		</th>
		<th className='cursor-pointer' onClick={() => handleSort('has_login_form')}>
			<p className='m-0'>Login Form</p>
			{getSortIcon('has_login_form')}
		</th>
		<th className='cursor-pointer' onClick={() => handleSort('status')}>
			<p className='m-0'>Status</p>
			{getSortIcon('status')}
		</th>
		<th className='align-middle'>Action</th>
	</tr>
)

export default TableHeader
