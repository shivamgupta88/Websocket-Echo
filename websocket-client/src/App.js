import React, { useState } from "react";
import useWebSocket, { ReadyState } from "react-use-websocket";
import MessageList from "./components/MessageList";
import MessageInput from "./components/MessageInput";
import ModeToggle from "./components/ModeToggle";

function App() {
  const [messageHistory, setMessageHistory] = useState([]);
  const [lastFiveMessages, setLastFiveMessages] = useState([]);
  const [isNightMode, setIsNightMode] = useState(false);

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

  const handleGetLastFiveMessages = async () => {
    try {
      const res = await fetch("http://localhost:8080/history");
      if (!res.ok) {
        throw new Error("Failed to fetch history");
      }
      const data = await res.json();
      setLastFiveMessages(data);
    } catch (error) {
      console.error("Error fetching last 5 messages:", error);
    }
  };

  return (
    <div className={`flex flex-col items-center p-6 min-h-screen ${isNightMode ? "bg-gray-900" : "bg-gray-100"}`}>
      <header className="flex justify-between items-center w-full max-w-md mb-6">
        <h1 className={`text-3xl font-bold ${isNightMode ? "text-white" : "text-gray-800"}`}>WebSocket Chat</h1>
        <ModeToggle isNightMode={isNightMode} setIsNightMode={setIsNightMode} />
      </header>

      <p className={`text-lg mb-4 ${isNightMode ? "text-gray-400" : "text-gray-700"}`}>
        Status: <span className={`font-semibold ${readyState === ReadyState.OPEN ? "text-green-400" : "text-red-400"}`}>{connectionStatus}</span>
      </p>

      <div className={`w-full max-w-md p-4 shadow-lg rounded-lg flex flex-col ${isNightMode ? "bg-gray-800" : "bg-white"}`}>
        <MessageList messages={messageHistory} isNightMode={isNightMode} />
        <MessageInput sendMessage={sendMessage} isNightMode={isNightMode} />

        <button
          onClick={handleGetLastFiveMessages}
          className={`mt-4 px-4 py-2 rounded-lg transition-all duration-300 ease-in-out ${isNightMode ? "bg-gray-700 text-white hover:bg-gray-600" : "bg-gray-500 text-white hover:bg-gray-600"}`}
        >
          Get Last 5 Messages
        </button>

        <h2 className={`text-xl font-semibold mt-4 mb-2 ${isNightMode ? "text-white" : "text-gray-800"}`}>Last 5 Messages</h2>
        <MessageList messages={lastFiveMessages} isNightMode={isNightMode} />
      </div>
    </div>
  );
}

export default App;