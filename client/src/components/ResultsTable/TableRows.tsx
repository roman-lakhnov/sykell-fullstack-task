import React from 'react'
import type { Link } from './types'

interface TableRowsProps {
	links: Link[]
	handleAction: (id: number, newStatus: string) => Promise<void>
}

const TableRows: React.FC<TableRowsProps> = ({ links, handleAction }) => (
	<>
		{links.length > 0 ? (
			links.map(link => (
				<tr key={link.id} className='align-middle'>
					<td>
						<a
							href={link.url}
							target='_blank'
							rel='noopener noreferrer'
							className='align-middle d-inline-block'
							style={{ maxWidth: '150px' }}
						>
							{link.url}
						</a>
					</td>
					<td style={{ maxWidth: '150px' }}>{link.title || 'N/A'}</td>
					<td>{link.html_version || 'N/A'}</td>
					<td>{link.headings_count?.h1 || 0}</td>
					<td>{link.headings_count?.h2 || 0}</td>
					<td>{link.headings_count?.h3 || 0}</td>
					<td>{link.headings_count?.h4 || 0}</td>
					<td>{link.headings_count?.h5 || 0}</td>
					<td>{link.headings_count?.h6 || 0}</td>
					<td>{link.internal_links || 0}</td>
					<td>{link.external_links || 0}</td>
					<td>{link.inaccessible_links || 0}</td>
					<td>{link.has_login_form ? 'Yes' : 'No'}</td>
					<td>
						<span
							className={`badge ${
								link.status === 'created'
									? 'bg-primary'
									: link.status === 'pending'
									? 'bg-warning'
									: link.status === 'checked'
									? 'bg-success'
									: 'bg-danger'
							}`}
						>
							{link.status}
						</span>
					</td>
					<td>
						{link.status === 'stop' ||
						link.status === 'checked' ||
						link.status === 'error' ? (
							<button
								className='btn btn-sm btn-primary'
								onClick={() => handleAction(link.id, 'created')}
							>
								analyze
							</button>
						) : (
							<button
								className='btn btn-sm btn-danger'
								onClick={() => handleAction(link.id, 'stop')}
								disabled={link.status === 'pending'}
							>
								stop
							</button>
						)}
					</td>
				</tr>
			))
		) : (
			<tr>
				<td colSpan={15} className='text-center'>
					No results found
				</td>
			</tr>
		)}
	</>
)

export default TableRows
