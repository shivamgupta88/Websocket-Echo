import React from "react";

const ModeToggle = ({ isNightMode, setIsNightMode }) => {
  return (
    <button
      onClick={() => setIsNightMode(!isNightMode)}
      className={`px-3 py-1 rounded-lg transition-all duration-300 ease-in-out ${
        isNightMode ? "bg-red-500 text-white hover:bg-red-600" : "bg-green-500 text-white hover:bg-green-600"
      }`}
      style={{ padding: '0.25rem 0.75rem' }} // Adjusted inline style for precise control
    >
      Switch to {isNightMode ? "Day" : "Night"} Mode
    </button>
  );
};

export default ModeToggle;