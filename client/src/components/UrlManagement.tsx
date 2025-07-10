import React, { useState } from 'react'
import { toast } from 'react-toastify'
import { isValidUrl } from '../utils'

const UrlManagement = () => {
	const [formData, setFormData] = useState<string>('')
	const [links, setLinks] = useState<string[]>([])

	const handleSubmit = (e: React.FormEvent) => {
		e.preventDefault()
		if (formData.trim()) {
			const urlList = formData
				.split(',')
				.map(url => url.trim())
				.filter(url => url !== '')
				.filter(url => isValidUrl(url))
			if (urlList.length > 0) {
				setLinks(prevLinks => [...prevLinks, ...urlList])
				setFormData('')
			}
			console.log('Submitted URLs:', links)
			toast.info('Links added!')
		}
	}

	const sendToServer = async () => {
		try {
			const response = await fetch('http://localhost:8080/links', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ urls: links })
			})
			if (response.ok) {
				setLinks([])
				toast.success('Successfully submitted!')
			}
		} catch (error) {
			console.error('Error submitting links:', error)
			toast.error('Failed to submit links. Please try again.')
		}
	}

	return (
		<div className='container mt-5'>
			<div className='row justify-content-center'>
				<div className='col-md-12'>
					<div className='card'>
						<div className='card-header bg-secondary text-white'>
							<h3 className='mb-0'>Links Management</h3>
						</div>
						<div className='card-body d-flex flex-column gap-3'>
							<form onSubmit={handleSubmit}>
								<div className='input-group'>
									<input
										type='text'
										className='form-control'
										value={formData}
										onChange={e => setFormData(e.target.value)}
										placeholder='Enter URLs separated by commas'
										required
									/>
									<button className='btn btn-primary' type='submit'>
										+
									</button>
								</div>
							</form>
							{links.length > 0 && (
								<div className='d-flex flex-column gap-3'>
									<div className='d-flex justify-content-between align-items-center'>
										<h4 className='m-0'>Links preview:</h4>
										<button className='btn btn-success' onClick={sendToServer}>
											Submit for Analysis
										</button>
									</div>
									<ul className='list-group'>
										<li className='list-group-item'>
											<div className='d-flex flex-wrap gap-2'>
												{links.map((url, index) => (
													<div key={index} className='border rounded p-2'>
														<a
															href={url}
															target='_blank'
															rel='noopener noreferrer'
															className='text-decoration-none'
														>
															{url}
														</a>
													</div>
												))}
											</div>
										</li>
									</ul>
								</div>
							)}
						</div>
					</div>
				</div>
			</div>
		</div>
	)
}

export default UrlManagement
