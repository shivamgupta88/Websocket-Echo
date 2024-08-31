import React, { useState } from "react";
import useWebSocket, { ReadyState } from "react-use-websocket";

function App() {
  const [messageHistory, setMessageHistory] = useState([]);
  const [inputValue, setInputValue] = useState("");
  const [lastFiveMessages, setLastFiveMessages] = useState([]);

  const { sendMessage, readyState } = useWebSocket("ws://localhost:8080/ws", {
    onOpen: () => console.log("Connected to WebSocket"),
    onMessage: (msg) => {
      setMessageHistory((prev) => [...prev, msg.data]);
    },
    shouldReconnect: (closeEvent) => true,
  });

  const connectionStatus = {
    [ReadyState.CONNECTING]: "Connecting",
    [ReadyState.OPEN]: "Open",
    [ReadyState.CLOSING]: "Closing",
    [ReadyState.CLOSED]: "Closed",
    [ReadyState.UNINSTANTIATED]: "Uninstantiated",
  }[readyState];

  const sendInputMessage = () => {
    if (inputValue.trim() !== "") {
      sendMessage(inputValue);
      setInputValue("");
    }
  };

  const handleGetLastFiveMessages = async () => {
    try {
      const res = await fetch("http://localhost:8080/history");
      if (!res.ok) {
        throw new Error("Failed to fetch history");
      }
      const data = await res.json();
      setLastFiveMessages(data); // Assuming the server sends JSON with last 5 messages
    } catch (error) {
      console.error("Error fetching last 5 messages:", error);
    }
  };

  return (
    <div className="flex flex-col items-center p-6 bg-gray-100 min-h-screen">
      <h1 className="text-2xl font-bold mb-4">WebSocket Client</h1>
      <p className="text-lg mb-4">Status: <span className={`font-semibold ${readyState === ReadyState.OPEN ? 'text-green-500' : 'text-red-500'}`}>{connectionStatus}</span></p>

      <div className="flex space-x-2 mb-4">
        <input
          type="text"
          value={inputValue}
          onChange={(e) => setInputValue(e.target.value)}
          placeholder="Type a message"
          className="p-2 border border-gray-300 rounded"
        />
        <button
          onClick={sendInputMessage}
          className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
        >
          Send
        </button>
      </div>

      <div className="w-full max-w-md bg-white p-4 shadow rounded">
        <h2 className="text-xl font-semibold mb-2">Messages</h2>
        <div className="space-y-2 mb-4">
          {messageHistory.map((message, idx) => (
            <p key={idx} className="p-2 bg-gray-200 rounded">{message}</p>
          ))}
        </div>

        <button
          onClick={handleGetLastFiveMessages}
          className="px-4 py-2 bg-gray-500 text-white rounded hover:bg-gray-600"
        >
          Get Last 5 Messages
        </button>

        <h2 className="text-xl font-semibold mt-4 mb-2">Last 5 Messages</h2>
        <div className="space-y-2">
          {lastFiveMessages.map((msg, idx) => (
            <p key={idx} className="p-2 bg-gray-200 rounded">{msg}</p>
          ))}
        </div>
      </div>
    </div>
  );
}

export default App;
