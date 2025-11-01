
import './App.css'
import { Link, Route, Routes } from 'react-router-dom'
import ProductList from './components/productList'

function App() {

  return (
    <>
    <div>
      <nav style={{ padding: "1rem", background: "#eee" }}>
        <Link to = "/products">Products</Link>
        {/* <Link to = "/cart">Cart</Link> */}
      </nav>
      <Routes>
        <Route path="/products" element={<ProductList />} />
      </Routes>
    </div>    
  </>
       
  )
}

export default App
