import { useState, useEffect } from 'react';

const API = "http://a2007186fbf734e00bf5ab798a4a0cc0-6175012.us-east-2.elb.amazonaws.com";

function App() {
  const [message, setMessage] = useState("");
  const [results, setResults] = useState([]);

  useEffect(() => {
    const timer = setInterval(async () => {
      const res = await fetch(`${API}/poll`);
      if (res.status === 200) {
        const task = await res.json();
        setResults(prev => [...prev, task]);
      }
    }, 50);

    return () => clearInterval(timer);
  }, []);

  const handleGo = async () => {
    await fetch(`${API}/tasks`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ message: message })
    });
  };

  const workers = [...new Set(results.map(r => r.workerID))];

  return (
    <div>
      <input 
        value={message} 
        onChange={e => setMessage(e.target.value)} 
        style={{ fontSize: '24px', padding: '12px', width: '400px' }}
      />
      <button onClick={handleGo} style={{ fontSize: '24px', padding: '12px 32px' }}>
        print
      </button>
      <div style={{ display: 'flex', gap: '24px' }}>
        {workers.map(w =>
          <div key={w} style={{ flex: 1 }}>
            <h3>{w}</h3>
            {results.filter(r => r.workerID === w).map((r, i) =>
              <p key={i}>ticket {r.ticketID}: {r.message}</p>)}
          </div>)}
      </div>
    </div>
  );
}

export default App;