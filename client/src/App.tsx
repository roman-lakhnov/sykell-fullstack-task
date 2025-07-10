import 'bootstrap/dist/css/bootstrap.min.css'
import { ToastContainer } from 'react-toastify'
import 'react-toastify/dist/ReactToastify.css' // Make sure styles are imported
import UrlManagement from './components/UrlManagement'
import ResultsDashboard from './components/ResultsDashboard'

function App() {
	return (
		<div className='App'>
			<UrlManagement />
			<ResultsDashboard />
			<ToastContainer
				limit={2} // Limit the maximum number of toasts to 2
				position='top-right'
				autoClose={3000}
				hideProgressBar={false}
				newestOnTop={true}
				closeOnClick
				rtl={false}
				pauseOnFocusLoss
				draggable
				pauseOnHover
			/>
		</div>
	)
}

export default App
