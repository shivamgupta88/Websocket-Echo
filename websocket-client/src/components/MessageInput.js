import React, { useState } from "react";

const MessageInput = ({ sendMessage, isNightMode }) => {
  const [inputValue, setInputValue] = useState("");

  const handleSendMessage = () => {
    if (inputValue.trim() !== "") {
      sendMessage(inputValue);
      setInputValue("");
    }
  };

  return (
    <div className="flex space-x-2 mb-4">
      <input
        type="text"
        value={inputValue}
        onChange={(e) => setInputValue(e.target.value)}
        placeholder="Type a message"
        className={`flex-1 p-3 border rounded-lg focus:outline-none focus:ring-2 ${isNightMode ? 'border-gray-700 bg-gray-700 text-white focus:ring-blue-400' : 'border-gray-300 bg-white text-gray-800 focus:ring-blue-400'}`}
      />
      <button
        onClick={handleSendMessage}
        className={`px-4 py-3 rounded-lg transition ${isNightMode ? 'bg-blue-500 text-white hover:bg-blue-600' : 'bg-blue-500 text-white hover:bg-blue-600'}`}
      >
        Send
      </button>
    </div>
  );
};

export default MessageInput;