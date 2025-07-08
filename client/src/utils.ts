export const isValidUrl = (url: string): boolean => {
	try {
		const urlObj = new URL(url)
		return urlObj.protocol === 'http:' || urlObj.protocol === 'https:'
	} catch (error) {
		console.error('Invalid URL:', url, error)
		return false
	}
}
