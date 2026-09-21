import { useState, useEffect } from 'react';

function App() {
  const [message, setMessage] = useState("");
  const [results, setResults] = useState([]);

  useEffect(() => {
    const timer = setInterval(async () => {
      const res = await fetch("http://a2007186fbf734e00bf5ab798a4a0cc0-6175012.us-east-2.elb.amazonaws.com/poll");
      if (res.status === 200) {
        const task = await res.json();
        setResults(prev => [...prev, task]);
      }
    }, 50);

    return () => clearInterval(timer);
  }, []);

  const handleGo = async () => {
  await fetch("http://a2007186fbf734e00bf5ab798a4a0cc0-6175012.us-east-2.elb.amazonaws.com/tasks", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ message: message })
  });
};

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
      <div>
        {results.map((r, i) => 
          <p key={i}>
            worker: {r.workerID} | ticket: {r.ticketID} | message: {r.message}
          </p>)}
      </div>
    </div>
  );
}

export default App;