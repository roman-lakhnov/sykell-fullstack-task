import 'bootstrap/dist/css/bootstrap.min.css'
import { ToastContainer } from 'react-toastify'
// import ResultsDashboard from '../../trash/ResultsDashboard'
import UrlManagement from './components/UrlManagement'

function App() {
	return (
		<div className='App'>
			<UrlManagement />
			{/* <ResultsDashboard /> */}
			<ToastContainer />
		</div>
	)
}

export default App
