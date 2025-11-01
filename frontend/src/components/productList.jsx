import React, { useEffect } from "react";

export default function ProductList() {

    const [products, setProducts] = React.useState([]);
    const [cart, setCart] = React.useState([]);
    const [showCheckout, setShowCheckout] = React.useState(false);
    const [sessionId, setSessionId] = React.useState(() => {
        const existing = localStorage.getItem("sessionId");
        if (existing) return existing;
        const newId = Math.floor(Math.random() * 1000000).toString();
        localStorage.setItem("sessionId", newId);
        return newId;
      });
      

    useEffect(() => {
        fetch('http://localhost:8080/products').then(res => res.json()).then(setProducts).catch(console.error); 
    },[])

    const total = cart.reduce(
        (sum, item) => sum + (item.priceCents / 100) * item.quantity,
        0
      );
      

    const handleCheckout = () => {
        setShowCheckout(true);
      };
    
      const closeModal = () => {
        setShowCheckout(false);
        clearCart();
      };

      const clearCart = () => {
        fetch(`http://localhost:8080/cart/${sessionId}`, {
          method: "DELETE",
        }).then(fetchCart).catch(console.error);
      };
    const fetchCart = () => {
        fetch(`http://localhost:8080/cart/${sessionId}`)
          .then(res => res.json())
          .then(data=>setCart(data || []))
          .catch(console.error);
      };

    // const fetchProducts = () => {
    //     fetch('http://localhost:8080/products').then(res => res.json()).then(setProducts).catch(console.error); 
    // }

    const addToCart = (product) => {
        fetch("http://localhost:8080/cart/items", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
              sessionid: sessionId,
              productid: product.id,
              quantity: 1,
            }),
          }).then(fetchCart).catch(console.error);
    }
      
    return (
        <div>
            <h2>ProductList</h2>   
            <div>     
            <ul>
                {products.map(product => (
                    <li key={product.id}>
                        <strong>{product.name}</strong>
                            &nbsp;&nbsp; <span>{product.description}</span> &nbsp;&nbsp; <span>${(product.priceCents/100)}</span>&nbsp;&nbsp;
                        <button onClick={() => addToCart(product)}>Add to cart</button> 
                    </li>
                ))}
            </ul>
            </div>
            <div>
            <h3>Cart</h3>
                <ul>
                    {cart.length === 0 ?(
                        <li>Cart is empty</li>
                    ):(
                        cart.map(item=>(<li key={item.productId}>{item.productName} x {item.quantity} -- $ {(item.priceCents/100 * item.quantity)}</li>)))}
                </ul>
            </div>
            {cart.length > 0 && (
        <button onClick={handleCheckout}>Checkout</button>
      )}
           {/* MODAL */}
      {showCheckout && (
        <div
          style={{
            position: "fixed",
            top: 0,
            left: 0,
            width: "100vw",
            height: "100vh",
            backgroundColor: "rgba(0,0,0,0.6)",
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
          }}
          onClick={closeModal}
        >
          <div
            style={{
              background: "#000",
              padding: "20px",
              borderRadius: "10px",
              width: "400px",
              textAlign: "center",
            }}
            onClick={(e) => e.stopPropagation()}
          >
            <h3>Order Placed 🎉</h3>
            <p>Here is your order summary:</p>
            <ul style={{ textAlign: "left" }}>
              {cart.map((item) => (
                <li key={item.productId}>
                  {item.productName} x {item.quantity} — $
                  {(item.priceCents / 100) * item.quantity}
                </li>
              ))}
            </ul>
            <p>
              <strong>Total: ${total.toFixed(2)}</strong>
            </p>
            <button onClick={closeModal}>Close</button>
          </div>
        </div>
      )}
        </div>
    );
}