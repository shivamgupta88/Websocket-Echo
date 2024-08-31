import React from "react";

const MessageList = ({ messages, isNightMode }) => {
  return (
    <div className="flex-1 overflow-y-auto mb-4">
      <h2 className={`text-xl font-semibold mb-2 ${isNightMode ? 'text-white' : 'text-gray-800'}`}>Messages</h2>
      <div className="space-y-2">
        {messages.map((message, idx) => (
          <div key={idx} className={`p-3 rounded-lg text-gray-800 ${isNightMode ? 'bg-blue-600' : 'bg-blue-100'}`}>
            {message}
          </div>
        ))}
      </div>
    </div>
  );
};

export default MessageList;