import React from 'react'

interface PaginationProps {
	currentPage: number
	totalPages: number
	pageSize: number
	setCurrentPage: React.Dispatch<React.SetStateAction<number>>
	setPageSize: React.Dispatch<React.SetStateAction<number>>
	fetchLinks: () => void
}

const Pagination: React.FC<PaginationProps> = ({
	currentPage,
	totalPages,
	pageSize,
	setCurrentPage,
	setPageSize,
	fetchLinks
}) => {
	const getPageNumbers = () => {
		let pages = []

		if (totalPages <= 5) {
			pages = Array.from({ length: totalPages }, (_, i) => i + 1)
		} else {
			pages.push(1)

			const startPage = Math.max(2, currentPage - 1)
			const endPage = Math.min(totalPages - 1, currentPage + 1)

			if (startPage > 2) {
				pages.push(-1)
			}
			for (let i = startPage; i <= endPage; i++) {
				pages.push(i)
			}
			if (endPage < totalPages - 1) {
				pages.push(-2)
			}

			if (totalPages > 1) {
				pages.push(totalPages)
			}
		}

		return pages
	}

	return (
		<div className='d-flex justify-content-between align-items-center flex-wrap gap-2'>
			<div>
				<select
					className='form-select form-select-sm'
					value={pageSize}
					onChange={e => {
						setPageSize(Number(e.target.value))
						setCurrentPage(1)
					}}
				>
					<option value='5'>5 per page</option>
					<option value='10'>10 per page</option>
					<option value='25'>25 per page</option>
					<option value='50'>50 per page</option>
				</select>
			</div>

			{totalPages > 0 && (
				<nav>
					<ul className='pagination mb-0'>
						<li className={`page-item ${currentPage === 1 ? 'disabled' : ''}`}>
							<button
								className='page-link'
								onClick={() => setCurrentPage(prev => Math.max(prev - 1, 1))}
								disabled={currentPage === 1}
							>
								Previous
							</button>
						</li>

						{getPageNumbers().map(pageNum => {
							if (pageNum < 0) {
								return (
									<li key={`ellipsis${pageNum}`} className='page-item disabled'>
										<span className='page-link'>...</span>
									</li>
								)
							}

							return (
								<li
									key={pageNum}
									className={`page-item ${
										currentPage === pageNum ? 'active' : ''
									}`}
								>
									<button
										className='page-link'
										onClick={() => setCurrentPage(pageNum)}
									>
										{pageNum}
									</button>
								</li>
							)
						})}

						<li
							className={`page-item ${
								currentPage === totalPages ? 'disabled' : ''
							}`}
						>
							<button
								className='page-link'
								onClick={() =>
									setCurrentPage(prev => Math.min(prev + 1, totalPages))
								}
								disabled={currentPage === totalPages}
							>
								Next
							</button>
						</li>
					</ul>
				</nav>
			)}
			<button className='btn btn-primary' onClick={fetchLinks}>
				Refresh Data
			</button>
		</div>
	)
}

export default Pagination
